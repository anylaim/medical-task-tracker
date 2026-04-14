package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceDaily    RecurrenceType = "daily"
	RecurrenceMonthly  RecurrenceType = "monthly"
	RecurrenceParity   RecurrenceType = "parity"
	RecurrenceSpecific RecurrenceType = "specific"
	RecurrenceWeekly   RecurrenceType = "weekly"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`

	RecurrenceType  *RecurrenceType `json:"recurrence_type,omitempty"`
	RecurrenceValue *int            `json:"recurrence_value,omitempty"`

	SpecificDates []time.Time `json:"specific_dates,omitempty"`
	ParityType    *string     `json:"parity_type,omitempty"` // even / odd

	DueDate  *time.Time `json:"due_date,omitempty"`
	ParentID *int64     `json:"parent_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
