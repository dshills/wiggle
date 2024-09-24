# Wigle the Multi-Node LLM Processing Framework

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
- [Coordinator](#Coordinator)
- [Error Handling](#Error-Handling)
- [Resource Management](#Resource-Management)
- [Advanced Usage](#Advanced-Usage)
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

func main() {
    logger := &ConsoleLogger{}
    coordinator := NewSimpleCoordinator(10 * time.Second)

    node1 := NewSimpleNode("node1", ActionPrint)
    node2 := NewSimpleNode("node2", ActionPrint)

    signal := Signal{
        NodeID:  "node1",
        Data:    MessageData{Message: "Some input data"},
        Context: "Example context",
    }

    go func() {
        node1.InputCh() <- signal
    }()

    err := coordinator.WaitForCompletion(node1, node2)
    if err != nil {
        fmt.Println("Error waiting for nodes:", err)
    } else {
        fmt.Println("All nodes completed successfully.")
    }
}
```

In this example:

- We define two simple nodes that print the signal they receive.
- The Coordinator ensures the nodes process the signal within a specified timeout.

## Core Concepts

### Node

A Node represents a processing unit in the system. Each node processes a Signal, potentially transforming it, querying an LLM, or passing it to other nodes.

```go
type Node interface {
    Clone() Node
    Connect(Node) error
    ID() string
    InputCh() chan Signal
    SetID(string)
    SetLogger(Logger)
    SetErrorHandler(ErrorHandler)
    SetCoordinator(Coordinator)
    SetResourceManager(ResourceManager)
    SetHooks(NodeHooks)
    SetStateManager(StateManager)
    Wait()
}
```

### Signal

A Signal carries the data, context, and metadata across nodes. It represents the input/output of a node.

```go
type Signal struct {
    NodeID   string
    Data     DataCarrier
    Context  string
    Meta     []Meta
    History  HistoryManager
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

A PartitionerNode splits input data into smaller tasks, distributing them across child nodes for parallel processing.

```go
type PartitionerNode interface {
    Node
    SetPartitionFunc(partitionFunc PartitionFn)
    SetChildNodes(nodes ...Node)
}
```

### IntegratorNode

An IntegratorNode gathers results from child nodes and combines them into a single coherent output.

```go
type IntegratorNode interface {
    Node
    SetIntegratorFunc(integratorFunc IntegratorFn)
    SetChildNodes(nodes ...Node)
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

### Error Handling

The ErrorHandler defines how errors are handled during node processing. It allows you to decide whether the workflow should continue or halt when errors occur.

```go
type ErrorHandler interface {
    HandleError(Signal, error) bool
}
```

### Resource Management

The ResourceManager controls resource usage (e.g., rate limiting) to prevent overwhelming external systems like LLM APIs or databases.

```go
type ResourceManager interface {
    RateLimit(Signal) error
}
```

## Advanced Usage

You can easily extend this framework by implementing additional Node, Guidance, or other interfaces, adding custom logic as needed. For instance, you could build more complex workflows using partitioners, integrators, and external vector databases for context retrieval.

### Partitioning and Integration Example

Here’s an example using partitioners and integrators to process and combine results from multiple tasks.

```go
func main() {
    partitioner := NewSimplePartitionerNode("partitioner", SentenceSplittingWithOverlap)
    integrator := NewSimpleIntegratorNode("integrator", SimpleConcatenationIntegrator)

    // Set up child nodes
    child1 := NewSimpleNode("child1", ActionPrint)
    child2 := NewSimpleNode("child2", ActionPrint)
    child3 := NewSimpleNode("child3", ActionPrint)

    // Configure partitioner and integrator
    partitioner.SetChildNodes(child1, child2, child3)
    integrator.SetChildNodes(child1, child2, child3)

    signal := Signal{
        NodeID:  "partitioner",
        Data:    MessageData{Message: "Some complex input data to partition"},
    }

    // Send signal to partitioner
    go func() {
        partitioner.InputCh() <- signal
    }()

    // Wait for completion
    err := integrator.WaitForCompletion(child1, child2, child3)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Integration complete")
    }
}
```

## Contributing

We welcome contributions! If you’d like to improve this framework or report issues, feel free to create a pull request or open an issue.

## License

This project is licensed under the MIT License. See the LICENSE file for details.