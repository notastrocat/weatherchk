package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	greetingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Italic(true).Background(lipgloss.Color("#7D56F4"))
	ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#adfabb"))
	WarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffdf12"))
	InfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))
)

func main() {
	colorTerm, ok := os.LookupEnv("COLORTERM")
	if !ok || colorTerm != "truecolor" {
		fmt.Println(WarnStyle.Render("⚠ COLORTERM variable not set to 'truecolor'. Please set it for proper color rendering."))

		os.Setenv("COLORTERM", "truecolor")
	}

	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		// log in Red
		fmt.Println(ErrStyle.Render("❌ API_KEY variable not set"))
		os.Exit(1)
	}

	fmt.Println(greetingStyle.Render("Welcome to Weatherchk!"))
	fmt.Println()

	client := WeatherClient(apiKey)
	redisClient := Connect()
	
    // Ping the Redis server
    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ Couldn't connect to Redis: %v", err)))
    }
    fmt.Println(SuccessStyle.Render("✔ Connected to Redis successfully!"))

	// First, get user input
	if _, err := tea.NewProgram(InitialTextModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	// Then show spinner while fetching weather
	p := tea.NewProgram(InitialModel(client, redisClient))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %s\n", err)
		os.Exit(1)
	}

	var startTime = time.Now()
	GetVal(redisClient, LocationInput)
	var timeTaken = time.Since(startTime)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	fmt.Println(helpStyle.Render(fmt.Sprintf("Took %v (cached)\n", timeTaken.String())))
}
