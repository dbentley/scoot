package sagalog

import "errors"
import "fmt"
import "github.com/scootdev/scoot/messages"

type inMemorySagaLog struct {
	/*
	 * In memory dictionary of SagaId to SagaLog Messages
	 * element 0 is always StartSaga message
	 * last element is always EndSaga message
	 */
	sagas    map[string][]SagaMessage
	sagaJobs map[string]messages.Job
}

func InMemorySagaFactory() saga {

	inMemLog := inMemorySagaLog{
		sagas:    make(map[string][]SagaMessage),
		sagaJobs: make(map[string]messages.Job),
	}

	return saga{
		log: &inMemLog,
	}
}

//Reuse all of this and just implement LogMessage for different implementations(?)
func (log *inMemorySagaLog) LogMessage(msg SagaMessage) error {
	sagaId := msg.sagaId
	sagaLog, ok := log.sagas[sagaId]

	if ok {
		log.sagas[sagaId] = append(sagaLog, msg)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Cannot Log Saga Message %i, Never Started Saga %s", msg.msgType, sagaId))
	}
}

/*
 * Log a Start Saga Message to the log.  Returns
 * an error if it fails.
 */
func (log *inMemorySagaLog) StartSaga(sagaId string, job messages.Job) error {
	// TODO: Initialize saga log size to hold successful saga number of messages
	// StartSaga, EndSaga, 2 * numTasks (StartTask & EndTask)
	logSize := 2 + 2*len(job.Tasks)
	sagaLog := make([]SagaMessage, logSize, 1)
	sagaLog[0] = SagaMessage{
		sagaId:  sagaId,
		msgType: StartSaga,
	}

	log.sagas[sagaId] = sagaLog
	log.sagaJobs[sagaId] = job
	return nil

	//TODO: check if it already exists if true return already exists error
}
