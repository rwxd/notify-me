# notify me

Small command line tool to notify myself through various services.

## Installation

Go to [Releases](https://github.com/rwxd/notify-me/releases) and download the latest release.

Or install it from source:

```bash
go install github.com/rxwxd/notify-me
```

## Usage

### Uptime Kuma

```bash
# send up notification
notify-me uptime-kuma -i uptime.com -t "<token>" --up -m "<message>"

# send down notification
notify-me uptime-kuma -i uptime.com -t "<token>" --down -m "<message>"

# send ping value
notify-me uptime-kuma -i uptime.com -t "<token>" --down -m "<message>" --ping "23"
```

#### Wrap a command

Runs a command and sends the up notification if the command exits with 0, otherwise sends a down notification.

```bash
# wrap ping command
notify-me uptime-kuma wrap -i uptime.com -t "<token>" -- ping -c 1 google.com

# only send down notification
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --error -- ping -c 1 google.com

# only send success notification
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --success -- ping -c 1 google.com

# send down notification when ping succeeds
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --reverse -- ping -c 1 google.com

# send custom message
notify-me uptime-kuma wrap -i uptime.com -t "<token>" -m "<message>" -- ping -c 1 google.com
```
