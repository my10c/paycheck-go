// BSD 3-Clause License
//
// Copyright (c) 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:	0.1
//

package main

import (
	"fmt"
	"os"
	"time"

	// local
	"calculate"
	"configurator"

	// on github
	"github.com/my10c/packages-go/print"
	"github.com/my10c/packages-go/spinner"
)

func main() {
	s := spinner.New(10)
	p := print.New()
	c := configurator.Configurator()

	// get given parameters
	c.InitializeArgs(p)

	// set Calculation settings
	c.SetCalculationSettings(p)

	// tinny sleeper :-)
	go s.Run()
	time.Sleep(1 * time.Second)
	s.Stop()

	fmt.Printf(print.ClearScreen)
	// get the configurations values
	// calculate.ShowCalc(c, p)
	calculate.BracketCalc(c, p)

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 64))
	// p.TheEnd()
	// fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 6j0))

	os.Exit(0)
}
