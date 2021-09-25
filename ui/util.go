package ui

import "unicode"

// NewWordIter docs here.
func NewWordIter(str string) WordIter {
	return WordIter{str, 0, 0}
}

// WordIter docs here.
type WordIter struct {
	str         string
	lowi, highi int
}

// Next docs here.
func (iter *WordIter) Next() bool {
	if iter.highi == len(iter.str) {
		return false
	}

	iter.lowi = iter.highi
	for iter.lowi < len(iter.str) &&
		unicode.IsSpace(rune(iter.str[iter.lowi])) {
		iter.lowi++
	}

	iter.highi = iter.lowi
	for iter.highi < len(iter.str) &&
		(unicode.IsLetter(rune(iter.str[iter.highi])) ||
			unicode.IsPunct(rune(iter.str[iter.highi]))) {
		iter.highi++
	}

	return iter.lowi != iter.highi
}

// Value docs here.
func (iter *WordIter) Value() string { return iter.str[iter.lowi:iter.highi] }
