package gosrvdir

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// Credentials maps usernames to bcrypt-hashed passwords.
type Credentials map[string]string

// ParseHtpasswd parses an htpasswd file (bcrypt format).
func ParseHtpasswd(path string) (Credentials, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	creds := make(Credentials)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		creds[parts[0]] = parts[1]
	}
	return creds, scanner.Err()
}

// CheckPassword verifies a plaintext password against a bcrypt hash.
func CheckPassword(creds Credentials, user, password string) bool {
	hash, ok := creds[user]
	if !ok {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// RunHtpasswd implements the htpasswd subcommand.
func RunHtpasswd(file, username string) error {
	fmt.Print("Password: ")
	pw1, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("reading password: %w", err)
	}

	fmt.Print("Confirm password: ")
	pw2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("reading password: %w", err)
	}

	if string(pw1) != string(pw2) {
		return fmt.Errorf("passwords do not match")
	}

	hash, err := bcrypt.GenerateFromPassword(pw1, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	newLine := fmt.Sprintf("%s:%s", username, hash)

	// Read existing file if it exists
	var lines []string
	replaced := false
	if data, err := os.ReadFile(file); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 && parts[0] == username {
				lines = append(lines, newLine)
				replaced = true
			} else {
				lines = append(lines, trimmed)
			}
		}
	}

	if !replaced {
		lines = append(lines, newLine)
	}

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	if replaced {
		fmt.Printf("Updated user %q in %s\n", username, file)
	} else {
		fmt.Printf("Added user %q to %s\n", username, file)
	}
	return nil
}
