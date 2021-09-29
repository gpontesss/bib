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

// Ref docs here.
func (bk *Book) Ref() Ref {
	return Ref{
		BookName: bk.Name,
	}
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

// Verses docs here.
func (bk *Book) Verses() []*Verse {
	total := 0
	for i := range bk.Chapters {
		total += len(bk.Chapters[i].Verses)
	}

	verses := make([]*Verse, total)
	versei := 0
	for i := range bk.Chapters {
		chap := &bk.Chapters[i]
		for j := range chap.Verses {
			verses[versei] = &chap.Verses[j]
			versei++
		}
	}
	return verses
}

// Chapter docs here.
type Chapter struct {
	// Book is a reference to the book which the chapter belongs to.
	Book   *Book
	Number int
	Verses []Verse
}

// Ref docs here.
func (chap *Chapter) Ref() Ref {
	return Ref{
		BookName:   chap.Book.Name,
		ChapterNum: chap.Number,
	}
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

// VerseRange docs here.
func (chap *Chapter) VerseRange(from, to int) []*Verse {
	verses := make([]*Verse, 0, to-from+1)
	for i := from; i <= to; i++ {
		if verse := chap.GetVerse(i); verse != nil {
			verses = append(verses, verse)
		}
	}
	return verses
}

// Next docs here.
func (chap *Chapter) Next() *Chapter {
	book := chap.Book
	for i := range book.Chapters {
		if chapptr := &book.Chapters[i]; chapptr == chap {
			if i == len(book.Chapters)-1 {
				return nil
			}
			return &book.Chapters[i+1]
		}
	}
	return nil
}

// Previous docs here.
func (chap *Chapter) Previous() *Chapter {
	book := chap.Book
	for i := range book.Chapters {
		if chapptr := &book.Chapters[i]; chapptr == chap {
			if i == 0 {
				return nil
			}
			return &book.Chapters[i-1]
		}
	}
	return nil
}

// LastVerse docs here.
func (chap *Chapter) LastVerse() *Verse {
	if len(chap.Verses) <= 0 {
		return nil
	}

	last := &chap.Verses[0]
	for i := 1; i < len(chap.Verses); i++ {
		verse := &chap.Verses[i]
		if last.Number < verse.Number {
			last = verse
		}
	}
	return last
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
	ref := vrs.Ref()
	return fmt.Sprintf(
		"%s %s", ref.String(), vrs.Text)
}

// Ref docs here.
func (vrs *Verse) Ref() Ref {
	return Ref{
		BookName:   vrs.Chapter.Book.Name,
		ChapterNum: vrs.Chapter.Number,
		VerseNum:   vrs.Number,
	}
}
