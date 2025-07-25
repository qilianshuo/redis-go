package sortedset

import (
	"math/bits"
	"math/rand"
)

const (
	maxLevel = 16
)

// Element is a key-score pair
type Element struct {
	Member string
	Score  float64
}

// Level aspect of a node
type Level struct {
	forward *node // forward node has greater score
	span    int64
}

type node struct {
	Element
	backward *node
	level    []*Level // level[0] is base level
}

type SkipList struct {
	header *node
	tail   *node
	length int64
	level  int16
}

func makeNode(level int16, score float64, member string) *node {
	n := &node{
		Element: Element{
			Score:  score,
			Member: member,
		},
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = new(Level)
	}
	return n
}

func makeSkipList() *SkipList {
	return &SkipList{
		level:  1,
		header: makeNode(maxLevel, 0, ""),
	}
}

func randomLevel() int16 {
	total := uint64(1)<<uint64(maxLevel) - 1
	k := rand.Uint64() % total
	return maxLevel - int16(bits.Len64(k+1)) + 1
}

func (skipList *SkipList) insert(member string, score float64) *node {
	update := make([]*node, maxLevel) // link new node with node in `update`
	rank := make([]int64, maxLevel)

	// find position to insert
	node := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		if i == skipList.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] // store rank that is crossed to reach the insert position
		}
		if node.level[i] != nil {
			// traverse the skip list
			for node.level[i].forward != nil &&
				(node.level[i].forward.Score < score ||
					(node.level[i].forward.Score == score && node.level[i].forward.Member < member)) { // same score, different key
				rank[i] += node.level[i].span
				node = node.level[i].forward
			}
		}
		update[i] = node
	}

	level := randomLevel()
	// extend SkipList level
	if level > skipList.level {
		for i := skipList.level; i < level; i++ {
			rank[i] = 0
			update[i] = skipList.header
			update[i].level[i].span = skipList.length
		}
		skipList.level = level
	}

	// make node and link into SkipList
	node = makeNode(level, score, member)
	for i := int16(0); i < level; i++ {
		node.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = node

		// update span covered by update[i] as node is inserted here
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// increment span for untouched levels
	for i := level; i < skipList.level; i++ {
		update[i].level[i].span++
	}

	// set backward node
	if update[0] == skipList.header {
		node.backward = nil
	} else {
		node.backward = update[0]
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node
	} else {
		skipList.tail = node
	}
	skipList.length++
	return node
}

/*
 * param node: node to delete
 * param update: backward node (of target)
 */
func (skipList *SkipList) removeNode(node *node, update []*node) {
	for i := int16(0); i < skipList.level; i++ {
		if update[i].level[i].forward == node {
			update[i].level[i].span += node.level[i].span - 1
			update[i].level[i].forward = node.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		skipList.tail = node.backward
	}
	for skipList.level > 1 && skipList.header.level[skipList.level-1].forward == nil {
		skipList.level--
	}
	skipList.length--
}

/*
 * return: has found and removed node
 */
func (skipList *SkipList) remove(member string, score float64) bool {
	/*
	 * find backward node (of target) or last node of each level
	 * their forward need to be updated
	 */
	update := make([]*node, maxLevel)
	node := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			(node.level[i].forward.Score < score ||
				(node.level[i].forward.Score == score &&
					node.level[i].forward.Member < member)) {
			node = node.level[i].forward
		}
		update[i] = node
	}
	node = node.level[0].forward
	if node != nil && score == node.Score && node.Member == member {
		skipList.removeNode(node, update)
		// free x
		return true
	}
	return false
}

/*
 * return: 1 based rank, 0 means member not found
 */
func (skipList *SkipList) getRank(member string, score float64) int64 {
	var rank int64 = 0
	x := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score &&
					x.level[i].forward.Member <= member)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		/* x might be equal to zsl->header, so test if obj is non-NULL */
		if x.Member == member {
			return rank
		}
	}
	return 0
}

/*
 * 1-based rank
 */
func (skipList *SkipList) getByRank(rank int64) *node {
	var i int64 = 0
	n := skipList.header
	// scan from top level
	for level := skipList.level - 1; level >= 0; level-- {
		for n.level[level].forward != nil && (i+n.level[level].span) <= rank {
			i += n.level[level].span
			n = n.level[level].forward
		}
		if i == rank {
			return n
		}
	}
	return nil
}

func (skipList *SkipList) hasInRange(min Border, max Border) bool {
	if min.isIntersected(max) { //是有交集的，则返回false
		return false
	}

	// min > tail
	n := skipList.tail
	if n == nil || !min.less(&n.Element) {
		return false
	}
	// max < head
	n = skipList.header.level[0].forward
	if n == nil || !max.greater(&n.Element) {
		return false
	}
	return true
}

func (skipList *SkipList) getFirstInRange(min Border, max Border) *node {
	if !skipList.hasInRange(min, max) {
		return nil
	}
	n := skipList.header
	// scan from top level
	for level := skipList.level - 1; level >= 0; level-- {
		// if forward is not in range than move forward
		for n.level[level].forward != nil && !min.less(&n.level[level].forward.Element) {
			n = n.level[level].forward
		}
	}
	/* This is an inner range, so the next node cannot be NULL. */
	n = n.level[0].forward
	if !max.greater(&n.Element) {
		return nil
	}
	return n
}

func (skipList *SkipList) getLastInRange(min Border, max Border) *node {
	if !skipList.hasInRange(min, max) {
		return nil
	}
	n := skipList.header
	// scan from top level
	for level := skipList.level - 1; level >= 0; level-- {
		for n.level[level].forward != nil && max.greater(&n.level[level].forward.Element) {
			n = n.level[level].forward
		}
	}
	if !min.less(&n.Element) {
		return nil
	}
	return n
}

// RemoveRange return removed elements
func (skipList *SkipList) RemoveRange(min Border, max Border, limit int) (removed []*Element) {
	update := make([]*node, maxLevel)
	removed = make([]*Element, 0)
	// find backward nodes (of target range) or last node of each level
	node := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil {
			if min.less(&node.level[i].forward.Element) { // already in range
				break
			}
			node = node.level[i].forward
		}
		update[i] = node
	}

	// node is the first one within range
	node = node.level[0].forward

	// remove nodes in range
	for node != nil {
		if !max.greater(&node.Element) { // already out of range
			break
		}
		next := node.level[0].forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		skipList.removeNode(node, update)
		if limit > 0 && len(removed) == limit {
			break
		}
		node = next
	}
	return removed
}

// RemoveRangeByRank 1-based rank, including start, exclude stop
func (skipList *SkipList) RemoveRangeByRank(start int64, stop int64) (removed []*Element) {
	var i int64 = 0 // rank of iterator
	update := make([]*node, maxLevel)
	removed = make([]*Element, 0)

	// scan from top level
	node := skipList.header
	for level := skipList.level - 1; level >= 0; level-- {
		for node.level[level].forward != nil && (i+node.level[level].span) < start {
			i += node.level[level].span
			node = node.level[level].forward
		}
		update[level] = node
	}

	i++
	node = node.level[0].forward // first node in range

	// remove nodes in range
	for node != nil && i < stop {
		next := node.level[0].forward
		removedElement := node.Element
		removed = append(removed, &removedElement)
		skipList.removeNode(node, update)
		node = next
		i++
	}
	return removed
}
