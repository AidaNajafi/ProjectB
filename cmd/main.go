package main

import (
	"authentication/cmd/app"
	"fmt"
	"os"
)

func main() {
	if err := app.AppStart(); err != nil {
		fmt.Println(err)
		fmt.Println("Failed to run server!")
		os.Exit(1)
	}

}
