package nlib

import (
	"strings"

	"github.com/dshills/wiggle/llm"
)

/*
	Integration
	•	Start with Simple Concatenation as the default integration strategy. It works for
	both small and large models.
	•	If using smaller models, you can optionally introduce a Coherence
	Rewriting step with a larger model to ensure that the final output is clean and coherent.
*/

/*
SimpleConcatenationIntegrator combines the outputs of partitioned tasks
by concatenating them into a single coherent result. This straightforward
integration strategy is effective when the partitioned chunks are semantically
meaningful and can be recombined without the need for further processing or
transformation, ensuring a seamless and consistent output.

Use Case:

  - Works well for both large models and smaller models. Since the partitioning
    strategy ensured the chunks were semantically meaningful, concatenation will
    produce a coherent result. For frontier models, it is straightforward; for local models,
    this ensures the results are combined correctly after processing smaller chunks.
*/
func SimpleConcatenationIntegrator(parts []string) (string, error) {
	return strings.Join(parts, "\n"), nil
}

/*
CoherenceRewritingIntegrator combines the outputs of partitioned tasks
and then passes the combined text to a large language model for coherence rewriting.
This integration strategy ensures that the final output is coherent and well-structured,
particularly useful when integrating results from smaller models with limited context windows.
The large model refines and polishes the text to improve readability and consistency.

Use Case:

  - For smaller models like local models, this strategy is very effective. After partitioning
    and processing smaller chunks, a large model like frontier models can handle the final integration
    step, ensuring the entire outpt is coherent.
  - For large models like frontier models, this step is often unnecessary, but it can still help in
    particularly complex workflows.
*/
func CoherenceRewritingIntegrator(lm llm.LLM, parts []string) (string, error) {
	// Combine the parts
	combined := strings.Join(parts, "\n")

	// Send the combined text to a large model for coherence rewriting
	response, err := lm.GenerateResponse(combined, "Please rewrite the following text for coherence.")
	if err != nil {
		return "", err
	}

	return response, nil
}

/*
SummaryBasedIntegrator processes partitioned outputs by summarizing each part
individually using a language model. The summaries are then combined to form
a final, concise result. This integration strategy is particularly useful when
dealing with large or complex outputs, allowing for efficient condensation of
information while maintaining the key points from each partition.

Use Case:

  - Smaller models benefit greatly from this strategy. Instead of trying to integrate
    large amounts of data directly, breaking down and summarizing allows smaller
    models to handle the integration phase more effectively.
  - Larger models like frontier models can handle this, but it may not be necessary unless the
    data is truly massive.
*/
func SummaryBasedIntegrator(llm llm.LLM, parts []string) (string, error) {
	var summaries []string
	for _, part := range parts {
		// Summarize each part individually
		summary, err := llm.GenerateResponse(part, "Summarize the following text.")
		if err != nil {
			return "", err
		}
		summaries = append(summaries, summary)
	}

	// Combine the summaries
	finalSummary := strings.Join(summaries, "\n")
	return finalSummary, nil
}
