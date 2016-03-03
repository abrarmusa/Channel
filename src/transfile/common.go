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
