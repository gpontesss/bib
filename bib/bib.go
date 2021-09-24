package bib

import (
	"fmt"
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

// GetVerse docs here.
func (chap *Chapter) GetVerse(num int) *Verse {
	for i := range chap.Verses {
		if verse := &chap.Verses[i]; verse.Number == num {
			return verse
		}
	}
	return nil
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
