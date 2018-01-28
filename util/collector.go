package util

import "sync"

// ErrorCollector is a thread save utility class to collect multiple errors into one slice
type ErrorCollector struct {
	errs      []error
	errsMutex sync.RWMutex
	wg        sync.WaitGroup
}

// Add adds an error to the ErrorCollector
func (ec *ErrorCollector) Add(err error) {
	ec.wg.Add(1)
	go func() {
		defer ec.wg.Done()
		ec.errsMutex.Lock()
		defer ec.errsMutex.Unlock()

		if ec.errs == nil {
			ec.errs = make([]error, 0, 1)
		}
		ec.errs = append(ec.errs, err)
	}()
}

// Wait waits for all pending error insertions to complete
func (ec *ErrorCollector) Wait() {
	ec.wg.Wait()
}

// Err returns an MultiError containing all errors added using Add, with baseMessage as base message
// If no errors were added, it returns nil
func (ec *ErrorCollector) Err(baseMessage string) error {
	ec.Wait()
	ec.errsMutex.RLock()
	defer ec.errsMutex.RUnlock()
	if len(ec.errs) == 0 {
		return nil
	}
	return NewMultiError(baseMessage, append([]error{}, ec.errs...))
}
