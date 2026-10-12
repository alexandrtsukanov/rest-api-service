package database

import (
	"database/sql"
	"fmt"
	"restapi-tasks/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

const DB_QUERY_ERROR string = "database query error: "
const NOT_FOUND_ERROR string = "task with id %d not found"

type TaskService struct {
	db *sqlx.DB
}

func NewTaskService(db *sqlx.DB) *TaskService {
	return &TaskService{db}
}

func (ts *TaskService) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
		SELECT id, title, descriptionn, completed, created_at, updated_at
		FROM tasks
		ordered by created_at desc;
	`

	err := ts.db.Select(&tasks, query)

	if err != nil {
		return nil, fmt.Errorf(DB_QUERY_ERROR+"%w", err)
	}

	return tasks, nil
}

func (ts *TaskService) GetTaskById(id int) (*models.Task, error) {
	var task models.Task

	query := `
		SELECT id, title, descriptionn, completed, created_at, updated_at
		FROM tasks
		where id = $1;
	`

	err := ts.db.Get(&task, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf(NOT_FOUND_ERROR, id)
	}

	if err != nil {
		return nil, fmt.Errorf(DB_QUERY_ERROR+"%w", err)
	}

	return &task, nil
}

func (ts *TaskService) CreateTask(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
		INSERT INTO tasks
		(title, descriptionn, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		returning id, title, descriptionn, completed, created_at, updated_at;
	`

	now := time.Now()

	err := ts.db.QueryRowx(
		query,
		input.Title,
		input.Description,
		input.Completed,
		now,
		now,
	).StructScan(&task)

	if err != nil {
		return nil, fmt.Errorf(DB_QUERY_ERROR+"%w", err)
	}

	return &task, nil
}

func (ts *TaskService) UpdateTask(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := ts.GetTaskById(id)

	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, descriptionn = $2, completed = $3, updated_at = $4
		where id = %5
		returning id, title, descriptionn, completed, created_at, updated_at;
	`

	var updatedTask models.Task

	dbUpdateErr := ts.db.QueryRowx(
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.UpdatedAt,
		id,
	).StructScan(&updatedTask)

	if dbUpdateErr != nil {
		return nil, fmt.Errorf(DB_QUERY_ERROR+"%w", dbUpdateErr)
	}

	return &updatedTask, nil
}

func (ts *TaskService) DeleteTask(id int) error {
	query := "DELETE FROM tasks where id = $1"

	result, err := ts.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf(DB_QUERY_ERROR+"%w", err)
	}

	rowsNum, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf(DB_QUERY_ERROR+"%w", err)
	}

	if rowsNum == 0 {
		return fmt.Errorf(NOT_FOUND_ERROR, id)
	}

	return nil
}
