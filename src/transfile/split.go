package transfile

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func getPrefix(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	myarray := strings.Split(s, ".")
	repeat := len(myarray) - 1
	for i := 0; i < repeat; i++ {
	}
	return "whatever"
}

func SplitFile(filename string, num int) {
	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()
	fileinfo, _ := file.Stat()
	var filesize uint64 = uint64(fileinfo.Size())
	fmt.Println("Filesize:", filesize)
	totalpartnum := uint64(num)
	part := uint64(math.Ceil(float64(filesize) / float64(totalpartnum)))
	for i := uint64(1); i < totalpartnum+1; i++ {
		partsize := int(math.Min(float64(part), float64(filesize - (i-1)*part)))
		buf := make([]byte, partsize)
		file.Read(buf)
		partname := "doggy" + strconv.FormatUint(i, 10)
		_, err := os.Create(partname)
		checkError(err)
		ioutil.WriteFile(partname, buf, os.ModeAppend)
		fmt.Println("split to :", partname, ":", partsize)
	}
}
