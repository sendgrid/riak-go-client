// +build integration

package riak

import (
	"net"
	"testing"
	"time"
)

func TestConnectionManagerDoesNotExpirePastMinConnections(t *testing.T) {
	minConnections := uint16(10)

	o := &testListenerOpts{
		test: t,
		host: "127.0.0.1",
		port: 13340,
	}
	tl := newTestListener(o)
	tl.start()
	defer tl.stop()

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:13340")

	cmopts := &connectionManagerOptions{
		addr:                   addr,
		minConnections:         minConnections,
		maxConnections:         20,
		idleExpirationInterval: time.Millisecond * 500,
		idleTimeout:            time.Millisecond * 10,
	}

	cm, err := newConnectionManager(cmopts)
	if err != nil {
		t.Fatal(err)
	}
	err = cm.start()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 250)
		if actual, expected := cm.connectionCounter.count(), minConnections; actual != expected {
			t.Errorf("got: %v, expected: %v", actual, expected)
		}
	}

	err = cm.stop()
	if err != nil {
		t.Error(err)
	}
}
