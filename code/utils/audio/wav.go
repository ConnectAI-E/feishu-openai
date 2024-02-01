package audio

import (
	"encoding/binary"
	"io"
)

type Encoder struct {
	Output          io.WriteSeeker
	SampleRate      int
	BitDepth        int
	totalBytes      uint32
	isHeaderWritten bool
}

func (e *Encoder) WriteHeader() error {
	if err := writeLe(e.Output, []byte("RIFF")); err != nil {
		return err
	}

	if err := writeLe(e.Output, uint32(0)); err != nil { // Placeholder for file size
		return err
	}

	if err := writeLe(e.Output, []byte("WAVE")); err != nil {
		return err
	}

	if err := writeLe(e.Output, []byte("fmt ")); err != nil {
		return err
	}
	if err := writeLe(e.Output, uint32(16)); err != nil {
		return err
	}

	if err := writeLe(e.Output, uint16(1)); err != nil { // Audio format: PCM
		return err
	}
	if err := writeLe(e.Output, uint16(1)); err != nil { // Number of channels: 1 (mono)
		return err
	}
	if err := writeLe(e.Output, uint32(e.SampleRate)); err != nil {
		return err
	}

	if err := writeLe(e.Output, uint32(e.SampleRate*e.BitDepth/8)); err != nil {
		return err
	}

	if err := writeLe(e.Output, uint16(e.BitDepth/8)); err != nil {
		return err
	}
	if err := writeLe(e.Output, uint16(e.BitDepth)); err != nil {
		return err
	}

	if err := writeLe(e.Output, []byte("data")); err != nil {
		return err
	}

	if err := writeLe(e.Output, uint32(0)); err != nil { //Placeholder for data size
		return err
	}
	e.isHeaderWritten = true
	return nil
}

func writeLe[T []byte | uint32 | uint16 | uint8](w io.Writer, data T) error {
	return binary.Write(w, binary.LittleEndian, data)
}

func (e *Encoder) Write(data []byte) error {
	if !e.isHeaderWritten {
		e.WriteHeader()
	}
	n, err := e.Output.Write(data)
	if err != nil {
		return err
	}
	e.totalBytes += uint32(n)
	return nil
}

func (e *Encoder) Close() error {
	if _, err := e.Output.Seek(4, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(e.Output, binary.LittleEndian, uint32(36+e.totalBytes)); err != nil {
		return err
	}
	if _, err := e.Output.Seek(40, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(e.Output, binary.LittleEndian, e.totalBytes); err != nil {
		return err
	}
	return nil
}

func NewEncoder(w io.WriteSeeker, sampleRate int, bitDepth int) *Encoder {
	return &Encoder{
		SampleRate:      sampleRate,
		Output:          w,
		BitDepth:        bitDepth,
		isHeaderWritten: false,
	}
}
