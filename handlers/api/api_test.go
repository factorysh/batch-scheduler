package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/factorysh/batch-scheduler/runner"
	"github.com/factorysh/batch-scheduler/scheduler"
	"github.com/factorysh/batch-scheduler/store"
	"github.com/factorysh/batch-scheduler/task"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "scheduler-")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	s := scheduler.New(scheduler.NewResources(4, 16*1024), runner.New(dir), store.NewMemoryStore())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Start(ctx)
	key := "plop"
	mux := MuxAPI(s, key)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	c, err := newClient(ts.URL, key)
	assert.NoError(t, err)

	var r []interface{}
	res, err := c.Do("GET", "/api/schedules", nil, nil, &r)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Len(t, r, 0)

	h := make(http.Header)
	h.Set("content-type", "application/json")
	b := bytes.NewReader([]byte(`{
		"cpu": 2,
		"ram": 128,
		"max_execution_time": "120s",
		"action": {
			"compose": {
				"version": "3",
				"services": {
					"hello": {
						"image":"busybox:latest",
						"command": "echo World"
					}
				}
			}
		}
	}`))
	var ta task.Task
	res, err = c.Do("POST", "/api/schedules", h, b, &ta)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)
	assert.Len(t, r, 0)

	// FIXME test schedule creation with a file upload
}

type testClient struct {
	root          string
	client        *http.Client
	authorization string
}

func newClient(root, key string) (*testClient, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"owner": "bob",
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	blob, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	return &testClient{
		root:          root,
		client:        &http.Client{},
		authorization: fmt.Sprintf("Bearer %s", blob),
	}, nil
}

// Do a request.
// value is a pointer for unmarshaled JSON response
func (t *testClient) Do(method, url string, header http.Header, body io.Reader, value interface{}) (*http.Response, error) {
	r, err := http.NewRequest(method, t.root+url, body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		r.Header = header
	}
	r.Header.Set("Authorization", t.authorization)
	res, err := t.client.Do(r)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()
	ct := res.Header.Get("content-type")
	if ct != "application/json" {
		return res, fmt.Errorf("Wrong content-type : %s", ct)
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	fmt.Println("raw", string(raw))
	err = json.Unmarshal(raw, value)
	return res, err
}
