package api

import (
	"fmt"
)

// PanamaxAdapter encapulates the CRUD operations for Services
type PanamaxAdapter interface {
	GetServices() ([]*Service, *Error)
	GetService(string) (*Service, *Error)
	CreateServices([]*Service) ([]*Service, *Error)
	UpdateService(*Service) *Error
	DestroyService(string) *Error
}

// Service structure with nested elements
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

type VolumesFrom struct {
	Name string `json:"name"`
}

// Metadata contains informational data about the current adapter.
type Metadata struct {
	Version string `json:"version"`
	Type    string `json:"type"`
}

// Error is a serializable Error structure.
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
