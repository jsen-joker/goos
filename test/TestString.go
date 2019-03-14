package main

import (
	"flag"
	"fmt"
)

func main() {
	realPath := flag.String("path", "/home/default_dir", "static resource path")
	flag.Parse()
	fmt.Println(*realPath)
}


