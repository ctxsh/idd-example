package deployment

import (
	"testing"
	"time"
)

func TestNewRequest(t *testing.T) {
	now := time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC)
	input := NewRequestInput{
		Service:     "payments",
		Environment: "staging",
		Version:     "v1.2.3",
		RequestedBy: "alice",
		References: []Reference{
			{Kind: ReferenceKindIssue, Value: "#2"},
			{Kind: ReferenceKindPullRequest, Value: "#8"},
		},
		Risk:         "low risk deployment",
		RollbackPlan: "redeploy v1.2.2",
	}

	request, err := NewRequest(input, now)
	if err != nil {
		t.Fatalf("NewRequest() error = %v, want nil", err)
	}

	if request.Service != input.Service {
		t.Errorf("Service = %q, want %q", request.Service, input.Service)
	}
	if request.Environment != input.Environment {
		t.Errorf("Environment = %q, want %q", request.Environment, input.Environment)
	}
	if request.Version != input.Version {
		t.Errorf("Version = %q, want %q", request.Version, input.Version)
	}
	if request.RequestedBy != input.RequestedBy {
		t.Errorf("RequestedBy = %q, want %q", request.RequestedBy, input.RequestedBy)
	}
	if request.Status != StatusPendingApproval {
		t.Errorf("Status = %q, want %q", request.Status, StatusPendingApproval)
	}
	if request.ApprovedBy != "" {
		t.Errorf("ApprovedBy = %q, want empty", request.ApprovedBy)
	}
	if len(request.References) != len(input.References) {
		t.Fatalf("len(References) = %d, want %d", len(request.References), len(input.References))
	}
	for i := range input.References {
		if request.References[i] != input.References[i] {
			t.Errorf("References[%d] = %#v, want %#v", i, request.References[i], input.References[i])
		}
	}
	if request.Risk != input.Risk {
		t.Errorf("Risk = %q, want %q", request.Risk, input.Risk)
	}
	if request.RollbackPlan != input.RollbackPlan {
		t.Errorf("RollbackPlan = %q, want %q", request.RollbackPlan, input.RollbackPlan)
	}
	if !request.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", request.CreatedAt, now)
	}
	if !request.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", request.UpdatedAt, now)
	}
}

func TestNewRequestStartsPendingApproval(t *testing.T) {
	now := time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC)
	input := validNewRequestInput()

	request, err := NewRequest(input, now)
	if err != nil {
		t.Fatalf("NewRequest() error = %v, want nil", err)
	}

	if request.Status != StatusPendingApproval {
		t.Errorf("Status = %q, want %q", request.Status, StatusPendingApproval)
	}
}

func TestNewRequestRequiresFields(t *testing.T) {
	tests := []struct {
		name  string
		input NewRequestInput
	}{
		{
			name: "service",
			input: NewRequestInput{
				Environment: "staging",
				Version:     "v1.2.3",
				RequestedBy: "alice",
			},
		},
		{
			name: "blank service",
			input: NewRequestInput{
				Service:     " ",
				Environment: "staging",
				Version:     "v1.2.3",
				RequestedBy: "alice",
			},
		},
		{
			name: "environment",
			input: NewRequestInput{
				Service:     "payments",
				Version:     "v1.2.3",
				RequestedBy: "alice",
			},
		},
		{
			name: "blank environment",
			input: NewRequestInput{
				Service:     "payments",
				Environment: " ",
				Version:     "v1.2.3",
				RequestedBy: "alice",
			},
		},
		{
			name: "version",
			input: NewRequestInput{
				Service:     "payments",
				Environment: "staging",
				RequestedBy: "alice",
			},
		},
		{
			name: "blank version",
			input: NewRequestInput{
				Service:     "payments",
				Environment: "staging",
				Version:     " ",
				RequestedBy: "alice",
			},
		},
		{
			name: "requested_by",
			input: NewRequestInput{
				Service:     "payments",
				Environment: "staging",
				Version:     "v1.2.3",
			},
		},
		{
			name: "blank requested_by",
			input: NewRequestInput{
				Service:     "payments",
				Environment: "staging",
				Version:     "v1.2.3",
				RequestedBy: " ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRequest(tt.input, time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC))
			if err == nil {
				t.Fatal("NewRequest() error = nil, want non-nil")
			}
		})
	}
}

func TestNewRequestCopiesReferences(t *testing.T) {
	input := validNewRequestInput()
	input.References = []Reference{{Kind: ReferenceKindIssue, Value: "#2"}}

	request, err := NewRequest(input, time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("NewRequest() error = %v, want nil", err)
	}

	input.References[0] = Reference{Kind: ReferenceKindPullRequest, Value: "#8"}

	if request.References[0].Kind != ReferenceKindIssue {
		t.Errorf("References[0].Kind = %q, want %q", request.References[0].Kind, ReferenceKindIssue)
	}
	if request.References[0].Value != "#2" {
		t.Errorf("References[0].Value = %q, want %q", request.References[0].Value, "#2")
	}
}

func TestNewRequestAcceptsSupportedReferenceKinds(t *testing.T) {
	tests := []struct {
		name string
		kind ReferenceKind
	}{
		{name: "issue", kind: ReferenceKindIssue},
		{name: "pull request", kind: ReferenceKindPullRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validNewRequestInput()
			input.References = []Reference{{Kind: tt.kind, Value: "#2"}}

			_, err := NewRequest(input, time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC))
			if err != nil {
				t.Fatalf("NewRequest() error = %v, want nil", err)
			}
		})
	}
}

func TestNewRequestRejectsUnsupportedReferenceKinds(t *testing.T) {
	input := validNewRequestInput()
	input.References = []Reference{{Kind: ReferenceKind("chat"), Value: "deploy thread"}}

	_, err := NewRequest(input, time.Date(2026, 6, 4, 12, 30, 0, 0, time.UTC))
	if err == nil {
		t.Fatal("NewRequest() error = nil, want non-nil")
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{status: StatusPendingApproval, want: "pending_approval"},
		{status: StatusApproved, want: "approved"},
		{status: StatusRejected, want: "rejected"},
		{status: StatusStarted, want: "started"},
		{status: StatusSucceeded, want: "succeeded"},
		{status: StatusFailed, want: "failed"},
		{status: StatusRolledBack, want: "rolled_back"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("string(%T(%q)) = %q, want %q", tt.status, tt.status, tt.status, tt.want)
		}
	}
}

func TestReferenceKindConstants(t *testing.T) {
	tests := []struct {
		kind ReferenceKind
		want string
	}{
		{kind: ReferenceKindIssue, want: "issue"},
		{kind: ReferenceKindPullRequest, want: "pull_request"},
	}

	for _, tt := range tests {
		if string(tt.kind) != tt.want {
			t.Errorf("string(%T(%q)) = %q, want %q", tt.kind, tt.kind, tt.kind, tt.want)
		}
	}
}

func validNewRequestInput() NewRequestInput {
	return NewRequestInput{
		Service:     "payments",
		Environment: "staging",
		Version:     "v1.2.3",
		RequestedBy: "alice",
	}
}
