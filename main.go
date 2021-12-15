package main

import (
	"flag"
	"os"
	"path"
)

func main() {
	modePtr := flag.String("mode", "encode", "sets program mode, possible modes: encode, decode")
	imagePtr := flag.String("image", "", "input image path, both modes, required")
	payloadPtr := flag.String("payload", "", "input payload path, required for encode mode")
	outputPtr := flag.String("output", "", "output path, both modes (in encode '.png' will always be appended), required")

	flag.Parse()

	switch *modePtr {
	case "encode":
		if *imagePtr == "" || *payloadPtr == "" || *outputPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		Hide(path.Clean(*imagePtr), path.Clean(*payloadPtr), path.Clean(*outputPtr))
	case "decode":
		if *imagePtr == "" || *outputPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		Uncover(path.Clean(*imagePtr), path.Clean(*outputPtr))
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
