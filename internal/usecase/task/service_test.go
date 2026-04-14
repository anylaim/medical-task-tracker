package task

import (
	"context"
	"testing"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type mockRepo struct {
	tasks []*taskdomain.Task
}

func (m *mockRepo) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	task.ID = int64(len(m.tasks) + 1)
	m.tasks = append(m.tasks, task)
	return task, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	return nil, nil
}

func (m *mockRepo) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	return nil, nil
}

func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockRepo) List(ctx context.Context) ([]taskdomain.Task, error) {
	return nil, nil
}

func TestCreate_Daily(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	recurrence := "daily"
	value := 2

	_, err := service.Create(context.Background(), CreateInput{
		Title:           "test",
		Status:          taskdomain.StatusNew,
		RecurrenceType:  &recurrence,
		RecurrenceValue: &value,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.tasks) != dailyHorizonDays {
		t.Fatalf("expected %d tasks, got %d", dailyHorizonDays, len(repo.tasks))
	}
}

func TestCreate_SpecificDates(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	recurrence := "specific"

	dates := []time.Time{
		time.Now(),
		time.Now().AddDate(0, 0, 1),
	}

	_, err := service.Create(context.Background(), CreateInput{
		Title:          "test",
		Status:         taskdomain.StatusNew,
		RecurrenceType: &recurrence,
		SpecificDates:  dates,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.tasks) != len(dates) {
		t.Fatalf("expected %d tasks, got %d", len(dates), len(repo.tasks))
	}
}

func TestCreate_ParityEven(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	recurrence := "parity"
	parity := "even"

	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	service.now = func() time.Time { return start }

	_, err := service.Create(context.Background(), CreateInput{
		Title:          "test",
		Status:         taskdomain.StatusNew,
		RecurrenceType: &recurrence,
		ParityType:     &parity,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, task := range repo.tasks {
		if task.DueDate.Day()%2 != 0 {
			t.Fatalf("found odd day in even parity task")
		}
	}
}

func TestCreate_InvalidDaily(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	recurrence := "daily"
	value := 0

	_, err := service.Create(context.Background(), CreateInput{
		Title:           "test",
		Status:          taskdomain.StatusNew,
		RecurrenceType:  &recurrence,
		RecurrenceValue: &value,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreate_DeduplicateSpecificDates(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	recurrence := "specific"

	date := time.Now()

	dates := []time.Time{
		date,
		date,
	}

	_, err := service.Create(context.Background(), CreateInput{
		Title:          "test",
		Status:         taskdomain.StatusNew,
		RecurrenceType: &recurrence,
		SpecificDates:  dates,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.tasks) != 1 {
		t.Fatalf("expected 1 task after deduplication, got %d", len(repo.tasks))
	}
}
