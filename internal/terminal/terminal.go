package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func GetDBPassword() ([]byte, error) {
	fmt.Print("Password: ")
	fd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("read password from terminal: %w", err)
	}

	return fd, nil
}
