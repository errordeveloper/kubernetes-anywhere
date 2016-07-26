package util

import (
	"fmt"

	"github.com/renstrom/dedent"
)

// TODO no-op for non-interactive shells

func PrintMessage(m string) {
	fmt.Printf(dedent.Dedent(m))
}

func PrintOkay() {
	fmt.Println("   ✓")
}
