package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// initContext create a context that will cancel in the event of a sigint or sigterm.
func InitContext() (context.Context, context.CancelFunc) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()
	return ctx, cancel
}
