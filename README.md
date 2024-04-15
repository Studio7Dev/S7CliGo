
# S7 Beta CLI (Command Line Interface) [Developer Version]


A powerful and customizable command line interface with various features and integrations.

Features
--------

* Authentication system
* Merlin AI (GPT 3) integration
* Hugging Face AI integration
* BlackBox Programming AI Chat integration
* Searx Search Engine integration
* TMDB (The Movie Database) integration
* Bing AI integration
* AI text to image

Getting Started
---------------

1. Clone the repository
2. Install the required dependencies
3. Run the `main.go` file

Commands
--------

| Command | Description |
| --- | --- |
| `exit` | Exits the application |
| `merlin` | Interact with Merlin AI (GPT 3) |
| `hug` | Interact with Hugging Face AI |
| `blackbox` | Interact with BlackBox Programming AI Chat |
| `searx` | Use Searx Search Engine |
| `movie` | Search for a movie with TMDB |
| `bingai`| Interact with the Bing AI from `bing.com/chat`|
| `youai` | Interact with `you.com` AI (YouAI) |
| `img-gen` | Generate images from text prompts, using AI |
| `enable-tcp` | Enable TCP API Usage |
| `enable-http`| Enable Web API Usage |
#### Available Models for HuggingFace

| Model Name | Description |
| --- | --- |
| `google/gemma-7b-it` | Google AI |
| `mistralai/Mixtral-8x7B-Instruct-v0.1` | Mixtral Chat AI v0.1 |
| `mistralai/Mistral-7B-Instruct-v0.2` | Mixtral Chat AI v0.2 |
| `meta-llama/Llama-2-70b-chat-hf` | Facebook (Meta) Llama AI |
| `NousResearch/Nous-Hermes-2-Mixtral-8x7B-DPO` | NousResearch x Mixtral-8x7B |
| `codellama/CodeLlama-70b-Instruct-hf` | CodeLlama (Programming Assistant AI) |
| `openchat/openchat-3.5-0106` | OpenChat 3.5 (GPT 3.5 Turbo) |

## Run?
```bash
go get .
go run main.go
```
## Build?
```bash
go get .
go build -o CLI main.go
```
## Wanna use the SshServer?
```bash
cd SshServer
./sshServer
```
## SshServer Config
```json
{
    "host":"0.0.0.0",
    "port":"22",
    "userdb":"./SshServer/users.txt",
    "working_dir":"../",
    "server_key":"/root/.ssh/id_rsa",
    "cli_binary":"go",
    "cli_cmd_args":["/usr/local/go/bin/go", "run", "main.go"]
}
```
## Cli Settings
```json
{
  "goliath_auth_token":"",
  "merlin_auth_token": "",
  "huggingface_cookie": "",
  "blackbox_cookie": "",
  "youai_cookie":"",
  "username": "admin",
  "password": "admin",
  "tcphost":"0.0.0.0:1337",
  "httphost":"0.0.0.0:8080"
}
```
# Stuff
#### Merlin Chrome `get the auth token from the headers of the chat request. ( Chrome Dev Tools ) `

#### The Auth token is the Authorization header in the chat request, copy it and paste it into `merlin_auth_token` in settings.json


#### https://chromewebstore.google.com/detail/merlin-1-click-access-to/camppjleccjaphfdbohjdohecfnoikec

---
#### Black Box AI
#### Go to https://www.blackbox.ai , open dev tools, send a chat to the ai, find the request for the chat in the network tab, look for the cookie section in the request and copy the entire thing as one line and paste it into the config, `blackbox_cookie` in settings.json
#### Hugging Face
#### Go to https://huggingface.co/chat, login and then get the cookie, im too lazy to explain just figure it out using the last step lmao
#### YouAI
#### Go to https://you.com/chat, open dev tools, send a chat to the ai, find the request for the chat in the network tab, look for the cookie section in the request and copy the entire thing as one line and paste it into the config, `youai_cookie` in settings.json
#### Goliath AI (tune studio)
#### Go to https://studio.tune.app/ and make an api key for `goliath_auth_token`
---
## WebAPI Implementations
---
### Merlin
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"hi"}' http://localhost:8080/merlin
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"hi\"}" http://localhost:8080/merlin
```
### HuggingFace
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"hi how are you", "model":"mistralai/Mistral-7B-Instruct-v0.2"}' http://localhost:8080/hug
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"hi how are you\", \"model\":\"mistralai/Mistral-7B-Instruct-v0.2\"}" http://localhost:8080/hug
```
### BingAI
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"hi how are you"}' http://localhost:8080/bing
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"hi how are you\"}" http://localhost:8080/bing
```
### BlackBox
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"def hello()"}' http://localhost:8080/blackbox
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"def hello()\"}" http://localhost:8080/blackbox
```
### YouAI
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"def hello()"}' http://localhost:8080/youai
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"def hello()\"}" http://localhost:8080/youai
```
### GoliathAI
#### Linux
```curl
curl -N -X POST -H "Content-Type: application/json" --data '{"message":"def hello()"}' http://localhost:8080/goliath
```
#### Windows
```curl
curl -N -X POST -H "Content-Type: application/json" --data "{\"message\":\"def hello()\"}" http://localhost:8080/goliath
```
---

## TCP Implementations
---
### Merlin
```
ai merlin hi how are you? can you code?
```
### HuggingFace
```
ai hug hi how are you? can you code?
```
### BingAI
```
ai bing hi how are you? can you code?
```
### BlackBox
```
ai blackbox hi how are you? can you code?
```
### YouAI (you.com)
```
ai youai hi how are you? can you code?
```
### GoliathAI
```
ai goliath hi how are you? can you code?
```