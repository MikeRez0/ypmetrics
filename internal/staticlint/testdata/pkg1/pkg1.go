package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit call from main is not recommended"
}
