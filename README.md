# Wiggle the Multi-Node LLM Processing Framework

This Go project provides a flexible and modular framework for chaining multiple Language Learning Models (LLMs), integrating context from various sources like vector databases, and efficiently processing large or complex data by partitioning tasks across nodes and integrating results. The framework is designed to support both large models (e.g., GPT-4) and smaller models (e.g., LLaMA 3.1), ensuring scalability, modularity, and efficiency.

Table of Contents

- [Overview](#overview)
- [Features](#Features)
- [Architecture](#Architecture)
- [Getting Started](#Getting-Started)
- [Installation](#Installation)
- [Basic Example](#Basic-Example)
- [Core Concepts](#Core-Concepts)
- [Node](#Node)
- [Signal](#Signal)
- [Guidance](#Guidance)
- [PartitionerNode](#PartitionerNode)
- [IntegratorNode](#IntegratorNode)
- [Set](#Set)
- [Hooks](#Hooks)
- [Coordinator](#Coordinator)
- [Resource Management](#Resource-Management)
- [Context Management](#Context-Management)
- [DataCarrier](#DataCarrier)
- [Example Workflow](#Example-Workflow)
- [Contributing](#Contributing)
- [License](#license)

## Overview

This framework is built around Nodes that can process data, transform it, or query external systems (such as LLMs). It supports partitioning data into manageable tasks and then integrating the results, providing a robust mechanism to handle both simple and complex workflows. The system is highly extensible and can integrate large LLMs like GPT-4, as well as smaller models like LLaMA.

## Features

- Modular Node-based Architecture: Chain nodes together to process and transform data in a structured way.
- Partitioning and Integration: Split large tasks into smaller units, process them independently, and aggregate results.
- ontext Integration: Automatically manage and update context across multiple nodes, ensuring relevance at each step.
- Error Handling: Gracefully handle errors, with configurable policies to continue or halt workflows.
- Rate Limiting: Prevent system overloads with rate-limited processing, ensuring efficient use of resources.
- LLM Integration: Easily connect to LLMs and vector databases, integrating AI into your workflows.

## Architecture

The system consists of several key components:

1. Nodes: Individual units of work that process data in the form of Signals. Nodes can be action nodes, query nodes, partitioner nodes, or integrator nodes.
2. Signals: Data structures passed between nodes, containing data, context, metadata, and history.
3. Guidance: Interface for generating processing instructions based on the signal’s data and context.
4. Partitioning: Split large data into smaller chunks that are processed independently by multiple nodes.
5. Integration: Combine the results of partitioned tasks into a coherent final result.

## Getting Started

### Installation

To get started, clone this repository and install dependencies:

```bash
git clone https://github.com/dshills/wiggle
cd wiggle
go mod tidy
```

### Basic Example

Here’s a simple example that demonstrates sending a signal to multiple nodes, each processing a message using a large language model (LLM).

```go
package main

import (
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
)

const envURL = "OPENAI_API_URL"
const envKey = "OPENAI_API_KEY"
const model = "gpt-4o"

func main() {
	// Setup LLM
	lm := openai.New(os.Getenv(envURL), model, os.Getenv(envKey), nil)

	// Create a Logger
	logger := nlib.NewSimpleLogger(log.Default())

	// Create State Manager
	stateMgr := nlib.NewSimpleStateManager()

	// Define output writer
	writer := os.Stdout

	// Create Nodes
	firstNode := nlib.NewAINode(lm, logger, stateMgr, "AI Node")
	outNode := nlib.NewOutputStringNode(writer, logger, stateMgr, "Output Node")
	firstNode.Connect(outNode)

	// Send it
	firstNode.InputCh() <- nlib.NewDefaultSignal(firstNode, "Why is the sky blue?")

	// Wait for the output node to print the result
	stateMgr.WaitFor(outNode)
}
```

In this example:

- We setup an LLM to receive queries
- Created an AI Node and Output Node
- Connected the nodes
- Sent the task to the first node
- Wait for the last node (output) to complete

## Core Concepts

### Node

A Node is the core processing unit in Wiggle. It processes incoming signals, executes actions (such as querying a model or transforming data), and forwards the processed signal to connected nodes. The interface is modular, allowing different node types to be chained together for flexible workflows.

```go
type Node interface {
    Connect(Node)
    ID() string
    InputCh() chan Signal
    SetGuidance(Guidance)
    SetHooks(Hooks)
    SetID(string)
    SetLogger(Logger)
    SetResourceManager(ResourceManager)
    SetStateManager(StateManager)
}
```

### Signal

A Signal represents the data structure passed between nodes, containing the main data being processed, its context, and any metadata or history required. It ensures smooth and coherent propagation of data and context across the entire node chain.

```go
type Signal struct {
	NodeID   string
	Data     DataCarrier
	Response DataCarrier
	Context  ContextManager
	Meta     []Meta
	History  HistoryManager
	Err      error
	Status   string
}
```
### Guidance

Guidance generates structured instructions for processing based on the signal’s content and metadata. It typically interfaces with an LLM.

```go
type Guidance interface {
    Generate(signal Signal) (Signal, error)
}
```

### PartitionerNode

A PartitionerNode splits large or complex tasks into smaller chunks using a partitioning function (PartitionerFn), enabling parallel processing by downstream nodes. This design allows for efficient handling of large-scale data processing.

```go
type PartitionerNode interface {
    Node
    SetPartitionFunc(partitionFunc PartitionerFn)
    SetChildNodes(nodes ...Node)
}
```

### IntegratorNode

The IntegratorNode aggregates the results from partitioned tasks using an integrator function (IntegratorFn). This ensures that all the partitioned results are combined into a single, coherent output, maintaining data consistency throughout the workflow.

```go
type IntegratorNode interface {
    Node
    SetIntegratorFunc(integratorFunc IntegratorFn)
    SetChildNodes(nodes ...Node)
}
```

### Set

A Set represents a collection of nodes that form a processing pipeline. It organizes nodes into a structured chain and manages the flow of data between them. A set allows you to define a complex workflow with multiple interconnected nodes.

```go
type Set interface {
    Node
    SetStartNode(Node)
    SetCoordinator(Coordinator)
}
```

### Hooks

The Hooks interface allows custom pre- and post-processing logic to be executed before or after a node processes a signal. This is useful for validation, logging, or modifying results before the data moves forward.

```go
type Hooks interface {
    BeforeAction(Signal) (Signal, error)
    AfterAction(Signal) (Signal, error)
}
```

### Coordinator

The Coordinator manages the execution flow of multiple nodes, ensuring tasks complete within a specified timeout and handling errors.

```go
type Coordinator interface {
    WaitForCompletion(nodes ...Node) error
    CancelOnTimeout(duration time.Duration)
}
```

### Resource Management

The ResourceManager controls resource usage (e.g., rate limiting) to prevent overwhelming external systems like LLM APIs or databases.

```go
type ResourceManager interface {
    RateLimit(Signal) error
}
```

### Context Management

A ContextManager manages the contextual data passed between nodes, ensuring that relevant information is consistent as the signal flows through the node chain.

```go
type ContextManager interface {
    SetContext(key string, data DataCarrier)
    RemoveContext(key string)
    GetContext(key string) (DataCarrier, error)
}
```

### DataCarrier

The DataCarrier provides an abstraction for handling different types of data, such as strings, JSON, or vectors, within a signal. It ensures flexibility in how data is passed and processed across the workflow.

```go
type DataCarrier interface {
    Vector() []float32
    String() string
    JSON() []byte
}
```

## Example Workflow

1. A PartitionerNode receives a signal containing a large chunk of data and uses its partition function to split the data into smaller chunks.
2. The smaller chunks are distributed to downstream Action Nodes or other specialized nodes for processing.
3. After processing, the results are aggregated by an IntegratorNode, which combines the partitioned data into a single output.
4. The Set orchestrates this entire workflow, ensuring that nodes are connected and coordinated.

## Contributing

Wiggle is under heavy development and we welcome ideas and contributions.  If you’d like to improve this framework or report issues, feel free to create a pull request or open an issue.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
