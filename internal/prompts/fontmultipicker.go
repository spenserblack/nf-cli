package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spenserblack/nf-cli/pkg/fonts"
)

// FPMultiModel is the model for the font multi-picker.
type fpmultimodel struct {
	fonts  []fonts.Font
	cursor int
	done   bool
	quit   bool
	// selected is a set of selected indices.
	selected map[int]struct{}
}

func (m fpmultimodel) Init() tea.Cmd {
	return nil
}

func (m fpmultimodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			m.done = true
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
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m fpmultimodel) View() string {
	if m.done {
		return ""
	}

	view := "select the font(s) you want to install (use arrow keys or j/k to navigate, Space to select):\n\n"

	// NOTE We don't show all fonts because the list could be very long.
	lowerBound := m.cursor - 5
	upperBound := m.cursor + 5
	// NOTE Theoretically there can be some awkward behavior if the list is very small,
	//		but that shouldn't happen in practice.
	if lowerBound < 0 {
		upperBound += -lowerBound
		lowerBound = 0
	}
	if upperBound > len(m.fonts) {
		lowerBound -= upperBound - len(m.fonts)
		upperBound = len(m.fonts)
	}

	for i := lowerBound; i < upperBound; i++ {
		font := m.fonts[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		check := " "
		if _, ok := m.selected[i]; ok {
			check = "O"
		}

		view += fmt.Sprintf(
			"%s %s %s (%s)\n",
			cursor,
			check,
			font.PatchedName,
			font.UnpatchedName,
		)
	}

	view += "\nPress q to quit, Enter to confirm.\n"

	return view
}

// MultiPromptForFonts prompts the user to select one or more fonts from the given choices.
func MultiPromptForFonts(choices []fonts.Font) (selected []fonts.Font, err error) {
	m := fpmultimodel{
		fonts:    choices,
		cursor:   0,
		done:     false,
		selected: make(map[int]struct{}, len(choices)),
	}
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := res.(fpmultimodel)

	if result.quit || len(result.selected) == 0 {
		return nil, ErrNoFontSelected
	}

	selected = make([]fonts.Font, 0, len(result.selected))
	for i := range result.selected {
		selected = append(selected, result.fonts[i])
	}
	return selected, nil
}
