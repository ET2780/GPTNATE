package agents

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/natexcvi/go-llm/engines"
	"github.com/natexcvi/go-llm/tools"
)

type Representable interface {
	Encode() string
	Schema() string
}

type Example[T Representable, S any] struct {
	Input             T
	IntermediarySteps []*engines.ChatMessage
	Answer            S
}

type Task[T Representable, S any] struct {
	Description  string
	Examples     []Example[T, S]
	AnswerParser func(string) (S, error)
}

func (task *Task[T, S]) Compile(input T, tools map[string]tools.Tool) *engines.ChatPrompt {
	prompt := &engines.ChatPrompt{
		History: []*engines.ChatMessage{
			{
				Role: engines.ConvRoleSystem,
				Text: fmt.Sprintf("You are a smart, autonomous agent given the task below. "+
					"You will be given input from the user in the following format:\n\n"+
					"%s\n\n Complete the task step-by-step, "+
					"reasoning about your solution steps by sending a message beginning "+
					"with `%s: `.\n\nTask description: %s",
					input.Schema(), ThoughtCode, task.Description),
			},
		},
	}
	task.enrichPromptWithTools(tools, prompt)
	task.enrichPromptWithExamples(prompt)
	prompt.History = append(prompt.History, &engines.ChatMessage{
		Role: engines.ConvRoleSystem,
		Text: fmt.Sprintf("Now, you will be given the input. "+
			"It's very important that every message you send begins with either "+
			"`%s: `, `%s: `, or `%s: `, and end with `%s`.",
			ThoughtCode, ActionCode, AnswerCode, EndMarker),
	})
	prompt.History = append(prompt.History, &engines.ChatMessage{
		Role: engines.ConvRoleUser,
		Text: input.Encode(),
	})
	return prompt
}

func (task *Task[T, S]) enrichPromptWithExamples(prompt *engines.ChatPrompt) {
	prompt.History = append(prompt.History, &engines.ChatMessage{
		Role: engines.ConvRoleSystem,
		Text: "Here are some examples of how you might solve this task:",
	})
	for _, example := range task.Examples {
		prompt.History = append(prompt.History, &engines.ChatMessage{
			Role: engines.ConvRoleUser,
			Text: example.Input.Encode(),
		})
		for _, step := range example.IntermediarySteps {
			prompt.History = append(prompt.History, step)
		}
		marshalledAnswer, err := json.Marshal(example.Answer)
		if err != nil {
			panic(err)
		}
		prompt.History = append(prompt.History, &engines.ChatMessage{
			Role: engines.ConvRoleSystem,
			Text: fmt.Sprintf(MessageFormat, AnswerCode, string(marshalledAnswer)),
		})
	}
}

func (*Task[T, S]) enrichPromptWithTools(tools map[string]tools.Tool, prompt *engines.ChatPrompt) {
	if len(tools) > 0 {
		toolsList := make([]string, 0, len(tools))
		for name, tool := range tools {
			toolsList = append(toolsList, fmt.Sprintf("%s(%s) # %s", name, tool.ArgsSchema(), tool.Description()))
		}
		prompt.History = append(prompt.History, &engines.ChatMessage{
			Role: engines.ConvRoleSystem,
			Text: fmt.Sprintf("Here are some tools you can use. To use a tool, "+
				"send a message in the form of `%s: tool_name(args)`, "+
				"where `args` is a valid one-line JSON representation of the arguments"+
				" to the tool, as specified for it (using JSON schema). You will get "+
				"the output in "+
				"a message beginning with `%s: `, or an error message beginning "+
				"with `%s: `.\n\nTools:\n%s",
				ActionCode, ObservationCode, ErrorCode, strings.Join(toolsList, "\n")),
		})
	}
}
