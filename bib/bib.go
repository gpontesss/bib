package bib

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Version docs here.
// TODO: probably some lazy loading in the future.
type Version struct {
	Name  string
	Books []Book
}

// GetBook docs here.
// TODO: probably some fuzzyfinding in the future.
func (vsr *Version) GetBook(name string) *Book {
	for i := range vsr.Books {
		if book := &vsr.Books[i]; book.Name == name {
			return book
		}
	}
	return nil
}

// Book docs here.
type Book struct {
	// Version is a reference to the version which the book belongs to.
	Version  *Version
	Number   int
	Name     string
	Chapters []Chapter
}

// GetChapter docs here.
func (bk *Book) GetChapter(num int) *Chapter {
	for i := range bk.Chapters {
		if chap := &bk.Chapters[i]; chap.Number == num {
			return chap
		}
	}
	return nil
}

// Chapter docs here.
type Chapter struct {
	// Book is a reference to the book which the chapter belongs to.
	Book   *Book
	Number int
	Verses []Verse
}

// Verse docs here.
type Verse struct {
	// Chapter is a reference to the chapter which the verse belongs to.
	Chapter *Chapter
	Number  int
	Text    string
}

// String implements the fmt.Stringer interface.
func (vrs *Verse) String() string {
	return fmt.Sprintf(
		"%s %d:%d %s",
		vrs.Chapter.Book.Name,
		vrs.Chapter.Number,
		vrs.Number,
		vrs.Text)
}

// VersionFromTSV docs here.
// TODO: something about metadata in TSV.
// TODO: helpful errors.
func VersionFromTSV(path string) (Version, error) {
	file, err := os.Open(path)
	if err != nil {
		return Version{}, err
	}
	version := Version{
		// TODO: include metadata in file and extract it from there.
		Name:  file.Name(),
		Books: []Book{},
	}

	// not the most efficient way of doing it.
	var bs []byte
	if bs, err = io.ReadAll(file); err != nil {
		return Version{}, err
	}

	lines := strings.Split(string(bs), "\n")
	for i := range lines {
		line := lines[i]
		// ignores empty lines
		if len(line) <= 0 {
			continue
		}

		lineparts := strings.Split(line, "\t")
		bookname, ref, text := lineparts[0], lineparts[1], lineparts[2]

		var book *Book
		if book = version.GetBook(bookname); book == nil {
			version.Books = append(version.Books, Book{
				Version:  &version,
				Name:     bookname,
				Chapters: []Chapter{},
			})
			book = &version.Books[len(version.Books)-1]
		}

		refparts := strings.Split(ref, ":")
		chapnumstr, versenumstr := refparts[0], refparts[1]

		var chapnum, versenum int
		if chapnum, err = strconv.Atoi(chapnumstr); err != nil {
			return Version{}, err
		} else if versenum, err = strconv.Atoi(versenumstr); err != nil {
			return Version{}, err
		}

		var chap *Chapter
		if chap = book.GetChapter(chapnum); chap == nil {
			book.Chapters = append(book.Chapters, Chapter{
				Book:   book,
				Number: chapnum,
				Verses: []Verse{},
			})
			chap = &book.Chapters[len(book.Chapters)-1]
		}

		chap.Verses = append(chap.Verses, Verse{
			Chapter: chap,
			Number:  versenum,
			Text:    text,
		})
	}

	return version, nil
}
