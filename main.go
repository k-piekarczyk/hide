package main

func main() {
	Hide("./resources/input.png", "./resources/toHide.copy.txt", "result")
	Uncover("./result.png", "result.txt")
}
