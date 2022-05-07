package main

import (
	"fmt"
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
		if vsr, err := LoadVersion(filename); err != nil {
			log.Fatal(err)
		} else {
			vsrs[i] = &vsr
		}
	}

	ui := ui.UI{Versions: vsrs}
	defer func() {
		err := recover()
		// forces UI to be ended before logging, to avoid terminal bugging
		// glitch.
		ui.End()
		if err != nil {
			panic(err)
		}
	}()

	if err := ui.Init(); err != nil {
		panic(fmt.Errorf("ui.Init: %v", err))
	}

	cancelchan := make(chan os.Signal, 1)
	signal.Notify(cancelchan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-ui.AsyncLoop():
		if err != nil {
			panic(fmt.Errorf("Exited with error: %v", err))
		}
	case sig := <-cancelchan:
		panic(fmt.Errorf("Caught exitin signal: %v", sig))
	}
}

// LoadVersion docs here.
func LoadVersion(filename string) (bib.Version, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return bib.Version{}, fmt.Errorf("os.Open: %s", err)
	}
	version, err := tsv.Decode(file, filename)
	if err != nil {
		return bib.Version{}, fmt.Errorf("tsv.Decode: %s", err)
	}
	return version, nil
}
