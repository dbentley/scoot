package sagalog

import "fmt"
import "testing"
import "github.com/golang/mock/gomock"

import msg "github.com/scootdev/scoot/messages"

//import mock "./mock"

func TestStartSaga(t *testing.T) {

	id := "testSaga"
	job := msg.Job{
		Id:      "1",
		Jobtype: "testJob",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().StartSaga(id, job)

	s := saga{
		log: sagaLogMock,
	}

	err := s.StartSaga(id, job)
	if err != nil {
		t.Error(fmt.Sprintf("Expected StartSaga to not return an error"))
	}
}

func TestEndSaga(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: EndSaga,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.EndSaga(entry.sagaId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected EndSaga to not return an error"))
	}
}

func TestAbortSaga(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: AbortSaga,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.AbortSaga(entry.sagaId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected AbortSaga to not return an error"))
	}
}

func TestStartTask(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: StartTask,
		taskId:  "task1",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.StartTask(entry.sagaId, entry.taskId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected StartTask to not return an error"))
	}
}

func TestEndTask(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: EndTask,
		taskId:  "task1",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.EndTask(entry.sagaId, entry.taskId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected EndTask to not return an error"))
	}
}

func TestStartCompTask(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: StartCompTask,
		taskId:  "task1",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.StartCompensatingTask(entry.sagaId, entry.taskId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected StartCompensatingTask to not return an error"))
	}
}

func TestEndCompTask(t *testing.T) {
	entry := SagaMessage{
		sagaId:  "1",
		msgType: EndCompTask,
		taskId:  "task1",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sagaLogMock := NewMockSagaLog(mockCtrl)
	sagaLogMock.EXPECT().LogMessage(entry)

	s := saga{
		log: sagaLogMock,
	}

	err := s.EndCompensatingTask(entry.sagaId, entry.taskId)
	if err != nil {
		t.Error(fmt.Sprintf("Expected EndCompensatingTask to not return an error"))
	}
}
