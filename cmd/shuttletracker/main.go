// Package main bundles together all of shuttletracker's subpackages
// to create, configure, and run the shuttle tracker.
package main

import (
	"github.com/thepoly/shuttletracker/cmd"
)

func main() {
	cmd.Execute()
}
