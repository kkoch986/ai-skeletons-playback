package common

import (
	"context"

	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

func Logger(ctx context.Context, level string) (context.Context, func(), error) {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		l = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	zapctx.LogLevel = l
	logger := zapctx.Logger(ctx)
	ctx = zapctx.WithLogger(ctx, logger)
	return ctx, func() {
		_ = logger.Sync()
	}, nil
}
