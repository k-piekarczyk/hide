package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
)

func decodePayload(img *image.Image) []uint8 {
	endPoint := (*img).Bounds().Max

	width := endPoint.X
	height := endPoint.Y

	bitCounter := 0
	lengthBits := make([]uint8, 64)
	var payloadBits []uint8

	fmt.Println(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			originalColor := (*img).At(x, y)

			r, g, b, _ := originalColor.RGBA()

			rgbArray := [3]uint8{uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101)} // De-premultiplication (Go automaticaly alpha premultiplifies)

			for i := 0; i < 3; i++ {
				if bitCounter < 64 {
					cc := rgbArray[i]
					lengthBits[bitCounter] = cc & 0b00000001
					bitCounter++

					if bitCounter == 64 {
						lengthBytes := GlueBits(&lengthBits)

						lengthBuffer := bytes.NewBuffer(lengthBytes)
						payloadLength, err := binary.ReadUvarint(lengthBuffer)

						fmt.Println("Encoded payload length:", payloadLength, "byte(s)")

						PanicOnError(err)

						payloadBits = make([]uint8, payloadLength*8)
					}
				} else if (bitCounter - 64) < len(payloadBits) {
					cc := rgbArray[i]
					payloadBits[bitCounter-64] = cc & 0b00000001
					bitCounter++
				} else {
					return payloadBits
				}
			}
		}
	}
	panic("We shouldn't have gotten this far")
}

func Uncover() {
	var pathToImage = "./result.png"

	fmt.Println("Opening a file ")

	img := GetImageFromFilePath(pathToImage)

	payloadBits := decodePayload(&img)

	payloadBytes := GlueBits(&payloadBits)

	fmt.Println(payloadBytes)
}
