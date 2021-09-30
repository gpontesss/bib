package bib

import (
	"fmt"
	"regexp"
	"strconv"
)

var refre = regexp.MustCompile("^(\\w+)\\s+(\\d+):(\\d+)$")

// ParseRef docs here.
// TODO: parse partial references that only contains book/chapter, etc., and
// ranges.
func ParseRef(str string) (Ref, error) {
	matches := refre.FindAllStringSubmatch(str, -1)[0]
	if len(matches) != 3 {
		return Ref{}, fmt.Errorf("Invalid reference '%s'", str)
	}
	bookname := matches[1]
	chapnum, _ := strconv.Atoi(matches[2])
	versenum, _ := strconv.Atoi(matches[3])
	return Ref{
		BookName:   bookname,
		ChapterNum: chapnum,
		VerseNum:   versenum,
	}, nil
}

// Ref docs here.
type Ref struct {
	BookName   string
	ChapterNum int
	VerseNum   int
	Offset     int
}

// TODO: deal with the nil pointers
// Book docs here.
func (ref *Ref) Book(vsr *Version) *Book { return vsr.GetBook(ref.BookName) }

// Chapter docs here.
func (ref *Ref) Chapter(vsr *Version) *Chapter {
	return vsr.GetBook(ref.BookName).GetChapter(ref.ChapterNum)
}

// NormOffset docs here.
func (ref *Ref) NormOffset(vsr *Version) {
	chap := ref.Chapter(vsr)
	if chap == nil {
		return
	}
	lastversenum := chap.LastVerse().Number
	if maxoffset := lastversenum - ref.VerseNum; ref.Offset > maxoffset {
		ref.Offset = maxoffset
	}
}

// String implements the fmt.Stringer interface.
func (ref *Ref) String() string {
	if ref.ChapterNum == 0 {
		return ref.BookName
	} else if ref.VerseNum == 0 {
		return fmt.Sprintf("%s %d", ref.BookName, ref.ChapterNum)
	} else if ref.Offset <= 0 {
		return fmt.Sprintf("%s %d:%d",
			ref.BookName, ref.ChapterNum, ref.VerseNum)
	} else {
		return fmt.Sprintf("%s %d:%d-%d",
			ref.BookName, ref.ChapterNum, ref.VerseNum, ref.VerseNum+ref.Offset)
	}
}

// Verses docs here.
// TODO: account for when ranges are weird.
func (ref *Ref) Verses(vsr *Version) []*Verse {
	book := vsr.GetBook(ref.BookName)
	if book == nil {
		return []*Verse{}
	}
	chap := book.GetChapter(ref.ChapterNum)
	if chap == nil {
		// it means all verses in the book.
		return book.Verses()
	}

	if ref.VerseNum <= 0 {
		// it means all verses in the chapter.
		return chap.VerseRange(1, chap.LastVerse().Number)
	}

	ref.NormOffset(vsr)
	return chap.VerseRange(ref.VerseNum, ref.VerseNum+ref.Offset)
}
