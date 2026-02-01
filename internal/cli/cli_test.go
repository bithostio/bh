package cli

import (
	"testing"

	"github.com/gsamokovarov/assert"
)

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
