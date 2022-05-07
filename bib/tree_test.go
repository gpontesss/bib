package bib

import (
	"fmt"
	"testing"
)

func strp(str string) *string {
	strp := new(string)
	*strp = str
	return strp
}

func AssertSearchNil[K comparable, V any](
	t *testing.T, tree *Branch[K, V], key []K) {
	if val := tree.Search(key); val != nil {
		t.Errorf("Expected nil, but got %v", *val)
	}
}

func AssertSearch[K comparable, V comparable](
	t *testing.T, tree *Branch[K, V], key []K, expect V) {
	val := tree.Search(key)
	if val == nil || *val != expect {
		valstr := "<nil>"
		if val != nil {
			valstr = fmt.Sprintf("%+v", *val)
		}
		t.Errorf("Expected %+v, got %s", expect, valstr)
	}
}

func TestNode(t *testing.T) {

	tree := NewTree[rune, string]()
	tree.Insert([]rune("hello"), strp("world"))
	tree.Insert([]rune("he"), strp("she"))
	tree.Insert([]rune("goodbye"), strp("ma'am"))

	AssertSearchNil(t, tree, []rune(""))
	AssertSearchNil(t, tree, []rune("none"))
	AssertSearch(t, tree, []rune("hello"), "world")
	AssertSearch(t, tree, []rune("he"), "she")
	AssertSearch(t, tree, []rune("hell"), "world")
	AssertSearchNil(t, tree, []rune("h"))

	word, expect := "goodbye", "ma'am"
	for i := 1; i <= len(word); i++ {
		AssertSearch(t, tree, []rune(word)[:i], expect)
	}
}
