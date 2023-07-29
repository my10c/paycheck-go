//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package configurator

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"vars"

	"github.com/my10c/packages-go/format"
	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"

	"github.com/akamensky/argparse"
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		ConfigFile			string
		State				string
		Salary				float64
		MaxSalary			float64
		CostHouse			float64
		CostCar				float64
		Insurance			map[string]float64
		Federal				map[string]float64
		FederalBracket		[][]float64
		StateBracket		[][]float64
		StatedDeduction		float64
		Adjustment			float64
		ExtraIncome			float64
	}

	State struct {
		StatedDeduction		float64	`toml:"StatedDeduction,omitempty"`
		PersonalExemption	float64	`toml:"PersonalExemption,omitempty"`
	}

	Bracket struct {
		TaxBracket			[][]float64	`toml:"TaxBracket"`
	}

	tomlConfig struct {
		Location	map[string]string	`toml:"location"`
		Base		map[string]float64	`toml:"base"`
		Insurance	map[string]float64	`toml:"insurance"`
		State		map[string]State	`toml:"state"`
		Federal		map[string]float64	`toml:"federal"`
		Bracket		map[string]Bracket	`toml:"bracket"`
		Adjustment	map[string]float64	`toml:"adjustment"`
	}
)

var (
	configFileSet	bool
	salarySet		bool
	maxSalarySet	bool
	stateSet		bool
	houseSet		bool
	carSet			bool
	noInsuranceSet	bool
	adjustmentSet	bool
	extraIncomeSet	bool
	insuranceSet	map[string]bool
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config{
		Insurance:	make(map[string]float64),
	}
}

func (c *Config) InitializeArgs(p *print.Print) {
	insuranceSet = make(map[string]bool)

	i := is.New()

	parser := argparse.NewParser(vars.MyProgname, vars.MyDescription)
	configFile := parser.String("c", "configFile",
		&argparse.Options{
			Required: false,
			Help:		"Configuration file to be use",
			Default:	vars.ConfigFile,
		})

	salary := parser.String("S", "salary",
		&argparse.Options{
			Required: false,
			Help:		"The yearly salary before tax, required if not set in the configuration file",
		})

	maxSalary := parser.String("m", "maxsalary",
		&argparse.Options{
			Required:	false,
			Help:		"The maximum allowed salary value",
		})

	state := parser.String("s", "state",
		&argparse.Options{
			Required:	false,
			Help:		"The state where taxes is collected, required if not set in the configuration file",
		})

	costHouse := parser.String("H", "house",
		&argparse.Options{
			Required:	false,
			Help:		"The monthly house rent/mortgage",
		})

	costCar := parser.String("C", "car",
		&argparse.Options{
			Required:	false,
			Help:		"The monthly cars payment",
		})

	costMedical := parser.String("M", "medical",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly medical insurance cost",
		})

	costPension := parser.String("P", "pension",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly 401k contribution",
		})

	costVision := parser.String("V", "vision",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly vision insurance cost",
		})

	costDental := parser.String("D", "dental",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly dental insurance cost",
		})

	costLife := parser.String("L", "life",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly life insurance cost",
		})

	costLongTerm := parser.String("T", "longterm",
		&argparse.Options{
			Required:	false,
			Help:		"Bi-weekly long term disability insurance cost",
		})

	noInsurance := parser.Flag("N", "noinsurance",
		&argparse.Options{
			Required:	false,
			Help:		"No insurance cost nor contibution to 401k",
		})

	aproxAdjustment := parser.Float("A", "adjustment",
		&argparse.Options{
			Required:	false,
			Help:		"Adjustment to the calculation in %, (suggestion 2.0 - 3.0)",
		})

	extraIncome := parser.Float("E", "extraincome",
		&argparse.Options{
			Required:	false,
			Help:		"extra imcome per month, should be ater tax!" ,
		})

	showVersion := parser.Flag("v", "version",
		&argparse.Options{
		Required:	false,
		Help:		"Show version",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *showVersion {
		p.ClearScreen()
		p.PrintYellow(vars.MyProgname + " version: " + vars.MyVersion + "\n")
		os.Exit(0)
	}

	configFileSet	= parser.GetArgs()[1].GetParsed()
	salarySet		= parser.GetArgs()[2].GetParsed()
	maxSalarySet	= parser.GetArgs()[3].GetParsed()
	stateSet		= parser.GetArgs()[4].GetParsed()
	houseSet		= parser.GetArgs()[5].GetParsed()
	carSet			= parser.GetArgs()[6].GetParsed()

	insuranceSet["Medical"]		= parser.GetArgs()[7].GetParsed()
	insuranceSet["401k"]		= parser.GetArgs()[8].GetParsed()
	insuranceSet["Vision"]		= parser.GetArgs()[9].GetParsed()
	insuranceSet["Dental"]		= parser.GetArgs()[10].GetParsed()
	insuranceSet["Life"]		= parser.GetArgs()[11].GetParsed()
	insuranceSet["LongTerm"]	= parser.GetArgs()[12].GetParsed()

	noInsuranceSet	= *noInsurance // 13th position
	adjustmentSet	= parser.GetArgs()[14].GetParsed()
	extraIncomeSet	= parser.GetArgs()[15].GetParsed()

	if _, ok, _ := i.IsExist(*configFile, "file"); !ok {
		p.PrintRed("Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.Salary, _		= strconv.ParseFloat(*salary, 64)
	c.MaxSalary, _	= strconv.ParseFloat(*maxSalary, 64)
	c.CostHouse, _	= strconv.ParseFloat(*costHouse, 64)
	c.CostCar, _	= strconv.ParseFloat(*costCar, 64)
	c.ConfigFile	= *configFile
	c.State			= strings.ToLower(*state)
	// insurances
	c.Insurance["Medical"], _	= strconv.ParseFloat(*costMedical, 64)
	c.Insurance["Dental"], _	= strconv.ParseFloat(*costDental, 64)
	c.Insurance["Vision"], _	= strconv.ParseFloat(*costVision, 64)
	c.Insurance["401k"], _		= strconv.ParseFloat(*costPension, 64)
	c.Insurance["LongTerm"], _	= strconv.ParseFloat(*costLongTerm, 64)
	c.Insurance["Life"], _		= strconv.ParseFloat(*costLife, 64)
	// adjustment
	c.Adjustment = *aproxAdjustment
	c.ExtraIncome = *extraIncome
}

// function to add the values to the Config object from the configuration file
func (c *Config) SetCalculationSettings(p *print.Print) {
	var errMsg string

	var configValues tomlConfig

	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		p.PrintRed("\tError reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// set state first we will need for next checks
	if !stateSet {
		c.State = strings.ToLower(configValues.Location["State"])
	}

	// need state so check after state was set
	if configValues.Location["State"] == "" && !stateSet {
		p.PrintRed("\tConfiguration file error: State is required\n")
		os.Exit(1)
	}

	// set salary first we will need for next checks
	if !salarySet {
		c.Salary = configValues.Base["Salary"]
	}
	// need salary in configuration file if not given on the command line
	if configValues.Base["Salary"] == 0  && !salarySet {
		p.PrintRed("\tConfiguration file error: Salary is required\n")
		os.Exit(1)
	}

	// set max salary first we will need for next checks
	if !maxSalarySet {
		c.MaxSalary = configValues.Base["MaxSalary"]
	}
	// need max salary in configuration file
	if configValues.Base["MaxSalary"] == 0  && !maxSalarySet {
		p.PrintRed("\tConfiguration file error: MaxSalary is required\n")
		os.Exit(1)
	}

	if c.Salary > c.MaxSalary {
		errMsg = fmt.Sprintf("\tMax Salary is configured to max $%v\n",
			format.Format(int64(c.MaxSalary)))
		p.PrintRed(errMsg)
		os.Exit(1)
	}

	// check we have the required values
	if len(configValues.Bracket["federal"].TaxBracket) == 0 {
		p.PrintRed("\tConfiguration file error: Federal tax brackets is required\n")
		os.Exit(1)
	}

	// need state so check after state was set
	if len(configValues.Bracket[c.State].TaxBracket) == 0 {
		errMsg = fmt.Sprintf("\tConfiguration file error: state %v tax brackets is required\n",
			strings.ToUpper(c.State))
		p.PrintRed(errMsg)
		os.Exit(1)
	}

	// set values from command line
	if !houseSet {
		c.CostHouse = configValues.Base["CostHouse"]
	}

	if !carSet {
		c.CostCar = configValues.Base["CostCar"]
	}

	// make sure federal is set properly
	var erroed bool = false
	for _, field := range vars.CalcTax {
		if configValues.Federal[field] == 0 {
			errMsg = fmt.Sprintf("\tConfiguration file error: federal %v is required\n", field)
			p.PrintRed(errMsg)
			erroed = true
		}
	}
	if erroed {
		os.Exit(1)
	}

	// set the insurances cost
	for field, _ := range configValues.Insurance {
		if !noInsuranceSet {
			if !insuranceSet[field] {
				c.Insurance[field] = configValues.Insurance[field]
			}
		} else {
			c.Insurance[field] = 0
		}
	}

	// set adjustmentm if given in the CLI we ignore the one from the configuration file
	if !adjustmentSet {
		c.Adjustment = configValues.Adjustment["adjustment"]
	}
	// overwrite extraIncome if it was given by cli
	// if not set in both configuration and cli it will default to 0
	if !extraIncomeSet {
		c.ExtraIncome = configValues.Adjustment["extraIncome"]
	}

	// get the state Standard Deduction and PersonalExemption
	c.StatedDeduction = configValues.State[c.State].StatedDeduction
	c.StatedDeduction += configValues.State[c.State].PersonalExemption

	// use bracket for tax calculation
	c.Federal = configValues.Federal
	c.FederalBracket = configValues.Bracket["federal"].TaxBracket
	c.StateBracket = configValues.Bracket[c.State].TaxBracket
}
