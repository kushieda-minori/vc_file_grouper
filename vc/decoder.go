// Original Author: Kellindil Maendellyn
// http://valkyriecrusade.wikia.com/wiki/Thread:119497#19
// Converted to go from java
package vc

import (
	"bytes"
	"encoding/binary"
	"errors"
	//"io"
	"io/ioutil"
	"os"
)

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

func toInt32(data []byte, sliceStart int) (ret int32) {
	buf := bytes.NewBuffer(data[sliceStart:(sliceStart + 4)])
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}
