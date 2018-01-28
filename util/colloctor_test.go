package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorCollector(t *testing.T) {
	var ec ErrorCollector
	assert.Nil(t, ec.Err("base message"), "newly created ErrorCollector did not return nil es Err()")

	for i := 0; i < 16; i++ {
		ec.Add(fmt.Errorf("error-%02d", i))
		ec.Wait()
		for j := 0; j <= i; j++ {
			assert.Equal(t, fmt.Sprintf("error-%02d", j), ec.errs[j].Error(), "after adding %d errors to ErrorCollector, error at index %d is incorrect", i+1, j)
		}
		assert.NotNil(t, ec.Err("base message"), "after adding %d errors to ErrorCollector, Err() returned nil as Err()")
	}
}
