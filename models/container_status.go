package models

type ContainerStatus string

const (
	Configured ContainerStatus = "configured"
	Created    ContainerStatus = "created"
	Running    ContainerStatus = "running"
	Stopped    ContainerStatus = "stopped"
	Paused     ContainerStatus = "paused"
	Exited     ContainerStatus = "exited"
	Removing   ContainerStatus = "removing"
	Unknown    ContainerStatus = "unknown"
	BadState   ContainerStatus = "bad_state"
)

var statusList = []ContainerStatus{
	Configured, Created, Running, Stopped, Paused, Exited, Removing, Unknown, BadState,
}

func NewContainerStatus(s string) ContainerStatus {
	for _, v := range statusList {
		if v.String() == s {
			return v
		}
	}

	return BadState
}

func (p ContainerStatus) String() string {
	return string(p)
}
