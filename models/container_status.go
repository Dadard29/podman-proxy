package models

type ContainerStatus int

const (
	Configured ContainerStatus = iota
	Created
	Running
	Stopped
	Paused
	Exited
	Removing

	Unknown
	BadState
)

func NewContainerStatus(s string) ContainerStatus {
	switch s {
	case "configured":
		return Configured
	case "created":
		return Created
	case "running":
		return Running
	case "stopped":
		return Stopped
	case "paused":
		return Paused
	case "exited":
		return Exited
	case "removing":
		return Removing
	case "unknown":
		return Unknown
	}
	return BadState
}

func (p ContainerStatus) String() string {
	return []string{"configured", "created", "running", "stopped", "paused", "exited", "removing", "unknown", "bad state"}[p]
}
