//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package calculate

import (
	"fmt"
	// "strconv"

	"configurator"
	"vars"

	"github.com/my10c/packages-go/format"
	"github.com/my10c/packages-go/print"
)

func ShowCalc(c *configurator.Config, p *print.Print) {
	var takeHone float32
	var taxTotal float32
	var insuranceTotal int
	//var taxItem string
	var taxValue float32

	monthlyCost := (c.CostHouse + c.CostCar)

	fmt.Printf("\t%s Before taxes: %s\n",
		p.PrintLine(print.Purple, 23),
		p.PrintLine(print.Purple, 22),
	)

	p.PrintGreen(fmt.Sprintf("\tYearly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary))),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary/12))),
	)
	p.PrintGreen(fmt.Sprintf("\tBi-weekly salary\t  $%12s\n",
		format.Format(int64(c.Salary/24))),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly Cost\t\t  $%12s\n",
		format.Format(int64(monthlyCost))),
	)

	fmt.Printf("\t%s Federal, State and Medicare: %s\n",
		p.PrintLine(print.Purple, 15),
		p.PrintLine(print.Purple, 15),
	)
	// Taxes
	for _, item := range vars.CalcTax {
		// taxItem = fmt.Sprintf("%.2f", (c.Tax[item] / 100) * (float32(c.Salary/24)))
		taxValue = (c.Tax[item] / 100) * (float32(c.Salary/24))
		p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
			item, c.Tax[item],	// taxItem),
			format.Format(int64(taxValue))),
		)
		taxTotal += (c.Tax[item] / 100) * (float32(c.Salary/24))
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s\n", format.Format(int64(taxTotal))))
	
	// CalcInsurance
	fmt.Printf("\t%s Insurance and 401K: %s\n",
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

	takeHone =	float32(c.Salary/24) - taxTotal - float32(insuranceTotal)
	afterCost := takeHone - float32(monthlyCost/ 2)

	p.PrintGreen(fmt.Sprintf("\tBring home salary ~$%s bi-weekly, ~$%s monthly\n",
		format.Format(int64(takeHone)),
		format.Format(int64(takeHone * 2))),
	)

	p.PrintGreen(fmt.Sprintf("\tAfter costs ~$%s bi-weekly, ~$%s monthly\n",
			format.Format(int64(afterCost)),
			format.Format(int64(afterCost * 2))),
	)
}
