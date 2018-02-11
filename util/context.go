package util

import "context"

// IsCanceled returned whether the context is canceled
func IsCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
