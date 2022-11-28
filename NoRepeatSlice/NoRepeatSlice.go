package NoRepeatSlice

import "fmt"

//NoRepeatSlice 是一个在入值时就避免重复的slice
type NoRepeatSlice struct {
	theMap map[string]byte
	items  []interface{}
}

func (nrs *NoRepeatSlice) Append(item interface{}) {
	key := fmt.Sprintf("%s", item)
	if _, ok := nrs.theMap[key]; !ok {
		nrs.theMap[key] = 0
		nrs.items = append(nrs.items, item)
	}
}

func (nrs NoRepeatSlice) Do(f func(interface{}) bool) {
	for _, item := range nrs.items {
		if !f(item) {
			break
		}
	}
}

func (nrs *NoRepeatSlice) Clear() {
	nrs.theMap = map[string]byte{}
	nrs.items = []interface{}{}
}

func NewNoRepeatSlice() *NoRepeatSlice {
	nrs := NoRepeatSlice{
		theMap: map[string]byte{},
		items:  []interface{}{},
	}
	return &nrs
}
