# Ezcron
Create cron expression easily.

## Feature
- Translate cron expression to English
- Show the next execute time
- [WIP] Creating cron expression easily by prompt

## CRON Expression Format
Only support 5 space-separated fields.
```
Field name     Mandatory?   Allowed values    Allowed special characters
----------     ----------   --------------    --------------------------
Minutes        Yes          0-59              * / , -
Hours          Yes          0-23              * / , -
Day of month   Yes          1-31              * / , -
Month          Yes          1-12 or JAN-DEC   * / , -
Day of week    Yes          0-6 or SUN-SAT    * / , -
```

## Installing

```
go get -u github.com/rueyaa332266/ezcron
```

## Usage
```
Usage:
  ezcron [flags]
  ezcron [command]

Available Commands:
  help        Help about any command
  next        return next execute time

Flags:
  -h, --help   help for ezcron

Use "ezcron [command] --help" for more information about a command.
```

## Example

### Translate cron expression
Pass the cron expression with pipe

```
$ echo "* * * * *" | ezcron
At every minute
```

### Show next execute time
```
$ ezcron next "* * * * *"
Next execute time: 2020-05-10 22:35:00 +0900 JST
```