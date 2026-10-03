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

const (
	LINUX_DESKTOP_ENTRY = `[Desktop Entry]
Type=Application
Name=Clipster-Desktop
Comment=A multi Platform Cloud Clipboard - Desktop Client (Go)
Exec=PLACEHOLDER
Terminal=false
`
	WIN_CREATE_SHORTCUT = `$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$HOME\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\clipster.lnk")
$Shortcut.TargetPath = 'PLACEHOLDER'
$Shortcut.Save()`
)

// desktopEntryExec quotes a path for the Exec key of a desktop entry
func desktopEntryExec(path string) string {
	r := strings.NewReplacer(`\`, `\\\\`, `"`, `\\"`, "`", "\\\\`", `$`, `\\$`)
	return `"` + r.Replace(path) + `"`
}

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

// enableAutostartLinux checks if  autostart folder exists on Linux
// if it does, creates a clipster.desktop file there for auto startup on X-Session
// pointing to current executable
func enableAutostartLinux() {
	startup_dir, startup_file := getAutostartDirAndFile()
	exec_path, err := os.Executable()
	if err != nil {
		log.Panicln("Error", err)
	}
	log.Println("Executable is: ", exec_path)
	entry := strings.Replace(LINUX_DESKTOP_ENTRY, "PLACEHOLDER", desktopEntryExec(exec_path), 1)

	if fileExists(startup_dir) {
		log.Println("Config file folder exists", startup_dir)
		if err := os.WriteFile(startup_file, []byte(entry), 0664); err != nil {
			log.Println("Error: could not write autostart file", err)
		} else {
			log.Println("Ok: written autostart file", startup_file)
			ShowNotification("Clipster", "Added Clipster to autostart by creating "+
				startup_file+"\nWhich points to "+exec_path)
		}
	} else {
		// Probabily no supported session manager
		log.Println("Error: No autostart folder exists")
		ShowNotification("Clipster", "Could not add Clipster to autostart. Folder "+
			startup_dir+" does not exist.")
	}
}

// enableAutostartWin creates a shortcut to clipster in the shell:startup folder
func enableAutostartWin() {
	_, startup_file := getAutostartDirAndFile()
	ps := New()
	exec_path, err := os.Executable()
	if err != nil {
		log.Panicln("Error", err)
	}
	log.Println("Executable is: ", exec_path)
	script := strings.Replace(WIN_CREATE_SHORTCUT, "PLACEHOLDER",
		strings.ReplaceAll(exec_path, "'", "''"), 1)
	_, stdErr, err := ps.execute(script)
	if err != nil {
		log.Println("Error: could not create shortcut in startup folder", err, stdErr)
		ShowNotification("Clipster", "Could not add Clipster to autostart:\n"+err.Error())
	} else {
		log.Println("Ok: shortcut created")
		ShowNotification("Clipster", "Added Clipster to autostart by creating "+
			startup_file+"\nWhich points to "+exec_path)
	}
}

// disableAutostartFile removes autostart file (Linux and Windows) and show status in Notification
func disableAutostartFile() {
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
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		disableAutostartFile()
	} else if runtime.GOOS == "darwin" {
		log.Println("RemoveAutostart not implemented on MacOS")
	}
}

// IsAutostartEnabled checks if autostart is enabled on different OSes
func IsAutostartEnabled() bool {
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		ok, _ := isAutostartEnabled()
		return ok
	} else {
		return false
	}
}

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

// execute command on powershell and return outputs
func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
