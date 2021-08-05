package cmd

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rossi1/coding-challenge/operation"
)

func Execute(db *sql.DB, ops operation.TokenReadWriter) error {

	var CMD = &cobra.Command{
		Use:     "start",
		Short:   "start up program",
		Long:    "program CLI",
		Version: "v0.1.0",
	}

	writeTokenToFileCMD, err := writeToCmd("token", ops)

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	//readTokenStdoutCMD, err := readTokenToStdoutCmd("token", ops)

	//if err != nil {
	//	return fmt.Errorf("error: %w", err)
	//}

	readTokenDatabaseCMD, err := readTokenToDatabaseCmd("token", db, ops)

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	CMD.AddCommand(writeTokenToFileCMD)
	//CMD.AddCommand(readTokenStdoutCMD)
	CMD.AddCommand(readTokenDatabaseCMD)

	if err := CMD.Execute(); err != nil {
		return err
	}
	return nil
}

func writeToCmd(filePath string, ops operation.TokenReadWriter) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "create token",
		Short: "create token",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := ops.WriteTokenToFile(filePath, 10000000, operation.GenerateToken); err != nil {
			return err
		}
		return nil

	}
	return startCmd, nil

}

/*

func readTokenToStdoutCmd(filePath string, ops operation.TokenReadWriter) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "read token stdout",
		Short: "read token stdout",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := ops.TokenToStdout(filePath)

		if err != nil {
			return err
		}
		return nil
	}

	return startCmd, nil
}
*/

func readTokenToDatabaseCmd(filePath string, db *sql.DB, ops operation.TokenReadWriter) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "read token database",
		Short: "read token database",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := ops.TokenToDatabase(filePath, db)

		if err != nil {
			return err
		}
		return nil
	}

	return startCmd, nil
}
