package vhaline

import (
	"fmt"

	"github.com/glycerine/rbtree"
)

// wtree is a sorted, ordered map
type wtree struct {
	*rbtree.Tree
}

// an waituntil lives inside an wtree.

func (t *wtree) insert(j *waituntil) {
	t.Insert(j)
}

func (t *wtree) del(matchme int64) {
	t.DeleteWithKey(&waituntil{matchme: matchme})
}

func (t *wtree) get(matchme int64) *waituntil {
	item := t.Get(&waituntil{matchme: matchme})
	if item == nil {
		return nil
	}
	return item.(*waituntil)
}

func (t *wtree) has(matchme int64) bool {
	item := t.Get(&waituntil{matchme: matchme})
	if item == nil {
		return false
	}
	if item.(*waituntil) == nil {
		return false
	}
	return true
}

func newWtree() *wtree {

	return &wtree{rbtree.NewTree(
		func(a1, b2 rbtree.Item) int {
			a := a1.(*waituntil)
			b := b2.(*waituntil)
			return int(a.matchme - b.matchme)
		})}
}

func (t *wtree) deleteThrough(x int64, callme func(goner *waituntil, through int64)) {
	for it := t.Min(); !it.Limit(); {
		cur := it.Item().(*waituntil)
		if cur.matchme <= x {
			//v("mtree2.deleteThrough deletes %#v\n", cur)

			next := it.Next()
			t.DeleteWithIterator(it)
			if callme != nil {
				callme(cur, x)
			}
			it = next
		} else {
			//fmt.Printf("delete pass ignores %#v\n", cur)
			break // we can stop scanning now.
		}
	}
}

func (t *wtree) String() string {
	s := ""
	for it := t.Min(); !it.Limit(); it = it.Next() {
		cur := it.Item().(*waituntil)
		s += fmt.Sprintf("%s\n", cur)
	}
	return s
}
