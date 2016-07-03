package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func makeSnapshotCmd(c *CliClient) *cobra.Command {
	snaps := &cobra.Command{
		Use:   "snapshot",
		Short: "interact with snapshots",
	}
	create := &cobra.Command{
		Use:   "create",
		Short: "create a snapshot",
		RunE:  c.snapshotCreate,
	}
	create.Flags().String("from", "", "directory to create snapshot from")
	snaps.AddCommand(create)
	return snaps
}

func (c *CliClient) snapshotCreate(cmd *cobra.Command, args []string) error {
	fromDir, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	if fromDir == "" {
		return fmt.Errorf("--from must be set to create a snapshot")
	}
	log.Println("Calling server", fromDir)

	conn, err := c.comms.Dial()
	if err != nil {
		return err
	}

	id, err := conn.SnapshotCreate(fromDir)
	if err != nil {
		return err
	}
	fmt.Printf("%v", id)
	return nil
}
