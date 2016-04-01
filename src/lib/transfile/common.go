package transfile

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error string: ", err)
		os.Exit(-1)
	}
}

func StayLive() {
	var temp string
	fmt.Scanln(&temp)
}

func concat(addition string, to *[]byte) {
	temp := []byte(addition)
	for i := 0; i < len(temp); i++ {
		*to = append(*to, temp[i])
	}
}

// Separate chunk num from the filename

func GetChunkNum(s string) (string, int) {
	arr := strings.Split(s, ".")
	var temp []byte
	repeat := len(arr) - 1
	for i := 0; i < repeat; i++ {
		concat(arr[i], &temp)
		if i < repeat-1 {
			concat(".", &temp)
		}
	}
	var ver int
	ver, _ = strconv.Atoi(arr[len(arr)-1])
	return string(temp), ver
}
