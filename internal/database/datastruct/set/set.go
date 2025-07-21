package set

// Set is a set of elements based on hash table
type Set interface {
	Add(val string) int
	Remove(val string) int
	Has(val string) bool
	Len() int
	ToSlice() []string
	ForEach(func(member string) bool)
}

// Intersect intersects two sets
func Intersect(sets ...Set) Set {
	result := NewSequentialSet()
	if len(sets) == 0 {
		return result
	}

	countMap := make(map[string]int)
	for _, set := range sets {
		set.ForEach(func(member string) bool {
			countMap[member]++
			return true
		})
	}
	for k, v := range countMap {
		if v == len(sets) {
			result.Add(k)
		}
	}
	return result
}

// Union adds two sets
func Union(sets ...Set) Set {
	result := NewSequentialSet()
	for _, set := range sets {
		set.ForEach(func(member string) bool {
			result.Add(member)
			return true
		})
	}
	return result
}
