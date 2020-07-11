# Ezcron

[![Go Report Card](https://goreportcard.com/badge/github.com/rueyaa332266/ezcron)](https://goreportcard.com/report/github.com/rueyaa332266/ezcron)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)

Ezcron is a CLI tool, helping you deal with cron expression easier.

## Feature
- Creating cron expression with prompts
- Translate cron expression into human-friendly language
- Show the next execute time
- And more (keep working) ...

## DEMO

Create cron expression like a boss 😎

![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/daily_schedule.gif)

See more DEMO at [example](#Example)

## TODO
- Add number option for Next command
- Refactor the code

## CRON Expression Format
Only support 5 fields.
```
# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
# │ │ ┌───────────── day of month (1 - 31)
# │ │ │ ┌───────────── month (1 - 12)
# │ │ │ │ ┌───────────── day of  week (0 - 6) (Sunday to Saturday)
# │ │ │ │ │
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
------------------------------------------------------------------------
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

### Create cron expression

```shell
$ ezcron expression
```
Five types of schedule are available.

- #### Time schedule:

    Create a schedule at specific time or time interval.
    - at {HH:MM}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/time_schedule_1.gif)
    - every_miniute {X_MINUTE}
    - every_hour {X_HOUR}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/time_schedule_2.gif)

- #### Daily schedule:

    Create a daily schedule at specific time.
    - every_day
    - every_day at {HH:MM}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/daily_schedule.gif)

- #### Weekly schedule:

    Create a weekly schedule on specific weekday at specific time.
    - on_every {WEEKDAY}
    - on_every {WEEKDAY} at {HH:MM}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/weekly_schedule.gif)

- #### Monthly schedule

    Create a monthly schedule on specific monthday at specific time.
    - on {MONTHDAY} of_every_month
    - on {MONTHDAY} of_every_month at {HH:MM}
    - on {MONTHDAY} of_every {X_MONTH}
    - on {MONTHDAY} of_every {X_MONTH} at {HH:MM}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/monthly_schedule.gif)

- #### Yearly schedule

    Create a yearly schedule in specific date at specific time.
    - in_every {MONTH} {MONTHDAY}
    - in_every {MONTH} {MONTHDAY} at {HH:MM}

    ![demo](https://github.com/rueyaa332266/assets/raw/master/ezcron/yearly_schedule.gif)


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