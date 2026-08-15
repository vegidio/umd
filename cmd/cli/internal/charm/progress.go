package charm

import (
	"fmt"
	"shared"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
	saktime "github.com/vegidio/go-sak/time"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type downloadMsg struct {
	resp *fetch.Response
}

type downloadDone struct{}

// downloadFinished signals that a single download has run to completion, either successfully or with
// an error.
type downloadFinished struct{}

func downloadCmd(ch <-chan *fetch.Response) tea.Cmd {
	return func() tea.Msg {
		if resp, ok := <-ch; ok {
			return downloadMsg{resp}
		}

		return downloadDone{}
	}
}

// waitForFinish blocks until the download is over and reports it back to the model. It must not rely
// on the progress reaching 100%, because failed downloads never get there.
func waitForFinish(resp *fetch.Response) tea.Cmd {
	return func() tea.Msg {
		_ = resp.Error()
		return downloadFinished{}
	}
}

type progressModel struct {
	progress      progress.Model
	result        <-chan *fetch.Response
	responses     []*fetch.Response
	queue         *shared.Queue
	total         int
	completed     int
	percent       float64
	done          bool
	fullBarShown  bool
	startTime     time.Time
	lastEtaUpdate time.Time
	eta           time.Duration
}

func (m *progressModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		downloadCmd(m.result),
	)
}

func (m *progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgValue := msg.(type) {
	case tickMsg:
		if m.fullBarShown {
			return m, tea.Quit
		}

		m.percent = 1
		if m.total > 0 {
			m.percent = float64(m.completed) / float64(m.total)
		}

		// Every download that was handed to us has run to completion; note that failed downloads count
		// as finished too, otherwise the bar would never reach 100% and we would never quit. Give the
		// finished bar one frame on screen before quitting.
		m.fullBarShown = m.done && m.completed >= len(m.responses)

		return m, tickCmd()

	case downloadMsg:
		m.queue.Add(msgValue.resp)
		m.responses = append(m.responses, msgValue.resp)
		return m, tea.Batch(downloadCmd(m.result), waitForFinish(msgValue.resp))

	case downloadFinished:
		m.completed++
		return m, nil

	case downloadDone:
		m.done = true
		m.eta = time.Duration(0)
		return m, nil

	case tea.KeyMsg:
		switch msgValue.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *progressModel) View() string {
	width := len(strconv.Itoa(m.total))

	percent := m.percent * 100
	intPart := int(percent)
	fracPart := int(percent*10) % 10
	percentStr := fmt.Sprintf("%3d.%1d%%", intPart, fracPart)

	now := time.Now()
	if now.Sub(m.lastEtaUpdate) >= time.Second {
		elapsed := time.Since(m.startTime)
		m.eta = saktime.CalculateEta(m.total, m.completed, elapsed)
		m.lastEtaUpdate = now
	}

	c := bold.Render(fmt.Sprintf("%0*d", width, m.completed))
	t := bold.Render(fmt.Sprintf("%d", m.total))

	eta := m.eta.Truncate(time.Second)
	if m.eta < 10*time.Second {
		eta = m.eta.Truncate(time.Second / 10)
	}

	return fmt.Sprintf("\nDownloading   %s%s%s%s%s  %s  %s   %s\n%s\n",
		gray.Render("["), c, gray.Render("/"), t, gray.Render("]"),
		m.progress.ViewAs(m.percent),
		green.Render(percentStr),
		magenta.Render(fmt.Sprintf("ETA %v", eta)),

		printLastFive(m.queue.Items()),
	)
}

func printLastFive(downloads []*fetch.Response) string {
	return lo.Reduce(downloads, func(acc string, r *fetch.Response, _ int) string {
		mType := shared.GetMediaType(r.Request.FilePath)
		index := fmt.Sprintf("%3.0f%%", r.Progress*100)

		var prefix string
		if mType == "video" {
			prefix = fmt.Sprintf("🎬 [%s] Downloading %s", blue.Render(index), orange.Render(mType))
		} else if mType == "image" {
			prefix = fmt.Sprintf("📸 [%s] Downloading %s", blue.Render(index), magenta.Render(mType))
		} else {
			prefix = fmt.Sprintf("✖️ [%s] Downloading %s", blue.Render(index), gray.Render(mType))
		}

		return acc + fmt.Sprintf("\n%s %s",
			prefix,
			cyanUnderline.Render(r.Request.Url),
		)
	}, "")
}

func initProgressModel(result <-chan *fetch.Response, total int) *progressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithWidth(50),
	)

	return &progressModel{
		progress:      p,
		result:        result,
		responses:     make([]*fetch.Response, 0),
		queue:         shared.NewQueue(5),
		total:         total,
		startTime:     time.Now(),
		lastEtaUpdate: time.Now(),
		eta:           time.Duration(0),
	}
}

func StartProgress(result <-chan *fetch.Response, total int) ([]*fetch.Response, error) {
	model, err := tea.NewProgram(initProgressModel(result, total)).Run()
	if err != nil {
		return nil, err
	}

	m := model.(*progressModel)
	return m.responses, nil
}
