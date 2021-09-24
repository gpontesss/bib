package main

import (
	"log"
	"os"

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

	ui.Loop()
}
