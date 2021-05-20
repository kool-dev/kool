package parser

import (
	"errors"
	"testing"
)

func TestSimilars(t *testing.T) {
	err := &ErrPossibleTypo{[]string{"cmd1", "cmd2"}}

	similars := err.Similars()

	if len(similars) != 2 {
		t.Error("failed to get the typo error similars")
	}

	if similars[0] != "cmd1" || similars[1] != "cmd2" {
		t.Error("failed to get the typo error similars")
	}
}

func TestSetSimilars(t *testing.T) {
	err := &ErrPossibleTypo{[]string{}}

	err.SetSimilars([]string{"cmd1", "cmd2"})

	similars := err.Similars()

	if len(similars) != 2 {
		t.Error("failed to get the typo error similars")
	}

	if similars[0] != "cmd1" || similars[1] != "cmd2" {
		t.Error("failed to get the typo error similars")
	}
}

func TestIsPossibleTypoError(t *testing.T) {
	err := &ErrPossibleTypo{[]string{}}

	isTypoError := IsPossibleTypoError(err)

	if !isTypoError {
		t.Error("failed to assert that error was a typo one")
	}

	isTypoError = IsPossibleTypoError(errors.New("error"))

	if isTypoError {
		t.Error("failed to assert that error was not a typo one")
	}
}
