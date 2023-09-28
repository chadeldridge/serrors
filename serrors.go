package serrors

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"time"
)

// SErrors contains the highest slog.Level and a list of errors (slog.Record)
type SErrors struct {
	// buf used for SErrors.handler to return records instead of logging them
	buf *bytes.Buffer
	// json flag to use JSON instead of text
	json bool
	// logger handler for writing logs
	logger slog.Handler
	// logger handler for writing to SErrors.buf
	handler slog.Handler
	// Level shows the highest slog.Level of errors added
	Level slog.Level
	// Errors is a list of slog.Record
	Errors []slog.Record
}

// UpperCaseKey converts slog.Attr.Key to upper case and returns the new slog.Attr
func UpperCaseKey(_ []string, a slog.Attr) slog.Attr {
	a.Key = strings.ToUpper(a.Key)
	return a
}

// UpperCaseKey converts slog.Attr.Key to lower case and returns the new slog.Attr
func LowerCaseKey(_ []string, a slog.Attr) slog.Attr {
	a.Key = strings.ToLower(a.Key)
	return a
}

// New creates a new SErrors struct using the slog.JSONHandler
func New(logWriter io.Writer, opts *slog.HandlerOptions) SErrors {
	return NewJSONHandler(logWriter, opts)
}

// NewJSONHandler creates a new SErrors struct which uses the slog.JSONHandler
func NewJSONHandler(logWriter io.Writer, opts *slog.HandlerOptions) SErrors {
	b := bytes.NewBuffer(nil)

	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return SErrors{
		buf:     b,
		json:    true,
		logger:  slog.NewJSONHandler(logWriter, opts),
		handler: slog.NewJSONHandler(b, opts),
		Errors:  []slog.Record{},
	}
}

// NewTextHandler creates a new SErrors struct which uses the slog.TextHandler
func NewTextHandler(logWriter io.Writer, opts *slog.HandlerOptions) SErrors {
	b := bytes.NewBuffer(nil)

	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return SErrors{
		buf:     b,
		json:    false,
		logger:  slog.NewTextHandler(logWriter, opts),
		handler: slog.NewTextHandler(b, opts),
		Errors:  []slog.Record{},
	}
}

// Add creates a new slog.Record and adds it to SErrors.Errors from slog.Attr(s).
func (e *SErrors) Add(t time.Time, l slog.Level, msg string, attrs ...slog.Attr) {
	r := slog.NewRecord(t, l, msg, 0)
	r.AddAttrs(attrs...)
	e.Errors = append(e.Errors, r)
	if l > e.Level {
		e.Level = l
	}
}

// Add creates a new slog.Record and adds it to SErrors.Errors from generics.
// args are grouped into key-value pairs.
func (e *SErrors) AddAny(t time.Time, l slog.Level, msg string, args ...any) {
	r := slog.NewRecord(time.Now(), l, msg, 0)
	r.Add(args...)
	e.Errors = append(e.Errors, r)
	if l > e.Level {
		e.Level = l
	}
}

// Debug creates a new Debug Level slog.Record and adds it to SErrors.Errors from slog.Attr(s)
func (e *SErrors) Debug(t time.Time, msg string, attrs ...slog.Attr) {
	e.Add(t, slog.LevelDebug, msg, attrs...)
}

// Debug adds a new Debug Level slog.Record and adds it to SErrors.Errors from generics
// args are grouped into key-value pairs.
func (e *SErrors) DebugAny(t time.Time, msg string, args ...any) {
	e.AddAny(t, slog.LevelDebug, msg, args...)
}

// Info adds a new Info Level slog.Record and adds it to SErrors.Errors from slog.Attr(s)
func (e *SErrors) Info(t time.Time, msg string, attrs ...slog.Attr) {
	e.Add(t, slog.LevelInfo, msg, attrs...)
}

// Info adds a new Info Level slog.Record and adds it to SErrors.Errors from generics
// args are grouped into key-value pairs.
func (e *SErrors) InfoAny(t time.Time, msg string, args ...any) {
	e.AddAny(t, slog.LevelInfo, msg, args...)
}

// Warn adds a new Warn Level slog.Record and adds it to SErrors.Errors from slog.Attr(s)
func (e *SErrors) Warn(t time.Time, msg string, attrs ...slog.Attr) {
	e.Add(t, slog.LevelWarn, msg, attrs...)
}

// Warn adds a new Warn Level slog.Record and adds it to SErrors.Errors from generics
// args are grouped into key-value pairs.
func (e *SErrors) WarnAny(t time.Time, msg string, args ...any) {
	e.AddAny(t, slog.LevelWarn, msg, args...)
}

// Error adds a new Error Level slog.Record and adds it to SErrors.Errors from slog.Attr(s)
func (e *SErrors) Error(t time.Time, msg string, attrs ...slog.Attr) {
	e.Add(t, slog.LevelError, msg, attrs...)
}

// Error adds a new Error Level slog.Record and adds it to SErrors.Errors from generics
// args are grouped into key-value pairs.
func (e *SErrors) ErrorAny(t time.Time, msg string, args ...any) {
	e.AddAny(t, slog.LevelError, msg, args...)
}

// Stack adds the arguement to the beginning of e.Errors and sets e.Level to the highest Level between the two
func (e *SErrors) Stack(errs SErrors) {
	if e.Level < errs.Level {
		e.Level = errs.Level
	}

	e.Errors = append(errs.Errors, e.Errors...)
}

// Append appends arguement to e.Errors and sets e.Level to the highest Level between the two
func (e *SErrors) Append(errs SErrors) {
	if e.Level < errs.Level {
		e.Level = errs.Level
	}

	e.Errors = append(e.Errors, errs.Errors...)
}

func (e SErrors) IsEmpty() bool { return len(e.Errors) < 1 }

// String returns all e.Errors as a sing string
func (e SErrors) String() string {
	var s string
	for _, r := range e.Errors {
		s += e.RtoString(r)
	}

	return s
}

// RtoString converst a slog.Record to a string
func (e SErrors) RtoString(r slog.Record) string {
	if err := e.handler.Handle(context.Background(), r); err != nil {
		return err.Error()
	}

	s := e.buf.String()
	e.buf.Reset()
	return s
}

// First returns the first slog.Record added to SErrors.Errors
func (e SErrors) First() slog.Record { return e.Errors[0] }

// Last returns the last slog.Record added to SErrors.Errors
func (e SErrors) Last() slog.Record { return e.Errors[len(e.Errors)-1] }

// ToArray returns SErrors.Errors as []string and an error
func (e SErrors) ToArray() ([]string, error) {
	s := make([]string, len(e.Errors))
	for i, r := range e.Errors {
		str := e.RtoString(r)
		str = strings.TrimSuffix(str, "\n")
		s[i] = str
	}
	return s, nil
}

// Log writes all SErrors.Errors using the SErrors.logger handler
func (e SErrors) Log() error {
	for _, r := range e.Errors {
		if err := e.logger.Handle(context.Background(), r); err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON converts SErrors.Errors to a JSON array
func (e SErrors) MarshalJSON() ([]byte, error) {
	s := "["
	l := len(e.Errors) - 1

	a, err := e.ToArray()
	if err != nil {
		return nil, err
	}

	for c, r := range a {
		if c < l {
			r += ","
		}
		s += r
	}

	s += "]"

	return []byte(s), nil
}
