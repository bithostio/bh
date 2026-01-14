package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	Green  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
	Red    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render
	Cyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render
	Gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render
	Bold   = lipgloss.NewStyle().Bold(true).Render
)

func PrintTable(rows [][]string) {
	if len(rows) == 0 {
		return
	}

	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			if w := lipgloss.Width(col); i < len(widths) && w > widths[i] {
				widths[i] = w
			}
		}
	}

	columns := make([]table.Column, len(rows[0]))
	for i, title := range rows[0] {
		columns[i] = table.Column{Title: title, Width: widths[i]}
	}

	tableRows := make([]table.Row, 0, len(rows)-1)
	for _, row := range rows[1:] {
		tableRows = append(tableRows, table.Row(row))
	}

	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().Bold(true).Padding(0, 1, 0, 0)
	s.Cell = lipgloss.NewStyle().Padding(0, 1, 0, 0)
	s.Selected = s.Cell

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(false),
		table.WithStyles(s),
	)
	fmt.Println(t.View())
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func FormatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func FormatMemory(v string) string  { return formatWithUnit(v, "MB") }
func FormatStorage(v string) string { return formatWithUnit(v, "GB") }

func formatWithUnit(v, unit string) string {
	u := strings.ToUpper(strings.TrimSpace(v))
	for _, s := range []string{"B", "KB", "MB", "GB", "TB", "PB"} {
		if strings.HasSuffix(u, s) {
			return v
		}
	}
	return v + " " + unit
}

func ShowProgress(message string) (stop, cancel func()) {
	done, clear := make(chan bool), make(chan bool)

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		for i := 0; ; i = (i + 1) % len(frames) {
			select {
			case <-done:
				fmt.Printf("\r%s %s\n", Green("✓"), message)
				return
			case <-clear:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r%s %s", frames[i], message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	wait := func(ch chan bool) func() {
		return func() { ch <- true; time.Sleep(200 * time.Millisecond) }
	}
	return wait(done), wait(clear)
}

type Paginator interface {
	GetPage() int
	HasPrevious() bool
	HasNext() bool
}

func DisplayPagination(p Paginator) {
	if p.GetPage() == 1 && !p.HasNext() {
		return
	}

	var hints []string
	if p.HasPrevious() {
		hints = append(hints, Gray(fmt.Sprintf("--page %d for previous", p.GetPage()-1)))
	}
	if p.HasNext() {
		hints = append(hints, Gray(fmt.Sprintf("--page %d for next", p.GetPage()+1)))
	}

	fmt.Println()
	if len(hints) > 0 {
		fmt.Printf("Page %d (%s)\n", p.GetPage(), strings.Join(hints, ", "))
	} else {
		fmt.Printf("Page %d\n", p.GetPage())
	}
}
