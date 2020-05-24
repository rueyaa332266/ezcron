# Ezcron
Create cron expression easily.

## Feature
- Creating cron expression easily by prompt
- Translate cron expression into human-friendly language
- Show the next execute time

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
  expression  Create a cron expression
  help        Help about any command
  next        Return next execute time
  translate   Translate into human-friendly language

Flags:
  -h, --help   help for ezcron

Use "ezcron [command] --help" for more information about a command.
```

## Example

### Translate cron expression

```
$ ezcron translate "* * * * *"
At every minute
```

It also works when passing the cron expression by pipe.
```
$ echo "* * * * *" | ezcron
At every minute
```

### Show next execute time
```
$ ezcron next "* * * * *"
Next execute time: 2020-05-10 22:35:00 +0900 JST
```