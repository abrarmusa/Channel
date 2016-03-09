package transfile

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func SplitFile(filename string, num int) {
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()
	fileinfo, _ := file.Stat()
	var filesize uint64 = uint64(fileinfo.Size())
	fmt.Println("Original file:", filename, ":", filesize)
	totalpartnum := uint64(num)
	part := uint64(math.Ceil(float64(filesize) / float64(totalpartnum)))
	for i := uint64(1); i < totalpartnum+1; i++ {
		partsize := int(math.Min(float64(part), float64(filesize - (i-1)*part)))
		buf := make([]byte, partsize)
		file.Read(buf)
		partname := filename + "." + strconv.FormatUint(i, 10)
		_, err := os.Create(partname)
		checkError(err)
		ioutil.WriteFile(partname, buf, os.ModeAppend)
		fmt.Println("split to :", partname, ":", partsize)
	}
}
