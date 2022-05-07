package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gpontesss/bib/ui"
)

func main() {
	vsrs, err := GetVersions()
	if err != nil {
		log.Fatal(err)
	}
	ui := ui.UI{Versions: vsrs}
	defer ui.Close()
	if err := ui.Init(); err != nil {
		panic(fmt.Errorf("ui.Init: %v", err))
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-ui.AsyncLoop():
		if err != nil {
			panic(fmt.Errorf("Exited with error: %v", err))
		}
	case sig := <-sigchan:
		panic(fmt.Errorf("Caught exit signal: %v", sig))
	}
}
