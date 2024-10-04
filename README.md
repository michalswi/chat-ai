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

- keep context (conversation continuity) between queries in the same session
- queries (answer+question) are kept in the local file

```
export API_KEY=<>

./chat-ai -p chatgpt
```

### Gemini

- **[not implemented]** keep context (conversation continuity) between queries in the same session
- **[not implemented]** queries (answer+question) are kept in the local file

```
export VAI_PROJECT_ID=<>
export API_KEY=<>>.json
export VAI_REGION=<>

./chat-ai -p gemini
```

### example

```
> chatgpt chat-ai

$ ./chat-ai -p chatgpt
Thu, 03 Oct 2024 16:16:05 UTC [chat-ai:chatgpt]: hello, my name is john
> Waiting for ChatGPT..
Hello, John! How can I assist you today?
Thu, 03 Oct 2024 16:16:10 UTC [chat-ai:chatgpt]: what is my name?
> Waiting for ChatGPT..
Your name is John. How can I help you today?
Thu, 03 Oct 2024 16:16:21 UTC [chat-ai:chatgpt]: what was my previous question?
> Waiting for ChatGPT..
Your previous question was, "what is my name?" If you have any more questions or need assistance, feel free to ask!
Thu, 03 Oct 2024 16:16:41 UTC [chat-ai:chatgpt]: q
Exiting chat. bye!


> gemini chat-ai

$ ./chat-ai -p gemini
Fri, 02 Aug 2024 13:59:58 UTC [chat-ai:gemini]: hello
> Waiting for Gemini..
Hello! 👋  How can I help you today? 😊

Fri, 02 Aug 2024 14:00:02 UTC [chat-ai:gemini]: q
Exiting chat. bye!
```