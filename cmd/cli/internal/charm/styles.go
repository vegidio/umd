package charm

import (
	"github.com/charmbracelet/lipgloss"
)

var bold = lipgloss.NewStyle().Bold(true)
var underline = lipgloss.NewStyle().Underline(true)

var blue = lipgloss.NewStyle().Foreground(lipgloss.Color("#82aaff"))
var cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("#88ddff"))
var gray = lipgloss.NewStyle().Foreground(lipgloss.Color("#686868"))
var green = lipgloss.NewStyle().Foreground(lipgloss.Color("#00c202"))
var magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#c792e9"))
var orange = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcb6b"))
var pink = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5faf"))
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("#f44336"))
var yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#e5db08"))

var cyanUnderline = cyan.Inherit(underline)
var yellowBold = yellow.Inherit(bold)
var yellowUnderline = yellow.Inherit(underline)
