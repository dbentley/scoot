package client

import (
	"fmt"
	"log"
	"time"

	"github.com/scootdev/scoot/tests/testhelpers"
	"github.com/spf13/cobra"
)

type smokeTestCmd struct {
	numJobs  int
	numTasks int
	timeout  time.Duration
}

func (c *smokeTestCmd) registerFlags() *cobra.Command {
	r := &cobra.Command{
		Use:   "run_smoke_test",
		Short: "Smoke Test",
	}
	r.Flags().IntVar(&c.numJobs, "num_jobs", 100, "number of jobs to run")
	r.Flags().IntVar(&c.numTasks, "num_tasks", -1, "number of tasks per job, or random if -1")
	r.Flags().DurationVar(&c.timeout, "timeout", 180*time.Second, "how long to wait for the smoke test")
	return r
}

func (c *smokeTestCmd) run(cl *simpleCLIClient, cmd *cobra.Command, args []string) error {
	fmt.Println("Starting Smoke Test")
	runner := &smokeTestRunner{cl: cl}
	return runner.run(c.numJobs, c.numTasks, c.timeout)
}

type smokeTestRunner struct {
	cl *simpleCLIClient
}

func (r *smokeTestRunner) run(numJobs int, numTasks int, timeout time.Duration) error {
	jobs := make([]string, 0, numJobs)

	for i := 0; i < numJobs; i++ {
		for {
			id, err := testhelpers.GenerateAndStartJob(r.cl.scootClient, numTasks)
			if err == nil {
				jobs = append(jobs, id)
				break
			}
			// retry starting job until it succeeds.
			// this is useful for testing where we are restarting the scheduler
			log.Printf("Error Starting Job: Retrying %v", err)
		}
	}
	return testhelpers.WaitForJobsToCompleteAndLogStatus(
		jobs, r.cl.scootClient, timeout)
}
