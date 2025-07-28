package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Conversation represents a chat interaction
type Conversation struct {
	ID        int       `json:"conversation_id"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	Timestamp string    `json:"timestamp"`
	Tag       string    `json:"tag,omitempty"` // Optional tag field
}

// Message defines the structure for API messages
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool defines an MCP tool
type Tool struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// RequestPayload is the structure for Heroku API requests
type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
}

// Choice contains the model's response
type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

// ResponseData holds the API response
type ResponseData struct {
	Choices []Choice `json:"choices"`
}

// loadHistory reads the conversation history from conversations.json
func loadHistory() ([]Conversation, error) {
	if _, err := os.Stat("conversations.json"); os.IsNotExist(err) {
		return []Conversation{}, nil
	}
	data, err := os.ReadFile("conversations.json")
	if err != nil {
		return nil, err
	}
	var history []Conversation
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// saveConversation saves a new conversation to conversations.json
func saveConversation(prompt, response, tag string) error {
	history, err := loadHistory()
	if err != nil {
		return err
	}
	history = append(history, Conversation{
		ID:        len(history) + 1,
		Prompt:    prompt,
		Response:  response,
		Timestamp: time.Now().Format(time.RFC3339),
		Tag:       tag,
	})
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("conversations.json", data, 0644)
}

// callHeroku sends a prompt to the Heroku API
func callHeroku(prompt, tag string) (string, error) {
	inferenceURL := os.Getenv("INFERENCE_URL")
	if inferenceURL == "" {
		inferenceURL = "https://eu.inference.heroku.com"
	}
	inferenceKey := os.Getenv("INFERENCE_KEY")
	if inferenceKey == "" {
		return "", fmt.Errorf("INFERENCE_KEY not configured")
	}

	// Load history for the given tag
	history, err := loadHistory()
	if err != nil {
		return "", fmt.Errorf("failed to load history: %v", err)
	}
	var messages []Message
	if tag != "" {
		for _, conv := range history {
			if conv.Tag == tag {
				messages = append(messages,
					Message{Role: "user", Content: conv.Prompt},
					Message{Role: "assistant", Content: conv.Response},
				)
			}
		}
	}
	messages = append(messages, Message{Role: "user", Content: prompt})

	url := inferenceURL + "/v1/agents/heroku"
	payload := RequestPayload{
		Model:    "claude-4-sonnet",
		Messages: messages,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to create payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+inferenceKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-Proto", "https")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("response status %d: %s", resp.StatusCode, string(body))
	}

	var fullResponse strings.Builder
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read stream: %v", err)
		}
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			if data == "[DONE]" {
				break
			}
			var responseData ResponseData
			if err := json.Unmarshal([]byte(data), &responseData); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing line: %v\n", err)
				continue
			}
			if len(responseData.Choices) > 0 && responseData.Choices[0].Message.Content != "" {
				fullResponse.WriteString(responseData.Choices[0].Message.Content)
			}
		}
	}
	if fullResponse.String() == "" {
		return "", fmt.Errorf("empty response from model; check prompt or add-on configuration")
	}
	return fullResponse.String(), nil
}

// viewHistory displays conversations, optionally filtered by tag
func viewHistory(tag string) error {
	history, err := loadHistory()
	if err != nil {
		return err
	}
	if len(history) == 0 {
		color.Yellow("⚠️ No history found.")
		return nil
	}
	found := false
	for _, conv := range history {
		if tag == "" || conv.Tag == tag {
			color.Cyan("📜 Conversation %d (%s) [Tag: %s]:", conv.ID, conv.Timestamp, conv.Tag)
			fmt.Printf("  Prompt: %s\n", conv.Prompt)
			fmt.Printf("  Response: %s\n\n", conv.Response)
			found = true
		}
	}
	if !found && tag != "" {
		color.Yellow("⚠️ No conversations found with tag '%s'.", tag)
	}
	return nil
}

// navigateConversations allows interactive navigation through conversations
func navigateConversations(tag string) {
	history, err := loadHistory()
	if err != nil {
		color.Red("❌ Error loading history: %v", err)
		return
	}
	var filteredHistory []Conversation
	if tag == "" {
		filteredHistory = history
	} else {
		for _, conv := range history {
			if conv.Tag == tag {
				filteredHistory = append(filteredHistory, conv)
			}
		}
	}
	if len(filteredHistory) == 0 {
		if tag == "" {
			color.Yellow("⚠️ No history found.")
		} else {
			color.Yellow("⚠️ No conversations found with tag '%s'.", tag)
		}
		return
	}

	currentIndex := len(filteredHistory) - 1
	color.Green("🔍 Use 'next', 'previous', 'select <ID>', or 'back' to exit navigation.")
	for {
		if currentIndex >= 0 && currentIndex < len(filteredHistory) {
			conv := filteredHistory[currentIndex]
			color.Cyan("\n📜 Current Conversation %d (%s) [Tag: %s]:", conv.ID, conv.Timestamp, conv.Tag)
			fmt.Printf("  Prompt: %s\n", conv.Prompt)
			fmt.Printf("  Response: %s\n", conv.Response)
		}
		fmt.Print(color.MagentaString("Navigate (next/previous/select <ID>/back): "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch {
		case input == "back":
			return
		case input == "next":
			if currentIndex < len(filteredHistory)-1 {
				currentIndex++
			} else {
				color.Yellow("⚠️ You are at the most recent conversation.")
			}
		case input == "previous":
			if currentIndex > 0 {
				currentIndex--
			} else {
				color.Yellow("⚠️ You are at the first conversation.")
			}
		case strings.HasPrefix(input, "select "):
			idStr := strings.TrimPrefix(input, "select ")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				color.Red("❌ Invalid ID. Use a valid number.")
				continue
			}
			found := false
			for i, conv := range filteredHistory {
				if conv.ID == id {
					currentIndex = i
					found = true
					break
				}
			}
			if !found {
				color.Red("❌ Invalid ID for the selected tag. Use a number from the listed IDs.")
			}
		default:
			color.Red("❌ Invalid command. Use 'next', 'previous', 'select <ID>', or 'back'.")
		}
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "herochat <tag> <prompt>",
		Short: "A CLI to chat with Heroku's Claude-4-Sonnet model",
		Args:  cobra.MinimumNArgs(2), // Require tag and at least one word for prompt
		Run: func(cmd *cobra.Command, args []string) {
			tag := args[0]
			prompt := strings.Join(args[1:], " ") // Join all args after tag as prompt
			response, err := callHeroku(prompt, tag)
			if err != nil {
				color.Red("❌ Error: %v", err)
				return
			}
			color.Green("✅ Response: %s", response)
			if err := saveConversation(prompt, response, tag); err != nil {
				color.Red("❌ Error saving conversation: %v", err)
			} else {
				color.Green("✅ Conversation saved in conversations.json with tag '%s'", tag)
			}
		},
	}

	var historyCmd = &cobra.Command{
		Use:   "history [tag]",
		Short: "View conversation history, optionally filtered by tag",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := ""
			if len(args) > 0 {
				tag = args[0]
			}
			if err := viewHistory(tag); err != nil {
				color.Red("❌ Error displaying history: %v", err)
			}
		},
	}

	var navigateCmd = &cobra.Command{
		Use:   "navigate [tag]",
		Short: "Navigate through conversation history, optionally filtered by tag",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := ""
			if len(args) > 0 {
				tag = args[0]
			}
			navigateConversations(tag)
		},
	}

	rootCmd.AddCommand(historyCmd, navigateCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Red("❌ err: %v", err)
		os.Exit(1)
	}
}