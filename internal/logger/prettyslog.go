package logger

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
)

const (
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Reset   = "\033[0m"
)

func Colorize(color, text string) string {
	return color + text + Reset
}

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = Colorize(Magenta, level)
	case slog.LevelInfo:
		level = Colorize(Blue, level)
	case slog.LevelWarn:
		level = Colorize(Yellow, level)
	case slog.LevelError:
		level = Colorize(Red, level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := Colorize(Cyan, r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		Colorize(White, string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
