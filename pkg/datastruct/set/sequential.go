package set

import "github.com/mirage208/redis-go/pkg/datastruct/dict"

type SequentialSet struct {
	dict dict.Dict
}

// NewSequentialSet creates a new set
func NewSequentialSet(members ...string) *SequentialSet {
	set := &SequentialSet{
		dict: dict.NewSequentialDict(),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

// Add adds member into set
func (set *SequentialSet) Add(val string) int {
	return set.dict.Put(val, nil)
}

// Remove removes member from set
func (set *SequentialSet) Remove(val string) int {
	_, ret := set.dict.Remove(val)
	return ret
}

// Has returns true if the val exists in the set
func (set *SequentialSet) Has(val string) bool {
	if set == nil || set.dict == nil {
		return false
	}
	_, exists := set.dict.Get(val)
	return exists
}

// Len returns number of members in the set
func (set *SequentialSet) Len() int {
	if set == nil || set.dict == nil {
		return 0
	}
	return set.dict.Len()
}

// ToSlice convert set to []string
func (set *SequentialSet) ToSlice() []string {
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
func (set *SequentialSet) ForEach(consumer func(member string) bool) {
	if set == nil || set.dict == nil {
		return
	}
	set.dict.ForEach(func(key string, val any) bool {
		return consumer(key)
	})
}

// ShallowCopy copies all members to another set
func (set *SequentialSet) ShallowCopy() *SequentialSet {
	result := NewSequentialSet()
	set.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	return result
}

// RandomMembers randomly returns keys of the given number, may contain duplicated key
func (set *SequentialSet) RandomMembers(limit int) []string {
	if set == nil || set.dict == nil {
		return nil
	}
	return set.dict.RandomKeys(limit)
}

// RandomDistinctMembers randomly returns keys of the given number, won't contain duplicated key
func (set *SequentialSet) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
