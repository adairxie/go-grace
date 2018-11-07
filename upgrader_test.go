package grace

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strconv"
    "testing"
    "time"
)

type testUpgrader struct {
    *Upgrader
    procs chan *testProcess
}

func newTestUpgrader(opts Options) *testUpgrader {
    env, procs := testEnv()
    u, err := newUpgrader(env, opts)
    if err != nil {
        panic(err)
    }
    err = u.Ready()
    if err != nil {
        panic(err)
    }

    return &testUpgrader{
        Upgrader: u,
        procs: procs,
    }
}

func (tu *testUpgrader) upgradeAsync() <-chan error {
    ch := make(chan error, 1)
    go func() {
        ch <- tu.Upgrade()
    }()
    return ch
}

var names = []string{"zaphod", "beeblebrox"}

func TestMain(m *testing.M) {
    upg, err := New(Options{})
    if err != nil {
        panic(err)
    }

    if upg.parent == nil {
        // Execute test suite if there is no parent.
        os.Exit(m.Run())
    }

    pid, err := upg.Fds.File("pid")
    if err != nil {
        panic(err)
    }

    if pid != nil {
        buf := make([]byte, 8)
        binary.LittleEndian.PutUint64(buf, uint64(os.Getpid()))
        pid.Write(buf)
        pid.Close()
    }

    for _, name := range names {
        file, err := upg.Fds.File(name)
        if err != nil {
            panic(err)
        }
        if file == nil {
            continue
        }
        if _, err := io.WriteString(file, name); err != nil {
            panic(err)
        }
    }

    if err := upg.Ready(); err != nil {
        panic(err)
    }
}
