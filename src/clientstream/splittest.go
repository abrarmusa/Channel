package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type VidSegment struct {
	Id int
	Body []byte
}

func main() {
	filename := "dummy.txt"
	num := 4 // replication factor

	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()
	fileinfo, _ := file.Stat()
	var filesize uint64 = uint64(fileinfo.Size())
	// fmt.Println("Original file:", filename, ":", filesize)
	totalpartnum := uint64(num)
	part := uint64(math.Ceil(float64(filesize) / float64(totalpartnum)))
	for i := uint64(1); i < totalpartnum+1; i++ {

		// create a VidSegment instance
		vs := VidSegment {
			Id: int(i),
			Body: nil,
		}
		partsize := int(math.Min(float64(part), float64(filesize - (i-1)*part)))
		buf := make([]byte, partsize)
		file.Read(buf)
		fmt.Println(vs.Id)
		buf_j, err := json.Marshal(buf)
		checkError(err)
		vs.Body = buf_j

		// now the vs VidSegment is ready to use!

		// just to check
		var temp []byte
		json.Unmarshal(vs.Body, &temp)
		fmt.Println(string(temp))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error string: ", err)
		os.Exit(-1)
	}
}
