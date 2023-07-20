module main

go 1.20

require (
	calculate v0.0.0
	configurator v0.0.0
	vars v0.0.0
)

require (
	github.com/my10c/packages-go/is v0.0.0-20230717011209-51a83962742b
	github.com/my10c/packages-go/print v0.0.0-20230717011209-51a83962742b
	github.com/my10c/packages-go/spinner v0.0.0-20230717011209-51a83962742b
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/akamensky/argparse v1.4.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/my10c/packages-go/format v0.0.0-20230717011209-51a83962742b // indirect
)

replace calculate => ./mod/calculate

replace configurator => ./mod/configurator

replace vars => ./mod/vars
