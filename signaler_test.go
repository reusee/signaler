package signaler

import "testing"

func TestOneshot(t *testing.T) {
	signaler := NewSignaler()
	i := 0
	signaler.OnSignal("foo", func() bool {
		i++
		return true
	})
	signaler.SignalSynced("foo")
	signaler.SignalSynced("foo")
	signaler.SignalSynced("foo")
	if i != 1 {
		t.Fail()
	}
}
