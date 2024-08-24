package ipc

import (
	"sync"
	"testing"
	"time"
)

func TestCreateRuntimePipe(t *testing.T) {
	testStateDir := "/tmp"
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(2)

	t.Log("creating runtime pipe")
	err := CreateRuntimePipe(testStateDir)
	if err != nil {
		t.Log("err:", err)
	}
	t.Log("created runtime pipe")

	go func() {
		defer wg.Done()
		t.Log("creating runtime pipe reader")
		waitForReady, err := NewRuntimePipeReader(testStateDir, func() {
			t.Log("> write-id-mapping called", time.Since(start))
		})
		checkErr(t, err)
		t.Log("created runtime pipe reader")

		t.Log("> wait for ready", time.Since(start))
		<-waitForReady
		t.Log("> received ready", time.Since(start))
	}()

	go func() {
		defer wg.Done()
		pipe, err := NewRuntimePipeWriter(testStateDir)
		checkErr(t, err)

		t.Log("< sending id mapping", time.Since(start))
		err = pipe.SendIdMapping()
		checkErr(t, err)
		t.Log("< send id mapping", time.Since(start))

		time.Sleep(10 * time.Millisecond)

		t.Log("< sending ready", time.Since(start))
		err = pipe.SendReady()
		checkErr(t, err)
		t.Log("< send ready", time.Since(start))
	}()

	wg.Wait()
}
