# chat-ai

![](https://img.shields.io/github/stars/michalswi/chat-ai)
![](https://img.shields.io/github/issues/michalswi/chat-ai)
![](https://img.shields.io/github/forks/michalswi/chat-ai)
![](https://img.shields.io/github/last-commit/michalswi/chat-ai)
![](https://img.shields.io/github/release/michalswi/chat-ai)


**Terminal based AI chat powered by ChatGPT and Gemini.**

Go app to transform text files into AI-powered reviews is [here](https://github.com/michalswi/file-go-openai) .

```
$ ./chat-ai -h
  -p string
    	AI provider [chatgpt, gemini]
```

### ChatGPT (OpenAI)

```
export API_KEY=<>

./chat-ai -p chatgpt
```

### Gemini

```
export VAI_PROJECT_ID=<>
export API_KEY=<>>.json
export VAI_REGION=<>

./chat-ai -p gemini
```

### example

```
> chatgpt

$ ./chat-ai_macos_arm64 -p chatgpt
Fri, 02 Aug 2024 13:55:46 UTC [chat-ai:chatgpt]: hello
> Waiting for ChatGPT..
Hello! How can I assist you today?

Fri, 02 Aug 2024 13:56:53 UTC [chat-ai:chatgpt]: q
Exiting chat. bye!


> gemini

$ ./chat-ai_macos_arm64 -p gemini
Fri, 02 Aug 2024 13:59:58 UTC [chat-ai:gemini]: hello
> Waiting for Gemini..
Hello! 👋  How can I help you today? 😊

Fri, 02 Aug 2024 14:00:02 UTC [chat-ai:gemini]: q
Exiting chat. bye!
```