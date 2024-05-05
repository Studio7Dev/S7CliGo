package verif

import (
	"bufio"
	"encoding/json"
	"fmt"
	"guiv1/models/huggingface"
	"io"
	"strings"
)

var (
	baseprompt = `
	Extract the following information from the verification email and return it in JSON format:

	Verification code
	Account username or email address
	Instructions on how to verify the account
	Any expiration time or deadline for verification
	URL or link to complete the verification process
	Please provide the extracted information in a JSON object with the following keys:
	verification_code, account_username, instructions, expiration_time, and verification_url.
	
	Verification Email:

	`
	ResponseText = ""
)

type VerificationEmail struct {
	VerificationCode string `json:"verification_code"`
	AccountUsername  string `json:"account_username"`
	Instructions     string `json:"instructions"`
	ExpirationTime   string `json:"expiration_time"`
	VerificationURL  string `json:"verification_url"`
}

func ExtractVerificationInfo(emailContent string) (*VerificationEmail, error) {
	hugclient := huggingface.NewHug()
	message := baseprompt + emailContent
	// wrap message to one line
	message = strings.ReplaceAll(message, "\n", " ")
	ChatId := hugclient.ChangeModel("meta-llama/Meta-Llama-3-70B-Instruct")
	Id_ := hugclient.GetMsgUID(ChatId)
	err, Response := hugclient.SendMessage(message, ChatId, Id_, true)

	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(Response.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event["type"] == "finalAnswer" {
			ResponseText = event["text"].(string)
		}
	}
	Split1 := strings.Split(ResponseText, "```")
	if len(Split1) < 2 {
		return nil, fmt.Errorf("no JSON response found")
	}
	Split2 := strings.Split(Split1[1], "```")[0]

	var verificationEmail VerificationEmail
	err = json.Unmarshal([]byte(Split2), &verificationEmail)
	if err != nil {
		return nil, err
	}
	return &verificationEmail, nil

}
