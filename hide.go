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
	"math/rand"
	"os"
)

func encodePayload(original *image.Image, payload *[]byte) *image.RGBA {
	length := uint64(len(*payload))
	lengthBuffer := make([]byte, 8)
	binary.PutUvarint(lengthBuffer, length)
	lengthBits := SplitBufferIntoChunks(&lengthBuffer)

	fmt.Println("Encoding payload:", length, "byte(s)")
	payloadBits := SplitBufferIntoChunks(payload)

	preparedPayload := append(lengthBits, payloadBits...)
	nBits := len(preparedPayload)

	encoded := image.NewRGBA((*original).Bounds())

	endPoint := (*original).Bounds().Max

	width := endPoint.X
	height := endPoint.Y

	bitCapacity := width * height * 3

	if bitCapacity < nBits {
		panic(fmt.Errorf("provided image can store up to %d bytes, but the data to be hidden requires at least %d bytes", bitCapacity/8, nBits/8))
	}

	pixelMap := map[Coord]bool{}
	rg := rand.New(rand.NewSource(int64(width * height)))

	bitCounter := 0
	for bitCounter < nBits {
		x := rg.Intn(width)
		y := rg.Intn(height)

		xy := Coord{x: x, y: y}

		for pixelMap[xy] {
			x = rg.Intn(width)
			y = rg.Intn(height)

			xy = Coord{x: x, y: y}
		}

		pixelMap[xy] = true

		originalColor := (*original).At(x, y)
		oR, oG, oB, _ := originalColor.RGBA()
		rgbArray := [3]uint8{uint8(oR / 0x101), uint8(oG / 0x101), uint8(oB / 0x101)} // De-premultiplication (Go automaticaly alpha premultiplifies)

		for i := 0; i < 3; i++ {
			if bitCounter < nBits {
				cc := rgbArray[i]
				rgbArray[i] = (cc & 0b11111110) | preparedPayload[bitCounter]
				bitCounter++

				if bitCounter == 64 {
					rg = rand.New(rand.NewSource(int64(width * height * len(*payload))))
				}
			}
		}
		encoded.Set(x, y, color.RGBA{R: rgbArray[0], G: rgbArray[1], B: rgbArray[2], A: 0xff})
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xy := Coord{x: x, y: y}

			if !pixelMap[xy] {
				encoded.Set(x, y, (*original).At(x, y))
			}
		}
	}

	return encoded
}

func Hide(image string, payload string, output string) {
	img := GetImageFromFilePath(image)

	messageBuffer := GetBytesFromFilePath(payload)

	f, err := os.Create(output + ".png")
	PanicOnError(err)
	defer f.Close()

	encoded := encodePayload(&img, &messageBuffer)

	encoder := png.Encoder{CompressionLevel: -1}

	err = encoder.Encode(f, encoded)
	PanicOnError(err)
}
