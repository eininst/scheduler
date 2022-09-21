package main

import "fmt"

func main() {
	s := "weweewewwex"

	if len(s) > 5 {
		fmt.Println(s[0:5])
	}

}
