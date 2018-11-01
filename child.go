package grace

import (
    "encoding/gob"
    "fmt"
    "os"

    "github.com/pkg/errors"
)

type child struct {
    *env
    proc        process
}
