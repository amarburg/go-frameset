package frameset

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"io"
	"testing"
)

func TestFrameSourceGoodJson(t *testing.T) {

	set, err := LoadFrameSet(frameset_testfiles.GoodMultiMovJson)

	if err != nil {
		t.Errorf("Unable to load good.json: %s", err)
	}

	source, err := MakeFrameSetFrameSource(set)

	if err != nil {
		t.Errorf("Unable to make frame source from good.json: %s", err)
	}

	frames := 0
  done := false

	for done == false {
		_, err := source.Next()

		switch err {
		case io.EOF:
			done = true
		case nil:
		default:
			t.Fatalf("Error reading frame: %s", err)
		}

		frames++
    t.Logf("Got frame %d", frames)

		// TODO.  Check that frames are valid.
	}

	if frames != frameset_testfiles.GoodMultiMovJsonFrames {
		t.Errorf("Didn't get as many frames as I expected %d, rather than %d", frames, frameset_testfiles.GoodMultiMovJsonFrames)
	}

}
