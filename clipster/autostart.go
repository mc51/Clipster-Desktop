// Deal with enabling and disabling auto start of Clipter on Desktop startup
package clipster

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type PowerShell struct {
	powerShell string
}

var (
	LINUX_DESKTOP_ENTRY = `[Desktop Entry]
Type=Application
Name=Clipster-Desktop
Comment=A multi Platform Cloud Clipboard - Desktop Client (Go)
Exec=PLACEHOLDER
Terminal=false
`
	WIN_CREATE_SHORTCUT = `$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$HOME\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\clipster.lnk")
$Shortcut.TargetPath = "PLACEHOLDER"
$Shortcut.Save()`
)

// getAutostartDirAndFile returns the absolute path to auto startup directory and file
// for different OSes
func getAutostartDirAndFile() (string, string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("Error", err)
		return "", ""
	}
	log.Println("Homedir is: ", homedir)

	if runtime.GOOS == "linux" {
		path_dir := filepath.Join(homedir, ".config", "autostart")
		path_file := filepath.Join(path_dir, "clipster.desktop")
		return path_dir, path_file
	} else if runtime.GOOS == "windows" {
		path_dir := filepath.Join(homedir,
			"AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup")
		path_file := filepath.Join(path_dir, "clipster.lnk")
		return path_dir, path_file
	} else {
		return "", ""
	}
}

// isAutostartEnabled checks if an autostart file exists and returns its absolute path
func isAutostartEnabled() (bool, string) {
	_, file := getAutostartDirAndFile()
	if fileExists(file) {
		log.Println("Ok: Autostart file exists", file)
		return true, file
	} else {
		log.Println("Warning: No Autostart file exists", file)
		return false, ""
	}
}

// enableAutostartLinux checks if  autostart folder ~/.config/autostart exists on Linux
// if it does, creates a clipster.desktop there for auto startup on X-Session
func enableAutostartLinux() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("Error", err)
	}
	log.Println("Homedir is: ", homedir)
	exec_path, err := os.Executable()
	if err != nil {
		log.Panicln("Error", err)
	}
	log.Println("Executable is: ", exec_path)
	LINUX_DESKTOP_ENTRY = strings.Replace(LINUX_DESKTOP_ENTRY, "PLACEHOLDER", exec_path, 1)

	path := filepath.Join(homedir, ".config", "autostart")
	if fileExists(path) {
		log.Println("Config file folder exists", path)
		_, file := getAutostartDirAndFile()
		if err := os.WriteFile(file, []byte(LINUX_DESKTOP_ENTRY), 0664); err != nil {
			log.Println("Error: could not write autostart file", err)
		} else {
			log.Println("Ok: written autostart file", file)
			ShowNotification("Clipster", "Added Clipster to autostart by creating "+file)
		}
	} else {
		// Probabily no supported session manager
		log.Println("Error: No autostart folder exists")
		ShowNotification("Clipster", "Could not add Clipster to autostart. Folder "+
			path+" does not exist.")
	}
}

// enableAutostartWin creates a shortcut to clipster in the shell:startup folder
func enableAutostartWin() {
	ps := New()
	exec_path, err := os.Executable()
	if err != nil {
		log.Panicln("Error", err)
	}
	log.Println("Executable is: ", exec_path)
	WIN_CREATE_SHORTCUT = strings.Replace(WIN_CREATE_SHORTCUT, "PLACEHOLDER", exec_path, 1)

	stdOut, stdErr, err := ps.execute(WIN_CREATE_SHORTCUT)
	log.Printf("CreateShortcut:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s",
		strings.TrimSpace(stdOut), stdErr, err)
}

// disableAutostartLinux removes autostart file ~/.config/autostart/clipster.desktop
// and show status in Notification
func disableAutostartLinux() {
	if ok, file := isAutostartEnabled(); ok {
		if err := os.Remove(file); err != nil {
			log.Println("Error: could not remove autostart file", file)
			ShowNotification("Clipster", "Could not remove autostart file "+
				file+"\n"+err.Error())
		} else {
			log.Println("Ok: removed autostart file " + file)
			ShowNotification("Clipster", "Removed autostart file "+file)
		}
	}
}

// enableAutostart deals with autostart of Clipster on different OSes
func enableAutostart() {
	if runtime.GOOS == "linux" {
		enableAutostartLinux()
	} else if runtime.GOOS == "darwin" {
		ShowNotification("Clipster", "To autostart Clipster, right click on the dock"+
			" icon and select\n`Options --> Open at Login`.")
	} else if runtime.GOOS == "windows" {
		enableAutostartWin()
	}
}

// disableAutostart deals with disabling autostart of Clipster on different OSes
func disableAutostart() {
	if runtime.GOOS == "linux" {
		disableAutostartLinux()
	} else if runtime.GOOS == "darwin" {
		log.Println("RemoveAutostart not implemented on MacOS")
	} else if runtime.GOOS == "windows" {
		log.Println("RemoveAutostart not implemented in Windows")
	}
}

// IsAutostartEnabled checks if autostart is enabled on different OSes
func IsAutostartEnabled() bool {
	if runtime.GOOS == "linux" {
		ok, _ := isAutostartEnabled()
		return ok
	} else if runtime.GOOS == "windows" {
		ok, _ := isAutostartEnabled()
		return ok
	} else {
		return false
	}
}

// TODO: Implement for Windows as well

// ToggleAutostart enable or disable autostart
func ToggleAutostart() {
	if IsAutostartEnabled() {
		disableAutostart()
	} else {
		enableAutostart()
	}
}

// New create new session
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
