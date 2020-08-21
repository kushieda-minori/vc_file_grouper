package vc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// Original Author: Kellindil Maendellyn
// https://valkyriecrusade.fandom.com/wiki/Thread:119497#19
// Converted to go from java

// Decode File header : 16 bytes
// 4 bytes for the signature (CODE)
// 8 bytes of unknown data
// 4 bytes for one of the encoding's keys (the second key is a magic number
// known from the app, 0x45AF6E5D at the time of writing)
//
// The remainder of the file is encoded 4 bytes by 4 bytes, the last few
// bytes unencoded if the file's length is not a multiple of 4
//
// this method does not return io.EOF as an error.
func Decode(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(data[0:4], []byte("CODE")) {
		return nil, errors.New("File '" + file + "' is not encoded")
	}

	subMe := toInt32(data, 12)
	xorMe := int32(0x45AF6E5D)

	// We'll ignore the 16-bytes signature
	excessBytes := len(data) % 4
	encodedLength := len(data) - 16 - excessBytes
	result := make([]byte, 0, len(data)-16)

	for i := 0; i < encodedLength/4; i++ {
		decodedBytes := (toInt32(data, 16+(i*4)) ^ xorMe) - subMe

		buf := bytes.NewBuffer(make([]byte, 0, 4))
		binary.Write(buf, binary.LittleEndian, decodedBytes)
		result = append(result, buf.Bytes()[:]...)
	}

	if excessBytes > 0 {
		result = append(result, data[(16+encodedLength):]...)
	}

	return result, nil
}

// DecodeAndSave Decodes the file and saves the result in the same location as the coded file.
// This method does not return io.EOF as an error.
func DecodeAndSave(file string) (string, []byte, error) {
	data, err := Decode(file)
	if err != nil {
		return "", nil, err
	}

	var fileName string
	if bytes.Equal(data[:4], []byte{0x89, 'P', 'N', 'G'}) {
		fileName = file + ".png"
	} else {
		fileName = file + ".json"

		// remove trailing 0's from the end of the file.
		dataLen := len(data)
		for data[dataLen-1] == 0 {
			dataLen--
		}
		data = data[:dataLen]
	}

	err = ioutil.WriteFile(fileName, data, os.FileMode(0655))
	if err != nil {
		return "", nil, err
	}

	return fileName, data, nil
}

// IsFileEncoded opens a path and reads the first 4 bytes to determin if the file uses VC encoding.
// Does not return io.OEF as an error
func IsFileEncoded(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	b := make([]byte, 4)
	c, err := f.Read(b)
	if err != nil && err != io.EOF {
		return false, err
	}
	if c != 4 {
		return false, errors.New("Unable to read 'magic number' of file to determin encoding")
	}
	return bytes.Equal(b, []byte("CODE")), nil
}

// reads a single 32bit int from the data starting at sliceStart position.
func toInt32(data []byte, sliceStart int) (ret int32) {
	buf := bytes.NewBuffer(data[sliceStart:(sliceStart + 4)])
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}
