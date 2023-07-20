# callit
My Salary Calculator


```
./paycheck -h
usage: paycheck [-h|--help] [-c|--configFile "<value>"] [-S|--salary "<value>"]
                -s|--state "<value>" [-H|--house "<value>"] [-C|--car
                "<value>"] [-v|--version]

                Simple script to calculate bi-weekly salary before and after
                tax 😁

Arguments:

  -h  --help        Print help information
  -c  --configFile  Configuration file to be use. Default:
                    /usr/local/etc/paycheck/paycheck.conf
  -S  --salary      The yearly salary before tax
  -s  --state       The state where taxes is collected, required
  -H  --house       The monthly house cost/mortgage
  -C  --car         The monthly cars payment
  -v  --version     Show version
```

###
Since taxes is adjusted based on the salary, the configs should be in this format: **_(state)(salary)_** 

Example:
```
ca150
co150
```

###
Taxes numbers take from :

[Colorado](https://smartasset.com/taxes/colorado-paycheck-calculator)

and

[California](https://smartasset.com/taxes/colorado-paycheck-calculator)

