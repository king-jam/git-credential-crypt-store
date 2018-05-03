package dialogs

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// PasswordBox displays a dialog box, returning the entered value and a error
func PasswordBox(message string) (string, error) {
	out, err := exec.Command(
		"zenity", "--entry",
		"--title", "Encryption Key Password",
		"--text", message,
		"--entry-text", "--hide-text").Output()
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
			"zenity", "--forms", "--add-password", "'Password'",
			"--add-password", "'Confirm Password'",
			"--title", "'Encryption Key Password Creation'",
			"--text", "'Please create a localized passphrase to encrypt/decrypt local password'").Output()
		// NOTE: exit code 1 = cancel was pressed
		if err != nil {
			return "", fmt.Errorf("Failure to get user password")
		}
		log.Printf("%s\n", string(out))
		parts := strings.SplitN(string(out), "|", 2)
		log.Printf("%+v", parts)
		if strings.TrimSpace(parts[0]) != strings.TrimSpace(parts[1]) {
			out, err = exec.Command(
				"zenity",
				"--error",
				"--text",
				"\"Passwords Do Not Match\"",
			).Output()
			if err != nil {
				return "", fmt.Errorf("Failure to get user password")
			}
		} else {
			return strings.TrimSpace(parts[0]), nil
		}
	}
}
