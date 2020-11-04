package prover

import "fmt"

type typeID int
type varID int

func (v varID) isVar() bool {
	return true
}

func (v varID) isEl() bool {
	return false
}

func (v varID) getEl() *el {
	return nil
}

func (v varID) getVarID() varID {
	return v
}

const (
	eq = iota
	add
	or
	and
	not
)

type el struct {
	l, r child
	t    int
}

func (e *el) isVar() bool {
	return false
}

func (e *el) isEl() bool {
	return true
}

func (e *el) getEl() *el {
	return e
}

func (e *el) getVarID() varID {
	return -1
}

type rel struct {
	l, r child
}

type child interface {
	isVar() bool
	isEl() bool
	getEl() *el
	getVarID() varID
}

func newRel(t int, l, r child) *el {
	return &el{
		t: t,
		l: l,
		r: r,
	}
}

var notID varID

func newNot(c child) *el {
	if notID == 0 {
		notID = newVar()
	}
	return &el{
		t: not,
		l: notID,
		r: c,
	}
}

var cVarID varID

func newVar() varID {
	cVarID++
	return cVarID
}

type varMap map[varID]child

func equals(c1, c2 child) bool {
	if c1.isEl() != c2.isEl() {
		return false
	}
	if c1.isVar() {
		return c1.getVarID() == c2.getVarID()
	}
	el1, el2 := c1.getEl(), c2.getEl()
	if el1.t != el2.t {
		return false
	}
	return equals(el1.l, el2.l) && equals(el1.r, el2.r)
}

func match(pattern, matchIn child, vm varMap) bool {
	if pattern.isEl() && matchIn.isVar() {
		return false
	}
	if pattern.isEl() && matchIn.isEl() {
		patEl := pattern.getEl()
		matchEl := matchIn.getEl()
		if patEl.t != matchEl.t {
			return false
		}
		return match(patEl.l, matchEl.l, vm) && match(patEl.r, matchEl.r, vm)
	}
	if pattern.isVar() && matchIn.isEl() {
		if vm[pattern.getVarID()] == nil {
			vm[pattern.getVarID()] = matchIn.getEl()
			return true
		}
		return equals(vm[pattern.getVarID()], matchIn.getEl())
	}
	if vm[pattern.getVarID()] == nil {
		vm[pattern.getVarID()] = matchIn.getVarID()
		return true
	}
	return vm[pattern.getVarID()] == matchIn.getVarID()
}

func replace(pattern child, vm varMap) child {
	if pattern.isVar() {
		return vm[pattern.getVarID()]
	}
	patEl := pattern.getEl()
	return &el{
		t: patEl.t,
		l: replace(patEl.l, vm),
		r: replace(patEl.r, vm),
	}
}

func comp(c1, c2 child) int {
	if c1.isVar() && c2.isVar() {
		return int(c1.getVarID()) - int(c2.getVarID())
	}
	if c1.isEl() && c2.isEl() {
		c1El := c1.getEl()
		c2El := c2.getEl()
		if c1El.t != c2El.t {
			return int(c1El.t) - int(c2El.t)
		}
		compL := comp(c1El.l, c2El.l)
		if compL != 0 {
			return compL
		}
		return comp(c1El.r, c2El.r)
	}
	if c1.isVar() && c2.isEl() {
		return -1
	}
	return 1
}

func Test() {
	a := newVar()
	b := newVar()
	c := newVar()
	d := newVar()
	e := newVar()
	/* kommEq := &rel{
		l: newRel(eq, a, b),
		r: newRel(eq, b, a),
	}
	assAdd := &rel{
		l: newRel(add, newRel(add, a, b), c),
		r: newRel(add, a, newRel(add, b, c)),
	}
	kommAdd := &rel{
		l: newRel(add, a, b),
		r: newRel(add, b, a),
	}
	dNeg := &rel{
		l: newNot(newNot(a)),
		r: a,
	} */

	m := make(varMap)

	aEqB := newRel(eq, a, b)
	bEqA := newRel(eq, b, a)

	test := newRel(eq, c, newRel(add, d, e))

	fmt.Println(test)
	fmt.Println(match(aEqB, test, m))
	fmt.Println(m)
	fmt.Println(replace(bEqA, m))
}
