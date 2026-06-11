package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

func Info(format string, args ...any)    { color(blue, format, args...) }
func Success(format string, args ...any) { color(green, format, args...) }
func Warn(format string, args ...any)    { color(yellow, format, args...) }
func Error(format string, args ...any)   { color(red, format, args...) }

func Confirm(r io.Reader, w io.Writer, question string) bool {
	if r == nil {
		r = os.Stdin
	}
	if w == nil {
		w = os.Stdout
	}
	_, _ = fmt.Fprintf(w, "%s [y/N]: ", question)
	line, _ := bufio.NewReader(r).ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

func color(c, format string, args ...any) {
	fmt.Printf(c+format+reset+"\n", args...)
}
