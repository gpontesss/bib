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
	filenames := []string{"./kjv.tsv", "./vul.tsv", "./sep.tsv"}
	vsrs := make([]*bib.Version, len(filenames))

	for i, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("os.Open", err)
		}
		version, err := tsv.Decode(file, filename)
		if err != nil {
			log.Fatal("tsv.Decode", err)
		}
		vsrs[i] = &version
		file.Close()
	}

	ui := ui.UI{Versions: vsrs}
	if err := ui.Init(); err != nil {
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