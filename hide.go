package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

func encodePayload(original *image.Image, payload *[]byte) *image.RGBA {
	length := uint64(len(*payload))
	lengthBuffer := make([]byte, 8)
	binary.PutUvarint(lengthBuffer, length)
	lengthBits := SplitBufferIntoChunks(&lengthBuffer)

	fmt.Println("Payload length", length, "byte(s)")
	fmt.Println(*payload)
	payloadBits := SplitBufferIntoChunks(payload)
	fmt.Println(payloadBits)

	preparedPayload := append(lengthBits, payloadBits...)
	nBits := len(preparedPayload)

	fmt.Println(preparedPayload)

	encoded := image.NewRGBA((*original).Bounds())

	endPoint := (*original).Bounds().Max

	width := endPoint.X
	height := endPoint.Y

	bitCounter := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			originalColor := (*original).At(x, y)

			if bitCounter < nBits {
				oR, oG, oB, _ := originalColor.RGBA()

				rgbArray := [3]uint8{uint8(oR / 0x101), uint8(oG / 0x101), uint8(oB / 0x101)} // De-premultiplication (Go automaticaly alpha premultiplifies)

				for i := 0; i < 3; i++ {
					if bitCounter < nBits {
						cc := rgbArray[i]
						rgbArray[i] = (cc & 0b11111110) | preparedPayload[bitCounter]
						fmt.Print(rgbArray[i]&0b1, " ")
						bitCounter++
					}
				}

				encoded.Set(x, y, color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 0xff})
			} else {
				encoded.Set(x, y, originalColor)
			}
		}
	}

	return encoded
}

func Hide() {
	var pathToHide = "./resources/toHide.txt"
	var pathToImage = "./resources/input.png"

	fmt.Println("Opening a file ")

	img := GetImageFromFilePath(pathToImage)

	messageBuffer := GetBytesFromFilePath(pathToHide)

	f, err := os.Create("result.png")
	PanicOnError(err)
	defer f.Close()

	encoded := encodePayload(&img, &messageBuffer)

	encoder := png.Encoder{CompressionLevel: -1}

	err = encoder.Encode(f, encoded)
	PanicOnError(err)

	fmt.Println("Message hidden")
}
