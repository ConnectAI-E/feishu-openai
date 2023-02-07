package calc

import (
	"testing"
)

func TestCalc(t *testing.T) {

	out, err := CalcStr("1+1")
	if err != nil {
		t.Error(err)
	}

	if out != 2 {
		t.Error("1+1 should be 2")
	}
}

func TestCalc2(t *testing.T) {

	out, err := CalcStr("1+2")
	if err != nil {
		t.Error(err)
	}

	if out != 3 {
		t.Error("1+2 should be 3")
	}
}

func TestCalc3(t *testing.T) {
	//22*32
	out, err := CalcStr("22*32")
	if err != nil {
		t.Error(err)
	}

	if out != 704 {
		t.Error("22*32 should be 704")
	}
}
