package logic

import "fmt"

type Var struct {
	Name string
}

func (v Var) String() string {
	return string(v.Name)
}

func (v Var) Traverse(f func(Operator) Operator) Operator {
	return f(v)
}

type Operator interface {
	Traverse(func(Operator) Operator) Operator
}

func Simplify(op Operator) Operator {
	r := op.Traverse(func(o Operator) Operator {
		if v, ok := o.(BiImpl); ok {
			return And{
				Impl{
					v.A,
					v.B,
				},
				Impl{
					v.B,
					v.A,
				},
			}
		}
		return o
	})
	r = r.Traverse(func(o Operator) Operator {
		if v, ok := o.(Impl); ok {
			return And{
				Not{v.A},
				v.B,
			}
		}
		return o
	})
	return r
}

type And struct {
	A, B Operator
}

func (a And) Traverse(f func(Operator) Operator) Operator {
	return f(And{
		a.A.Traverse(f),
		a.B.Traverse(f),
	})
}

func (a And) String() string {
	return fmt.Sprintf("(%v & %v)", a.A, a.B)
}

type Or struct {
	A, B Operator
}

func (o Or) String() string {
	return fmt.Sprintf("(%v | %v)", o.A, o.B)
}

func (o Or) Traverse(f func(Operator) Operator) Operator {
	return f(Or{
		o.A.Traverse(f),
		o.B.Traverse(f),
	})
}

type Impl struct {
	A, B Operator
}

func (i Impl) Traverse(f func(Operator) Operator) Operator {
	return f(Impl{
		A: i.A.Traverse(f),
		B: i.B.Traverse(f),
	})
}

func (i Impl) String() string {
	return fmt.Sprintf("(%v -> %v)", i.A, i.B)
}

type BiImpl struct {
	A, B Operator
}

func (bi BiImpl) Traverse(f func(Operator) Operator) Operator {
	return f(BiImpl{
		bi.A.Traverse(f),
		bi.B.Traverse(f),
	})
}

func (bi BiImpl) String() string {
	return fmt.Sprintf("(%v <-> %v)", bi.A, bi.B)
}

type Not struct {
	A Operator
}

func (n Not) Traverse(f func(Operator) Operator) Operator {
	return f(Not{
		n.A.Traverse(f),
	})
}

func (n Not) String() string {
	return fmt.Sprintf("!%v", n.A)
}
