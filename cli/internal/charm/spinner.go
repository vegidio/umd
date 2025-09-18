package charm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vegidio/umd"
)

type queryDoneMsg struct{}

func waitForDone(doneCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		<-doneCh
		return queryDoneMsg{}
	}
}

type spinnerModel struct {
	spinner  spinner.Model
	typeName string
	name     string
	resp     *umd.Response
}

func initSpinnerModel(typeName string, name string, resp *umd.Response) *spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = pink

	return &spinnerModel{spinner: s, typeName: typeName, name: name, resp: resp}
}

func (m *spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForDone(m.resp.Done),
	)
}

func (m *spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgValue := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case queryDoneMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch msgValue.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *spinnerModel) View() string {
	typeName := strings.ToLower(m.typeName)
	numberStr := fmt.Sprintf("%d", len(m.resp.Media))

	return fmt.Sprintf("📝 Querying %s %s for media files... %s%s Found\n",
		typeName, bold.Render(m.name), m.spinner.View(), bold.Render(numberStr))
}

func StartSpinner(typeName string, name string, resp *umd.Response) error {
	_, err := tea.NewProgram(initSpinnerModel(typeName, name, resp)).Run()
	if err != nil {
		return err
	}

	return nil
}
