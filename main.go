package main

import "fmt"

func main() {
	ScrapeJSON()
	fmt.Printf("Choose a vendor:\n\t(1) NeoMacro\n\t(2) AceKBD\n\t(3) GenesisPC\n\t(4) Keebsmod\n")

	ReadJSON()
}
