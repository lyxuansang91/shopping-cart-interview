package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// WorkflowReferenceType represents the type of entity a workflow is associated with
type WorkflowReferenceType string

const (
	WorkflowReferenceTypeSubscription WorkflowReferenceType = "subscription"
	WorkflowReferenceTypeOrder        WorkflowReferenceType = "order"
	WorkflowReferenceTypeInvoice      WorkflowReferenceType = "invoice"
)

// WorkflowStatus represents the current status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusTimedOut  WorkflowStatus = "timed_out"
)

type WorkflowExecution struct {
	ID            uuid.ID
	WorkflowID    string
	WorkflowType  string
	ReferenceType WorkflowReferenceType
	ReferenceID   uuid.ID
	Status        WorkflowStatus
	StartedAt     time.Time
	EndedAt       *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
