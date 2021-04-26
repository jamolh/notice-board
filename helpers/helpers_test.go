package helpers_test

import (
	"testing"

	"github.com/jamolh/notice-board/helpers"
)

func TestGenerateRandomHash(t *testing.T) {
	str := "sf@f3 9#2Ln_r"
	trueResult := "sffLnr"
	if result := helpers.RemoveNonLetter(str); result != trueResult {
		t.Errorf("RemoveNonLetter failed test, exepted:%v got:%v", trueResult, result)
	}
}

func TestIsValidUUID(t *testing.T) {
	trueID := "5dcca8d2-a5d6-11eb-bcbc-0242ac130002"
	wrongID := "fs3a2f4-sfa-sa2wfsa-fan3-3234234fsa"

	if !helpers.IsValidUUID(trueID) {
		t.Errorf("IsValidUUID failed test, true uuid didnt passed test")
	}
	if helpers.IsValidUUID(wrongID) {
		t.Errorf("IsValidUUID failed test, wrong uuid passed test")
	}
}
