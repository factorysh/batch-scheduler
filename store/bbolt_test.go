package store

import (
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	var wg sync.WaitGroup

	// should success since the store is not open
	store, err := NewBoltStore("../tests/store.bolt")
	assert.NoError(t, err)

	wg.Add(1)

	go func() {
		// should fail since the store is already open
		_, err = NewBoltStore("../tests/store.bolt")
		assert.Error(t, err)

		wg.Done()
	}()

	wg.Wait()

	store.Db.Close()
}

func TestPut(t *testing.T) {
	store, err := NewBoltStore("../tests/store.bolt")
	assert.NoError(t, err)

	err = store.Put([]byte("put_test"), []byte("put_value"))
	assert.NoError(t, err)

	store.Db.Close()
}

func TestGet(t *testing.T) {
	store, err := NewBoltStore("../tests/store.bolt")
	assert.NoError(t, err)

	err = store.Put([]byte("answer"), []byte("42"))
	assert.NoError(t, err)

	v, err := store.Get([]byte("answer"))
	assert.NoError(t, err)
	assert.Equal(t, "42", string(v))

	store.Db.Close()
}

func TestDelete(t *testing.T) {
	store, err := NewBoltStore("../tests/store.bolt")
	assert.NoError(t, err)

	err = store.Put([]byte("test_delete"), []byte("delete_value"))
	assert.NoError(t, err)

	err = store.Delete([]byte("test_delete"))
	assert.NoError(t, err)

	v, err := store.Get([]byte("test_delete"))
	assert.Equal(t, "", string(v))

	store.Db.Close()
}

func TestBBoltForEach(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "bolt-")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	store, err := NewBoltStore(f.Name())
	assert.NoError(t, err)
	for _, name := range []string{"pim", "pam", "poum"} {
		err = store.Put([]byte(name), []byte{})
		assert.NoError(t, err)
	}
	names := make([]string, 0)
	err = store.ForEach(func(k, v []byte) error {
		names = append(names, string(k))
		return nil
	})
	assert.NoError(t, err)
	sort.Strings(names)
	assert.Equal(t, []string{"pam", "pim", "poum"}, names)
}
