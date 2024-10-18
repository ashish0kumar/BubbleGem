# ♊ GemGo

GemGo is a command-line interface (CLI) for [Google Gemini](https://gemini.google.com/app), built using the [Bubble Tea framework](https://github.com/charmbracelet/bubbletea) in Go. 
With GemGo, you can ask questions and interact with the Gemini AI in your terminal environment.

## Features

- **Command-line interface**: Interact with Gemini AI from your terminal.
- **Bubble Tea framework**: Uses the Bubble Tea TUI framework for a responsive and clean user experience.
- **Loading animation**: Animated loading indicator while waiting for the AI response.

## Installation

### Prerequisites

- Go (version 1.18 or higher)
- Gemini AI API Key (set as an environment variable `GEMINI_API_KEY`)

### Steps

1. Clone the repository:

```bash
git clone https://github.com/ashish0kumar/GemGo.git
cd GemGo
```

2. To use `GemGo`, you'll need an API key set in the `GEMINI_API_KEY` environment variable. If you don't already have one, create a key in [Google AI Studio](https://aistudio.google.com/app/apikey).

4. Set the `GEMINI_API_KEY` environment variable with your API key:

  - For Linux/macOS:

  ```bash
  export GEMINI_API_KEY="your-api-key-here"
  ```

  - For Windows (PowerShell):

  ```powershell
  $env:GEMINI_API_KEY="your-api-key-here"
  ```

3. Build and run the application:

```bash
go build
./GemGo
```

## Usage

Once you run the CLI, you can start typing questions. To submit a query, press Enter. To exit, press `Ctrl+C` or `Esc`.


## Dependencies

- [Charm's Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Fatih Color](https://github.com/fatih/color)

## Contributing

Feel free to open issues and contribute to this project. All contributions are welcome!

## License

This project is licensed under the MIT License.
