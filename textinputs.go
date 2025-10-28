// package main

// // A simple program demonstrating the text input component from the Bubbles
// // component library.

// import (
// 	"fmt"
// 	"log"

// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// )

// var textInputValue string

// func GetTextInput() string {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return textInputValue
// }

// type (
// 	errMsg error
// )

// type TextModel struct {
// 	textInput textinput.Model
// 	err       error
// }

// func initialModel() TextModel {
// 	ti := textinput.New()
// 	ti.Placeholder = "London"
// 	ti.Focus()
// 	ti.CharLimit = 156
// 	ti.Width = 20

// 	return TextModel{
// 		textInput: ti,
// 		err:       nil,
// 	}
// }

// func (m TextModel) Init() tea.Cmd {
// 	return textinput.Blink
// }

// func (m TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return m, tea.Quit
// 		case tea.KeyEnter:
// 			textInputValue = m.textInput.Value()
// 			return m, tea.Quit
// 		}

// 	// We handle errors just like any other message
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}

// 	m.textInput, cmd = m.textInput.Update(msg)
// 	return m, cmd
// }

// func (m TextModel) View() string {
// 	return fmt.Sprintf(
// 		"What's your City?\n\n%s\n\n%s",
// 		m.textInput.View(),
// 		"(esc to quit)",
// 	) + "\n"
// }

package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var (
	LocationInput string
	UnitTypeInput string
)

type textModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func InitialTextModel() textModel {
	m := textModel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Width = 100

		switch i {
		case 0:
			t.Placeholder = "Location (e.g., London)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Unit Type (e.g., metric/us/uk)"
			t.CharLimit = 8
		// case 2:
		// 	t.Placeholder = "Password"
		// 	t.EchoMode = textinput.EchoPassword
		// 	t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				LocationInput = m.inputs[0].Value()
				UnitTypeInput = m.inputs[1].Value()
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *textModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m textModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

// func GetTextInput() {
// 	if _, err := tea.NewProgram(InitialTextModel()).Run(); err != nil {
// 		fmt.Printf("could not start program: %s\n", err)
// 		os.Exit(1)
// 	}
// }
