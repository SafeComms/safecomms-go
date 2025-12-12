// Package safecomms provides a client for the SafeComms API.
//
// SafeComms is a powerful moderation API that helps you keep your community safe.
// This package allows you to easily integrate SafeComms into your Go applications.
package safecomms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const DefaultBaseURL = "https://api.safecomms.dev"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string, baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type ModerateTextRequest struct {
	Content             string `json:"content"`
	Language            string `json:"language,omitempty"`
	Replace             bool   `json:"replace,omitempty"`
	Pii                 bool   `json:"pii,omitempty"`
	ReplaceSeverity     string `json:"replaceSeverity,omitempty"`
	ModerationProfileId string `json:"moderationProfileId,omitempty"`
}

func (c *Client) ModerateText(req ModerateTextRequest) (map[string]interface{}, error) {
	if req.Language == "" {
		req.Language = "en"
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/moderation/text", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

type ModerateImageRequest struct {
	Image               string `json:"image"`
	Language            string `json:"language,omitempty"`
	ModerationProfileId string `json:"moderationProfileId,omitempty"`
}

func (c *Client) ModerateImage(req ModerateImageRequest) (map[string]interface{}, error) {
	if req.Language == "" {
		req.Language = "en"
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/moderation/image", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

type ModerateImageFileRequest struct {
	FilePath            string
	Language            string
	ModerationProfileId string
}

func (c *Client) ModerateImageFile(req ModerateImageFileRequest) (map[string]interface{}, error) {
	if req.Language == "" {
		req.Language = "en"
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filepath.Base(req.FilePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.WriteField("language", req.Language)
	if req.ModerationProfileId != "" {
		writer.WriteField("moderationProfileId", req.ModerationProfileId)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/moderation/image/upload", body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetUsage() (map[string]interface{}, error) {
	httpReq, err := http.NewRequest("GET", c.baseURL+"/usage", nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
}
