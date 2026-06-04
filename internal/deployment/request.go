// Package deployment contains the deployment request domain model.
package deployment

import (
	"fmt"
	"strings"
	"time"
)

// Status is the lifecycle state of a deployment request.
type Status string

const (
	// StatusPendingApproval means the deployment request is waiting for review.
	StatusPendingApproval Status = "pending_approval"
	// StatusApproved means the deployment request has been approved.
	StatusApproved Status = "approved"
	// StatusRejected means the deployment request has been rejected.
	StatusRejected Status = "rejected"
	// StatusStarted means the deployment request has started.
	StatusStarted Status = "started"
	// StatusSucceeded means the deployment request completed successfully.
	StatusSucceeded Status = "succeeded"
	// StatusFailed means the deployment request failed.
	StatusFailed Status = "failed"
	// StatusRolledBack means the deployment request was rolled back.
	StatusRolledBack Status = "rolled_back"
)

// ReferenceKind identifies the type of external item connected to a deployment request.
type ReferenceKind string

const (
	// ReferenceKindIssue references an issue.
	ReferenceKindIssue ReferenceKind = "issue"
	// ReferenceKindPullRequest references a pull request.
	ReferenceKindPullRequest ReferenceKind = "pull_request"
)

// Reference connects a deployment request to an issue or pull request.
type Reference struct {
	Kind  ReferenceKind
	Value string
}

// Request records intent to deploy one service version to one environment.
type Request struct {
	Service      string
	Environment  string
	Version      string
	Status       Status
	RequestedBy  string
	ApprovedBy   string
	References   []Reference
	Risk         string
	RollbackPlan string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewRequestInput contains the caller-provided fields for a new deployment request.
type NewRequestInput struct {
	Service      string
	Environment  string
	Version      string
	RequestedBy  string
	References   []Reference
	Risk         string
	RollbackPlan string
}

// NewRequest validates input and returns a new pending deployment request.
func NewRequest(input NewRequestInput, now time.Time) (Request, error) {
	if strings.TrimSpace(input.Service) == "" {
		return Request{}, fmt.Errorf("service is required")
	}
	if strings.TrimSpace(input.Environment) == "" {
		return Request{}, fmt.Errorf("environment is required")
	}
	if strings.TrimSpace(input.Version) == "" {
		return Request{}, fmt.Errorf("version is required")
	}
	if strings.TrimSpace(input.RequestedBy) == "" {
		return Request{}, fmt.Errorf("requested_by is required")
	}
	for i, reference := range input.References {
		if !validReferenceKind(reference.Kind) {
			return Request{}, fmt.Errorf("references[%d].kind %q is unsupported", i, reference.Kind)
		}
	}

	return Request{
		Service:      input.Service,
		Environment:  input.Environment,
		Version:      input.Version,
		Status:       StatusPendingApproval,
		RequestedBy:  input.RequestedBy,
		References:   append([]Reference(nil), input.References...),
		Risk:         input.Risk,
		RollbackPlan: input.RollbackPlan,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func validReferenceKind(kind ReferenceKind) bool {
	switch kind {
	case ReferenceKindIssue, ReferenceKindPullRequest:
		return true
	default:
		return false
	}
}
