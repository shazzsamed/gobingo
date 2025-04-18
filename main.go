package main

import (
	"fmt"
	"gobingo/cmd"
	"gobingo/helpers"
	"os"
	"os/signal"
	"syscall"
)

func enableAlternateScreen() {
	fmt.Print("\033[?1049h") 
}

func disableAlternateScreen() {
	helpers.ClearTerminal()
	fmt.Print("\033[?1049l")
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer func() {
		disableAlternateScreen()
	}()

	enableAlternateScreen()

	go func() {
		<-signalChan
		fmt.Println("\nExiting gracefully...")
		disableAlternateScreen()
		os.Exit(0)
	}()

	cmd.Execute()

	disableAlternateScreen()
}
