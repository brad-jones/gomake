package print

import (
	"fmt"

	"github.com/fatih/color"
)

// H1 prints a heading with bold seperators and a color of your choice.
func H1(text string, textColor color.Attribute) {
	defer color.Unset()
	color.Set(textColor)
	fmt.Println("")
	fmt.Println("================================================================================")
	fmt.Println(">>> " + text)
	fmt.Println("================================================================================")
	fmt.Println("")
}

// H2 prints a heading with slightly less emphasise than H1
func H2(text string) {
	fmt.Println("")
	fmt.Println(text)
	fmt.Println("--------------------------------------------------------------------------------")
}
