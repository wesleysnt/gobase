package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/you/gobase/internal/config"
	"github.com/you/gobase/internal/database"
	"github.com/you/gobase/internal/log"
	"github.com/you/gobase/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		// Override from flags if set
		if port, _ := cmd.Flags().GetInt("port"); port != 0 {
			cfg.Port = port
		}
		if env, _ := cmd.Flags().GetString("env"); env != "" {
			cfg.Env = env
		}

		logger := log.New(cfg)

		db, err := database.Connect(cfg)
		if err != nil {
			return fmt.Errorf("database: %w", err)
		}
		defer db.Close()

		router := server.NewRouter(cfg, db, logger)

		return server.ListenAndServe(cfg, logger, router)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 0, "Server port (default: from PORT env or 8080)")
	serveCmd.Flags().String("env", "", "Environment: development|production|test")
}
