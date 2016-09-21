package client

import (
	"fmt"
	"github.com/scootdev/scoot/scootapi/gen-go/scoot"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	jobStatusSleepSeconds time.Duration = 3 * time.Second
)

func makeWatchJobCmd(c Client) *cobra.Command {
	r := &cobra.Command{
		Use:   "watch_job",
		Short: "Watch job",
		RunE:  c.watchJob,
	}

	r.Flags().StringVar(&c.addr, "addr", "localhost:9090", "address to connect to")
	return r
}

func (c Client) watchJob(cmd *cobra.Command, args []string) error {

	log.Println("Checking Status for Scoot Job", args)
	client, err := c.Dial()

	if err != nil {
		return err
	}

	jobId := args[0]

	for {
		jobStatus, err := GetAndPrintStatus(jobId, client)

		if err != nil {
			return err
		}

		if *jobStatus == scoot.Status_COMPLETED || *jobStatus == scoot.Status_ROLLED_BACK {
			return nil
		}

		time.Sleep(jobStatusSleepSeconds)
	}

}

func GetAndPrintStatus(jobId string, thriftClient *scoot.CloudScootClient) (*scoot.Status, error) {

	status, err := thriftClient.GetStatus(jobId)

	if err != nil {
		switch err := err.(type) {
		case *scoot.InvalidRequest:
			return nil, fmt.Errorf("Invalid Request: %v", err.GetMessage())
		case *scoot.ScootServerError:
			return nil, fmt.Errorf("Error getting status: %v", err.Error())
		}
	}
	PrintJobStatus(status)

	return &status.Status, nil

}

func PrintJobStatus(jobStatus *scoot.JobStatus) {
	fmt.Printf(fmt.Sprintf("Job id: %s\n", jobStatus.ID))
	fmt.Printf(fmt.Sprintf("Job status: %s\n", jobStatus.Status.String()))
	for taskId, taskStatus := range jobStatus.TaskStatus {
		fmt.Printf(fmt.Sprintf("\tTask id: %s\n", taskId))
		fmt.Printf(fmt.Sprintf("\tTask status: %s\n", taskStatus.String()))
	}
}