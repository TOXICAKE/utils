package utils

import (
	"fmt"
	"testing"
)

type A struct {
	Arg1 string
	Arg2 int
	OK   bool
}
type B struct {
	Arg1 string
	Arg2 int
	OK   int
}

func TestCopyStructByName(t *testing.T) {
	a := A{
		Arg1: "aaa",
		Arg2: 3,
		OK:   false,
	}
	b := B{}
	CopyStructByName(a, &b)
	fmt.Println(b.Arg1, b.Arg2)
}

func TestCopyStructString(t *testing.T) {
	a := A{
		Arg1: "",
		Arg2: 3,
		OK:   false,
	}
	b := B{Arg1: "bb"}
	CopyStructString(a, &b, true)
	fmt.Println(b.Arg1, b.Arg2, b.OK)
}
