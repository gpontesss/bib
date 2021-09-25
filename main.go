package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gpontesss/bib/bib"
	"github.com/gpontesss/bib/encoding/tsv"
	"github.com/gpontesss/bib/ui"
)

func main() {
	file, err := os.Open("./kjv.tsv")
	if err != nil {
		log.Fatal("os.Open", err)
	}

	version, err := tsv.Decode(file)
	if err != nil {
		log.Fatal("tsv.Decode", err)
	}

	ui := ui.UI{Versions: []bib.Version{version, version}}
	if ui.Init() != nil {
		log.Fatal("ui.Init", err)
	}
	defer ui.End()

	cancelchan := make(chan os.Signal, 1)
	signal.Notify(cancelchan, syscall.SIGTERM, syscall.SIGINT)

	loopend := make(chan struct{})
	go func() {
		ui.Loop()
		loopend <- struct{}{}
	}()
	select {
	case <-loopend:
		log.Print("Exiting application...")
	case sig := <-cancelchan:
		log.Printf("Caught signal: %v", sig)
	}
}
