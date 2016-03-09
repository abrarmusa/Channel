package transfile

import (
	"fmt"
	"os"
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
