package cli

import (
	"errors"

	"github.com/charmbracelet/huh"
)

// ErrAborted is returned when the user cancels a prompt
var ErrAborted = errors.New("aborted")

// Confirm prompts the user for a yes/no confirmation
func Confirm(label string) (bool, error) {
	var confirmed bool

	err := huh.NewConfirm().
		Title(label).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()

	if errors.Is(err, huh.ErrUserAborted) {
		return false, ErrAborted
	}

	return confirmed, err
}

// InputOptions configures an input prompt
type InputOptions struct {
	Mask        rune
	Default     string
	Placeholder string
	Validate    func(string) error
}

// InputOption is a function that configures InputOptions
type InputOption func(*InputOptions)

// WithMask sets the mask character for password inputs
func WithMask(r rune) InputOption {
	return func(o *InputOptions) {
		o.Mask = r
	}
}

// WithDefault sets the default value
func WithDefault(s string) InputOption {
	return func(o *InputOptions) {
		o.Default = s
	}
}

// WithPlaceholder sets the placeholder text
func WithPlaceholder(s string) InputOption {
	return func(o *InputOptions) {
		o.Placeholder = s
	}
}

// WithValidation sets a validation function
func WithValidation(fn func(string) error) InputOption {
	return func(o *InputOptions) {
		o.Validate = fn
	}
}

// Input prompts the user for text input
func Input(label string, opts ...InputOption) (string, error) {
	options := &InputOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var value string
	if options.Default != "" {
		value = options.Default
	}

	input := huh.NewInput().
		Title(label).
		Value(&value)

	if options.Mask != 0 {
		input.EchoMode(huh.EchoModePassword)
	}

	if options.Placeholder != "" {
		input.Placeholder(options.Placeholder)
	}

	if options.Validate != nil {
		input.Validate(options.Validate)
	}

	err := input.Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return "", ErrAborted
	}

	return value, err
}
