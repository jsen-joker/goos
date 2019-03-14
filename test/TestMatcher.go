package main

import (
	"fmt"
	"github.com/jsen-joker/goos/core/utils"
)

func main() {
	s := []string{"tes.*", "dfds"}
	fmt.Println(utils.Matcher("test.json", s))
}
