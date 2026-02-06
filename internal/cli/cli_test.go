package cli

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gsamokovarov/assert"
)

func captureOutput(fn func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic("failed to create pipe: " + err.Error())
	}

	old := os.Stdout
	defer func() { os.Stdout = old }()

	os.Stdout = w
	fn()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func TestTruncate(t *testing.T) {
	t.Run("shorter than max", func(t *testing.T) {
		assert.Equal(t, "abc", Truncate("abc", 10))
	})

	t.Run("equal to max", func(t *testing.T) {
		assert.Equal(t, "abcdefghij", Truncate("abcdefghij", 10))
	})

	t.Run("longer than max", func(t *testing.T) {
		assert.Equal(t, "abcdefg...", Truncate("abcdefghijk", 10))
	})

	t.Run("empty string", func(t *testing.T) {
		assert.Equal(t, "", Truncate("", 10))
	})
}

func TestFormatMoney(t *testing.T) {
	t.Run("whole number", func(t *testing.T) {
		assert.Equal(t, "$10.00", FormatMoney(10))
	})

	t.Run("with cents", func(t *testing.T) {
		assert.Equal(t, "$10.50", FormatMoney(10.5))
	})

	t.Run("zero", func(t *testing.T) {
		assert.Equal(t, "$0.00", FormatMoney(0))
	})

	t.Run("rounds to two decimals", func(t *testing.T) {
		assert.Equal(t, "$10.57", FormatMoney(10.567))
	})
}

func TestFormatMemory(t *testing.T) {
	t.Run("adds MB when no unit", func(t *testing.T) {
		assert.Equal(t, "512 MB", FormatMemory("512"))
	})

	t.Run("preserves existing MB", func(t *testing.T) {
		assert.Equal(t, "512 MB", FormatMemory("512 MB"))
	})

	t.Run("preserves existing GB", func(t *testing.T) {
		assert.Equal(t, "2 GB", FormatMemory("2 GB"))
	})

	t.Run("case insensitive check", func(t *testing.T) {
		assert.Equal(t, "512 mb", FormatMemory("512 mb"))
	})
}

func TestPrintTable(t *testing.T) {
	t.Run("shows all rows", func(t *testing.T) {
		rows := [][]string{
			{"Slug", "Name"},
			{"digital_ocean", "DigitalOcean"},
			{"packet", "Packet"},
			{"linode", "Linode"},
			{"hetzner", "Hetzner"},
			{"vultr", "Vultr"},
		}

		output := captureOutput(func() {
			PrintTable(rows)
		})

		cleaned := stripANSI(output)

		for _, row := range rows {
			for _, cell := range row {
				assert.True(t, strings.Contains(cleaned, cell))
			}
		}
	})

	t.Run("empty input", func(t *testing.T) {
		output := captureOutput(func() {
			PrintTable(nil)
		})

		assert.Equal(t, "", output)
	})

	t.Run("header only", func(t *testing.T) {
		rows := [][]string{
			{"Name", "Value"},
		}

		output := captureOutput(func() {
			PrintTable(rows)
		})

		cleaned := stripANSI(output)
		assert.True(t, strings.Contains(cleaned, "Name"))
	})
}

func TestFormatStorage(t *testing.T) {
	t.Run("adds GB when no unit", func(t *testing.T) {
		assert.Equal(t, "100 GB", FormatStorage("100"))
	})

	t.Run("preserves existing GB", func(t *testing.T) {
		assert.Equal(t, "100 GB", FormatStorage("100 GB"))
	})

	t.Run("preserves existing TB", func(t *testing.T) {
		assert.Equal(t, "2 TB", FormatStorage("2 TB"))
	})
}
