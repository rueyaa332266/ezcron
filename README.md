# Ezcron
Create cron expression easily

## Feature
- Translate cron expression to English
- Show the next execute time
- [WIP] Creating cron expression easily by prompt

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