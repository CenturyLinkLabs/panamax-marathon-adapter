package api

import (
	"fmt"
)

// PanamaxAdapter encapulates the CRUD operations for Services.
// These methods must be implemented to fulfill the adapter contract.
type PanamaxAdapter interface {
	GetServices() ([]*Service, *Error)
	GetService(string) (*Service, *Error)
	CreateServices([]*Service) ([]*Service, *Error)
	UpdateService(*Service) *Error
	DestroyService(string) *Error
}

// A Service describes the information needed to deploy and
// scale a desired application. The Id is required and the
// actualState is used to provide status back to the remote
// agent.
type Service struct {
	Id          string         `json:"id"`
	Name        string         `json:"name,omitempty"`
	Source      string         `json:"source,omitempty"`
	Command     string         `json:"command,omitempty"`
	Links       []*Link        `json:"links,omitempty"`
	Ports       []*Port        `json:"ports,omitempty"`
	Expose      []uint16       `json:"expose,omitempty"`
	Environment []*Environment `json:"environment,omitempty"`
	Volumes     []*Volume      `json:"volumes,omitempty"`
	VolumesFrom []*VolumesFrom `json:"volumes_from,omitempty"`
	ActualState string         `json:"actualState,omitempty"`
	Deployment  Deployment     `json:"deployment,omitempty"`
}

// Deployment structure contains the deployment count
// for a service.
type Deployment struct {
	Count int `json:"count,omitempty"`
}

// A Link is equivalent to the docker link command.
// It contains the named of a service and the desired alias.
type Link struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// A Port is desribes a docker port mapping.
type Port struct {
	HostPort      uint16 `json:"hostPort,omitempty"`
	ContainerPort uint16 `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

// A Environment is a structure which contains environmental variables.
// They are equivalent to the -e "Name=Value" on the docker command line.
type Environment struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

// A Volume is used to mount a host directory into the container and is
// translated into a -v docker command.
type Volume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
}

type VolumesFrom struct {
	Name string `json:"name"`
}

// Metadata contains informational data about the current adapter.
type Metadata struct {
	Version string `json:"version"`
	Type    string `json:"type"`
}

// Error is an application specific error structure which
// encapsulates an error code and message.
type Error struct {
	Code    int
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error(%d): %s", e.Code, e.Message)
}

// NewError creates an error instance with the specified code and message.
func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg}
}
