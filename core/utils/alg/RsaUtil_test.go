package alg

import (
	"fmt"
	"testing"
)

func TestRsaEncrypt(t *testing.T) {
	if data, e := RsaEncrypt("hello"); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(data)
	}
}

func TestRsaDecrypt(t *testing.T) {
	if data, e := RsaEncrypt("hello"); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(data)
		if data, e := RsaDecrypt(data); e != nil {
			fmt.Println(e)
		} else {
			fmt.Println(data)
		}
	}

}