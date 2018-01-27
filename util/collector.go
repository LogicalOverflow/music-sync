package util

import "sync"

type ErrorCollector struct {
	errs      []error
	errsMutex sync.RWMutex
	wg        sync.WaitGroup
}

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

func (ec *ErrorCollector) Wait() {
	ec.wg.Wait()
}

func (ec *ErrorCollector) Err(baseMessage string) error {
	ec.Wait()
	ec.errsMutex.RLock()
	defer ec.errsMutex.RUnlock()
	if len(ec.errs) == 0 {
		return nil
	}
	return NewMultiError(baseMessage, append([]error{}, ec.errs...))
}
