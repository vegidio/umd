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

func downloadCmd(ch <-chan *fetch.Response) tea.Cmd {
	return func() tea.Msg {
		if resp, ok := <-ch; ok {
			return downloadMsg{resp}
		}

		return downloadDone{}
	}
}

type progressModel struct {
	progress      progress.Model
	result        <-chan *fetch.Response
	responses     []*fetch.Response
	queue         *shared.Queue
	total         int
	completed     int
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
		if m.progress.Percent() >= 1 && !m.progress.IsAnimating() && m.queue.Incompleted() == 0 {
			return m, tea.Quit
		}

		var percent float64
		if m.total > 0 {
			percent = float64(m.completed) / float64(m.total)
		}

		barCmd := m.progress.SetPercent(percent)
		return m, tea.Batch(tickCmd(), barCmd)

	case downloadMsg:
		m.queue.Add(msgValue.resp)

		go func() {
			_ = msgValue.resp.Track(func(_, _ int64, progress float64) {
				if progress >= 1 {
					m.completed++
				}
			})
		}()

		m.responses = append(m.responses, msgValue.resp)
		return m, downloadCmd(m.result)

	case downloadDone:
		barCmd := m.progress.SetPercent(1)
		m.eta = time.Duration(0)
		return m, barCmd

	case progress.FrameMsg:
		updated, cmd := m.progress.Update(msg)
		m.progress = updated.(progress.Model)
		return m, cmd

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

	percent := m.progress.Percent() * 100
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
		m.progress.View(),
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
