package executor

const (
	StatusPending Status = iota
	StatusSkipped
	StatusRunning
	StatusCancelled
	StatusSuccess
	StatusFailure
)

const (
	ContainerStatusUnknown ContainerStatus = iota
	ContainerStatusCreated
	ContainerStatusRunning
	ContainerStatusPausing
	ContainerStatusPaused
	ContainerStatusStopped
	ContainerStatusExited
)

var (
	StatusNames = map[Status]string{
		StatusPending:   "pending",
		StatusSkipped:   "skipped",
		StatusRunning:   "running",
		StatusCancelled: "cancelled",
		StatusSuccess:   "success",
		StatusFailure:   "failure",
	}

	ContainerStatusNames = map[ContainerStatus]string{
		ContainerStatusUnknown: "unknown",
		ContainerStatusCreated: "created",
		ContainerStatusRunning: "running",
		ContainerStatusPausing: "pausing",
		ContainerStatusPaused:  "paused",
		ContainerStatusStopped: "stopped",
		ContainerStatusExited:  "exited",
	}

	DefaultStatus          = StatusPending
	DefaultContainerStatus = ContainerStatusUnknown
)
