package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"math/rand"
	"os"
)

func decodePayload(img *image.Image) []uint8 {
	endPoint := (*img).Bounds().Max

	width := endPoint.X
	height := endPoint.Y

	bitCounter := 0
	lengthBits := make([]uint8, 64)
	var payloadBits []uint8

	pixelMap := map[Coord]bool{}
	rg := rand.New(rand.NewSource(int64(width * height)))

	for true {
		x := rg.Intn(width)
		y := rg.Intn(height)

		xy := Coord{x: x, y: y}

		for pixelMap[xy] {
			x = rg.Intn(width)
			y = rg.Intn(height)

			xy = Coord{x: x, y: y}
		}

		pixelMap[xy] = true

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

					fmt.Println("Decoded payload length:", payloadLength, "byte(s)")

					PanicOnError(err)

					payloadBits = make([]uint8, payloadLength*8)

					rg = rand.New(rand.NewSource(int64(width * height * int(payloadLength))))
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

	panic("We shouldn't have gotten this far")
}

func Uncover(input string, output string) {
	fmt.Println("Opening a file ")

	img := GetImageFromFilePath(input)

	payloadBits := decodePayload(&img)

	payloadBytes := GlueBits(&payloadBits)

	f, err := os.Create(output)
	PanicOnError(err)
	defer f.Close()

	_, err = f.Write(payloadBytes)
	PanicOnError(err)
}
