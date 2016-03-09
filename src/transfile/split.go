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

func Test() {
	filename := "dogg"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fileinfo, _ := file.Stat()
	var filesize int64 = fileinfo.Size()
	const part = 1 * (1 << 20) // 1 MB
	totalpartnum := uint64(math.Ceil(float64(filesize) / float64(part)))
	fmt.Println(totalpartnum)
	fmt.Printf("splitting to %d pieces.\n", totalpartnum)

	for i := uint64(1); i < totalpartnum+1; i++ {
		partsize := int(math.Min(part, float64(filesize - int64((i-1)*part))))
	fmt.Println(filesize)
	fmt.Println(partsize)
		buf := make([]byte, partsize)
		file.Read(buf)
		partname := "doggy" + strconv.FormatUint(i, 10)
		_, err := os.Create(partname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ioutil.WriteFile(partname, buf, os.ModeAppend)
		fmt.Println("split to : ", partname)
	}
}
