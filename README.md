# callit
My Salary Calculator

```
usage: paycheck [-h|--help] [-c|--configFile "<value>"] [-S|--salary "<value>"]
                [-s|--state "<value>"] [-H|--house "<value>"] [-C|--car
                "<value>"] [-v|--version]

                Simple script to calculate bi-weekly salary before and after
                tax 😁

Arguments:

  -h  --help        Print help information
  -c  --configFile  Configuration file to be use. Default:
                    /usr/local/etc/paycheck/paycheck.conf
  -S  --salary      The yearly salary before tax, required if no set in the
                    configuration file
  -s  --state       The state where taxes is collected, required if no set in
                    the configuration file
  -H  --house       The monthly house rent/mortgage
  -C  --car         The monthly cars payment
  -v  --version     Show version

```

### Example
Look in the example director of a working configuration for the states of Colorado and California 

The brackets were taken from: 
[Colorado](https://leg.colorado.gov/agencies/legislative-council-staff/individual-income-tax%C2%A0) 

and 

[California](https://www.ftb.ca.gov/forms/2022/2022-540-tax-rate-schedules.pdf) 
__** California is 2022, I have not found the oficial 2023 schedule **__

### Tax Brackets and Tax Calculation
- the configuration used the Married Filing Jointly Tax Brackets bracket
- Allowance is set to 1 for Federal, State and local
- Insurance is Pre-Tax
- Taxable income is use base on the Standard Deduction for Married Filing Jointly
- Taxable income is use to calculate the State taxes! __** could be wrong for other states then CO **__


