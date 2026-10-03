# Privacy Policy

Clipster contains no ads, analytics or trackers. It only talks to the Clipster server you configured: the public server at [clipster.cc](https://clipster.cc) or your own. It only does so when you log in, register or pick an action from its tray menu.

## Data sent to the server

- Your username and a hash of your password (used to log in, your password itself is never sent)
- Your clips, encrypted on your device. The server can't read them. A clip is only sent when you select `Share Clip`.
- For each clip: its format (text or image), the device name `desktop` and the time it was shared

The server keeps only your newest clips (5 on the public server) and deletes older ones.

## Data on your device

Your username, server address and the hashes derived from your password are stored in `.config/clipster/config.toml` in your home folder, readable only by your user. Delete that folder to remove them.

If you enable `Autostart Clipster`, Clipster creates a start entry: `~/.config/autostart/clipster.desktop` on Linux, a `clipster.lnk` shortcut in your Startup folder on Windows. Disabling autostart removes it again.
