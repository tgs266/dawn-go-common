package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type stack []uintptr
type Formatter string

// constants for Stack formatters.
const (
	DefaultFormatter    Formatter = "file.go:152\n"
	JavaLikeFormatter             = "at foo(file.go:152)\n"
	GoLikeFormatter               = "foo\n\tfile.go:152\n"
	PythonLikeFormatter           = "File file.go, line 152, in foo\n"
)

var tracerFormatter = DefaultFormatter

// ApplyFormatter specifies the formatter for errorTracer.
// Apply DefaultFormatter if not specified.
// The template of formatter: "foo" for function name, "file.go" for file name,
// "152" for line number.
func ApplyFormatter(formatter Formatter) {
	tracerFormatter = formatter
}

func (f Formatter) format(frame runtime.Frame) string {
	formatted := string(tracerFormatter)
	formatted = strings.Replace(formatted, "foo", frame.Function, 1)
	formatted = strings.Replace(formatted, "file.go", frame.File, 1)
	formatted = strings.Replace(formatted, "152", fmt.Sprint(frame.Line), 1)
	return formatted
}

func recordStack() *stack {
	s := make(stack, 64)
	n := runtime.Callers(3, s)
	s = s[:n]
	return &s
}

// Format formats the stack trace with the layout registered in FormatTracer.
func (s *stack) Format() (ret string) {
	// check nil.
	if s == nil {
		return
	}
	v := *s

	// handle if no frames.
	if len(v) == 0 {
		return
	}

	// get frames for current stack.
	frames := runtime.CallersFrames(v)
	for {
		frame, more := frames.Next()

		// apply the layout formatter.
		ret = ret + tracerFormatter.format(frame)

		if !more {
			return
		}
	}
}
