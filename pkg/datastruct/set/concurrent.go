package set

import "github.com/mirage208/redis-go/pkg/datastruct/dict"

type ConcurrentSet struct {
	dict dict.Dict
}

// NewConcurrentSet creates a new set which is concurrent safe
func NewConcurrentSet(members ...string) *ConcurrentSet {
	set := &ConcurrentSet{
		dict: dict.NewConcurrentDict(1),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

// Add adds member into set
func (set *ConcurrentSet) Add(val string) int {
	return set.dict.Put(val, nil)
}

// Remove removes member from set
func (set *ConcurrentSet) Remove(val string) int {
	_, ret := set.dict.Remove(val)
	return ret
}

// Has returns true if the val exists in the set
func (set *ConcurrentSet) Has(val string) bool {
	if set == nil || set.dict == nil {
		return false
	}
	_, exists := set.dict.Get(val)
	return exists
}

// Len returns number of members in the set
func (set *ConcurrentSet) Len() int {
	if set == nil || set.dict == nil {
		return 0
	}
	return set.dict.Len()
}

// ToSlice convert set to []string
func (set *ConcurrentSet) ToSlice() []string {
	slice := make([]string, set.Len())
	i := 0
	set.dict.ForEach(func(key string, val any) bool {
		if i < len(slice) {
			slice[i] = key
		} else {
			// set extended during traversal
			slice = append(slice, key)
		}
		i++
		return true
	})
	return slice
}

// ForEach visits each member in the set
func (set *ConcurrentSet) ForEach(consumer func(member string) bool) {
	if set == nil || set.dict == nil {
		return
	}
	set.dict.ForEach(func(key string, val any) bool {
		return consumer(key)
	})
}

// ShallowCopy copies all members to another set
func (set *ConcurrentSet) ShallowCopy() *ConcurrentSet {
	result := NewConcurrentSet()
	set.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	return result
}

// RandomDistinctMembers randomly returns keys of the given number, won't contain duplicated key
func (set *ConcurrentSet) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
