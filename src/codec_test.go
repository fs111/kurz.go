package main

import (
	"testing"
)

func TestEncode_1(t *testing.T) {
	if Encode(1) != "1" {
		t.Error("1 does Encode to '1' ")
	}
}

func TestEncode_16(t *testing.T) {
	if Encode(16) != "g" {
		t.Error("16 does Encode to 'g' ")
	}
}

func TestDecode_g(t *testing.T) {
	if Decode("g") != 16 {
		t.Error("g does Decode to '16' ")
	}
}
