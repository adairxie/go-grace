package grace

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestFdsListen(t *testing.T) {
	addrs := [][2]string{
		{"unix", ""},
		{"tcp", "localhost:0"},
	}

	fds := newFds(nil)

	for _, addr := range addrs {
		ln, err := fds.Listen(addr[0], addr[1])
		if err != nil {
			t.Fatal(err)
		}
		if ln == nil {
			t.Fatal("Missing listener", addr)
		}
		ln.Close()
	}
}

func TestFdsListener(t *testing.T) {
	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}

	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer tcp.Close()

	temp, err := ioutil.TempDir("", "tableflip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	socketPath := filepath.Join(temp, "socket")
	unix, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close()

	parent := newFds(nil)
	if err := parent.AddListener(addr.Network(), addr.String(), tcp); err != nil {
		t.Fatal("Can't add listener:", err)
	}
	tcp.Close()

	if err := parent.AddListener("unix", socketPath, unix.(Listener)); err != nil {
		t.Fatal("Can't add listener:", err)
	}
	unix.Close()

	if _, err := os.Stat(socketPath); err != nil {
		t.Error("Unix.Close() unlinked socketPath:", err)
	}

	child := newFds(parent.copy())
	ln, err := child.Listener(addr.Network(), addr.String())
	if err != nil {
		t.Fatal("Can't get listener:", err)
	}

	if ln == nil {
		t.Fatal("Missing listener")
	}
	ln.Close()

	child.closeInherited()
	if _, err := os.Stat(socketPath); err == nil {
		t.Error("closeInherited() did not unlink socketPath")
	}
}
