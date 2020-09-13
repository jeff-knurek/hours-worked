# hours-worked
A little tool that tracks if the user is actively logged in and no screensaver. ie, track how many hours worked (or just tracks how much time spent in front of the screen).

Developed specifically for ubuntu 18.04 _(soon to be tested on 20.01)_. Focused on the fact that multiple users might be running on the same machine without actually logging off, so it tracks the UI activity of the user.

Config values can be passed by cli args, and/or provided via `~/.hours-worked/config.toml`. If that files doesn't exist on first run, a default one will be created.

## Install

In order to show the menu icon, `github.com/getlantern/systray` requires a few libraries: _(some may already be installed)_

```
sudo apt install gcc libgtk-3-dev libappindicator3-dev gir1.2-appindicator3
```

The latest release can be downloaded directly from the [release page](/releases)

### systemctl

While you can use any option to run this, if splitting across multiple users, `systemctl` might be a good approach.

Save the binary to `/usr/local/bin` and create a file: `~/.local/share/systemd/user/hours-worked.service` with:

```
[Unit]
Description=track hours worked

[Service]
Type=simple
ExecStart=/usr/local/bin/hours-worked track
Restart=always
StandardOutput=journal

[Install]
WantedBy=default.target
```

And run:

```
systemctl --user enable hours-worked.service
```

Logs can be found: `journalctl --user | grep hours-worked`

## Reporting Time

Results of time tracked are saved in json file `~/.hours-worked/tracked.json`,

A simple text report can be generated in cli running: `hours-worked report`. _(future reporting output formats are to be developed)_

## Contributing

It should be fairly easy to port this solution to also work with Mac, but would require someone else to test/implement.

While the tool is simple in what it does and what it's used for, opening issues and new functionality PRs are most welcome.

data format of time tracked:
```
{
    <USER>: {
        "2020": {           //year
            "January": {    //month
                "1": 0,
                "2": 180,   // day: minutes active
            },
        }
    },
}
```