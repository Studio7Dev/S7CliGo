package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://lethiayellow.pythonanywhere.com/"

type HttpResponse struct {
	StatusCode int
	Data       interface{}
}
type AuthService struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (as *AuthService) PerformLogin(username, password string) (*HttpResponse, error) {
	url := baseUrl + "login"
	body := &LoginRequest{
		Username: username,
		Password: password,
	}
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("failed building request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling response: %v", err)
	}

	httpResp := &HttpResponse{
		StatusCode: int(result["statusCode"].(float64)),
	}
	if httpResp.StatusCode >= 200 && httpResp.StatusCode < 300 {
		httpResp.Data = result
	}

	return httpResp, nil
}

// func main() {
// 	authService := &AuthService{}
// 	response, err := authService.PerformLogin("admin", "admins")
// 	if err != nil {
// 		fmt.Printf("An error occurred: %+v\n", err)
// 		return
// 	}

// 	switch response.StatusCode {
// 	case 200:
// 		authenticated, _ := response.Data.(map[string]interface{})["authenticated"].(bool)
// 		fmt.Println("Authenticated:", authenticated)
// 		username, _ := response.Data.(map[string]interface{})["username"].(string)
// 		fmt.Println("Username:", username)
// 		role, _ := response.Data.(map[string]interface{})["role"].(string)
// 		fmt.Println("Role:", role)
// 	default:
// 		fmt.Println("Authentication failed with status code", response.StatusCode)
// 	}
// }
