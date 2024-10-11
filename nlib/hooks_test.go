package nlib_test

import (
	"errors"
	"testing"

	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/assert"
)

func TestSimpleNodeHooks_BeforeAction(t *testing.T) {
	// Mock signal and hook functions
	signal := node.Signal{}

	// Test case: BeforeAction should return the signal unmodified if no before hook is set
	hooks := nlib.NewSimpleNodeHooks(nil, nil)
	outSignal, err := hooks.BeforeAction(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, outSignal)

	// Test case: BeforeAction should call the before hook if provided
	mockBefore := func(s node.Signal) (node.Signal, error) {
		return s, nil
	}
	hooks = nlib.NewSimpleNodeHooks(mockBefore, nil)
	outSignal, err = hooks.BeforeAction(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, outSignal)

	// Test case: BeforeAction should return an error if the before hook fails
	mockBeforeWithError := func(s node.Signal) (node.Signal, error) {
		return s, errors.New("before hook failed")
	}
	hooks = nlib.NewSimpleNodeHooks(mockBeforeWithError, nil)
	_, err = hooks.BeforeAction(signal)
	assert.Error(t, err)
	assert.EqualError(t, err, "before hook failed")
}

func TestSimpleNodeHooks_AfterAction(t *testing.T) {
	// Mock signal and hook functions
	signal := node.Signal{}

	// Test case: AfterAction should return the signal unmodified if no after hook is set
	hooks := nlib.NewSimpleNodeHooks(nil, nil)
	outSignal, err := hooks.AfterAction(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, outSignal)

	// Test case: AfterAction should call the after hook if provided
	mockAfter := func(s node.Signal) (node.Signal, error) {
		return s, nil
	}
	hooks = nlib.NewSimpleNodeHooks(nil, mockAfter)
	outSignal, err = hooks.AfterAction(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, outSignal)

	// Test case: AfterAction should return an error if the after hook fails
	mockAfterWithError := func(s node.Signal) (node.Signal, error) {
		return s, errors.New("after hook failed")
	}
	hooks = nlib.NewSimpleNodeHooks(nil, mockAfterWithError)
	_, err = hooks.AfterAction(signal)
	assert.Error(t, err)
	assert.EqualError(t, err, "after hook failed")
}
