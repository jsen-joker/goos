package utils

import (
	"fmt"
	"testing"
)

func TestMatcher(t *testing.T) {
	//s := []string{"tes.*", "dfds"}
	//if Matcher("test.json", s) {
	//	t.Error("test failed")
	//} else {
	//	t.Log("succeed")
	//}
	s := []string{"/**"}
	//fmt.Println(Matcher("test.json", s))
	fmt.Println(Matcher("/api/config/configList", s))
}