package output

import (
	"context"

	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

// HardwareWriter is an interface for things that can
// send the sequence signals outward to the physical
// devices rendering the show
// Examples might be OSC, DMX, MQTT etc..
type HardwareWriter interface {
	Name() string
	Start(ctx context.Context) error
	End(ctx context.Context) error
	Send(ctx context.Context, v Value) error
}

var _ HardwareWriter = &LogWriter{}
var _ HardwareWriter = &Multiplexer{}

// LogWriter simply issues a debug log for each value
// sent to it
type LogWriter struct {
}

// Name returns "logwriter", assuming there should never
// be more than one of these so...
func (m *LogWriter) Name() string {
	return "logwriter"
}

// Start logs a start message when the sequence begins to stream
func (m *LogWriter) Start(ctx context.Context) error {
	logger := zapctx.Logger(ctx)
	logger.Debug("stream started")
	return nil
}

// End logs a message when the sequence stops streaming
func (m *LogWriter) End(ctx context.Context) error {
	logger := zapctx.Logger(ctx)
	logger.Debug("stream complete")
	return nil
}

// Send will call Send on each of the writers in m
// Since they are all fired asynchronously,
// errors from the underlying writers are discareded
func (m *LogWriter) Send(ctx context.Context, v Value) error {
	logger := zapctx.Logger(ctx)
	logger.Debug(
		"value received",
		zap.String("key", v.Key),
		zap.Float64("value", v.Value),
	)
	return nil
}

// Multiplexer is a HardwareWriter that writes
// the same message asynchronously to each device
// the writes are not treated as reliable and errors
// are logged and discarded
type Multiplexer struct {
	name    string
	writers []HardwareWriter
}

// NewMultiplexer creates and returns a Multiplexer with
// the given writes
func NewMultiplexer(n string, w []HardwareWriter) *Multiplexer {
	return &Multiplexer{n, w}
}

// Name returns the name this was given when it was created
func (m *Multiplexer) Name() string {
	return m.name
}

// Start calls start on all of the other writers
// since it is done asynchronously, errors are logged and disregarded
func (m *Multiplexer) Start(ctx context.Context) error {
	logger := zapctx.Logger(ctx)
	var err error
	for _, w := range m.writers {
		go func(ctx context.Context, w HardwareWriter) {
			err = w.Start(ctx)
			if err != nil {
				logger.Error("error writing to hardware", zap.String("name", w.Name()))
			}
		}(ctx, w)
	}
	return nil
}

// End calls end on all of the other writers
// since it is done asynchronously, errors are logged and disregarded
func (m *Multiplexer) End(ctx context.Context) error {
	logger := zapctx.Logger(ctx)
	var err error
	for _, w := range m.writers {
		go func(ctx context.Context, w HardwareWriter) {
			err = w.End(ctx)
			if err != nil {
				logger.Error("error writing to hardware", zap.String("name", w.Name()))
			}
		}(ctx, w)
	}
	return nil
}

// Send will call Send on each of the writers in m
// Since they are all fired asynchronously,
// errors from the underlying writers are discareded
func (m *Multiplexer) Send(ctx context.Context, v Value) error {
	logger := zapctx.Logger(ctx)
	var err error
	for _, w := range m.writers {
		go func(ctx context.Context, w HardwareWriter) {
			err = w.Send(ctx, v)
			if err != nil {
				logger.Error("error writing to hardware", zap.String("name", w.Name()))
			}
		}(ctx, w)
	}
	return nil
}
