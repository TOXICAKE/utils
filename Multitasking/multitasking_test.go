package Multitasking

import (
	"fmt"
	"testing"
)

func TestNewMultitasking(t *testing.T) {
	mt := NewMultitasking("TEST", 10)
	run, err := mt.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(run)
	}
}
