package list

// Expected check whether given item is equals to expected value
type Expected func(a any) bool

// Consumer traverses list.
// It receives index and value as params, returns true to continue traversal, while returns false to break
type Consumer func(i int, v any) bool

type List interface {
	Add(val any)
	Get(index int) (val any)
	Set(index int, val any)
	Insert(index int, val any)
	Remove(index int) (val any)
	RemoveLast() (val any)
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	ReverseRemoveByVal(expected Expected, count int) int
	Len() int
	ForEach(consumer Consumer)
	Contains(expected Expected) bool
	Range(start int, stop int) []any
}
