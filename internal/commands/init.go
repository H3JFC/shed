package commands

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented yet")

func Init(_ context.Context) error {
	return ErrNotImplemented
}
