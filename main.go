package main

import (
	"log"

	"github.com/gpontesss/bib/bib"
	"github.com/gpontesss/bib/ui"
)

func main() {
	version, err := bib.VersionFromTSV("./kjv.tsv")
	if err != nil {
		log.Fatal("bib.VersionFromTSV", err)
	}

	ui := ui.UI{Versions: []bib.Version{version, version}}
	if ui.Init() != nil {
		log.Fatal("ui.Init", err)
	}
	defer ui.End()

	ui.Loop()
}
