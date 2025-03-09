package usecases

import (
	"context"
	"go-clean-arch/internal/entities"
	"go-clean-arch/internal/repositories"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context, task *entities.Task) (string, error)
	GetTasks(ctx context.Context) ([]entities.Task, error)
	UpdateTask(ctx context.Context, id string, task *entities.Task) error
	DeleteTask(ctx context.Context, id string) error
}

type taskUseCase struct {
	repo repositories.TaskRepository
}

func NewTaskUseCase(repo repositories.TaskRepository) TaskUseCase {
	return &taskUseCase{repo: repo}
}

func (uc *taskUseCase) CreateTask(ctx context.Context, task *entities.Task) (string, error) {
	return uc.repo.Create(ctx, task)
}

func (uc *taskUseCase) GetTasks(ctx context.Context) ([]entities.Task, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *taskUseCase) UpdateTask(ctx context.Context, id string, task *entities.Task) error {
	return uc.repo.Update(ctx, id, task)
}

func (uc *taskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
