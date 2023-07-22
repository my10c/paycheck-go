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

	var takeHone float64 = 0
	var taxTotal float64 = 0

	var SocialSecurity float64 = 0
	var Medicare float64 = 0
	var insuranceTotal float64 = 0

	// use in the loops
	var idx int
	var over float64
	var overTax float64
	var totalTax float64
	var percentTax float64

	var overState float64
	var overTaxState float64
	var totalTaxState float64
	var percentTaxState float64

	monthlyCost := (c.CostHouse + c.CostCar)
	federalBracket := c.FederalBracket
	stateBracket   := c.StateBracket
	fedTaxableSalary := int64(c.Salary) - int64(c.Federal["StandardDeduction"])
	stateTaxableSalary := int64(c.Salary) - int64(c.StatedDeduction)

	fmt.Printf("\t%s Before taxes: %s\n",
		p.PrintLine(print.Purple, 23),
		p.PrintLine(print.Purple, 22),
	)

	p.PrintGreen(fmt.Sprintf("\tYearly salary\t\t\t  $%12s\n",
		format.Format(int64(c.Salary))),
	)
	p.PrintGreen(fmt.Sprintf("\tFederal Taxable salary\t\t  $%12s / -$%s\n",
		format.Format(fedTaxableSalary),
		format.Format(int64(c.Federal["StandardDeduction"]))),
	)
	p.PrintGreen(fmt.Sprintf("\tState Taxable salary\t\t  $%12s / -$%s\n",
		format.Format(stateTaxableSalary),
		format.Format(int64(c.StatedDeduction))),
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
			if	fedTaxableSalary > int64(federalBracket[idx][2]) &&
				fedTaxableSalary <  int64(federalBracket[idx][3]) {
				over = float64(c.Salary) - federalBracket[idx][2]
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
			if	stateTaxableSalary > int64(stateBracket[idx][2]) &&
				stateTaxableSalary <  int64(stateBracket[idx][3]) {
				overState = float64(c.Salary) - stateBracket[idx][2]
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

	if fedTaxableSalary > int64(c.Federal["SocialSecurityMax"]) {
		SocialSecurity =
			((float64(c.Federal["SocialSecurityMax"]) * c.Federal["SocialSecurity"]) / 100) / 24
	} else {
		SocialSecurity =
			((float64(fedTaxableSalary) * c.Federal["SocialSecurity"]) / 100) / 24
	}
	p.PrintYellow(fmt.Sprintf("\t%14s (%v%%)\t  $ %11s\n",
		"SocialSecurity", c.Federal["SocialSecurity"], format.Format(int64(SocialSecurity))),
	)

	Medicare =
		((float64(fedTaxableSalary) * c.Federal["Medicare"])  / 100) / 24
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
		p.PrintBlue(fmt.Sprintf("\t%18s\t  $ %11s\n",
			field, format.Format(int64(c.Insurance[field]))),
		)
		insuranceTotal += c.Insurance[field]
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s\n", format.Format(int64(insuranceTotal))))
	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))	

	taxTotal = float64(totalTax + totalTaxState + SocialSecurity + Medicare )
	takeHone =	float64(c.Salary/24) - taxTotal - float64(insuranceTotal)
	afterCost := takeHone - float64(monthlyCost/ 2)

	p.PrintGreen(fmt.Sprintf("\t(approx.) Bring home salary:\n\t\t\tbi-weekly : $%s\n\t\t\tmonthly   : $%s\n",
		format.Format(int64(takeHone)),
		format.Format(int64(takeHone * 2))),
	)

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))
	p.PrintGreen(fmt.Sprintf("\t(approx.) After house and car payment:\n\t\t\tbi-weekly : $%s\n\t\t\tmonthly   : $%s\n",
			format.Format(int64(afterCost)),
			format.Format(int64(afterCost * 2))),
	)
}
