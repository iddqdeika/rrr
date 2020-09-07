package rrr

import "context"

type Root interface {
	// configuration and DI logic
	Register() []error
	// main logic execution
	Resolve(ctx context.Context) error
	// resources releasing, shutdown messages sending etc
	Release() error
}
