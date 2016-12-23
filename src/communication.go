package main

import (
	"bytes"
	"github.com/tarm/serial"
	"io"
	"time"
)

func requestData(writeData []byte, device string) ([]byte, error) {
	c := &serial.Config{Name: device, Baud: 9600, ReadTimeout: time.Second * 3}

	s, err := serial.OpenPort(c)
	if err != nil {
		return []byte{}, err
	}

	_, err = s.Write(writeData)
	if err != nil {
		return []byte{}, err
	}

	time.Sleep(100 * time.Millisecond)

	byteBuffer := bytes.NewBuffer(nil)

	buf := make([]byte, 64)
	_, err = s.Read(buf)
	if err != nil {
		return []byte{}, err
	}
	r := bytes.NewReader(buf)
	io.Copy(byteBuffer, r)

	s.Flush()
	s.Close()

	return byteBuffer.Bytes(), nil
}

func idPos(id []byte, data []byte) int {
	bytePos := bytes.Index(data, id)
	if bytePos == -1 {
		return -1
	}

	return bytePos + len(id)
}

func readLEInt(data []byte, start int) int {
	t := 0

	if len(data) > start {
		t += int(data[start])
	}

	if len(data) > start+1 {
		t += (int(data[start+1]) << 8)
	}

	return t
}
