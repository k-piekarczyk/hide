package main

import (
	"fmt"
	"image"
	"os"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func SplitBufferIntoChunks(buffer *[]byte) []uint8 {
	nBytes := len(*buffer)
	nBits := nBytes * 8

	bitArray := make([]uint8, nBits)

	for i := 0; i < nBytes; i++ {
		currentByte := (*buffer)[i]

		for j := 0; j < 8; j++ {
			currentBit := currentByte & (0b10000000 >> j) >> (7 - j)

			bitArray[(8*i)+j] = currentBit
		}
	}

	return bitArray
}

func GlueBits(bits *[]uint8) []byte {
	nBits := len(*bits)
	nBytes := nBits / 8

	glued := make([]byte, nBytes)

	for i := 0; i < nBytes; i++ {
		currentByte := byte(0)

		for j := 0; j < 8; j++ {
			currentByte = currentByte | ((*bits)[(i*8)+j] << (7 - j))
		}

		glued[i] = currentByte
	}

	return glued
}

func GetImageFromFilePath(filePath string) image.Image {
	f, err := os.Open(filePath)
	PanicOnError(err)
	defer f.Close()

	img, imgType, err := image.Decode(f)
	PanicOnError(err)

	fmt.Println(imgType)

	return img
}

func GetBytesFromFilePath(filePath string) []byte {
	f, err := os.Open(filePath)
	PanicOnError(err)
	defer f.Close()

	info, err := f.Stat()
	PanicOnError(err)

	nBytes := info.Size()

	var buffer = make([]byte, nBytes)

	_, err = f.Read(buffer)
	PanicOnError(err)

	return buffer
}
