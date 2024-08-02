# chat-ai

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