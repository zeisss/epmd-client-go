package epmd

import "fmt"
import "net"
import "testing"
import "time"

func TestRegisterAndNames(t *testing.T) {
	SetPort("9910")

	conns := make([]net.Conn, 0, 10)
	for i := 0; i < 10; i++ {
		c, err := Register(
			uint16(9000+i), NODE_TYPE_ERLANG,
			5, 5,
			fmt.Sprintf("aaaa%vbbbb", i),
			"")
		if err != nil {
			t.Fatal(err)
		}
		conns = append(conns, c)
	}

	// Give EPMD some time
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 1000; i++ {
		names, err := Names("localhost")
		if err != nil {
			t.Fatal(err)
		}
		if len(names) != len(conns) {
			t.Fatalf("%v: Only received %v names: %v", i, len(names), names)
		}
	}
}
