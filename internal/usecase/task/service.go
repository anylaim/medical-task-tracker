package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
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

    
    if input.RecurrenceType == nil {
        return s.createOne(ctx, normalized, nil, nil)
    }

    var lastCreated *taskdomain.Task
    
    startDate := time.Now().UTC()
    if input.DueDate != nil {
        startDate = *input.DueDate
    }


    if *input.RecurrenceType == string(taskdomain.RecurrenceDaily) {
        interval := 1
        if input.RecurrenceValue != nil {
            interval = *input.RecurrenceValue
        }

        for i := 0; i < 10; i++ {
            currentDate := startDate.AddDate(0, 0, i*interval)
            
            created, err := s.createOne(ctx, normalized, &currentDate, nil)
            if err != nil {
                return nil, err
            }
            lastCreated = created
        }
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

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
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
