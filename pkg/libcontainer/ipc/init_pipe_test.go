package ipc

import (
	"sync"
	"testing"
	"time"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestInitPipe(t *testing.T) {
	testStateDir := "/tmp"
	//defer os.Remove(filepath.Join(testStateDir, initPipeFileName))
	var wg sync.WaitGroup
	wg.Add(1)
	err := CreateInitPipe(testStateDir)
	if err != nil {
		t.Log(err)
	}

	start := time.Now()
	go func() {
		defer wg.Done()
		r, err := NewInitPipeReader(testStateDir)
		checkErr(t, err)

		t.Log("waiting for start", time.Since(start))
		err = r.WaitForStart()
		checkErr(t, err)
		t.Log("received start", time.Since(start))
	}()

	w, err := NewInitPipeWriter(testStateDir)
	checkErr(t, err)

	time.Sleep(10 * time.Millisecond)

	t.Log("sending start", time.Since(start))
	err = w.SendStart()
	checkErr(t, err)
	t.Log("start send", time.Since(start))

	wg.Wait()
}
