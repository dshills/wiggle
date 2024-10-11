package nlib_test

import (
	"sync"
	"testing"

	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/assert"
)

func TestSimpleHistoryManager_AddHistory(t *testing.T) {
	hx := nlib.NewSimpleHistoryManager()

	// Create a mock signal to add
	sig1 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testData")}

	// Add the signal to history
	hx.AddHistory(sig1)

	// Check if the signal is added correctly
	history := hx.GetHistory()
	assert.Len(t, history, 1, "History should contain exactly one signal")
	assert.Equal(t, sig1, history[0], "The added signal should match the first element in history")
}

func TestSimpleHistoryManager_GetHistory(t *testing.T) {
	hx := nlib.NewSimpleHistoryManager()

	// Add multiple signals
	sig1 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testdata1")}
	sig2 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testdata2")}
	hx.AddHistory(sig1)
	hx.AddHistory(sig2)

	// Get the full history and validate
	history := hx.GetHistory()
	assert.Len(t, history, 2, "History should contain two signals")
	assert.Equal(t, sig1, history[0], "First signal should match sig1")
	assert.Equal(t, sig2, history[1], "Second signal should match sig2")
}

func TestSimpleHistoryManager_Filter(t *testing.T) {
	hx := nlib.NewSimpleHistoryManager()

	// Add multiple signals with different NodeIDs
	sig1 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testdata1")}
	sig2 := node.Signal{NodeID: "node2", Task: nlib.NewTextCarrier("testdata2")}
	sig3 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testdata3")}
	hx.AddHistory(sig1)
	hx.AddHistory(sig2)
	hx.AddHistory(sig3)

	// Filter by NodeID "node1" and validate
	filteredHistory := hx.Filter("node1")
	assert.Len(t, filteredHistory, 2, "Filter should return two signals with NodeID 'node1'")
	assert.Equal(t, sig1, filteredHistory[0], "First filtered signal should match sig1")
	assert.Equal(t, sig3, filteredHistory[1], "Second filtered signal should match sig3")
}

func TestSimpleHistoryManager_CompressHistory(t *testing.T) {
	hx := nlib.NewSimpleHistoryManager()

	// CompressHistory currently does nothing, so just ensure it returns nil
	err := hx.CompressHistory()
	assert.NoError(t, err, "CompressHistory should not return an error")
}

func TestSimpleHistoryManager_Concurrent(_ *testing.T) {
	hx := nlib.NewSimpleHistoryManager()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sig1 := node.Signal{NodeID: "node1", Task: nlib.NewTextCarrier("testdata1")}
			hx.AddHistory(sig1)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			hx.Filter("node1")
		}
	}(&wg)

	wg.Wait()
}
