# Clipster - Desktop Client (Go)

[![GitHub Actions Build Workflow](https://github.com/mc51/Clipster-Desktop/workflows/Build/badge.svg)](https://github.com/mc51/Clipster-Desktop/actions)  

Clipster is a multi platform cloud clipboard:  
Copy a text on your smartphone and paste it on your desktop, or vice versa.  
Easy, secure, open source.  
Supports Android, Linux, MacOS, Windows and all browsers.   

You can use the web front-end of the public server at [clipster.cc](https://clipster.cc).  
For the Android client see [Clipster-Android](https://github.com/mc51/Clipster-Android).  
To run your own server check [Clipster-Server](https://github.com/mc51/Clipster-Server).  
There is an alternative [Clipster-Desktop](https://github.com/mc51/Clipster-Desktop-Py) implementation written in Python.
  
![Clipster demo](assets/demo_01.gif)  
  
## Setup

### Linux 

Download [`clipster`](https://github.com/mc51/Clipster-Desktop/releases/latest/download/clipster) from the latest Linux release and start it. To have Clipster auto start, add it to `Application Autostart`.

Clipster depends on gtk-3.0. To install it (Ubuntu/Debian):
`sudo apt-get install libgtk-3-0`  

### Windows (coming soon...)

Download [`clipster.exe`](https://github.com/mc51/Clipster-Desktop/releases/latest/download/clipster.exe) from the latest Windows release and start it. To have Clipster auto start for the current user, open the startup folder by opening Explorer and typing `shell:startup`. Copy `clipster.exe` there. 

### MacOS (coming soon...)

Download [`clipster_mac.zip`](https://github.com/mc51/Clipster-Desktop/releases/latest/download/clipster_mac.zip) from the latest MacOS release, move it to `Applications` and start it via right-click -> open. You might get a warning message, that you need to ignore. If that fails:
Go to `System Preferences --> Security & Privacy`. In the `General` Tab the App will be listed and you can start it from there.  
  
To automatically start Clipster, right click on the icon in your Dock and click on `Options --> Open at Login`.  
  
Now, you can [use](#usage) clipster!  
  
## Usage

On the first startup, you can register a new account or enter your existing credentials for the login. Your credentials will be stored in your `HOMEPATH` in `./config/clipster/config`.  
Clipster will add an Icon to your system tray which you can click for opening up a menu with the following options:  
`Get last Clip` will fetch the last shared Clip from the server and put it into your clipboard.  
`Get all Clips` will fetch all shared Clips from the server and display them to you.  
`Share Clip` will share your current clipboard. Then, it's available for all your devices.  
`Edit Credentials` allows you to register a new account or change your login credentials.  
`Quit` will terminate the app.  

## Roadmap

- [x] Encrypt / Decrypt clipboard locally and only transmit encrypted data to server
- [x] Add clipboard history: share multiple Clips
- [x] Add PyPi package
- [ ] Support image sharing
  
## Contributions

Contributions are very welcome. If you come across a bug, please open an issue. The same thing goes for feature requests.

## Credits

- GUI based on [goey](https://pkg.go.dev/bitbucket.org/rj/goey)
- [Systray](https://pkg.go.dev/github.com/getlantern/systray) for tray icon and menu
- Notifications using [beep](https://github.com/gen2brain/beeep)
- Config by [Viper](https://github.com/spf13/viper)
- Crypto using [PBKDF2](https://pkg.go.dev/golang.org/x/crypto/pbkdf2) and [Fernet](https://github.com/fernet/fernet-go)