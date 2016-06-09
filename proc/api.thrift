service Proc {

  exception InvalidRequest {
    1: optional string message
  }

  exception CanNotScheduleNow {
    1: optional i64 retryAfterMs
  }

  enum JobType {
    UNKNOWN=1,
    IRON_TESTS=2,
  }

  struct Command {
    1: list<string> args
  }

  struct Task {
    1: optional string id
    2: required Command command,
    3: optional string snapshotId,
  }

  struct Job {
    1: required string id,
    2: required list<Task> tasks,
  }

  Job RunJob(
    1: required list<Task> tasks
    2: optional JobType jobType,
  ) throws (
    1: InvalidRequest ir,
    2: CanNotScheduleNow cnsn,
  )
}
