package bib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type refType int

const (
	// Signs that it is a book reference.
	BookRef refType = iota
	// Signs that it is a chapter reference.
	ChapterRef
	// Signs that the reference contains a single verse.
	SingleVerseRef
	// Signs that the reference contains a verse range.
	RangeVerseRef
	// Signs that the reference contains a list of verses.
	ListVerseRef
)

// Ref represents a generic reference of a verse, or a collection of verses,
// that can be indepently built and applied to versions.
type Ref struct {
	// Type indicates how the reference should be interpreted. See RefType and
	// constants of this type.
	Type refType
	// BookName references the book name. Relevant to all types, for it must
	// always be present.
	BookName string
	// ChapterNum references the chapter number. Relevant for type ChapterRef,
	// SingleVerseRef, RangeVerseRef and ListVerseRef.
	ChapterNum int
	// VerseNum references the lower index when the type is RangeVerseRef, or
	// the verse number when the type is SingleVerseRef.
	VerseNum int
	// EndVerseNum references the higher index when the type is RangeVerseRef.
	EndVerseNum int
	// VerseNums  references unrelated verse numbers when the type is
	// ListVerseRef.
	VerseNums []int
}

// Reference regular expression.
// It considers all possible cases of a reference, and allow spaces between
// parts where it is reasonable.
// TODO: allow unicode matching.
var refre = regexp.MustCompile(
	"^\\s*(\\d?\\s*[\\w]+\\.?)\\s*((\\d+)(:((\\d+)((-)(\\d+)|((,)\\s*\\d+)+)?))?\\s*)?$")

// ParseRef docs here.
// In the future, would be nice to do it with zero allocations.
func ParseRef(str string) (Ref, error) {
	result := refre.FindAllStringSubmatch(str, -1)
	if len(result) != 1 {
		return Ref{}, fmt.Errorf("%q is an invalid reference\n", str)
	}
	// since it matches the edges (^...$), it ever only produces one result
	// match set.
	matches := result[0]

	// Book name should always be matched.
	ref := Ref{BookName: matches[1], Type: BookRef}

	// Chapter number is included
	if chapternum := matches[3]; chapternum != "" {
		// regex ensures that number is valid.
		ref.ChapterNum, _ = strconv.Atoi(matches[3])
		ref.Type = ChapterRef
	}

	// Verse range reference.
	if matches[8] == "-" {
		// regex ensures that number is valid.
		lowi, _ := strconv.Atoi(matches[6])
		highi, _ := strconv.Atoi(matches[9])

		if lowi >= highi {
			return Ref{}, fmt.Errorf(
				"Invalid range: higher index (%d) must be higher than lower (%d)\n",
				highi, lowi)
		}
		ref.VerseNum = lowi
		ref.EndVerseNum = highi
		ref.Type = RangeVerseRef

		// Verse list reference.
	} else if matches[11] == "," {
		numstrs := strings.Split(
			strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return -1
				}
				return r
			}, matches[5]),
			",")

		versenums := make([]int, len(numstrs))
		for i := 0; i < len(versenums); i++ {
			// regex ensures that number is valid.
			versenums[i], _ = strconv.Atoi(numstrs[i])
		}
		ref.VerseNums = versenums
		ref.Type = ListVerseRef

		// Only verse number is included.
	} else if versenumstr := matches[5]; versenumstr != "" {
		// regex ensures that number is valid.
		ref.VerseNum, _ = strconv.Atoi(matches[5])
		ref.Type = SingleVerseRef
	}

	return ref, nil
}

// TODO: deal with the nil pointers
// Book docs here.
func (ref *Ref) Book(vsr *Version) *Book { return vsr.GetBook(ref.BookName) }

// Chapter docs here.
func (ref *Ref) Chapter(vsr *Version) *Chapter {
	return vsr.GetBook(ref.BookName).GetChapter(ref.ChapterNum)
}

// String implements the fmt.Stringer interface.
func (ref *Ref) String() string {
	switch ref.Type {
	case BookRef:
		return ref.BookName
	case ChapterRef:
		return fmt.Sprintf("%s %d", ref.BookName, ref.ChapterNum)
	case SingleVerseRef:
		return fmt.Sprintf("%s %d:%d", ref.BookName, ref.ChapterNum, ref.VerseNum)
	case ListVerseRef:
		// damn you, golang, and your lack of generics to allow mapping and
		// other sane things.
		versenums := make([]string, len(ref.VerseNums))
		for i := range versenums {
			versenums[i] = strconv.Itoa(ref.VerseNums[i])
		}
		return fmt.Sprintf("%s %d:%s",
			ref.BookName, ref.ChapterNum, strings.Join(versenums, ", "))
	case RangeVerseRef:
		return fmt.Sprintf("%s %d:%d-%d",
			ref.BookName, ref.ChapterNum, ref.VerseNum, ref.EndVerseNum)
	default:
		panic("Can't stringify ref type")
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

	if ref.VerseNums == nil {
		// it means all verses in the chapter.
		return chap.VerseRange(1, chap.LastVerse().Number)
	}

	// TODO: refactor
	verses := make([]*Verse, 0, len(ref.VerseNums))
	for _, versenum := range ref.VerseNums {
		if verse := chap.GetVerse(versenum); verse != nil {
			verses = append(verses, verse)
		}
	}
	return verses
}
