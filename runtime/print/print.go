package print

import (
	"fmt"

	"github.com/fatih/color"
)

// H1 prints a heading with bold seperators and a color of your choice.
//
// NOTE: A closure is returned thus making for nicer syntax when used with
// run.Serial/Parallel functions. If using standalone you will need to invoke
// the returned closure.
func H1(text string, textColor color.Attribute) func() error {
	return func() error {
		defer color.Unset()
		color.Set(textColor)
		fmt.Println("")
		fmt.Println("================================================================================")
		fmt.Println(">>> " + text)
		fmt.Println("================================================================================")
		fmt.Println("")
		return nil
	}
}

// H2 prints a heading with slightly less emphasise than H1
//
// NOTE: A closure is returned thus making for nicer syntax when used with
// run.Serial/Parallel functions. If using standalone you will need to invoke
// the returned closure.
func H2(text string) func() error {
	return func() error {
		fmt.Println("")
		fmt.Println(text)
		fmt.Println("--------------------------------------------------------------------------------")
		return nil
	}
}
