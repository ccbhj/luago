package state

func (l *luaState) PushNil()             { l.stack.push(nil) }
func (l *luaState) PushBoolean(b bool)   { l.stack.push(b) }
func (l *luaState) PushInteger(i int64)  { l.stack.push(i) }
func (l *luaState) PushNumber(n float64) { l.stack.push(n) }
func (l *luaState) PushString(s string)  { l.stack.push(s) }
