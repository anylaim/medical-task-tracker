package task

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

const (
	dailyHorizonDays  = 30
	monthlyHorizon    = 12
	parityHorizonDays = 30
	weeklyHorizon     = 10
	maxSpecificDates  = 100
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	if normalized.RecurrenceType == nil {
		return s.createOne(ctx, normalized, normalized.DueDate, nil)
	}

	var lastCreated *taskdomain.Task
	var parentID *int64

	startDate := s.now()
	if input.DueDate != nil {
		startDate = *input.DueDate
	}

	rtype := taskdomain.RecurrenceType(*normalized.RecurrenceType)

	createWithParent := func(date time.Time) error {
		task, err := s.createOne(ctx, normalized, &date, parentID)
		if err != nil {
			return err
		}

		if parentID == nil {
			parentID = &task.ID
		}

		lastCreated = task
		return nil
	}

	switch rtype {

	case taskdomain.RecurrenceDaily:
		interval := *normalized.RecurrenceValue
		for i := 0; i < dailyHorizonDays; i++ {
			date := startDate.AddDate(0, 0, i*interval)
			if err := createWithParent(date); err != nil {
				return nil, err
			}
		}

	case taskdomain.RecurrenceMonthly:
		for i := 0; i < monthlyHorizon; i++ {
			date := startDate.AddDate(0, i, 0)
			if err := createWithParent(date); err != nil {
				return nil, err
			}
		}

	case taskdomain.RecurrenceParity:
		isEven := *normalized.ParityType == "even"

		for i := 0; i < parityHorizonDays; i++ {
			date := startDate.AddDate(0, 0, i)
			day := date.Day()

			if (isEven && day%2 == 0) || (!isEven && day%2 != 0) {
				if err := createWithParent(date); err != nil {
					return nil, err
				}
			}
		}

	case taskdomain.RecurrenceSpecific:
		dates := deduplicateDates(normalized.SpecificDates)

		for _, date := range dates {
			if err := createWithParent(date); err != nil {
				return nil, err
			}
		}

	case taskdomain.RecurrenceWeekly:
		for i := 0; i < weeklyHorizon; i++ {
			date := startDate.AddDate(0, 0, i*7)
			if err := createWithParent(date); err != nil {
				return nil, err
			}
		}

	default:
		return s.createOne(ctx, normalized, &startDate, nil)
	}

	return lastCreated, nil
}

func (s *Service) createOne(ctx context.Context, input CreateInput, dueDate *time.Time, parentID *int64) (*taskdomain.Task, error) {
	model := &taskdomain.Task{
		Title:           input.Title,
		Description:     input.Description,
		Status:          input.Status,
		RecurrenceType:  (*taskdomain.RecurrenceType)(input.RecurrenceType),
		RecurrenceValue: input.RecurrenceValue,
		DueDate:         dueDate,
		ParentID:        parentID,
		ParityType:      input.ParityType,
		SpecificDates:   input.SpecificDates,
	}

	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	return s.repo.Create(ctx, model)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}

	return s.repo.Update(ctx, model)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType != nil {
		val := strings.ToLower(strings.TrimSpace(*input.RecurrenceType))
		*input.RecurrenceType = val

		rtype := taskdomain.RecurrenceType(val)

		switch rtype {
		case taskdomain.RecurrenceDaily:
			if input.RecurrenceValue == nil || *input.RecurrenceValue <= 0 {
				return CreateInput{}, fmt.Errorf("%w: daily requires positive recurrence_value", ErrInvalidInput)
			}
		case taskdomain.RecurrenceMonthly:
			if input.DueDate == nil {
				now := time.Now().UTC()
				input.DueDate = &now
			}
		case taskdomain.RecurrenceParity:
			if input.ParityType == nil {
				return CreateInput{}, fmt.Errorf("%w: parity_type required", ErrInvalidInput)
			}
			p := strings.ToLower(strings.TrimSpace(*input.ParityType))
			*input.ParityType = p
			if p != "even" && p != "odd" {
				return CreateInput{}, fmt.Errorf("%w: parity_type must be 'even' or 'odd'", ErrInvalidInput)
			}
		case taskdomain.RecurrenceSpecific:
			if len(input.SpecificDates) == 0 {
				return CreateInput{}, fmt.Errorf("%w: specific_dates cannot be empty", ErrInvalidInput)
			}
		case taskdomain.RecurrenceWeekly:
			break
		default:
			return CreateInput{}, fmt.Errorf("%w: unknown recurrence_type", ErrInvalidInput)
		}
	}
	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func deduplicateDates(dates []time.Time) []time.Time {
	seen := make(map[time.Time]struct{})
	result := make([]time.Time, 0, len(dates))

	for _, d := range dates {
		normalized := d.UTC().Truncate(24 * time.Hour)
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			result = append(result, normalized)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Before(result[j])
	})

	return result
}
