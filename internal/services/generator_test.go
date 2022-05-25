package services

import (
	"strings"
	"testing"
)

func TestGenShort(t *testing.T) {
	//Arrange
	alphabet = []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_`)
	requiredLength := 10
	//Act
	res := genShort()
	realLength := len(res)
	containWrongSymb, wrSymb := func(str string, alph []byte) (bool, rune) {
		for _, symb := range str {
			if !strings.Contains(string(alphabet), string(symb)) {
				return true, symb
			}
		}
		return false, 0
	}(res, alphabet)
	//Assert
	if requiredLength != realLength {
		t.Errorf("Incorrect length of short url: expected %d, got %d", requiredLength, realLength)
	}
	if containWrongSymb {
		t.Errorf("Incorrect symbol in short url: %s", string(wrSymb))
	}
}
