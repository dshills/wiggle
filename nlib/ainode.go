package nlib

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Node = (*AINode)(nil)

type AINode struct {
	EmptyNode
	lm llm.LLM
}

func NewAINode(lm llm.LLM, l node.Logger, sm node.StateManager, name string) node.Node {
	ai := AINode{lm: lm}
	ai.SetID(name)
	ai.SetStateManager(sm)
	ai.SetLogger(l)
	ai.MakeInputCh()

	go ai.listen()

	return &ai
}

func (n *AINode) listen() {
	var err error
	var ctx = context.TODO()
	select {
	case sig := <-n.inCh:
		n.LogInfo("Received signal")
		if n.RateLimit(sig) != nil {
			time.Sleep(1 * time.Second)
		}
		sig = node.SignalFromSignal(n.ID(), sig)
		sig, err = n.RunBeforeHook(sig)
		if err != nil {
			n.LogErr(err)
		}
		sig, err = n.GenGuidance(sig)
		if err != nil {
			n.LogErr(err)
		}
		n.LogInfo(fmt.Sprintf("Sending to llm: %v", sig.Data.String()))
		sig, err := n.callLLM(ctx, sig)
		if err != nil {
			n.LogErr(err)
		}
		sig, err = n.RunAfterHook(sig)
		if err != nil {
			n.LogErr(err)
		}
		n.UpdateState(sig)
		n.sendToConnected(sig)
	case <-n.doneCh:
		return
	}
}

func (n *AINode) sendToConnected(sig node.Signal) {
	for _, n := range n.nodes {
		n.InputCh() <- sig
	}
}

func (n *AINode) callLLM(ctx context.Context, sig node.Signal) (node.Signal, error) {
	msgList := llm.MessageList{llm.UserMsg(sig.Data.String())}
	msg, err := n.lm.Chat(ctx, msgList)
	if err != nil {
		return sig, err
	}
	sig.Response = NewStringData(msg.Content)
	return sig, nil
}
