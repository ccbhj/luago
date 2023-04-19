package state

func (l *luaState) GetTop() int {
	return l.stack.top
}

func (l *luaState) AbsIndex(idx int) int {
	return l.stack.absIndex(idx)
}

func (l *luaState) CheckStack(n int) bool {
	l.stack.check(n)
	return true
}

func (l *luaState) Pop(n int) {
	for ; n > 0; n-- {
		l.stack.pop()
	}
}

func (l *luaState) Copy(from, to int) {
	v := l.stack.get(from)
	l.stack.set(to, v)
}

// PushValue get the value from `idx` and push it to the top
func (l *luaState) PushValue(idx int) {
	val := l.stack.get(idx)
	l.stack.push(val)
}

// Replace pop the value on stack top and set it at `idx`
func (l *luaState) Replace(idx int) {
	val := l.stack.pop()
	l.stack.set(idx, val)
}

func (l *luaState) Insert(idx int) {
	l.Rotate(idx, 1)
}

func (l *luaState) Remove(idx int) {
	l.Rotate(idx, -1)
	l.Pop(1)
}

/*

insert(2) -> l.rotate(2, 1)

4 e    4*e    4*e    4 d
3*d    3 b    3 b    3 c
2 c -> 2 c -> 2 c -> 2 b
1#b    1 d    1#d    1 e
0 a    0 a    0 a    0 a

n = 1
t = 5 - 1 = 4
p = 2 - 1 = 1
m = t - 1 = 3

remove(2) -> l.rotate(2, -1)

4 e    4*e    4*c    4 b
3 d    3 d    3 d    3 e
2 c -> 2#c -> 2 e -> 2 d
1*b    1 b    1#b    1 c
0 a    0 a    0 a    0 a

idx = 2
n = -1
top = 5
t = 5 - 1 = 4
p = 2 - 1 = 1
m = 1 -(-1) - 1 = 1
*/

func (l *luaState) Rotate(idx, n int) {
	t := l.stack.top - 1
	p := l.stack.absIndex(idx) - 1

	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	l.reverse(p, m)
	l.reverse(m+1, t)
	l.reverse(p, t)
}

func (l *luaState) reverse(from, to int) {
	slots := l.stack.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

// SetTop set stack's top to `idx`
// if current top < idx, push nil to stack
// if current top > idx, pop top util current top is idx
func (l *luaState) SetTop(idx int) {
	newTop := l.AbsIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := l.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			l.stack.pop()
		}
		return
	}
	for i := 0; i > n; i-- {
		l.stack.push(nil)
	}
}
