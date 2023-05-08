package state

import . "luago/api"

func (l *luaState) PushGoFunction(fn GoFunction) {
	l.stack.push(newGoClosure(fn))
}

func (l *luaState) IsGoFunction(idx int) bool {
	val := l.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

func (l *luaState) ToGoFunction(idx int) GoFunction {
	val := l.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}

func (l *luaState) PushNil()            { l.stack.push(nil) }
func (l *luaState) PushBoolean(b bool)  { l.stack.push(b) }
func (l *luaState) PushInteger(i int64) { l.stack.push(i) }

func (l *luaState) PushNumber(f float64) { l.stack.push(f) }
func (l *luaState) PushString(s string)  { l.stack.push(s) }
