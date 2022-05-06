package bib

// NewTree returns a new tree with no branches.
func NewTree[K comparable, V any]() *Branch[K, V] {
	return &Branch[K, V]{subbranches: []*Branch[K, V]{}}
}

// Branch gathes a matching key, an optional value (ending branches should
// always have a value; root branch should never have a value) and its
// subbranches.
type Branch[K comparable, V any] struct {
	key         K
	val         *V
	subbranches []*Branch[K, V]
}

// Search searches for a key in a branch and its subbranches. If a unumbiguous value
// is found, it is returned; otherwise, yields nil. An unumbiguous value is only
// asserted with partial keys. A partial key doesn't contain all the values
// needed to match an entire key, but can match a value if there's only one
// subbranch in the branch that it reached.
// TODO: include a drawing illustrating this concept.
func (branch *Branch[K, V]) Search(key []K) *V {
	// no keys provided, or recursion exausted all given keys. search continues
	// to see if it's a branch without sub-branches. if it is, returns only
	// possible value; if not, returns nil.
	if !(len(key) > 0) {
		// no branches left, (len=0), reached the end;
		// there's multiple choices, (len>1) which should not yield an ambiguous
		// value.
		if len(branch.subbranches) != 1 {
			return branch.val
		}
		// there's only one branch and no intermediary value; continues search
		// in it.
		if branch.val == nil {
			return (branch.subbranches[0]).Search(key)
		}
		// there's only one subbranch, but with a value in an intermediary
		// subbranch. it should not be returned, since it would be an ambiguous value.
		return nil
	}
	// key is provided; search continues by matching first value from K, trying
	// to match it with a branch. if a match is found, search continues within
	// the branch with the rest of the values in the key
	first := key[0]
	for _, subbranch := range branch.subbranches {
		if subbranch.key == first {
			tail := key[1:]
			// ending subbranch has a value attached, so it can be returned.
			if len(tail) == 0 && subbranch.val != nil {
				return subbranch.val
			}
			// ending subbranch has no value, so it proceeds to check if there's
			// single subbranch left.
			return subbranch.Search(tail)
		}
	}
	// key value wasn't found in no branch.
	return nil
}

// Insert inserts a key in a branch, creating subbranches necessary to
// accomodate its path to given a value, which will be attached to the end of
// the branch.
func (branch *Branch[K, V]) Insert(key []K, val *V) {
	if !(len(key) >= 1) {
		panic("Insert: key requires at list one item")
	}
	// searches for an existent branch with a matching first key value.
	first := key[0]
	var selected *Branch[K, V]
	for _, subbranch := range branch.subbranches {
		if subbranch.key == first {
			selected = subbranch
			break
		}
	}
	// no branch with current key value, so one must be created and appended to
	// branch's subbranches.
	if selected == nil {
		selected = &Branch[K, V]{
			key:         first,
			subbranches: []*Branch[K, V]{}}
		branch.subbranches = append(branch.subbranches, selected)
	}
	// final key was reached; value should be set in final branch.
	if len(key) == 1 {
		selected.val = val
		return
	}
	// key has length greater than one, so insert must proceed in selected
	// branch.
	selected.Insert(key[1:], val)
}
