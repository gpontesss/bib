package bib

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRef(t *testing.T) {
	cases := []struct {
		input               string
		expectedRef         Ref
		expectedErrContains string
	}{
		{"Genesis", Ref{BookName: "Genesis", Type: BookRef}, ""},
		{"Isaiah 53", Ref{BookName: "Isaiah", ChapterNum: 53, Type: ChapterRef}, ""},
		{"John 1:1", Ref{BookName: "John", ChapterNum: 1, VerseNum: 1, Type: SingleVerseRef}, ""},
		{"1 Cor. 6:9-11", Ref{BookName: "1 Cor.", ChapterNum: 6, VerseNum: 9, EndVerseNum: 11, Type: RangeVerseRef}, ""},
		{"Apo. 21:3, 4", Ref{BookName: "Apo.", ChapterNum: 21, VerseNums: []int{3, 4}, Type: ListVerseRef}, ""},
		{"John 1:a", Ref{}, "invalid reference"},
		{"John 1:2-1", Ref{}, "higher index (1) must be higher than lower (2)"},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			ref, err := ParseRef(c.input)

			if c.expectedErrContains != "" &&
				!strings.Contains(err.Error(), c.expectedErrContains) {
				t.Fatalf("Expected error to contain %v", c.expectedErrContains)
			}
			if !cmp.Equal(c.expectedRef, ref) {
				t.Fatalf("Expected ref %v; received %v", c.expectedRef, ref)
			}
		})
	}
}
