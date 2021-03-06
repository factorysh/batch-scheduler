package scheduler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/factorysh/density/store"
	"github.com/factorysh/density/task"
	"github.com/google/uuid"
)

// JSONStore stores task.Task
type JSONStore struct {
	store store.Store
}

// Get a Task
func (j *JSONStore) Get(id uuid.UUID) (*task.Task, error) {
	v, err := j.store.Get([]byte(id.String()))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return parseTask(v)
}

func parseTask(v []byte) (*task.Task, error) {
	var t task.Task
	err := json.Unmarshal(v, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Put task.Task
func (j *JSONStore) Put(t *task.Task) error {
	if t.Id == uuid.Nil {
		return errors.New("Task wihtout id")
	}
	t.Mtime = time.Now()
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return j.store.Put([]byte(t.Id.String()), value)
}

// Delete a task
func (j *JSONStore) Delete(id uuid.UUID) error {
	return j.store.Delete([]byte(id.String()))
}

// Length of the store
func (j *JSONStore) Length() int {
	return j.store.Length()
}

// ForEach loops over kv
func (j *JSONStore) ForEach(fn func(t *task.Task) error) error {
	return j.store.ForEach(func(k, v []byte) error {
		// k is the UUID, serialized in v
		t, err := parseTask(v)
		if err != nil {
			return err
		}
		return fn(t)
	})
}

// DeleteWithClause batch delete
func (j *JSONStore) DeleteWithClause(fn func(t *task.Task) bool) error {
	return j.store.DeleteWithClause(func(k, v []byte) bool {
		t, err := parseTask(v)
		if err != nil {
			panic(err)
		}
		return fn(t)
	})
}
