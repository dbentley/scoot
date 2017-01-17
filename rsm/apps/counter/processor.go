package main

import (
	"log"
	"time"
)

func processor(q *Queue, agentID string) {
	for {
		log.Printf("%v: processing", agentID)
		taskID, err := q.Process(agentID)
		if err != nil {
			panic(err)
		}

		if taskID != INVALID {
			log.Printf("%v: processed %v", agentID, taskID)
		}
		time.Sleep(2000 * time.Millisecond)
	}
}
