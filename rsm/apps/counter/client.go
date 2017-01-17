package main

import (
	"log"
	"math/rand"
	"time"
)

func client(q *Queue, clientID string, priority int, weight int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		var taskID TaskID
		for {
			var err error
			log.Printf("%v: adding %d, %d", clientID, priority, weight)
			taskID, err = q.AddTask(priority, weight)
			if err == nil {
				log.Printf("%v: created task %v", clientID, taskID)
				break
			}
			if err, ok := err.(QueueBusyError); ok {
				log.Printf("%v: got busy signal %v", clientID, err)
			} else {
				panic(err)
			}
		}

		for {
			st, err := q.Status(taskID)
			if err != nil {
				panic(err)
			}
			if st.WeightAhead != -1 {
				log.Printf("%v: waiting on task %v (behind %v weight)", clientID, taskID, st.WeightAhead)
			} else {
				log.Printf("%v: task %v done by %v", clientID, taskID, st.Processor)
				break
			}
			time.Sleep(time.Duration(r.Intn(3000)) * time.Millisecond)
		}
	}
}
