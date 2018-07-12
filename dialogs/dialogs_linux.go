package dialogs

import (
	"fmt"
	"os/exec"
	"strings"
)

// PasswordBox displays a dialog box, returning the entered value and a error
func PasswordBox() (string, error) {
	out, err := exec.Command(
		"zenity", "--entry",
		"--title", "Decryption Password",
		"--text", "Please enter your decryption password",
		"--hide-text").Output()
	// NOTE: exit code 1 = cancel was pressed
	if err != nil {
		return "", fmt.Errorf("Failure to get user password")
	}
	return strings.TrimSpace(string(out)), nil
}

// PasswordCreationBox displays a password box with confirmation.
// This will only return if the user has entered matching passwords or hit the cancel box.
func PasswordCreationBox() (string, error) {
	for {
		out, err := exec.Command(
			"zenity", "--forms", "--add-password", "Password",
			"--add-password", "Confirm Password",
			"--title", "Encryption Password Creation",
			"--text", "Please create a localized passphrase to encrypt/decrypt local password").Output()
		// NOTE: exit code 1 = cancel was pressed
		if err != nil {
			return "", fmt.Errorf("Failure to get user password")
		}
		parts := strings.SplitN(string(out), "|", 2)
		if strings.TrimSpace(parts[0]) != strings.TrimSpace(parts[1]) {
			_, err = exec.Command(
				"zenity",
				"--error",
				"--text",
				"Passwords Do Not Match",
			).Output()
			if err != nil {
				return "", fmt.Errorf("Failure to get user password")
			}
		} else {
			return strings.TrimSpace(parts[0]), nil
		}
	}
}
