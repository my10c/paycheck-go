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

	"vars"

	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"

	"github.com/akamensky/argparse"
	"github.com/BurntSushi/toml"
)

type (
   Config struct {
		ConfigFile			string
		Salary				int
		CostHouse			int
		CostCar				int
		StandardDeduction	int
		State				string
		Tax					map[string]float32
		Insurance			map[string]int
		Federal				map[string]float32
		FederalBracket		[][]float32
		StateBracket		[][]float32
	}

	State struct {
		StateOtherTax		float32		`toml:"StateOtherTax,omitempty"`
	}

	Insurance struct {
		Medical		int					`toml:"Medical,omitempty"`
		Dental		int					`toml:"Dental,omitempty"`
		Vision		int					`toml:"Vision,omitempty"`
		Pension		int					`toml:"Pension,omitempty"`
		LongTerm	int					`toml:"LongTerm,omitempty"`
		Life		int					`toml:"Life,omitempty"`
	}

	Federal struct {
		SocialSecurity		float32		`toml:"SocialSecurity,omitempty"`
		SocialSecurityMax	float32		`toml:"SocialSecurityMax,omitempty"`
		Medicare			float32		`toml:"Medicare,omitempty"`
	}

	Bracket struct {
		TaxBracket			[][]float32	`toml:"TaxBracket,omitempty"`	
	}

	tomlConfig struct {
		Base		map[string]int		`toml:"base,omitempty"`
		State		map[string]State	`toml:"state,omitempty"`
		Insurance	map[string]int		`toml:"insurance,omitempty"`
		Federal		map[string]float32	`toml:"federal,omitempty"`
		Bracket		map[string]Bracket	`toml:"bracket,omitempty"`
	}
)

var (
	configFileSet	bool
	salarySet		bool
	stateSet		bool
	houseSet		bool
	carSet			bool
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config{
		Tax: 		make(map[string]float32),
		Insurance:	make(map[string]int),
	}
}

func (c *Config) InitializeArgs(p *print.Print) {

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
			Help:		"The yearly salary before tax",
		})

	state := parser.String("s", "state",
		&argparse.Options{
			Required:	false,
			Help:		"The state where taxes is collected, required",
		})

	costHouse := parser.String("H", "house",
		&argparse.Options{
			Required:	false,
			Help:		"The monthly house cost/mortgage",
		})

	costCar := parser.String("C", "car",
		&argparse.Options{
			Required:	false,
			Help:		"The monthly cars payment",
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
	stateSet		= parser.GetArgs()[3].GetParsed()
	houseSet		= parser.GetArgs()[4].GetParsed()
	carSet			= parser.GetArgs()[5].GetParsed()

	if !stateSet  {
		p.PrintBlue(vars.MyInfo)
		p.PrintRed("\n\tThe flags [-s|--state] is required\n\n")
		p.PrintGreen(parser.Usage(err))
		os.Exit(1)
	}

	if _, ok, _ := i.IsExist(*configFile, "file"); !ok {
		p.PrintRed("Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.Salary, _		= strconv.Atoi(*salary)
	c.CostHouse, _	= strconv.Atoi(*costHouse)
	c.CostCar, _	= strconv.Atoi(*costCar)
	c.ConfigFile	= *configFile
	c.State			= *state
}

// function to add the values to the Config object from the configuration file
func (c *Config) SetCalculationSettings(p *print.Print) {

	var configValues tomlConfig

	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		p.PrintRed("\tError reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Standard Deduction
	c.StandardDeduction = configValues.Base["StandardDeduction"]

	// ignore value from config the flag was given : -s, H and C
	if !salarySet {
		c.Salary = configValues.Base["Salary"]
	}
	
	if !houseSet {
		c.CostHouse = configValues.Base["CostHouse"]
	}

	if !carSet {
		c.CostCar = configValues.Base["CostCar"]
	}

	// set the insurances cost
	for _, field := range vars.CalcInsurance {
		c.Insurance[field] = configValues.Insurance[field]
	}

	// use bracket for tax calculation
	c.Tax["StateOtherTax"] = configValues.State[c.State].StateOtherTax
	c.Federal = configValues.Federal
	c.FederalBracket = configValues.Bracket["federal"].TaxBracket
	c.StateBracket = configValues.Bracket[c.State].TaxBracket
}
