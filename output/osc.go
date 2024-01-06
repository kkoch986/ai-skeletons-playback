package output

import (
	"context"

	"github.com/hypebeast/go-osc/osc"
)

// OSCWriter is a HardwareWriter that transmits the data
// as OSC messages
type OSCWriter struct {
	name   string
	client *osc.Client
}

func NewOSCWriter(n string, c *osc.Client) *OSCWriter {
	return &OSCWriter{n, c}
}

func (w *OSCWriter) Name() string {
	return w.name
}

func (w *OSCWriter) Start(ctx context.Context) error {
	return nil
}

func (w *OSCWriter) End(ctx context.Context) error {
	return nil
}

func (w *OSCWriter) Send(ctx context.Context, v Value) error {
	msg := osc.NewMessage(v.Key)
	msg.Append(float32(v.Value))
	return w.client.Send(msg)
}
