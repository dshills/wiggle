package node

import "strings"

// PartitionerFn is a function type that takes an input string and splits it
// into smaller parts or tasks. It is used by PartitionerNodes to divide
// large or complex data into manageable chunks, enabling parallel processing
// by multiple nodes in the chain.
type PartitionerFn func(string) ([]string, error)

/*
	Partitioning
	•	Use Sentence Splitting with Overlap as the default partitioning strategy.
	This ensures that even smaller models maintain some level of context between partitions.
	•	Keep chunk sizes small enough (e.g., 1000 characters or less) to be
	compatible with smaller models while still taking advantage of large models’
	capacity to handle more information.
*/

/*
SemanticChunkingPartition splits input text into semantically meaningful chunks,
such as paragraphs or sections, based on natural boundaries (e.g., newlines).
This partitioning ensures that each chunk maintains context and meaning, making
it suitable for both large and small language models to process efficiently.
For models with smaller context windows, the function further splits large chunks
to prevent overflow while preserving semantic integrity.

Use Case:

  - Frontier models can handle large chunks, so this strategy works directly without much chunking. It can process
    entire paragraphs or sections and maintain coherence.
  - Local models or smaller models have more limited context windows. For these models, the chunkText
    helper ensures that partitions are not too large. The chunks are smaller, but they
    remain meaningful by maintaining semantic structure.
*/
func SemanticChunkingPartition(input string) ([]string, error) {
	// Split input into chunks by semantic boundaries such as paragraphs or sentence groups.
	// For simplicity, we'll use paragraphs here.
	chunks := strings.Split(input, "\n\n") // Split by paragraphs

	// Ensure that chunks aren't too long for smaller models by further splitting if necessary
	maxChunkSize := 1000 // Character limit for smaller models
	var finalChunks []string
	for _, chunk := range chunks {
		if len(chunk) > maxChunkSize {
			subChunks := ChunkText(chunk, maxChunkSize)
			finalChunks = append(finalChunks, subChunks...)
		} else {
			finalChunks = append(finalChunks, chunk)
		}
	}

	return finalChunks, nil
}

// ChunkText is a helper function to split long text into smaller chunks.
func ChunkText(text string, size int) []string {
	var chunks []string
	for len(text) > size {
		chunks = append(chunks, text[:size])
		text = text[size:]
	}
	if len(text) > 0 {
		chunks = append(chunks, text)
	}
	return chunks
}

/*
SentenceSplittingWithOverlap partitions input text by splitting it into
sentences, with a configurable overlap between consecutive chunks. This overlap
helps maintain context across partitions, especially useful for models with smaller
context windows. The function ensures that important context from the end of one
chunk is carried over to the beginning of the next, improving coherence in downstream processing.

Use Case:

  - Frontier models can process large portions of text in one go, but this strategy ensures that even if frontier models
    breaks output into smaller chunks for parallel processing, overlapping sentences help maintain flow and coherence.
  - Local models benefits even more from overlap since it has a smaller context window.
    The overlapping sentences help ensure that context is carried across chunks.
*/
func SentenceSplittingWithOverlap(input string, overlap int) ([]string, error) {
	// Split input into sentences
	sentences := strings.Split(input, ". ") // Basic sentence splitting

	// Generate overlapping chunks
	var chunks []string
	for i := 0; i < len(sentences); i += overlap {
		end := i + overlap
		if end > len(sentences) {
			end = len(sentences)
		}
		chunk := strings.Join(sentences[i:end], ". ")
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

/*
TaskBasedPartitioning divides the input into logical sections or tasks,
such as "Introduction", "Body", and "Conclusion", based on predefined structure
(e.g., headers or sections). This partitioning strategy is ideal for processing
documents or multi-step tasks, ensuring that each part is handled independently
while maintaining the logical flow of the overall content.

Use Case:

  - Frontier models can process entire sections of a complex task or document, allowing it to handle
    large parts like introductions or conclusions.
  - Local models will still benefit, but more complex tasks may need to be split further
    based on sub-tasks to accommodate its smaller context window.
*/
func TaskBasedPartitioning(input string) ([]string, error) {
	// This assumes a structured input where sections can be identified, for example by headings
	sections := strings.Split(input, "\n# ") // Split by section headers
	for i, section := range sections {
		sections[i] = "# " + section // Re-add the header symbol for proper formatting
	}
	return sections, nil
}
