package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}

type Model struct {
	textinput textinput.Model
	viewport  viewport.Model
	responses []string
	err       error
	width     int
	height    int
	ready     bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Gemini something..."
	ti.Focus()

	vp := viewport.New(80, 20)
	vp.SetContent("")

	return Model{
		textinput: ti,
		viewport:  vp,
		responses: []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := m.textinput.Value()
			m.textinput.SetValue("")
			return m, handleGeminiQuery(query)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3
		m.updateViewportContent()
		m.ready = true

	case GeminiResponseMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.responses = []string{fmt.Sprintf("Error: %v", msg.Err)}
		} else {
			m.responses = []string{fmt.Sprintf(msg.Response)}
		}
		m.updateViewportContent()
	}

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateViewportContent() {
	content := strings.Join(m.responses, "\n\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	return fmt.Sprintf(
		"%s\n\n%s",
		m.textinput.View(),
		m.viewport.View(),
	)
}

type GeminiResponseMsg struct {
	Response string
	Err      error
}

func handleGeminiQuery(query string) tea.Cmd {
	return func() tea.Msg {
		userQuery := strings.Join([]string{query}, " ")
		apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
		if !ok {
			return GeminiResponseMsg{Err: fmt.Errorf("GEMINI_API_KEY environment variable not set")}
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return GeminiResponseMsg{Err: err}
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-1.5-pro")

		extraPrompt := "Please respond without using any bold or italic text formatting"
		resp, err := model.GenerateContent(ctx, genai.Text(userQuery+extraPrompt))
		if err != nil {
			return GeminiResponseMsg{Err: err}
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return GeminiResponseMsg{Err: fmt.Errorf("no response from Gemini")}
		}

		responseText := ""
		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseText += string(textPart)
			}
		}

		if responseText == "" {
			return GeminiResponseMsg{Err: fmt.Errorf("no text content in the response")}
		}

		return GeminiResponseMsg{Response: responseText}
	}
}
