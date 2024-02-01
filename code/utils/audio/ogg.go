package audio

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/pion/opus"
	"github.com/pion/opus/pkg/oggreader"
)

func OggToWavByPath(ogg string, wav string) error {
	input, err := os.Open(ogg)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(wav)
	if err != nil {
		return err
	}

	defer output.Close()
	return OggToWav(input, output)
}

func OggToWav(input io.Reader, output io.WriteSeeker) error {
	ogg, _, err := oggreader.NewWith(input)
	if err != nil {
		return err
	}

	out := make([]byte, 1920)

	decoder := opus.NewDecoder()
	encoder := NewEncoder(output, 44100, 16)

	for {
		segments, _, err := ogg.ParseNextPage()
		if errors.Is(err, io.EOF) {
			break
		} else if bytes.HasPrefix(segments[0], []byte("OpusTags")) {
			continue
		}

		if err != nil {
			panic(err)
		}

		for i := range segments {
			if _, _, err = decoder.Decode(segments[i], out); err != nil {
				panic(err)
			}
			encoder.Write(out)
		}
	}
	encoder.Close()
	return nil
}
