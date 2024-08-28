package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("main: Starting application")
	app := NewApp()
	go app.startup()

	// Keep the main function running
	for {
		time.Sleep(1 * time.Second)
	}
}
