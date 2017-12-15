
package framesource

import (
	"fmt"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-frameset/frameset"
	"image"
	"io"
)

type FrameSetFrameSource struct {
	*frameset.FrameSet
	Movie         lazyquicktime.MovieExtractor
	chunkIdx      int
	frameIdx      int
	segmentOffset uint64
}

func MakeFrameSetFrameSource(set *frameset.FrameSet) (*FrameSetFrameSource, error) {

	mm, err := set.MovieExtractor()

	if err != nil {
		return &FrameSetFrameSource{}, err
	}

	return &FrameSetFrameSource{
		FrameSet: set,
		Movie:    mm,
	}, nil
}

func (source *FrameSetFrameSource) Valid() error {
	if source.chunkIdx >= len(source.FrameSet.Chunks) {
		return io.EOF
	}

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	if chunk.HasFrames() {
		if source.frameIdx >= len(chunk.Frames) {
			return fmt.Errorf("Frame offset is off end of frame array (error) in chunk %d; %d >= %d", source.chunkIdx, source.frameIdx, len(chunk.Frames))
		}
	} else if (chunk.Start + source.segmentOffset) >= chunk.End {
			return fmt.Errorf("Segment offset is off end of segment (error) in chunk %d; %d >= %d", source.chunkIdx, (chunk.Start + source.segmentOffset), chunk.End)
		}

	return nil
}

func (source *FrameSetFrameSource) Advance() {
	source.frameIdx++
	source.segmentOffset++

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	if chunk.HasFrames() {

		if source.frameIdx >= len(chunk.Frames) {
			source.frameIdx = 0
			source.segmentOffset = 0
			source.chunkIdx++
		}
	} else if (chunk.Start + source.segmentOffset) >= chunk.End {
			source.frameIdx = 0
			source.segmentOffset = 0
			source.chunkIdx++
		}

}

func (source *FrameSetFrameSource) Next() (image.Image, error) {
	if err := source.Valid(); err != nil {
		return nil, err
	}

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	var frame uint64
	if chunk.HasFrames() {
		frame = chunk.Frames[source.frameIdx]
	} else {
		frame = chunk.Start + source.segmentOffset
	}

	defer source.Advance()

	return source.Movie.ExtractFrame(frame)
}