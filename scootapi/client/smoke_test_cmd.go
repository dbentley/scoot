package client

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/scootdev/scoot/scootapi/gen-go/scoot"
	"github.com/scootdev/scoot/tests/testhelpers"
	"github.com/spf13/cobra"
	"strconv"
)

func makeSmokeTestCmd(c *Client) *cobra.Command {
	r := &cobra.Command{
		Use:   "run_smoke_test",
		Short: "Smoke Test",
		RunE:  c.runSmokeTest,
	}

	r.Flags().StringVar(&c.addr, "addr", "localhost:9090", "address to connect to")
	return r
}

func (c *Client) runSmokeTest(cmd *cobra.Command, args []string) error {
	fmt.Println("Starting Smoke Test")

	numTasks := 100

	if (len(args)) > 0 {
		var err error
		numTasks, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}

	timeout := 10 * time.Second
	if (len(args)) > 1 {
		var err error
		timeout, err = time.ParseDuration(args[1])
		if err != nil {
			return err
		}
	}
	// run a bunch of concurrent jobs and track their status
	ch := make(chan jobAndStatus)
	timedOutJobs := make(chan error)
	var wg sync.WaitGroup
	errCh := make(chan error, numTasks)
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func() {
			err := c.generateAndRunJob(timeout, ch)
			if err != nil {
				timedOutJobs <- err
				errCh <- err
				fmt.Println(err)
			}
			wg.Done()
		}()
	}

	jobStatusMap := make(map[string]scoot.Status)
	ticker := time.NewTicker(time.Millisecond)
	timeouts := 0
Loop:
	for {
		select {
		case <-timedOutJobs:
			timeouts++
		case <-ticker.C:
			// job id's grouped by status
			statusJobMap := make(map[scoot.Status][]string)
			jobAndStatus := <-ch
			jobStatusMap[jobAndStatus.job] = jobAndStatus.status
			for job, status := range jobStatusMap {
				// populate statusJobMap
				statusJobMap[status] = append(statusJobMap[status], job)
			}
			fmt.Println(timeouts, "jobs have timed out")
			for status, jobs := range statusJobMap {
				sort.Sort(sort.StringSlice(jobs))
				fmt.Println(status, ":", jobs)
				// if all jobs are completed, break loop
				if len(statusJobMap[scoot.Status_COMPLETED])+len(statusJobMap[scoot.Status_ROLLED_BACK])+timeouts == numTasks {
					ticker.Stop()
					break Loop
				}
			}
		}
	}

	wg.Wait()

	// if any errors were logged return an error
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *Client) generateAndRunJob(timeout time.Duration, ch chan jobAndStatus) error {
	client, err := c.Dial()

	if err != nil {
		return err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// We just want the JobDefinition here Id doesn't matter
	job := testhelpers.GenJobDefinition(rng)
	jobId, err := client.RunJob(job)

	// Error Enqueuing Job
	if err != nil {
		switch err := err.(type) {
		case *scoot.InvalidRequest:
			return fmt.Errorf("Invalid Request: %v", err.GetMessage())
		default:
			return fmt.Errorf("Error running job: %v %T", err, err)
		}
	}
	fmt.Println("Successfully Scheduled Job", jobId.ID)

	// Check Job Status
	jobInProgress := true
	timeSpent := 0 * time.Second
	for jobInProgress && timeSpent < timeout {
		status, err := client.GetStatus(jobId.ID)
		if status.Status == scoot.Status_COMPLETED || status.Status == scoot.Status_ROLLED_BACK {
			jobInProgress = false
		}

		if err != nil {
			switch err := err.(type) {
			case *scoot.InvalidRequest:
				return fmt.Errorf("Invalid Request: %v", err.GetMessage())
			case *scoot.ScootServerError:
				return fmt.Errorf("Error getting status: %v", err.Error())
			}
		}
		// send it back with updated status
		ch <- jobAndStatus{job: jobId.ID, status: status.Status}
		time.Sleep(50 * time.Millisecond)
		timeSpent += 50 * time.Millisecond
	}

	if jobInProgress {
		return fmt.Errorf("Could Not Complete Jobs in Alloted Time %v", timeout)
	} else {
		return nil
	}
}

// struct for passing status to UI
type jobAndStatus struct {
	job    string
	status scoot.Status
}
