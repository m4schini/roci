package ipc

import (
	"bytes"
	"context"
	pb "roci/proto"
	"testing"
)

var testMessage = &pb.IdMapping{
	InsideId:  42,
	OutsideId: 69,
}

func TestListen(t *testing.T) {
	var buf bytes.Buffer
	var expected = testMessage

	go func() {
		err := write(&buf, expected)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("written message into buffer")
	}()

	ch := Listen(context.TODO(), &buf, func() *pb.IdMapping {
		return new(pb.IdMapping)
	})

	actual := <-ch
	t.Log("received message")

	if actual.InsideId != expected.InsideId {
		t.Error()
		t.FailNow()
	}
	if actual.OutsideId != expected.OutsideId {
		t.FailNow()
	}

	t.Log("expected:", expected.String())
	t.Log("  actual:", actual.String())
}
