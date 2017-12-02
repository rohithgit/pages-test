// Package security contains code to handle sceurity, e.g. authentication, authorization, encryption, etc.
package security

import (
	"fmt"
)

// init is an optional function that can be used by each source file to set up whatever state is required.
// init is called automatically by GO after all of the variable declarations in the package have evaluated their initializers,
// and those are evaluated only after all the imported packages have been initialized.
func init() {
	fmt.Println("security package init ..")
}
