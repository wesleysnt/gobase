package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/you/gobase/internal/auth"
	"github.com/you/gobase/internal/config"
)

var jwtGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a signed JWT",
	Example: `  gobase jwt generate --user-id=usr_123
  gobase jwt generate --user-id=usr_123 --expires-in=72h
  gobase jwt generate --user-id=usr_123 --claims='{"role":"admin"}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		userID, _ := cmd.Flags().GetString("user-id")
		if userID == "" {
			return fmt.Errorf("--user-id is required")
		}

		expiresIn, _ := cmd.Flags().GetDuration("expires-in")
		if expiresIn == 0 {
			expiresIn = 24 * time.Hour
		}

		claimsJSON, _ := cmd.Flags().GetString("claims")
		var extra map[string]interface{}
		if claimsJSON != "" {
			if err := json.Unmarshal([]byte(claimsJSON), &extra); err != nil {
				return fmt.Errorf("invalid --claims JSON: %w", err)
			}
		}

		secret := []byte(cfg.JWTSecret)
		token, _, err := auth.GenerateToken(userID, expiresIn, extra, secret)
		if err != nil {
			return fmt.Errorf("generate token: %w", err)
		}

		fmt.Println(token)
		return nil
	},
}

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT token utilities",
}

func init() {
	rootCmd.AddCommand(jwtCmd)
	jwtCmd.AddCommand(jwtGenerateCmd)

	jwtGenerateCmd.Flags().String("user-id", "", "Subject claim (user ID) — required")
	jwtGenerateCmd.Flags().Duration("expires-in", 24*time.Hour, "Token expiry duration")
	jwtGenerateCmd.Flags().String("claims", "", "Extra claims as JSON object")
}
