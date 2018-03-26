package ahocorasick

import (
	"container/list"
	"log"
	"unsafe"
)

type trieNode struct {
	count uint32
	fail  *trieNode
	child [256]*trieNode
	index uint32
}

func (t *trieNode) print(offset int) {

	id := (uintptr)(unsafe.Pointer(t))
	prefix := ""
	for i := 0; i < offset; i++ {
		prefix += "\t"
	}
	log.Printf("%s[Node] %q count: %v index: %v", prefix, id, t.count, t.index)
	for _, c := range t.child {
		log.Printf("%schild", prefix)
		c.print(offset + 1)
	}
	log.Printf("%sfail: %v", prefix, (uintptr)(unsafe.Pointer(t.fail)))
}

func newTrieNode() *trieNode {
	return &trieNode{
		count: 0,
		fail:  nil,
		child: [256]*trieNode{},
		index: 0,
	}
}

type Matcher struct {
	root *trieNode
	size uint32
}

func NewMatcher() *Matcher {
	return &Matcher{
		root: newTrieNode(),
		size: 0,
	}
}

// initialize the ahocorasick
func (this *Matcher) Build(dictionary []string) {
	for i := range dictionary {
		this.insert(dictionary[i])
	}
	this.build()
}

// string match search
// return all strings matched as indexes into the original dictionary
func (this *Matcher) Match(s string) []uint32 {
	curNode := this.root
	mark := make([]bool, this.size)
	var p *trieNode = nil

	ret := make([]uint32, 0, this.size)

	for _, v := range []byte(s) {
		for curNode.child[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = this.root
		}

		p = curNode
		for p != this.root {
			if p.count > 0 && !mark[p.index] {
				mark[p.index] = true
				for i := uint32(0); i < p.count; i++ {
					ret = append(ret, p.index)
				}
			}
			p = p.fail
		}
	}

	return ret
}

// just return the number of len(Match(s))
func (this *Matcher) GetMatchResultSize(s string) int {
	return len(this.Match(s))
}

func (this *Matcher) build() {
	ll := list.New()
	ll.PushBack(this.root)
	for ll.Len() > 0 {
		temp := ll.Remove(ll.Front()).(*trieNode)
		var p *trieNode = nil

		for i, v := range temp.child {
			if v == nil {
				continue
			}

			if temp == this.root {
				v.fail = this.root
			} else {
				p = temp.fail
				for p != nil {
					if p.child[i] != nil {
						v.fail = p.child[i]
						break
					}
					p = p.fail
				}
				if p == nil {
					v.fail = this.root
				}
			}
			ll.PushBack(v)
		}
	}
}

func (this *Matcher) insert(s string) {
	curNode := this.root
	for _, v := range []byte(s) {
		if curNode.child[v] == nil {
			curNode.child[v] = newTrieNode()
		}
		curNode = curNode.child[v]
	}
	curNode.count++
	curNode.index = this.size
	this.size++
}
