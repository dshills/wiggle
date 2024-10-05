package nlib_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/nmock"
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	StatusSuccess   = "success"
	StatusInProcess = "in-process"
	StatusFail      = "fail"
)

// Helper function to create a basic signal
func createTestSignal(id string) node.Signal {
	return node.Signal{
		NodeID: id,
		Status: StatusInProcess,
		Meta:   []node.Meta{},
	}
}

func TestEmptyNode_SetID(t *testing.T) {
	node := &nlib.EmptyNode{}
	node.SetID("test-node")
	assert.Equal(t, "test-node", node.ID())
}

func TestEmptyNode_Connect(t *testing.T) {
	node := &nlib.EmptyNode{}
	childNode1 := &nlib.EmptyNode{}
	childNode2 := &nlib.EmptyNode{}

	node.Connect(childNode1, childNode2)

	assert.Equal(t, 2, len(node.Nodes()))
	assert.Equal(t, childNode1, node.Nodes()[0])
	assert.Equal(t, childNode2, node.Nodes()[1])
}

func TestEmptyNode_PreProcessSignal_Success(t *testing.T) {
	mockResMgr := new(nmock.MockResourceManager)
	mockResMgr.On("RateLimit", mock.Anything).Return(nil)

	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("ResourceManager").Return(mockResMgr)

	node := &nlib.EmptyNode{}
	node.SetStateManager(mockStateMgr)

	signal := createTestSignal("test-node")
	preProcessedSignal, err := node.PreProcessSignal(signal)

	assert.NoError(t, err)
	assert.Equal(t, "in-process", preProcessedSignal.Status)
}

func TestEmptyNode_PreProcessSignal_ExceedsRateLimit(t *testing.T) {
	mockResMgr := new(nmock.MockResourceManager)
	mockResMgr.On("RateLimit", mock.Anything).Return(errors.New("rate limit exceeded"))

	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("ResourceManager").Return(mockResMgr)

	node := &nlib.EmptyNode{}
	node.SetStateManager(mockStateMgr)

	signal := createTestSignal("test-node")
	_, err := node.PreProcessSignal(signal)

	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceeded rate limit, could not recover", err.Error())
	}
}

func TestEmptyNode_PostProcessSignal_Success(t *testing.T) {
	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("UpdateState", mock.Anything).Return()

	node := &nlib.EmptyNode{}
	node.SetStateManager(mockStateMgr)

	signal := createTestSignal("test-node")
	postProcessedSignal, err := node.PostProcessSignal(signal)

	assert.NoError(t, err)
	assert.Equal(t, StatusInProcess, postProcessedSignal.Status)
	mockStateMgr.AssertCalled(t, "UpdateState", signal)
}

func TestEmptyNode_SendToConnected_Success(t *testing.T) {
	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("Log", mock.Anything).Return()

	childNode := &nlib.EmptyNode{}
	childNode.MakeInputCh()
	childNode.SetStateManager(mockStateMgr)

	n := &nlib.EmptyNode{}
	n.MakeInputCh()
	n.SetStateManager(mockStateMgr)
	n.Connect(childNode)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	signal := createTestSignal("parent-node")
	go func() {
		<-childNode.InputCh()
	}()

	err := n.SendToConnected(ctx, signal)
	assert.NoError(t, err)
}

func TestEmptyNode_SendToConnected_ContextTimeout(t *testing.T) {
	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("Log", mock.Anything).Return()

	childNode := &nlib.EmptyNode{}
	childNode.MakeInputCh()
	childNode.SetStateManager(mockStateMgr)

	n := &nlib.EmptyNode{}
	n.MakeInputCh()
	n.SetStateManager(mockStateMgr)
	n.Connect(childNode)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	signal := createTestSignal("parent-node")
	err := n.SendToConnected(ctx, signal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context timeout or cancellation")
}

func TestEmptyNode_RunBeforeHook(t *testing.T) {
	mockHooks := new(nmock.MockHooks)
	signal := createTestSignal("test-node")
	mockHooks.On("BeforeAction", signal).Return(signal, nil)

	n := &nlib.EmptyNode{}
	n.SetOptions(node.Options{Hooks: mockHooks})

	beforeSignal, err := n.RunBeforeHook(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, beforeSignal)
	mockHooks.AssertCalled(t, "BeforeAction", signal)
}

func TestEmptyNode_RunAfterHook(t *testing.T) {
	mockHooks := new(nmock.MockHooks)
	signal := createTestSignal("test-node")
	mockHooks.On("AfterAction", signal).Return(signal, nil)

	n := &nlib.EmptyNode{}
	n.SetOptions(node.Options{Hooks: mockHooks})

	afterSignal, err := n.RunAfterHook(signal)
	assert.NoError(t, err)
	assert.Equal(t, signal, afterSignal)
	mockHooks.AssertCalled(t, "AfterAction", signal)
}

func TestEmptyNode_Fail(t *testing.T) {
	mockStateMgr := new(nmock.MockStateManager)
	mockStateMgr.On("UpdateState", mock.Anything).Return()
	mockStateMgr.On("Complete").Return()
	mockStateMgr.On("Log", mock.Anything).Return()

	n := &nlib.EmptyNode{}
	n.MakeInputCh()
	n.SetStateManager(mockStateMgr)

	signal := createTestSignal("test-node")
	err := fmt.Errorf("test error")
	n.Fail(signal, err)

	signal.Err = err.Error()
	signal.Status = StatusFail

	assert.Equal(t, StatusFail, signal.Status)
	mockStateMgr.AssertCalled(t, "UpdateState", signal)
	mockStateMgr.AssertCalled(t, "Complete")
}
