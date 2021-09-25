package tsv

import (
	"io"
	"strconv"
	"strings"

	"github.com/gpontesss/bib/bib"
)

// Decode docs here.
// TODO: something about metadata in TSV.
// TODO: helpful errors.
func Decode(rdr io.Reader, name string) (bib.Version, error) {
	version := bib.Version{
		// TODO: include metadata in file and extract it from there.
		Name:  name,
		Books: []bib.Book{},
	}

	// not the most efficient way of doing it.
	var err error
	var bs []byte
	if bs, err = io.ReadAll(rdr); err != nil {
		return bib.Version{}, err
	}

	// for now, it reads weirdly linebreaks. it should work for all cases,
	// though.
	lines := strings.Split(string(bs), "\r\n")
	for i := range lines {
		line := lines[i]
		// ignores empty lines
		if len(line) <= 0 {
			continue
		}

		lineparts := strings.Split(line, "\t")
		bookname, chapnumstr, versenumstr, text := lineparts[0], lineparts[1], lineparts[2], lineparts[3]

		var book *bib.Book
		if book = version.GetBook(bookname); book == nil {
			version.Books = append(version.Books, bib.Book{
				Version:  &version,
				Name:     bookname,
				Chapters: []bib.Chapter{},
			})
			book = &version.Books[len(version.Books)-1]
		}

		var chapnum, versenum int
		if chapnum, err = strconv.Atoi(chapnumstr); err != nil {
			return bib.Version{}, err
		} else if versenum, err = strconv.Atoi(versenumstr); err != nil {
			return bib.Version{}, err
		}

		var chap *bib.Chapter
		if chap = book.GetChapter(chapnum); chap == nil {
			book.Chapters = append(book.Chapters, bib.Chapter{
				Book:   book,
				Number: chapnum,
				Verses: []bib.Verse{},
			})
			chap = &book.Chapters[len(book.Chapters)-1]
		}

		chap.Verses = append(chap.Verses, bib.Verse{
			Chapter: chap,
			Number:  versenum,
			Text:    text,
		})
	}

	return version, nil
}
