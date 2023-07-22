//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package calculate

import (
	"fmt"
	"strings"

	"configurator"

	"github.com/my10c/packages-go/format"
	"github.com/my10c/packages-go/print"
)

// calculate based on brackets
func BracketCalc(c *configurator.Config, p *print.Print) {

	var takeHone float32 = 0
	var taxTotal float32 = 0

	var SocialSecurity float32 = 0
	var Medicare float32 = 0
	var insuranceTotal int = 0

	// use in the loops
	var idx int
	var over float32
	var overTax float32
	var totalTax float32
	var percentTax float32

	var overState float32
	var overTaxState float32
	var totalTaxState float32
	var percentTaxState float32

	monthlyCost := (c.CostHouse + c.CostCar)
	federalBracket := c.FederalBracket
	stateBracket   := c.StateBracket
	taxableSalary  := int64(c.Salary) - int64(c.Federal["StandardDeduction"])

	fmt.Printf("\t%s Before taxes: %s\n",
		p.PrintLine(print.Purple, 23),
		p.PrintLine(print.Purple, 22),
	)

	p.PrintGreen(fmt.Sprintf("\tYearly salary\t\t\t  $%12s\n",
		format.Format(int64(c.Salary))),
	)
	p.PrintGreen(fmt.Sprintf("\tTaxable salary\t\t\t  $%12s\n",
		format.Format(taxableSalary)),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly salary\t\t\t  $%12s\n",
		format.Format(int64(c.Salary/12))),
	)
	p.PrintGreen(fmt.Sprintf("\tBi-weekly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary/24))),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly Cost (house & car)\t  $%12s\n",
		format.Format(int64(monthlyCost))),
	)

	fmt.Printf("\t%s Federal, State (%s) and Medicare: %s\n",
		p.PrintLine(print.Purple, 13),
		strings.ToUpper(c.State),
		p.PrintLine(print.Purple, 13),
	)

	for idx, _ = range federalBracket {
		for _ = range federalBracket[idx]{
			if	taxableSalary > int64(federalBracket[idx][2]) &&
				taxableSalary <  int64(federalBracket[idx][3]) {
				over = float32(c.Salary) - federalBracket[idx][2]
				overTax = (over * federalBracket[idx][1]) / 100
				//bi-weekly
				totalTax = (overTax + federalBracket[idx][0]) / 24
				percentTax = federalBracket[idx][1]
				break
			}
		}
	}
	p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
		"FederalTax", percentTax, format.Format(int64(totalTax))),
	)

	for idx, _  = range stateBracket {
		for _ = range stateBracket[idx] {
			if	int64(c.Salary) > int64(stateBracket[idx][2]) &&
				int64(c.Salary) <  int64(stateBracket[idx][3]) {
				overState = float32(c.Salary) - stateBracket[idx][2]
				overTaxState = (overState * stateBracket[idx][1]) / 100
				//bi-weekly
				totalTaxState = (overTaxState + stateBracket[idx][0]) / 24
				percentTaxState = stateBracket[idx][1]
				break
			}
		}
	}
	p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
		"StateTax", percentTaxState, format.Format(int64(totalTaxState))),
	)

	if taxableSalary > int64(c.Federal["SocialSecurityMax"]) {
		SocialSecurity =
			((float32(c.Federal["SocialSecurityMax"]) * c.Federal["SocialSecurity"]) / 100) / 24
	} else {
		SocialSecurity =
			((float32(taxableSalary) * c.Federal["SocialSecurity"]) / 100) / 24
	}
	p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
		"SocialSecurity", c.Federal["SocialSecurity"], format.Format(int64(SocialSecurity))),
	)

	Medicare =
		((float32(taxableSalary) * c.Federal["Medicare"])  / 100) / 24
	p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
		"Medicare", c.Federal["Medicare"], format.Format(int64(Medicare))),
	)

	p.PrintRed(fmt.Sprintf("\tTotal $%s\n",
		format.Format(int64(totalTax + totalTaxState + SocialSecurity + Medicare ))))

	// Insurance
	fmt.Printf("\t%s Insurance and 401K: %s\n",
		p.PrintLine(print.Purple, 20),
		p.PrintLine(print.Purple, 20),
	)
	insuranceTotal = 0
	for field, _ := range c.Insurance {
		p.PrintBlue(fmt.Sprintf("\t%18s\t  $ %8d\n",
			field, c.Insurance[field]),
		)
		insuranceTotal += c.Insurance[field]
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s\n", format.Format(int64(insuranceTotal))))
	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))	

	taxTotal = float32(totalTax + totalTaxState + SocialSecurity + Medicare )
	takeHone =	float32(c.Salary/24) - taxTotal - float32(insuranceTotal)
	afterCost := takeHone - float32(monthlyCost/ 2)

	p.PrintGreen(fmt.Sprintf("\t(approx.) Bring home salary : $%s bi-weekly, $%s monthly\n",
		format.Format(int64(takeHone)),
		format.Format(int64(takeHone * 2))),
	)

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))
	p.PrintGreen(fmt.Sprintf("\t(approx.) After costs       : $%s bi-weekly, $%s monthly\n",
			format.Format(int64(afterCost)),
			format.Format(int64(afterCost * 2))),
	)
}
