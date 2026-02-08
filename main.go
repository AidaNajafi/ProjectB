package main

import (
	"authentication/internal/app"
	"fmt"
	"os"
)

func main() {
	flags := os.Args[1:]
	if err := app.AppStart(flags); err != nil {
		fmt.Println("Failed to run server!")
		os.Exit(1)
	}

}
