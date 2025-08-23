package input

import (
	"bufio"
	"io"
	"strings"

	"github.com/viniciuscg/vinux/internal/notify"
)

type InputReader interface {
	ReadInput(prompt string) string
}

type ConsoleReader struct {
	reader *bufio.Reader
}

func NewConsoleReader(
	r io.Reader,
) *ConsoleReader {
	return &ConsoleReader{
		reader: bufio.NewReader(r),
	}
}

func (c *ConsoleReader) ReadInput(prompt string) string {
	notify.Print(
		notify.TypeWrite,
		prompt,
		nil,
	)
	input, err := c.reader.ReadString('\n')
	if err != nil {
		notify.Print(
			notify.TypeError,
			notify.ErrorReadingInput,
			err,
		)
	}

	return strings.TrimSpace(input)
}
