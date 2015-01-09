package api

const (
	StartedStatus = "started"
	StoppedStatus = "stopped"
	ErrorStatus = "error"
	)

type PanamaxAdapter interface {
	GetServices() ([]*Service)
	GetService(string) (*Service)
	CreateServices([]*Service) ([]*Response)
	UpdateService(*Service) (bool)
	DestroyService(*Service) (bool)
}

type Service struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Source       string `json:"source"`
	Command      string `json:"command,omitempty"`
	Links        []*Link `json:"links,omitempty"`
	Ports        []*Port `json:"ports,omitempty"`
	Expose       []uint16 `json:"expose,omitempty"`
	Environment  []*Environment `json:"environment,omitempty"`
	Volumes      []*Volume `json:"volumes,omitempty"`
	DesiredState string `json:"desiredState,omitempty"`
	CurrentState string `json:"currentState,omitempty"`
}

type Link struct {
 	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type Port struct {
	HostPort      uint16 `json:"hostPort,omitempty"`
	ContainerPort uint16 `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

type Environment struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

type Volume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
}

type Response struct {
	Id	      string `json:"id"`
	CurrrentState string `json:"currentState,omitempty"`
}
