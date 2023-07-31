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

	var TakeHome float64 = 0
	var taxTotal float64 = 0

	var SocialSecurity float64 = 0
	var Medicare float64 = 0
	var insuranceTotal float64 = 0
	var percentOther float64

	// use in the loops
	var idx int
	var over float64
	var overTax float64
	var fedTotalTax float64
	var stateTotalTax float64
	var percentTax float64

	monthlyCost := (c.CostHouse + c.CostCar)
	federalBracket := c.FederalBracket
	stateBracket   := c.StateBracket
	fedTaxableSalary := int64(c.Salary) - int64(c.Federal["StandardDeduction"])
	stateTaxableSalary := int64(c.Salary) - int64(c.StatedDeduction)

	fmt.Printf("\n\t%s Before taxes: %s\n",
		p.PrintLine(print.Green, 25),
		p.PrintLine(print.Green, 24),
	)

	p.PrintGreen(fmt.Sprintf("\tYearly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary))),
	)
	p.PrintGreen(fmt.Sprintf("\tFederal Taxable salary\t  $%12s / -$%s\n",
		format.Format(fedTaxableSalary),
		format.Format(int64(c.Federal["StandardDeduction"]))),
	)
	p.PrintGreen(fmt.Sprintf("\tState Taxable salary\t  $%12s / -$%s\n",
		format.Format(stateTaxableSalary),
		format.Format(int64(c.StatedDeduction))),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly salary\t\t  $%12s\n",
		format.Format(int64(c.Salary/12))),
	)
	p.PrintGreen(fmt.Sprintf("\tBi-weekly salary\t  $%12s\n",
		format.Format(int64(c.Salary/24))),
	)
	p.PrintGreen(fmt.Sprintf("\tMonthly Cost house & car  $%12s (house: $%4s / car: $%4s) \n",
		format.Format(int64(monthlyCost)),
		format.Format(int64(c.CostHouse)),
		format.Format(int64(c.CostCar))),
	)

	p.PrintYellow(
		fmt.Sprintf("\n\t%s Federal, State (%s) and Medicare: %s\n\t\t\t\t  bi-weekly / monthly  / yearly\n",
		p.PrintLine(print.Yellow, 15),
		strings.ToUpper(c.State),
		p.PrintLine(print.Yellow, 15)),
	)

	for idx, _ = range federalBracket {
		for _ = range federalBracket[idx]{
			if	fedTaxableSalary > int64(federalBracket[idx][2]) &&
				fedTaxableSalary <  int64(federalBracket[idx][3]) {
				over = float64(c.Salary) - federalBracket[idx][2]
				overTax = (over * federalBracket[idx][1]) / 100
				//bi-weekly
				fedTotalTax = (overTax + federalBracket[idx][0]) / 24
				// percentTax = federalBracket[idx][1]
				break
			}
		}
	}
	// effective tax in present
	percentTax = float64((fedTotalTax * 24 ) / c.Salary) * 100

	p.PrintYellow(fmt.Sprintf("\t%14s (%.1f%%)\t  $ %6s  / $ %6s / $ %6s \n",
		"FederalTax", percentTax,
		format.Format(int64(fedTotalTax)),
		format.Format(int64(fedTotalTax * 2)),
		format.Format(int64(fedTotalTax * 24))),
	)

	for idx, _  = range stateBracket {
		for _ = range stateBracket[idx] {
			if	stateTaxableSalary > int64(stateBracket[idx][2]) &&
				stateTaxableSalary <  int64(stateBracket[idx][3]) {
				over = float64(c.Salary) - stateBracket[idx][2]
				overTax = (over * stateBracket[idx][1]) / 100
				//bi-weekly
				stateTotalTax = (overTax + stateBracket[idx][0]) / 24
				// percentTax = stateBracket[idx][1]
				break
			}
		}
	}
	// effective tax in present
	percentTax = float64((stateTotalTax * 24 ) / c.Salary) * 100

	p.PrintYellow(fmt.Sprintf("\t%14s (%.1f%%)\t  $ %6s  / $ %6s / $ %6s \n",
		"StateTax", percentTax,
		format.Format(int64(stateTotalTax)),
		format.Format(int64(stateTotalTax * 2)),
		format.Format(int64(stateTotalTax * 24))),
	)

	if fedTaxableSalary > int64(c.Federal["SocialSecurityMax"]) {
		SocialSecurity =
			((float64(c.Federal["SocialSecurityMax"]) * c.Federal["SocialSecurity"]) / 100) / 24
	} else {
		SocialSecurity =
			((float64(fedTaxableSalary) * c.Federal["SocialSecurity"]) / 100) / 24
	}
	percentOther = float64((SocialSecurity * 24) / c.Salary) * 100
	p.PrintYellow(fmt.Sprintf("\t%14s (%.1f%%)\t  $ %6s  / $ %6s / $ %6s \n",
		"SocialSecurity", percentOther,
		format.Format(int64(SocialSecurity)),
		format.Format(int64(SocialSecurity * 2)),
		format.Format(int64(SocialSecurity * 24))),
	)

	Medicare =
		((float64(fedTaxableSalary) * c.Federal["Medicare"])  / 100) / 24
	percentOther = float64((Medicare * 24) / c.Salary) * 100
	p.PrintYellow(fmt.Sprintf("\t%14s (%.1f%%)\t  $ %6s  / $ %6s / $ %6s \n",
		"Medicare", percentOther,
		format.Format(int64(Medicare)),
		format.Format(int64(Medicare * 2)),
		format.Format(int64(Medicare * 23))),
	)

	totalAll := fedTotalTax + stateTotalTax + SocialSecurity + Medicare
	p.PrintRed(fmt.Sprintf("\tTotal $%s / $%s / $%s\n",
		format.Format(int64(totalAll)),
		format.Format(int64(totalAll * 2)),
		format.Format(int64(totalAll * 24))),
	)

	// Insurance
	fmt.Printf("\n\t%s Insurance and 401K: %s\n\t\t\t\t  bi-weekly / monthly  / yearly\n",
		p.PrintLine(print.Blue, 22),
		p.PrintLine(print.Blue, 22),
	)
	insuranceTotal = 0
	var insuranceCost int64
	for field, _ := range c.Insurance {
		insuranceCost = int64(c.Insurance[field])
		p.PrintBlue(fmt.Sprintf("\t%18s\t  $  %6s / $ %6s / $ %6s \n",
			field,
			format.Format(insuranceCost),
			format.Format(insuranceCost * 2),
			format.Format(insuranceCost * 24)),
		)
		// format.Format(int64(c.Insurance[field] * 24))),
		insuranceTotal += c.Insurance[field]
	}
	p.PrintRed(fmt.Sprintf("\tTotal $%s / $%s / $%s\n",
		format.Format(int64(insuranceTotal)),
		format.Format(int64(insuranceTotal * 2 )),
		format.Format(int64(insuranceTotal * 24))),
	)
	taxTotal = float64(fedTotalTax + stateTotalTax + SocialSecurity + Medicare )
	TakeHome =	float64(c.Salary/24) - taxTotal - float64(insuranceTotal)
	// adjustment and possible extra income
	TakeHome = (TakeHome * (100 + c.Adjustment)) / 100
	// after adjustment
	afterCost := TakeHome - float64(monthlyCost/ 2)

	// to make adjustment easier
	fmt.Printf("\n\t%s\n", p.PrintLine(print.Purple, 64))	
	biWeekly := fmt.Sprintf("\t$ %7s  / $ %7s",
		format.Format(int64(TakeHome + (c.ExtraIncome/2))),
		format.Format(int64(afterCost +(c.ExtraIncome/2))),
	)
	monthly := fmt.Sprintf("\t$ %7s  / $ %7s",
		format.Format(int64((TakeHome * 2) + c.ExtraIncome)),
		format.Format(int64((afterCost * 2) + c.ExtraIncome)),
	)
	yearly := fmt.Sprintf("\t$ %7s  / $ %7s",
		format.Format(int64((TakeHome * 24) + (c.ExtraIncome * 12))),
		format.Format(int64((afterCost * 24) + (c.ExtraIncome * 12))),
	)
	p.PrintYellow(fmt.Sprintf("\tAdjust by +%.2f%% and extre income $%s monthly, $%s yearly\n",
		c.Adjustment,
		format.Format(int64(c.ExtraIncome)),
		format.Format(int64(c.ExtraIncome * 12))),
	)
	p.PrintYellow(fmt.Sprintf("\tinclude the extra income\n"))
	p.PrintGreen(fmt.Sprintf("\t(approx.) Bring home salary: \t   / After house and car payment:\n"))
	p.PrintGreen(fmt.Sprintf("\t\t bi-weekly : %s\n", biWeekly))
	p.PrintGreen(fmt.Sprintf("\t\t monthly   : %s\n", monthly))
	p.PrintGreen(fmt.Sprintf("\t\t yearly    : %s\n", yearly))
}
