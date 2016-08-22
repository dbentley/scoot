package runner

// Helper functions to create ProcessStatus

func AbortStatus(runId RunId) (r ProcessStatus) {
	r.RunId = runId
	r.State = ABORTED
	return r
}

func TimeoutStatus(runId RunId) (r ProcessStatus) {
	r.RunId = runId
	r.State = TIMEDOUT
	return r
}

func ErrorStatus(runId RunId, err error) (r ProcessStatus) {
	r.RunId = runId
	r.State = FAILED
	r.Error = err.Error()
	return r
}

func BadRequestStatus(runId RunId, err error) (r ProcessStatus) {
	r.RunId = runId
	r.State = BADREQUEST
	r.Error = err.Error()
	return r
}

func RunningStatus(runId RunId) (r ProcessStatus) {
	r.RunId = runId
	r.State = RUNNING
	return r
}

func CompleteStatus(runId RunId, stdoutRef string, stderrRef string, exitCode int) (r ProcessStatus) {
	r.RunId = runId
	r.State = COMPLETE
	r.StdoutRef = stdoutRef
	r.StderrRef = stderrRef
	r.ExitCode = exitCode
	return r
}

func PreparingStatus(runId RunId) (r ProcessStatus) {
	return ProcessStatus{
		RunId: runId,
		State: PREPARING,
	}
}
