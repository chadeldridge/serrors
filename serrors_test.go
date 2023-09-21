package serrors

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

type Param struct {
	level slog.Level
	msg   string
	attrs []slog.Attr
	want  string
}

type Test struct {
	name   string
	opts   slog.HandlerOptions
	params []Param
}

var (
	testTime           = time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	testAttrParamsText = []Test{
		{
			"default",
			slog.HandlerOptions{},
			[]Param{
				{
					slog.LevelDebug,
					"m",
					[]slog.Attr{slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2})},
					`time=2000-01-02T03:04:05.000Z level=DEBUG msg=m a=1 m=map[b:2]`,
				},
				{
					slog.LevelInfo,
					"m",
					[]slog.Attr{slog.Int("a", 2), slog.Any("m", map[string]int{"b": 2})},
					`time=2000-01-02T03:04:05.000Z level=INFO msg=m a=2 m=map[b:2]`,
				},
				{
					slog.LevelWarn,
					"m",
					[]slog.Attr{slog.Int("a", 3), slog.Any("m", map[string]int{"b": 2})},
					`time=2000-01-02T03:04:05.000Z level=WARN msg=m a=3 m=map[b:2]`,
				},
				{
					slog.LevelError,
					"m",
					[]slog.Attr{slog.Int("a", 4), slog.Any("m", map[string]int{"b": 2})},
					`time=2000-01-02T03:04:05.000Z level=ERROR msg=m a=4 m=map[b:2]`,
				},
			},
		},
		{
			"upperCaseKeys",
			slog.HandlerOptions{ReplaceAttr: UpperCaseKey},
			[]Param{
				{
					slog.LevelDebug,
					"m",
					[]slog.Attr{slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2})},
					`TIME=2000-01-02T03:04:05.000Z LEVEL=DEBUG MSG=m A=1 M=map[b:2]`,
				},
				{
					slog.LevelInfo,
					"m",
					[]slog.Attr{slog.Int("a", 2), slog.Any("m", map[string]int{"b": 2})},
					`TIME=2000-01-02T03:04:05.000Z LEVEL=INFO MSG=m A=2 M=map[b:2]`,
				},
				{
					slog.LevelWarn,
					"m",
					[]slog.Attr{slog.Int("a", 3), slog.Any("m", map[string]int{"b": 2})},
					`TIME=2000-01-02T03:04:05.000Z LEVEL=WARN MSG=m A=3 M=map[b:2]`,
				},
				{
					slog.LevelError,
					"m",
					[]slog.Attr{slog.Int("a", 4), slog.Any("m", map[string]int{"b": 2})},
					`TIME=2000-01-02T03:04:05.000Z LEVEL=ERROR MSG=m A=4 M=map[b:2]`,
				},
			},
		},
	}
	testAttrParamsJSON = []Test{
		{
			"default",
			slog.HandlerOptions{},
			[]Param{
				{
					slog.LevelDebug,
					"m",
					[]slog.Attr{slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2})},
					`{"time":"2000-01-02T03:04:05Z","level":"DEBUG","msg":"m","a":1,"m":{"b":2}}`,
				},
				{
					slog.LevelInfo,
					"m",
					[]slog.Attr{slog.Int("a", 2), slog.Any("m", map[string]int{"b": 2})},
					`{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"m","a":2,"m":{"b":2}}`,
				},
				{
					slog.LevelWarn,
					"m",
					[]slog.Attr{slog.Int("a", 3), slog.Any("m", map[string]int{"b": 2})},
					`{"time":"2000-01-02T03:04:05Z","level":"WARN","msg":"m","a":3,"m":{"b":2}}`,
				},
				{
					slog.LevelError,
					"m",
					[]slog.Attr{slog.Int("a", 4), slog.Any("m", map[string]int{"b": 2})},
					`{"time":"2000-01-02T03:04:05Z","level":"ERROR","msg":"m","a":4,"m":{"b":2}}`,
				},
			},
		},
		{
			"upperCaseKeys",
			slog.HandlerOptions{ReplaceAttr: UpperCaseKey},
			[]Param{
				{
					slog.LevelDebug,
					"m",
					[]slog.Attr{slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2})},
					`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"DEBUG","MSG":"m","A":1,"M":{"b":2}}`,
				},
				{
					slog.LevelInfo,
					"m",
					[]slog.Attr{slog.Int("a", 2), slog.Any("m", map[string]int{"b": 2})},
					`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"INFO","MSG":"m","A":2,"M":{"b":2}}`,
				},
				{
					slog.LevelWarn,
					"m",
					[]slog.Attr{slog.Int("a", 3), slog.Any("m", map[string]int{"b": 2})},
					`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"WARN","MSG":"m","A":3,"M":{"b":2}}`,
				},
				{
					slog.LevelError,
					"m",
					[]slog.Attr{slog.Int("a", 4), slog.Any("m", map[string]int{"b": 2})},
					`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"ERROR","MSG":"m","A":4,"M":{"b":2}}`,
				},
			},
		},
	}
)

func TestSErrorsNew(t *testing.T) {
	for _, test := range testAttrParamsJSON {
		t.Run(test.name, func(t *testing.T) {
			var want string
			e := New(nil, &test.opts)

			for _, p := range test.params {
				e.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + "\n"
			}

			got := e.String()
			if got != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}

func TestSErrorsNewTextHandler(t *testing.T) {
	for _, test := range testAttrParamsText {
		t.Run(test.name, func(t *testing.T) {
			var want string
			e := NewTextHandler(nil, &test.opts)

			for _, p := range test.params {
				e.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + "\n"
			}

			got := e.String()
			if got != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}

func TestSErrorsMarshalJSON(t *testing.T) {
	for _, test := range testAttrParamsJSON {
		t.Run(test.name, func(t *testing.T) {
			e := New(os.Stdout, &test.opts)
			want := "["

			for _, p := range test.params {
				e.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + ","
			}
			want = strings.TrimSuffix(want, ",")
			want += "]"

			// got, err := e.MarshalJSON()
			got, err := json.Marshal(e)
			if err != nil {
				t.Fatalf("\ngot  %s\nwant nil", err.Error())
			}

			if string(got) != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}

func TestSErrorsStructFormatting(t *testing.T) {
	for _, test := range testAttrParamsJSON {
		t.Run(test.name, func(t *testing.T) {
			s := struct {
				String string  `json:"string"`
				Int    int     `json:"int"`
				Errors SErrors `json:"errors,omitempty"`
			}{
				"m",
				1,
				New(os.Stdout, &test.opts),
			}

			want := `{"string":"m","int":1,"errors":[`

			for _, p := range test.params {
				s.Errors.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + ","
			}
			want = strings.TrimSuffix(want, ",")
			want += "]}"

			// got, err := e.MarshalJSON()
			got, err := json.Marshal(s)
			if err != nil {
				t.Fatalf("\ngot  %s\nwant nil", err.Error())
			}

			if string(got) != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}

func TestSErrorsLogJSON(t *testing.T) {
	for _, test := range testAttrParamsJSON {
		t.Run(test.name, func(t *testing.T) {
			var want string
			got := bytes.NewBuffer(nil)
			e := New(got, &test.opts)

			for _, p := range test.params {
				e.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + "\n"
			}

			err := e.Log()
			if err != nil {
				t.Fatalf("\ngot  %s\nwant nil", err.Error())
			}

			if got.String() != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}

func TestSErrorsLogText(t *testing.T) {
	for _, test := range testAttrParamsText {
		t.Run(test.name, func(t *testing.T) {
			var want string
			got := bytes.NewBuffer(nil)
			e := NewTextHandler(got, &test.opts)

			for _, p := range test.params {
				e.Add(testTime, p.level, p.msg, p.attrs...)
				want += p.want + "\n"
			}

			err := e.Log()
			if err != nil {
				t.Fatalf("\ngot  %s\nwant nil", err.Error())
			}

			if got.String() != want {
				t.Fatalf("\ngot  %s\nwant %s", got, want)
			}
		})
	}
}
