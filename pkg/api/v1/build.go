package v1

import (
	// TODO: feels a little weird to have to import this here. should type definitions under pkg/system be moved into pkg/types?
	"github.com/mlab-lattice/lattice/pkg/definition/tree"
	"time"
)

type (
	BuildID    string
	BuildState string
)

const (
	BuildStatePending   BuildState = "pending"
	BuildStateRunning   BuildState = "running"
	BuildStateSucceeded BuildState = "succeeded"
	BuildStateFailed    BuildState = "failed"
)

type Build struct {
	ID    BuildID    `json:"id"`
	State BuildState `json:"state"`

	StartTimestamp      *time.Time `json:"startTimestamp,omitempty"`
	CompletionTimestamp *time.Time `json:"completionTimestamp,omitempty"`

	Version SystemVersion `json:"version"`
	// Services maps service paths (e.g. /foo/bar/buzz) to the
	// status of the build for that service in the Build.
	Services map[tree.NodePath]ServiceBuild `json:"services"`
}

type (
	ServiceBuildID    string
	ServiceBuildState string
)

const (
	ServiceBuildStatePending   ServiceBuildState = "pending"
	ServiceBuildStateRunning   ServiceBuildState = "running"
	ServiceBuildStateSucceeded ServiceBuildState = "succeeded"
	ServiceBuildStateFailed    ServiceBuildState = "failed"
)

type ServiceBuild struct {
	State ServiceBuildState `json:"state"`

	StartTimestamp      *time.Time `json:"startTimestamp,omitempty"`
	CompletionTimestamp *time.Time `json:"completionTimestamp,omitempty"`

	// Components maps the component name to the build for that component.
	Components map[string]ComponentBuild `json:"components"`
}

type (
	ComponentBuildID    string
	ComponentBuildState string
	ComponentBuildPhase string
)

const (
	ComponentBuildPhasePullingGitRepository ComponentBuildPhase = "pulling git repository"
	ComponentBuildPhasePullingDockerImage   ComponentBuildPhase = "pulling docker image"
	ComponentBuildPhaseBuildingDockerImage  ComponentBuildPhase = "building docker image"
	ComponentBuildPhasePushingDockerImage   ComponentBuildPhase = "pushing docker image"

	ComponentBuildStatePending   ComponentBuildState = "pending"
	ComponentBuildStateQueued    ComponentBuildState = "queued"
	ComponentBuildStateRunning   ComponentBuildState = "running"
	ComponentBuildStateSucceeded ComponentBuildState = "succeeded"
	ComponentBuildStateFailed    ComponentBuildState = "failed"
)

type ComponentBuild struct {
	State ComponentBuildState `json:"state"`

	StartTimestamp      *time.Time `json:"startTimestamp,omitempty"`
	CompletionTimestamp *time.Time `json:"completionTimestamp,omitempty"`

	LastObservedPhase *ComponentBuildPhase `json:"lastObservedPhase,omitempty"`
	FailureMessage    *string              `json:"failureMessage,omitempty"`
}

type ComponentBuildFailureInfo struct {
	Message  string `json:"message"`
	Internal bool   `json:"internal"`
}
