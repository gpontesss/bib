package main

import (
	"fmt"

	"github.com/gpontesss/bib/bib"
)

func main() {
	version, err := bib.VersionFromTSV("./kjv.tsv")
	if err != nil {
		panic(err)
	}

	fmt.Println(&version.Books[0].Chapters[0].Verses[0])
}
