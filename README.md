# hours-worked
track if the linux user is actively logged in and no screensaver (ie, track how many hours worked)
developed specifically for ubuntu 18.04 (and 20.01) and focused on the fact that multiple users might be running on the same machine without actually logging off.

config values can be passed by cli args, or provided via `~/.hours-worked/config`

results are saved in json file `~/.hours-worked/tracked.json`,

and report can be generated `(TODO)`

------

data format

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

-----
In order to show the menu icon, `github.com/getlantern/systray` requires a few libraries: _(some may already be installed)_

```
sudo apt install gcc libgtk-3-dev libappindicator3-dev gir1.2-appindicator3
```
