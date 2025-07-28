
# 💬 Heroku Go Chat

Heroku Go Chat is a command-line chat client written in Go that lets you interact with Heroku's Claude-4-Sonnet model using a simple and efficient interface. It supports real-time chat, persistent conversation history, tagging, and easy navigation of past interactions. The tool is designed for quick local use or seamless deployment to Heroku, and is ideal for developers, experimenters, and anyone who wants a fast, scriptable AI chat experience.


## 🚀 Getting Started

### 1. Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) (for deployment)

### 2. Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/alexandrespindola/heroku-go-chat.git
cd heroku-go-chat
go mod tidy
```

### 3. Environment & Alias Setup (Heroku Managed Inference and Agents)

Set up the required environment variables and alias for convenient usage. Run the following commands in your terminal:

```bash
echo "alias hc='noglob ~/path/to/herochat'" >> ~/.zshrc
echo 'export INFERENCE_URL=https://your-inference-url.heroku.com' >> ~/.zshrc
echo 'export INFERENCE_KEY=inf-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' >> ~/.zshrc
source ~/.zshrc
```

These commands will:

- Add an alias `hc` to run the chat client easily
- Set the `INFERENCE_URL` and `INFERENCE_KEY` environment variables required for Heroku Managed Inference
- Reload your shell configuration

**Note:**
- Replace `~/path/to/herochat` with the actual path to your `herochat` binary.
- Replace `https://your-inference-url.heroku.com` with your actual INFERENCE_URL from your Heroku app.
- Replace `inf-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` with your actual INFERENCE_KEY from your Heroku app dashboard.


## 🛠️ Main Functionalities & Commands

- **Chat with AI**: Start a new conversation or continue a tagged thread directly from your terminal.
    - Usage: `hc <tag> <prompt>`
    - Example: `hc projectX "Summarize the last meeting notes"`
    - The `<tag>` allows you to group related conversations for context continuity.

- **View History**: List all previous conversations, optionally filtered by tag.
    - Usage: `hc history [tag]`
    - Example: `hc history projectX`

- **Navigate Conversations**: Interactively browse through your conversation history, moving forward and backward, or selecting by ID.
    - Usage: `hc navigate [tag]`
    - Example: `hc navigate projectX`

- **Persistent Storage**: All conversations are saved in `conversations.json` for later review or analysis.

### Example Workflow

1. Start a chat: `hc mytag "How do I deploy a Go app to Heroku?"`
2. View history: `hc history mytag`
3. Navigate interactively: `hc navigate mytag`

See the section below for alias and environment setup to use the `hc` command conveniently.

## ✨ Features

- ⚡ Real-time chat between multiple users
- 🧹 Simple and clean codebase
- 🚀 Easy deployment to Heroku
- 💾 JSON-based conversation history

## 🚀 Getting Started

### 🛠️ Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) (for deployment)

### 📦 Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/alexandrespindola/heroku-go-chat.git
   cd heroku-go-chat
   ```
2. Install dependencies:

   ```bash
   go mod tidy
   ```

### 🏃 Running Locally

```bash
go run main.go
```


### ☁️ Deployment to Heroku

To use this chat client with Heroku, you must:

1. **Install the "Heroku Managed Inference and Agents" add-on** from the [Heroku Marketplace](https://elements.heroku.com/addons/managed-inference-agents) in your Heroku app dashboard.
2. **Select the Claude Sonnet 4 model** ("claude-4-sonnet") as your inference model in the add-on configuration.
3. Create a new Heroku app (if you haven't already):

   ```bash
   heroku create
   ```
4. Deploy your app:

   ```bash
   git push heroku main
   ```
5. Open your app in the browser:

   ```bash
   heroku open
   ```

## 📁 Project Structure

- `main.go` - Main application entry point
- `herochat/` - Core chat logic and handlers
- `conversations.json` - Stores chat history

## 🪪 License

This project is released into the public domain under the [Creative Commons Zero v1.0 Universal (CC0 1.0) Public Domain Dedication](https://creativecommons.org/publicdomain/zero/1.0/). You can copy, modify, distribute and perform the work, even for commercial purposes, all without asking permission.
