package sagalog

import "github.com/scootdev/scoot/messages"

/*
 *  SagaLog Interface, Implemented
 */
type SagaLog interface {
	StartSaga(sagaId string, job messages.Job) error
	LogMessage(message SagaMessage) error
}
