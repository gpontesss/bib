package bib

import "fmt"

// ParseRef docs here.
// TODO
func ParseRef(str string) (Ref, error) {
	return Ref{}, nil
}

// Ref docs here.
type Ref struct {
	BookName   string
	ChapterNum int
	VerseNum   int
	Offset     int
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
	verses := make([]*Verse, 0, ref.Offset)
	for i := ref.VerseNum; i <= ref.VerseNum+ref.Offset; i++ {
		verse := vsr.
			GetBook(ref.BookName).
			GetChapter(ref.ChapterNum).
			GetVerse(i)
		if verse != nil {
			verses = append(verses, verse)
		}
	}
	return verses
}
