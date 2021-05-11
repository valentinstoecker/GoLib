package turing

import (
	"fmt"
	"strings"
)

// Char character of turing machine
type Char string

// State of turing machine
type State string

// Dirs of TM
const (
	Left    = "L"
	Right   = "R"
	Neutral = "N"
)

// Dir of TM
type Dir string

// TMState combined state of turing machine
type TMState struct {
	State
	Char
}

// TMAct action of turing machine
type TMAct struct {
	State
	Char
	Dir
}

// Trans transition function
type Trans map[TMState]TMAct

type Band map[int]Char

// TM turing machine
type TM struct {
	States []State
	Alph   []Char
	BAlph  []Char
	Trans  Trans
	SState State
	Empty  Char
	Akt    []State
	band   Band
	bInd   int
	cState State
	lb     int
	ub     int
}

func NewMachine(s []State, a []Char, ba []Char, t Trans, ss State, e Char, ac []State) *TM {
	return &TM{
		States: s,
		Alph:   a,
		BAlph:  ba,
		Trans:  t,
		SState: ss,
		Empty:  e,
		Akt:    ac,
		bInd:   0,
		band:   Band{},
		cState: ss,
		lb:     -10,
		ub:     10,
	}
}

func (tm *TM) Init(band Band) {
	tm.band = band
	tm.bInd = 0
	tm.cState = tm.SState
}

func (tm *TM) GetBChar() Char {
	return tm.GetChar(tm.bInd)
}

func (tm *TM) GetChar(i int) Char {
	if tm.band[i] == "" {
		tm.band[i] = tm.Empty
	}
	return tm.band[i]
}

func (tm *TM) Step() bool {
	st := TMState{
		State: tm.cState,
		Char:  tm.GetBChar(),
	}
	act := tm.Trans[st]

	eAct := TMAct{}

	if act == eAct {
		return false
	}

	tm.band[tm.bInd] = act.Char
	tm.cState = act.State
	switch act.Dir {
	case Left:
		tm.bInd--
	case Right:
		tm.bInd++
	default:
	}
	return true
}

func (tm *TM) printHead() {
	str := ""
	for i := tm.lb; i < tm.ub; i++ {
		if i == tm.bInd {
			str += "!"
		} else {
			str += " "
		}
	}
	fmt.Println(str)
}

func (tm *TM) printBand() {
	str := ""
	for i := tm.lb; i < tm.ub; i++ {
		str += string(tm.GetChar(i))
	}
	fmt.Println(str)
}

func (tm *TM) Print() {
	fmt.Println("------------------------------------------")
	tm.printHead()
	tm.printBand()
	fmt.Println("Current State:", tm.cState)
}

func parseLine(l string) (TMState, TMAct) {
	strs := strings.Split(l, ";")
	if len(strs) != 5 {
		fmt.Println(strs)
		panic("Hilfe falsches format")
	}
	tms := TMState{}
	tms.State = State(strs[0])
	tms.Char = Char(strs[1])

	tma := TMAct{}
	tma.State = State(strs[2])
	tma.Char = Char(strs[3])
	switch strs[4] {
	case "L":
		tma.Dir = Left
	case "R":
		tma.Dir = Right
	default:
		tma.Dir = Neutral
	}
	return tms, tma
}

func ParseTrans(s string) Trans {
	trans := make(Trans)
	strs := strings.Split(s, "\n")
	for _, str := range strs {
		s, a := parseLine(str)
		trans[s] = a
	}
	return trans
}

// Match tests string
func (tm *TM) Match(str string, print bool) bool {
	band := make(Band)
	for i, c := range str {
		band[i] = Char([]rune{c})
	}
	tm.Init(band)
	if print {
		tm.Print()
	}
	fmt.Scanln()
	for tm.Step() {
		if print {
			tm.Print()
		}
		fmt.Scanln()
	}

	for _, ze := range tm.Akt {
		if ze == tm.cState {
			return true
		}
	}
	return false
}
