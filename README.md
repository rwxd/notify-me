# notify me

Small command line tool to notify myself through various services.

## Installation

Go to [Releases](https://github.com/rwxd/notify-me/releases) and download the latest release.

Or install it from source:

```bash
go install github.com/rwxd/notify-me@latest
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

# send with a custom message before the command output
notify-me uptime-kuma wrap -i uptime.com -t "<token>" -m "<message>" -- ping -c 1 google.com

# only send down notification
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --fail -- ping -c 1 google.com

# only send success notification
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --success -- ping -c 1 google.com

# send down notification when ping succeeds
notify-me uptime-kuma wrap -i uptime.com -t "<token>" --reverse -- ping -c 1 google.com

# send custom message
notify-me uptime-kuma wrap -i uptime.com -t "<token>" -m "<message>" -- ping -c 1 google.com
```

### ntfy


```bash
# send notification
notify-me ntfy -t "<topic>" -m "<message>"

# send notification to self hosted instance with username and password
notify-me ntfy -i "<instance>" -u "<user>" -p "<pass>" -t "<topic>" -m "<message>"

# send notification with token for authentication
notify-me ntfy -i "<instance>" --token "<token>" -t "<topic>" -m "<message>"

# add title
notify-me ntfy -t "<topic>" -m "<message>" -T "title"

# set priority to high
notify-me ntfy -t "<topic>" -m "<message>" -p "high"

# add tags
notify-me ntfy -t "<topic>" -m "<message>" --Tags "tag1,tag2"

# add url to open when clicking on notification
notify-me ntfy -t "<topic>" -m "<message>" --url "https://example.com"

# more options
notify-me ntfy --help
```

#### Wrap a command

Runs a command and sends a notification with the output.

Supports all the options from the ntfy command.

```bash
# wrap ping command
notify-me ntfy wrap -t "<topic>" -- ping -c 1 google.com

# send with a custom message before the command output
notify-me ntfy wrap -t "<topic>" -m "<message>" -- ping -c 1 google.com

# send notification to self hosted instance with username and password
notify-me ntfy wrap -i "<instance>" -u "<user>" -p "<pass>" -t "<topic>" -- ping -c 1 google.com

# only send notification when command fails
notify-me ntfy wrap -t "<topic>" --fail -- ping -c 1 google.com

# only send notification when command succeeds
notify-me ntfy wrap -t "<topic>" --success -- ping -c 1 google.com

# more options
notify-me ntfy wrap --help
```
