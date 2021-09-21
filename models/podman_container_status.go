package models

type PodmanContainerStatus int

const (
	Configured PodmanContainerStatus = iota
	Created
	Running
	Stopped
	Paused
	Exited
	Removing

	Unknown
	BadState
)

func NewPodmanContainerStatus(s string) PodmanContainerStatus {
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

func (p PodmanContainerStatus) String() string {
	return []string{"configured", "created", "running", "stopped", "paused", "exited", "removing", "unknown", "bad state"}[p]
}
