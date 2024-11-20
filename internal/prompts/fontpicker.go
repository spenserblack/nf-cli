package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spenserblack/nf-cli/pkg/fonts"
)

// FPModel is the model for the font picker.
type fpmodel struct {
	fonts    []fonts.Font
	cursor   int
	quit     bool
	selected bool
}

func (m fpmodel) Init() tea.Cmd {
	return nil
}

func (m fpmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.fonts) - 1
			}
		case "down", "j":
			if m.cursor < len(m.fonts)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			m.quit = true
			// NOTE cursor != -1 if user has moved the cursor at least once
			m.selected = m.cursor != -1
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m fpmodel) View() string {
	if m.quit {
		return ""
	}

	view := "Select a font (use arrow keys or j/k to navigate, Enter to select):\n\n"

	for i, font := range m.fonts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s (%s)\n", cursor, font.PatchedName, font.UnpatchedName)
	}

	view += "\nPress q to quit.\n"

	return view
}

// PromptForFont prompts the user to select a font from the given choices.
func PromptForFont(choices []fonts.Font) (font fonts.Font, err error) {
	m := fpmodel{
		fonts:    choices,
		cursor:   -1,
		quit:     false,
		selected: false,
	}
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return font, err
	}

	result := res.(fpmodel)

	if result.selected {
		return result.fonts[result.cursor], nil
	}

	return font, ErrNoFontSelected
}
