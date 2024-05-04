# S7 Gui V1

S7 Gui V1 is a simple command-line interface (CLI) tool for interacting with various AI models.

## Features

- Supports multiple AI models, including Merlin, Bing, Hugging Face, Black Box, Tune App, and YouAI.
- Allows users to input messages and receive responses from the selected AI model.
- Provides a user-friendly interface for selecting AI models and inputting messages.
- Supports saving and loading settings.

## Getting Started

To use S7 Gui V1, simply run the `main.go` file using a Go compiler.

```bash
go run main.go
```

Once the application is running, you will be presented with a prompt to input a message. After inputting a message, press enter to receive a response from the currently selected AI model.

## Installation

S7 Gui V1 can be installed by following these steps:

1. Ensure that you have Go installed on your system.
2. Clone the repository to your local machine.
3. Navigate to the repository directory in your terminal.
4. Run `go build` to build the binary.
5. Run the binary to start the application.

## Usage

Here are some usage examples for S7 Gui V1:

### Inputting a Message

To input a message, simply type your message into the input field at the bottom of the application and press enter.

### Selecting an AI Model

To select an AI model, use the "Select AI Provider" button in the top-right corner of the application. This will open a modal where you can select the desired AI model.

### Saving and Loading Settings

To save settings, use the "Save" button in the settings modal. To load settings, use the "Load" button in the settings modal.

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
---
## Contributing

Contributions to S7 Gui V1 are welcome! To contribute, please submit a pull request or open an issue.

## License

S7 Gui V1 is released under the MIT License. See the `LICENSE` file for more information.

## Support

If you encounter any issues or have any questions about S7 Gui V1, please open an issue in the repository.

## Acknowledgements

S7 Gui V1 was built using the following open-source libraries and tools:

- Fyne GUI library
- Go programming language

Thank you to the maintainers and contributors of these projects for their hard work and dedication.

## Disclaimer

S7 Gui V1 is a research project and is not intended for use in production environments. The maintainers and contributors of S7 Gui V1 are not responsible for any damage or loss caused by the use of this software.

