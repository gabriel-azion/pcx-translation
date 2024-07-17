package gptintegration

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func GptIntegration(option, filepath string) error {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	client := openai.NewClient(apiKey)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "you are a helpful chatbot",
			},
		},
	}

	done := make(chan bool, 2) 

	fmt.Println("Translations take some time, take a chill pill and a cup of coffee")

	go func() {
		defer func() { done <- true }() 

		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: `Text: ` + string(data) + `
			---

			- Translate this text to Brazilian portuguese.
			- Don't change the context
			- The output format must be markdown
			- Never translate the following terms:
				Edge Application or Edge Applications
				Edge Functions
				Edge Storage
				Object Storage
				Azion CLI
				(and the variations in lower case, plural and singular)
			`,
		})

		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			panic(err)
		}
		response := resp.Choices[0].Message.Content

		// Write the response to a new file
		err = os.WriteFile("ptversion.md", []byte(response), 0644)
		if err != nil {
			panic(err)
		}

	}()

	go func() {
		defer func() { done <- true }() // Send a signal to the channel when this goroutine finishes

		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: `Text: ` + string(data) + `
		---

		- Translate this text to latam spanish.
		- Don't change the context
		- The output format must be markdown
		- Actúa como un traductor y editor experto en español localizado para LATAM, que sea lo más neutro posible y evite regionalismos.
		- Quiero que priorices estructuras y tiempos verbales simples, y evites tiempos compuestos que puedan aumentar la complejidad de la versión traducida.
		- El contenido está enfocado en el área de computación y tecnología, por lo que debe priorizar el uso de términos de uso común y relacionados.
		- No traduzcas los siguientes términos: edge computing, edge, Azion, edge application, edge function, edge firewall [FALTA COMPLETAR COM PLANILHA DE TERMOS]
		- Usa artículos en femenino para: edge computing, [FALTA COMPLETAR COM PLANILHA DE TERMOS]
		- Usa artículos en masculino para: edge [FALTA COMPLETAR COM PLANILHA DE TERMOS]
		- Siempre debes traducir request como solicitud.
		- No cambies el contexto del contenido.
		- Traduce siempre usando como segundo persona: tú y sus conjugaciones.
		`,
		})

		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			panic(err)
		}
		response := resp.Choices[0].Message.Content

		// Write the response to a new file
		err = os.WriteFile("./spanishversion.md", []byte(response), 0644)
		if err != nil {
			panic(err)
		}
	}()

	<-done 
	<-done

	print("Done - feito o brique")

	return nil
}
