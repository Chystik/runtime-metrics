package main

import "os"

func main() {
	os.Exit(0) // want "calling os.Exit in main function is not allowed"
}
