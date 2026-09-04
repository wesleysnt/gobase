package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/database"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrate("up", cmd)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrate("down", cmd)
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file pair",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.CreateMigration(args[0])
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		// For now: run up with steps=0 to show pending status
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		db, err := database.Connect(cfg)
		if err != nil {
			return fmt.Errorf("database: %w", err)
		}
		defer db.Close()

		// Use golang-migrate to get version info
		version, dirty, err := database.MigrationStatus(db, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
		return nil
	},
}

func runMigrate(direction string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer db.Close()

	steps, _ := cmd.Flags().GetInt("steps")
	return database.RunMigrations(db, cfg.DatabaseURL, direction, steps)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	migrateUpCmd.Flags().Int("steps", 0, "Max migrations to apply (0 = all)")
	migrateDownCmd.Flags().Int("steps", 1, "Number of migrations to roll back")
}
