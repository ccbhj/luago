package state

import (
	"math"

	"luago/number"
)

type luaTable struct {
	metatable *luaTable
	arr       []luaValue
	_map      map[luaValue]luaValue
	keys      map[luaValue]luaValue
	changed   bool
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := new(luaTable)
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}

	return t
}

func (t *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(t.arr)) {
			return t.arr[idx-1]
		}
	}

	return t._map[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (t *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(t.arr))
		if idx <= arrLen {
			t.arr[idx-1] = val
			if idx == arrLen && val == nil {
				t._shrinkArray()
			}
			return
		}

		if idx == arrLen+1 {
			delete(t._map, key)
			if val != nil {
				t.arr = append(t.arr, val)
				t._growArray()
			}
			return
		}
	}

	if val == nil {
		delete(t._map, key)
		return
	}
	if t._map == nil {
		t._map = make(map[luaValue]luaValue, 1<<3)
	}
	t._map[key] = val
}

func (t *luaTable) _shrinkArray() {
	for i := len(t.arr) - 1; i >= 0; i-- {
		if t.arr[i] == nil {
			t.arr = t.arr[:i]
		}
	}
}

func (t *luaTable) _growArray() {
	// move value from t._map to idx
	for idx := int64(len(t.arr)) + 1; true; idx++ {
		val, found := t._map[idx]
		if !found {
			break
		}
		delete(t._map, idx)
		t.arr = append(t.arr, val)
	}
}

func (t *luaTable) len() int {
	return len(t.arr)
}

func (t *luaTable) hasMetafield(fieldName string) bool {
	return t.metatable != nil && t.metatable.get(fieldName) != nil
}

func (t *luaTable) nextKey(key luaValue) luaValue {
	if t.keys == nil || key == nil {
		t.initKeys()
		t.changed = true
	}

	return t.keys[key]
}

func (t *luaTable) initKeys() {
	var key luaValue

	t.keys = make(map[luaValue]luaValue)
	for i, v := range t.arr {
		if v != nil {
			t.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}

	for k, v := range t._map {
		if v != nil {
			t.keys[key] = k
			key = k
		}
	}
}
