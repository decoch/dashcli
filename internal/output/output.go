package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Options struct {
	JSON   bool
	Stdout io.Writer
	Stderr io.Writer
}

type Output struct {
	json   bool
	stdout io.Writer
	stderr io.Writer
}

func New(options Options) *Output {
	return &Output{
		json:   options.JSON,
		stdout: options.Stdout,
		stderr: options.Stderr,
	}
}

func (out *Output) JSONEnabled() bool {
	if out == nil {
		return false
	}
	return out.json
}

func (out *Output) Print(value any) error {
	if out == nil {
		return nil
	}
	if out.json {
		encoder := json.NewEncoder(out.stdout)
		return encoder.Encode(value)
	}
	_, err := fmt.Fprintln(out.stdout, value)
	return err
}

func (out *Output) PrintText(format string, args ...any) error {
	if out == nil {
		return nil
	}
	_, err := fmt.Fprintf(out.stdout, format, args...)
	return err
}

func (out *Output) Errorf(format string, args ...any) {
	if out == nil {
		return
	}
	_, _ = fmt.Fprintf(out.stderr, format+"\n", args...)
}
