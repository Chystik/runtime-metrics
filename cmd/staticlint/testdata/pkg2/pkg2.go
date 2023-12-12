package main

import "os"

func main() {
	Exit(0)
}

func Exit(code int) {
	os.Exit(0) // don't want analyzer err here, because os.Exit not in the main function
}
