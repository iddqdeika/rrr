package rrr

import (
	"context"
	"sync"
)

type Service interface {
	Run(ctx context.Context) error
}

// runs services and returns all not-nil errs
func RunServices(ctx context.Context, svcs... Service) []error{
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch := make(chan error, len(svcs))
	wg := sync.WaitGroup{}
	wg.Add(len(svcs))
	for _, svc := range svcs {
		go func(svc Service) {
			defer wg.Done()
			defer cancel()
			err := svc.Run(ctx)
			if err != nil {
				ch <- err
			}
		}(svc)
	}
	wg.Wait()
	return collectErr(ch)
}

func collectErr(ch chan error) []error {
	var errs []error
	for {
		select {
		case err := <-ch:
			errs = append(errs, err)
		default:
			return errs
		}
	}
}