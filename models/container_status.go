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

var toString = map[ContainerStatus]string{
	Configured: "configured",
	Created:    "created",
	Running:    "running",
	Stopped:    "stopped",
	Paused:     "paused",
	Exited:     "exited",
	Removing:   "removing",
	Unknown:    "unknown",
}

func NewContainerStatus(s string) ContainerStatus {
	for k, v := range toString {
		if v == s {
			return k
		}
	}

	return BadState
}

func (p ContainerStatus) String() string {
	return toString[p]
}
