package publicip

import "context"

// Get some Haskell in to you.
func race[R any](ctx context.Context, fs ...func(context.Context) (R, error)) (r R, errs []error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rs := make(chan R)
	errChan := make(chan error, len(fs))
	for _, _f := range fs {
		f := _f
		go func() {
			r, err := f(ctx)
			if err == nil {
				select {
				case rs <- r:
				case <-ctx.Done():
				}
			} else {
				errChan <- err
			}
		}()
	}
	errs = make([]error, 0, len(fs))
	for range fs {
		select {
		case <-ctx.Done():
			errs = append(errs, ctx.Err())
			return
		case r = <-rs:
			errs = nil
			return
		case err := <-errChan:
			errs = append(errs, err)
		}
	}
	return
}
