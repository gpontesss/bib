package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gpontesss/bib/bib"
	"github.com/gpontesss/bib/encoding/tsv"
)

// FlagArray docs here.
type FlagArray []string

// Value docs here.
func (arrflag *FlagArray) String() string {
	return strings.Join(*arrflag, ",")
}

// Set docs here.
func (arrflag *FlagArray) Set(val string) error {
	*arrflag = append((*arrflag), val)
	return nil
}

// GetVersionPaths docs here.
func GetVersionPaths() []string {
	var vsrpaths FlagArray
	flag.Var(&vsrpaths, "version", "Path to version to be loaded (.tsv)")
	flag.Parse()
	return []string(vsrpaths)
}

// GetVersions docs here.
func GetVersions() ([]*bib.Version, error) {
	filenames := GetVersionPaths()
	vsrs := make([]*bib.Version, len(filenames))

	for i, filename := range filenames {
		if vsr, err := LoadVersion(filename); err != nil {
			return nil, err
		} else {
			vsrs[i] = &vsr
		}
	}
	return vsrs, nil
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
