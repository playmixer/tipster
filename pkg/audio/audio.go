package audio

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-audio/wav"
)

type WavMetadata struct {
	Duration time.Duration
}

func ReadWavMetadata(data []byte) (*WavMetadata, error) {
	buf := bytes.NewReader(data)

	dec := wav.NewDecoder(buf)
	dec.ReadMetadata()
	if err := dec.Err(); err != nil {
		return nil, fmt.Errorf("failed read metadata: %w", err)
	}
	// if dec.Metadata == nil {
	// 	return nil, errors.New("no metadata present")
	// }

	duration := (float32(len(data)) / float32(dec.AvgBytesPerSec)) * float32(time.Second)
	meta := &WavMetadata{
		Duration: time.Duration(duration),
	}

	return meta, nil
}
