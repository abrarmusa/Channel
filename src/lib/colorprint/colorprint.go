package colorprint

import (
	"fmt"
	"github.com/fatih/color"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// COLORPRINT METHODS
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Prints message into console in magenta
// -------------------
// INSTRUCTIONS:
// -------------------
// call colorprint.Debug("{YOUR CONSOLE OUTPUT STRING}"")
func Debug(str string) {
	color.Set(color.FgMagenta)
	fmt.Println(str)
	color.Unset()
}

// Prints message into console in yellow
// -------------------
// INSTRUCTIONS:
// -------------------
// call colorprint.Warning("{YOUR CONSOLE OUTPUT STRING}"")
func Warning(str string) {
	color.Set(color.FgYellow)
	fmt.Println(str)
	color.Unset()
}

// Prints message into console in magenta
func Alert(str string) {
	color.Set(color.FgRed)
	fmt.Println(str)
	color.Unset()
}

// Prints message into console in blue
// -------------------
// INSTRUCTIONS:
// -------------------
// call colorprint.Blue("{YOUR CONSOLE OUTPUT STRING}"")
func Blue(str string) {
	color.Set(color.FgBlue)
	fmt.Println(str)
	color.Unset()
}

// Prints message into console in green
// -------------------
// INSTRUCTIONS:
// -------------------
// call colorprint.Info("{YOUR CONSOLE OUTPUT STRING}"")
func Info(str string) {
	color.Set(color.FgGreen)
	fmt.Println(str)
	color.Unset()
}
