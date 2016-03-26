package colorprint

// package clientstream

import (
	"fmt"
	"github.com/fatih/color"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// COLORPRINT METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// debug(str string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints message into console in magenta
func Debug(str string) {
	color.Set(color.FgMagenta)
	fmt.Println(str)
	color.Unset()
}

// warning(str string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints message into console in yellow
func Warning(str string) {
	color.Set(color.FgYellow)
	fmt.Println(str)
	color.Unset()
}

// debug(str string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints message into console in magenta
func Alert(str string) {
	color.Set(color.FgRed)
	fmt.Println(str)
	color.Unset()
}

// blue(str string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints message into console in blue
func Blue(str string) {
	color.Set(color.FgBlue)
	fmt.Println(str)
	color.Unset()
}

// info(str string)
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints message into console in green
func Info(str string) {
	color.Set(color.FgGreen)
	fmt.Println(str)
	color.Unset()
}
