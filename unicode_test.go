package encoder_test

import (
	"slices"
	"testing"

	"github.com/msultra/encoder"
)

func TestStrToUTF16(t *testing.T) {
	str := "Test\x00"
	exp := []byte{'T', '\x00', 'e', '\x00', 's', '\x00', 't', '\x00', '\x00', '\x00'}
	res := encoder.StrToUTF16(str)
	if slices.Compare(res, exp) != 0 {
		t.Errorf("%v, expected: %v", res, exp)
	}

	str2 := ""
	res2 := encoder.StrToUTF16(str2)
	if res2 != nil {
		t.Errorf("%v, expected: %v", res2, nil)
	}

	i := len(res2)
	if i != 0 {
		t.Errorf("should be 0")
	}
}

func TestUTF16ToStr(t *testing.T) {
	bst := []byte{'T', '\x00', 'e', '\x00', 's', '\x00', 't', '\x00'}
	exp := "Test"
	res := encoder.UTF16ToStr(bst)
	if res != exp {
		t.Errorf("%v, expected: %v", res, exp)
	}
}
