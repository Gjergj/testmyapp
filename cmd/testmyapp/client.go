package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gjergj/testmyapp/pkg/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type CustomHTTPClient struct {
	Client       *http.Client
	Host         string
	Token        string
	RefreshToken string
}

func NewCustomHTTPClient(host, token string, refreshToken string) *CustomHTTPClient {
	return &CustomHTTPClient{
		Client: &http.Client{
			Transport: &http.Transport{
				DisableCompression: false,
			},
			Timeout: 5 * time.Second,
		},
		Host:         host,
		Token:        token,
		RefreshToken: refreshToken,
	}
}

func (c *CustomHTTPClient) Get(url string) (*http.Response, error) {
	return c.doRequest(http.MethodGet, url, "", nil)
}

func (c *CustomHTTPClient) Post(url string, body interface{}) (*http.Response, error) {
	return c.doRequest(http.MethodPost, url, "application/json", body)
}

func (c *CustomHTTPClient) Delete(url string, body interface{}) (*http.Response, error) {
	return c.doRequest(http.MethodDelete, url, "application/json", body)
}

func (c *CustomHTTPClient) Put(url string, body interface{}) (*http.Response, error) {
	return c.doRequest(http.MethodPut, url, "application/json", body)
}

func (c *CustomHTTPClient) Upload(url string, files []string) (*http.Response, error) {
	// Create a buffer to store the request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Add files to the request
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer file.Close()

		// by adding the directory as the form field name, the server will create the directory if it doesn't exist
		dir := filepath.Dir(fileName)
		part, err := writer.CreateFormFile(dir, fileName)
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return nil, err
		}

		// Copy file content to the form data
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println("Error copying file content:", err)
			return nil, err
		}
	}
	// Close the multipart writer
	writer.Close()

	return c.doRequest(http.MethodPost, url, writer.FormDataContentType(), &requestBody)
}

func (c *CustomHTTPClient) doRequest(method, url string, contentType string, body interface{}) (*http.Response, error) {
	var buf io.ReadWriter
	var req *http.Request
	var err error
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req, err = http.NewRequest(method, url, v)
		default:
			buf = new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, url, buf)
		}
	} else {
		req, err = http.NewRequest(method, url, buf)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		// Refresh the token
		err = c.refreshToken()
		if err != nil {
			return nil, err
		}

		// Retry the request with the new token
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		res, err = c.Client.Do(req)
	}
	return res, err
}

func (c *CustomHTTPClient) refreshToken() error {
	// URL to send the POST request
	serverURL := c.Host + "/v1/refresh_token"

	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("refresh-token", c.RefreshToken)
	req.Header.Add("Accept-Encoding", "gzip")
	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status
	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
		// Read the response body
		responseBody, _ := io.ReadAll(response.Body)
		if len(responseBody) > 0 {
			fmt.Printf("Response body: %s\n", responseBody)
		}

		return errors.New("failed to refresh token")
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Parse the response body to get the new token
	// Parse the response body to get the JWT and refresh token
	apiResp := models.LoginResponse{}
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return err
	}
	// Update the client's token
	c.Token = apiResp.Token
	c.RefreshToken = apiResp.RefreshToken
	return nil
}

func (c *CustomHTTPClient) Login(username, password string) (string, string, string, error) {
	// URL to send the POST request to
	serverURL := apiHost + "/v1/login"

	// Credentials
	creds := models.LoginRequest{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(creds)
	if err != nil {
		return "", "", "", err
	}

	// Make the POST request
	response, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", "", "", err
	}
	defer response.Body.Close()

	// Check the response status
	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
		return "", "", "", err
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", "", "", err
	}

	// Parse the response body to get the JWT and refresh token
	apiResp := models.LoginResponse{}
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return "", "", "", err
	}

	c.Token = apiResp.Token
	c.RefreshToken = apiResp.RefreshToken
	return c.Token, c.RefreshToken, apiResp.UserID, nil
}
