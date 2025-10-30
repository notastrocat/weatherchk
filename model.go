package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

    "github.com/nitishm/go-rejson/v4"
)

// Model represents the application state
type Model struct {
	spinner   spinner.Model
	client    *Client
	result    map[string]interface{}
	timeTaken time.Duration
	loading   bool
	err       error
	quitting  bool
	isCached  bool
}

// weatherMsg is used to deliver weather data
type weatherMsg struct {
	result    map[string]interface{}
	timeTaken time.Duration
	err       error
}

var rejsonClient *rejson.Handler

func InitialModel(client *Client, rh *rejson.Handler) Model {
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	rejsonClient = rh
	return Model{
		spinner: s,
		client:  client,
		loading: false,
	}
}

func (m Model) Init() tea.Cmd {
	m.loading = true
	// Start fetching weather data immediately
	return tea.Batch(
		m.spinner.Tick,
		fetchWeather(m.client, m),
	)
}

// fetchWeather returns a command that fetches weather data
func fetchWeather(client *Client, m Model) tea.Cmd {
	return func() tea.Msg {
		result, timeTaken, err := client.GetCurrentWeather()
		m.loading = false
		return weatherMsg{result: result, timeTaken: timeTaken, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case weatherMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.result = msg.result
		m.timeTaken = msg.timeTaken

		// goofy error handling here.
		var err error
		m.isCached, err = SetWeatherData(rejsonClient, LocationInput, m.result)
		if err != nil {
			// do something about it
			// m.err = err
		}
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	if m.err != nil {
		return fmt.Sprintf("❌ Error: %v\n", m.err.Error())
	}

	if m.loading {
		return m.spinner.View() + " Fetching weather data...\n\n"
	}

	if m.result != "" {
		timeStr   := fmt.Sprintf("\nTook %v\n\n", m.timeTaken.String())

		return m.result + helpStyle.Render(timeStr)
	}

	return m.spinner.View() + " Starting...\n\n"
}
