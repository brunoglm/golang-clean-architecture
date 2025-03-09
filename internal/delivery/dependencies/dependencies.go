package dependencies

import (
	"go-clean-arch/internal/infra"
	"go-clean-arch/internal/interfaces/handlers"
	"go-clean-arch/internal/repositories"
	"go-clean-arch/internal/usecases"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

func Setup() *dig.Container {
	container := dig.New()

	err := container.Provide(infra.NewMongoDatabase)
	if err != nil {
		log.Fatalf("Erro ao registrar MongoDB: %v", err)
	}

	err = container.Provide(func(db *mongo.Database) repositories.TaskRepository {
		return repositories.NewTaskRepository(db)
	})
	if err != nil {
		log.Fatalf("Erro ao registrar TaskRepository: %v", err)
	}

	err = container.Provide(usecases.NewTaskUseCase)
	if err != nil {
		log.Fatalf("Erro ao registrar TaskUseCase: %v", err)
	}

	err = container.Provide(handlers.NewTaskHandler)
	if err != nil {
		log.Fatalf("Erro ao registrar TaskHandler: %v", err)
	}

	return container
}
