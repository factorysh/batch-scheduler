package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/factorysh/density/owner"
	"github.com/factorysh/density/scheduler"
	"github.com/factorysh/density/task"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleGetTask will retreive a task, convert it into a resp and return the data
func HandleGetTask(schd *scheduler.Scheduler, u *owner.Owner, _ http.ResponseWriter, r *http.Request) (interface{}, error) {
	fmt.Println(u)

	vars := mux.Vars(r)
	rawID, ok := vars[task.UUID]
	if !ok {
		return nil, errors.New("No uuid in request")
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, err
	}

	t, err := schd.GetTask(id)
	if err != nil {
		return nil, err
	}

	return t.ToTaskResp(), nil
}
