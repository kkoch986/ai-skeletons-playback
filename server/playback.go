package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"io"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/juju/zaputil/zapctx"
	"github.com/kkoch986/ai-skeletons-playback/output"
)

// PlaybackRequest is the struct that can parse the POST request sent to /generate
type PlaybackRequest struct {
	RawAudioData string          `json:"audio"`
	SequenceData output.Sequence `json:"sequence"`
}

// PlaybackResponse contains the respost to POST /generate calls
type PlaybackResponse struct {
}

// HandlePlayback handles the POST /generate calls
// It will configure beep to play the audio and trigger the output modules to play the provided sequencing events
func HandlePOSTPlayback(
	ctx context.Context,
	r *PlaybackRequest,
	w output.HardwareWriter,
) (*PlaybackResponse, error) {
	logger := zapctx.Logger(ctx)

	// reverse the hex encoding on the mp3 audio and base64 decode it to get
	// the raw mp3 data
	data, err := hex.DecodeString(r.RawAudioData)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	_, err = base64.StdEncoding.Decode(dst, data)
	if err != nil {
		return nil, err
	}
	rc := io.NopCloser(bytes.NewReader(dst))

	// create the streamer for the mp3 file
	streamer, format, err := mp3.Decode(rc)
	if err != nil {
		return nil, err
	}
	logger.Debug("audio decoded")

	// append watchers for each element in the sequence provided
	watchStreamer := beep.Watch(streamer)
	watchStreamer.StartedAsync(func(pos int) {
		w.Start(ctx)
	})
	watchStreamer.EndedAsync(func(pos int) {
		w.End(ctx)
	})
	for pos, values := range r.SequenceData {
		for _, v := range values {
			watchStreamer.AtAsync(
				format.SampleRate.N(time.Duration(pos*float64(time.Second))),
				func(pos int) {
					w.Send(ctx, v)
				},
			)
		}
	}
	logger.Debug("sequence data addeed")

	// play the streamer
	err = speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)
	if err != nil {
		return nil, err
	}
	logger.Debug("starting stream")
	speaker.Play(watchStreamer)

	return &PlaybackResponse{}, nil
}
