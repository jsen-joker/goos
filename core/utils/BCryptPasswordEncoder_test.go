package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	fmt.Println(os.Getenv("HOME"))
	password := "nacos"
	if encode, err := Encode(password); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(encode)

		if e := Matches(encode, password + "s"); e != nil {
			fmt.Println(e)
		} else {
			fmt.Println("succeed")
		}
	}
}