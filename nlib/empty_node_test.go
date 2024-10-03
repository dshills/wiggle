package nlib_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/nmock"
	"github.com/dshills/wiggle/node"
)

func TestEmptyNode_Init(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := &nmock.MockStateManager{}
	n := &nlib.EmptyNode{}

	n.Init(logger, stateMgr, "test-node")

	if n.ID() != "test-node" {
		t.Errorf("expected n ID to be 'test-node', got '%s'", n.ID())
	}

	if n.Logger() != logger {
		t.Errorf("expected logger to be set")
	}

	if n.StateManager() != stateMgr {
		t.Errorf("expected state manager to be set")
	}
}

func TestEmptyNode_Close(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := nmock.NewMockStateManager()
	n := &nlib.EmptyNode{}

	n.Init(logger, stateMgr, "test-node")

	n.Close()

	select {
	case <-n.DoneCh():
		// Pass
	default:
		t.Errorf("expected done channel to be closed")
	}

	select {
	case _, open := <-n.InputCh():
		if open {
			t.Errorf("expected input channel to be closed")
		}
	default:
		t.Errorf("expected input channel to be closed")
	}
}

func TestEmptyNode_SendToConnected_Success(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := nmock.NewMockStateManager()
	n := &nlib.EmptyNode{}
	n.Init(logger, stateMgr, "test-node")

	mockNode := &nlib.EmptyNode{}
	mockNode.Init(logger, stateMgr, "connected-node")
	n.Connect(mockNode)

	ctx := context.Background()
	signal := node.Signal{NodeID: "test-node"}

	err := n.SendToConnected(ctx, signal)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	select {
	case sig := <-mockNode.InputCh():
		if sig.NodeID != "connected-node" {
			t.Errorf("expected NodeID to be 'connected-node', got '%s'", sig.NodeID)
		}
	default:
		t.Errorf("expected signal to be sent")
	}
}

func TestEmptyNode_SendToConnected_Timeout(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := nmock.NewMockStateManager()
	n := &nlib.EmptyNode{}
	n.Init(logger, stateMgr, "test-node")

	mockNode := &nlib.EmptyNode{}
	mockNode.Init(logger, stateMgr, "connected-node")
	n.Connect(mockNode)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	signal := node.Signal{NodeID: "test-node"}
	err := n.SendToConnected(ctx, signal)
	if err == nil {
		t.Errorf("expected timeout error, got nil")
		return
	}

	if err.Error() != "context timeout or cancellation while sending signal to n connected-node: context deadline exceeded" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEmptyNode_Logger(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := &nmock.MockStateManager{}
	n := &nlib.EmptyNode{}
	n.Init(logger, stateMgr, "test-node")

	n.LogInfo("This is a test log")

	logEntries := logger.Entries()
	if len(logEntries) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(logEntries))
	}
}

func TestEmptyNode_ValidateSignal(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := &nmock.MockStateManager{}
	n := &nlib.EmptyNode{}
	n.Init(logger, stateMgr, "test-node")

	// Test valid signal
	validSignal := node.Signal{NodeID: "valid-node"}
	if err := n.ValidateSignal(validSignal); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test invalid signal
	invalidSignal := node.Signal{}
	if err := n.ValidateSignal(invalidSignal); err == nil {
		t.Errorf("expected error for invalid signal, got nil")
	}

	expectedError := "invalid signal missing ID"
	if err := n.ValidateSignal(invalidSignal); err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEmptyNode_PreProcessSignal(t *testing.T) {
	logger := &nmock.MockLogger{}
	stateMgr := &nmock.MockStateManager{}
	n := &nlib.EmptyNode{}
	n.Init(logger, stateMgr, "test-node")

	// Test valid signal pre-processing
	sig := node.Signal{NodeID: "test-node"}
	processedSig, err := n.PreProcessSignal(sig)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if processedSig.NodeID != "test-node" {
		t.Errorf("expected signal NodeID to remain 'test-node', got '%s'", processedSig.NodeID)
	}

	// Test invalid signal
	invalidSig := node.Signal{}
	_, err = n.PreProcessSignal(invalidSig)
	if err == nil {
		t.Errorf("expected error for invalid signal, got nil")
	}

	expectedError := "invalid signal missing ID"
	if err != nil && err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEmptyNode_ErrorGuidance(t *testing.T) {
	n := &nlib.EmptyNode{}

	// Test default behavior without setting error guidance
	if n.ErrorAction(errors.New("test")) != node.ErrGuideFail {
		t.Errorf("expected default error action to be ErrGuideFail")
	}

	// Mock error guidance and set it on the n
	mockErrGuide := &nmock.MockErrorGuidance{RetryErr: nmock.ErrMockRetry}
	n.SetErrorGuidance(mockErrGuide)

	// Test error guidance action
	if n.ErrorAction(nmock.ErrMockRetry) != node.ErrGuideRetry {
		t.Errorf("expected error guidance to return ErrGuideRetry")
	}
}
