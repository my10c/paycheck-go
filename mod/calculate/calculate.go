//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package calculate

import (
	"fmt"

	"configurator"
	"vars"

	"github.com/my10c/packages-go/format"
	"github.com/my10c/packages-go/print"
)

func ShowCalc(c *configurator.Config, p *print.Print) {
	var takeHone float32
	var taxTotal float32
	var insuranceTotal int
	var taxItem string

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))	
	p.PrintGreen(fmt.Sprintf("\tYearly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary))),
	)
	p.PrintGreen(fmt.Sprintf("\tBi-Weekly salary\t  $%12s\n",
		format.Format(int64(c.Salary/24))),
	)

	fmt.Printf("\n\t%s Federal, State and Medicare: %s\n",
		p.PrintLine(print.Purple, 15),
		p.PrintLine(print.Purple, 15),
	)
	// Taxes
	for _, item := range vars.CalcTax {
		taxItem = fmt.Sprintf("%.2f", (c.Tax[item] / 100) * (float32(c.Salary/24)))
		p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %12s\n",
			item, c.Tax[item], taxItem),
		)
		taxTotal += (c.Tax[item] / 100) * (float32(c.Salary/24))
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s\n", format.Format(int64(taxTotal))))
	
	// CalcInsurance
	fmt.Printf("\n\t%s Insurance and 401K: %s\n",
		p.PrintLine(print.Purple, 20),
		p.PrintLine(print.Purple, 20),
	)
	for _, item := range vars.CalcInsurance {
		p.PrintBlue(fmt.Sprintf("\t%18s\t  $ %8d\n",
			item, c.Insurance[item]),
		)
		insuranceTotal += c.Insurance[item]
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s\n", format.Format(int64(insuranceTotal))))
	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))	

	monthlyCost := (c.CostHouse + c.CostCar)/ 2 
	takeHone =	float32(c.Salary/24) - taxTotal - float32(insuranceTotal)
	afterCost := (takeHone *2) - float32(monthlyCost)

	p.PrintGreen(fmt.Sprintf("\tBring home salary $%s bi-weekly\n", format.Format(int64(takeHone))))
	p.PrintGreen(fmt.Sprintf("\tAfter monthly costs $%s ($%s/monthly before cost)\n",
			format.Format(int64(afterCost)), format.Format(int64(takeHone * 2))))
}
