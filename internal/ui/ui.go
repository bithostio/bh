package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Gray   = color.New(color.FgHiBlack).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)

// PrintTable prints data in aligned columns
// First row is treated as headers (bolded)
func PrintTable(rows [][]string) {
	if len(rows) == 0 {
		return
	}

	widths := calculateWidths(rows)

	// Print header
	printRow(rows[0], widths, true)
	printDivider(widths)

	// Print data rows
	for _, row := range rows[1:] {
		printRow(row, widths, false)
	}
	fmt.Println()
}

func calculateWidths(rows [][]string) []int {
	if len(rows) == 0 {
		return nil
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			// Strip ANSI codes to get actual display width
			visibleLen := len(stripAnsi(col))
			if i < len(widths) && visibleLen > widths[i] {
				widths[i] = visibleLen
			}
		}
	}
	return widths
}

func printRow(row []string, widths []int, header bool) {
	var parts []string
	for i, col := range row {
		if i >= len(widths) {
			break
		}
		// Calculate padding based on visible length (without ANSI codes)
		visibleLen := len(stripAnsi(col))
		padding := max(widths[i]-visibleLen, 0)
		padded := col + strings.Repeat(" ", padding)

		if header {
			parts = append(parts, Bold(padded))
		} else {
			parts = append(parts, padded)
		}
	}
	fmt.Println(strings.Join(parts, "  "))
}

func printDivider(widths []int) {
	total := 0
	for _, w := range widths {
		total += w
	}
	total += (len(widths) - 1) * 2
	fmt.Println(strings.Repeat("─", total))
}

// Header returns a bold header
func Header(text string) string {
	return Bold(text)
}

// StepHeader returns a colored step header for wizards
func StepHeader(step int, text string) string {
	return Cyan(fmt.Sprintf("Step %d: %s", step, text))
}

// Divider returns a divider line
func Divider(length int) string {
	return strings.Repeat("─", length)
}

// Truncate truncates a string to maxLen, adding "..." if needed
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatMoney formats a float as currency
func FormatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// FormatEnabled returns colored "Enabled" or "Disabled"
func FormatEnabled(enabled bool) string {
	if enabled {
		return Green("Enabled")
	}
	return "Disabled"
}

// ShowProgress displays a spinner with a message
// Returns a stop function (shows checkmark) and cancel function (clears line)
func ShowProgress(message string) (stop func(), cancel func()) {
	done := make(chan bool)
	clear := make(chan bool)

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r%s %s\n", Green("✓"), message)
				return
			case <-clear:
				fmt.Printf("\r\033[K") // Clear line
				return
			default:
				fmt.Printf("\r%s %s", frames[i], message)
				i = (i + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	stopFn := func() {
		done <- true
		time.Sleep(200 * time.Millisecond)
	}

	cancelFn := func() {
		clear <- true
		time.Sleep(200 * time.Millisecond)
	}

	return stopFn, cancelFn
}

// ansiRegex matches ANSI escape codes
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsi removes ANSI escape codes from a string
func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// Pagination represents pagination metadata
type Pagination struct {
	Page            int
	HasPreviousPage bool
	HasNextPage     bool
}

// DisplayPagination shows pagination information if metadata is available
func DisplayPagination(page int, hasPrev, hasNext bool) {
	// Don't show pagination info if we're on page 1 and there's no next page
	if page == 1 && !hasNext {
		return
	}

	fmt.Println()

	// Show current page
	fmt.Printf("Page %d", page)

	// Show navigation hints
	var hints []string
	if hasPrev {
		hints = append(hints, Gray(fmt.Sprintf("--page %d for previous", page-1)))
	}
	if hasNext {
		hints = append(hints, Gray(fmt.Sprintf("--page %d for next", page+1)))
	}

	if len(hints) > 0 {
		fmt.Print(" (")
		for i, hint := range hints {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(hint)
		}
		fmt.Print(")")
	}

	fmt.Println()
}
