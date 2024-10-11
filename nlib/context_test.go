package nlib_test

import (
	"sync"
	"testing"

	"github.com/dshills/wiggle/nlib"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleContextManager(t *testing.T) {
	cm := nlib.NewSimpleContextManager()
	assert.NotNil(t, cm, "Expected NewSimpleContextManager to return a valid instance")
}

func TestSetContext(t *testing.T) {
	cm := nlib.NewSimpleContextManager()
	carrier := &nlib.Carrier{TextData: "TestData"}

	cm.SetContext("node1", carrier)

	result, err := cm.GetContext("node1")
	assert.NoError(t, err, "Expected no error when retrieving context after setting it")
	assert.Equal(t, carrier, result, "Expected the stored context data to match the set data")
}

func TestGetContext_NotFound(t *testing.T) {
	cm := nlib.NewSimpleContextManager()

	_, err := cm.GetContext("nonexistent-node")
	assert.Error(t, err, "Expected an error when retrieving context for nonexistent node")
	assert.EqualError(t, err, "not found", "Expected error message to be 'not found'")
}

func TestRemoveContext(t *testing.T) {
	cm := nlib.NewSimpleContextManager()
	carrier := &nlib.Carrier{TextData: "TestData"}

	cm.SetContext("node1", carrier)
	cm.RemoveContext("node1")

	_, err := cm.GetContext("node1")
	assert.Error(t, err, "Expected an error when retrieving context after removal")
	assert.EqualError(t, err, "not found", "Expected error message to be 'not found' after removal")
}

func TestThreadSafety(_ *testing.T) {
	cm := nlib.NewSimpleContextManager()
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			cm.SetContext("node1", &nlib.Carrier{TextData: "data"})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = cm.GetContext("node1")
		}
	}()

	wg.Wait()

	// Just ensuring the test completes without data races or deadlocks
}
