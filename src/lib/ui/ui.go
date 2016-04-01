package ui

import (
	"../../consts"
	"../colorprint"
	"fmt"
)

// Intro()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Prints application information
func Intro() {
	fmt.Println("==========================================================================================")
	colorprint.Warning("      ___          ___          ___          ___          ___          ___          ___ ")
	colorprint.Warning("     /\\  \\        /\\__\\        /\\  \\        /\\__\\        /\\__\\        /\\  \\        /\\__\\")
	colorprint.Warning("    /::\\  \\      /:/  /       /::\\  \\      /::|  |      /::|  |      /::\\  \\      /:/  /")
	colorprint.Warning("   /:/\\:\\  \\    /:/__/       /:/\\:\\  \\    /:|:|  |     /:|:|  |     /:/\\:\\  \\    /:/  / ")
	colorprint.Warning("  /:/  \\:\\  \\  /::\\  \\ ___  /::\\~\\:\\  \\  /:/|:|  |__  /:/|:|  |__  /::\\~\\:\\  \\  /:/  /  ")
	colorprint.Warning(" /:/__/ \\:\\__\\/:/\\:\\  /\\__\\/:/\\:\\ \\:\\__\\/:/ |:| /\\__\\/:/ |:| /\\__\\/:/\\:\\ \\:\\__\\/:/__/   ")
	colorprint.Warning(" \\:\\  \\  \\/__/\\/__\\:\\/:/  /\\/__\\:\\/:/  /\\/__|:|/:/  /\\/__|:|/:/  /\\:\\~\\:\\ \\/__/\\:\\  \\   ")
	colorprint.Warning("  \\:\\  \\           \\::/  /      \\::/  /     |:/:/  /     |:/:/  /  \\:\\ \\:\\__\\   \\:\\  \\  ")
	colorprint.Warning("   \\:\\  \\          /:/  /       /:/  /      |::/  /      |::/  /    \\:\\ \\/__/    \\:\\  \\ ")
	colorprint.Warning("    \\:\\__\\        /:/  /       /:/  /       /:/  /       /:/  /      \\:\\__\\       \\:\\__\\")
	colorprint.Warning("     \\/__/        \\/__/        \\/__/        \\/__/        \\/__/        \\/__/        \\/__/")
	fmt.Println()
	fmt.Println("==========================================================================================")
	colorprint.Blue(">>>>                                     CHANNEL                                      <<<<")
	colorprint.Blue(">>>>----------------------------------------------------------------------------------<<<<")
	fmt.Printf(">>>>  %-78s  <<<<\n", "BITTORRENT BASED P2P VIDEO STREAMING APPLICATION")
	fmt.Printf(">>>>  VERSION %-70s  <<<<\n", consts.VersionNum)
	fmt.Printf(">>>>  %-78s  <<<<\n", consts.Builders)
	colorprint.Blue(">>>>----------------------------------------------------------------------------------<<<<")
}

// Help()
// --------------------------------------------------------------------------------------------
// DESCRIPTION:
// -------------------
// Lists available commands
func Help() {
	Intro()
	colorprint.Warning(">>>>----------------------------------------------------------------------------------<<<<")
	colorprint.Blue("    AVAILABLE LIST OF COMMANDS")
	colorprint.Blue("    1. 'help' -  displays available list of commands")
	colorprint.Blue("    2. 'get {filename} {node-address}' - asks node at {node-adress} for a file. If the file is present, user is asked y/n for transfer of file")
	colorprint.Blue("    3. 'list' - lists all the available commands")
	colorprint.Warning(">>>>----------------------------------------------------------------------------------<<<<")
}
