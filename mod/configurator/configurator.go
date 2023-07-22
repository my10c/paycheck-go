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

	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"

	"github.com/akamensky/argparse"
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		ConfigFile			string
		State				string
		Salary				int
		CostHouse			int
		CostCar				int
		Insurance			map[string]int
		Federal				map[string]float32
		FederalBracket		[][]float32
		StateBracket		[][]float32
	}

	Bracket struct {
		TaxBracket			[][]float32	`toml:"TaxBracket"`
	}

	tomlConfig struct {
		Location	map[string]string	`toml:"location"`
		Base		map[string]int		`toml:"base"`
		Insurance	map[string]int		`toml:"insurance"`
		Federal		map[string]float32	`toml:"federal"`
		Bracket		map[string]Bracket	`toml:"bracket"`
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
			Help:		"The yearly salary before tax, required if not set in the configuration file",
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

	if _, ok, _ := i.IsExist(*configFile, "file"); !ok {
		p.PrintRed("Configuration file " + *configFile + " does not exist\n")
		os.Exit(1)
	}

	c.Salary, _		= strconv.Atoi(*salary)
	c.CostHouse, _	= strconv.Atoi(*costHouse)
	c.CostCar, _	= strconv.Atoi(*costCar)
	c.ConfigFile	= *configFile
	c.State			= strings.ToLower(*state)
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
	// need state so check after state was set
	if configValues.Base["Salary"] == 0  && !salarySet {
		p.PrintRed("\tConfiguration file error: Salary is required\n")
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
		c.Insurance[field] = configValues.Insurance[field]
	}

	// use bracket for tax calculation
	c.Federal = configValues.Federal
	c.FederalBracket = configValues.Bracket["federal"].TaxBracket
	c.StateBracket = configValues.Bracket[c.State].TaxBracket
}
