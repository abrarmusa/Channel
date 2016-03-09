package main

import (
	"transfile"
	"fmt"
)

func main() {
	transfile.SplitFile("dummy.txt", 6) // saves 6 chunks in ./
	str, n := transfile.GetChunkNum("hey0re-rte.ertwer.jpg.34")
	fmt.Println()
	fmt.Println("Filename:", str)
	fmt.Println("Chunk num:", n)
}
