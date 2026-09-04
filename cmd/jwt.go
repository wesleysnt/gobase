package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const placeholderSecret = "change-me-to-a-random-secret"

var jwtGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a JWT secret and write it to .env",
	Example: `  gobase jwt generate
  gobase jwt generate --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		secret, err := randomSecret(32)
		if err != nil {
			return fmt.Errorf("generate secret: %w", err)
		}

		if err := writeEnvSecret(".env", secret, force); err != nil {
			return err
		}

		fmt.Println("Generated JWT_SECRET and wrote it to .env")
		return nil
	},
}

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT utilities",
}

func init() {
	rootCmd.AddCommand(jwtCmd)
	jwtCmd.AddCommand(jwtGenerateCmd)

	jwtGenerateCmd.Flags().Bool("force", false, "Overwrite an existing JWT_SECRET")
}

// randomSecret returns a hex-encoded cryptographically secure random string of
// the given number of bytes (32 bytes = 256 bits, suitable for HMAC-SHA256).
func randomSecret(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// writeEnvSecret sets JWT_SECRET in the .env file at path, preserving all other
// lines (including comments) and the existing file mode. If a non-placeholder
// secret is already present, it refuses to overwrite unless force is true.
func writeEnvSecret(path, secret string, force bool) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	mode := os.FileMode(0o600)
	if err == nil {
		if info, statErr := os.Stat(path); statErr == nil {
			mode = info.Mode()
		}
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "JWT_SECRET=") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "JWT_SECRET="))
		if value != "" && value != placeholderSecret && !force {
			return fmt.Errorf("JWT_SECRET already set in %s — pass --force to overwrite", path)
		}
		lines[i] = "JWT_SECRET=" + secret
		replaced = true
	}

	if !replaced {
		lines = append(lines, "", "JWT_SECRET="+secret)
	}

	out := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(out), mode); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
