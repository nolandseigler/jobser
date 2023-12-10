package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/labstack/echo/v4"
)

type GetAnalyzeHandlerReq struct {
	AnalyzeText string `query:"analyze-text"`
	Summarize   string `query:"summarize"`
	Sentiment   string `query:"sentiment"`
	Keyword     string `query:"keyword"`
}

type supportedAPI string

const (
	supportedAPISummary   supportedAPI = "SUPPORTED_API_SUMMARY"
	supportedAPISentiment supportedAPI = "SUPPORTED_API_SENTIMENT"
	supportedAPIKeyword   supportedAPI = "SUPPORTED_API_KEYWORD"
)

type APIResponse struct {
	api  supportedAPI
	data []byte
	err  error
}

type SummaryAPIResp struct {
	Summary string `json:"summary"`
}

type SentimentAPIResp struct {
	Polarity string  `json:"polarity"`
	Score    float64 `json:"score"`
}

type Keyword struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

type KeywordAPIResp struct {
	Keywords []Keyword `json:"keywords"`
}

type AnalyzeData struct {
	OriginalText string
	Summary      *SummaryAPIResp
	Sentiment    *SentimentAPIResp
	Keywords     *KeywordAPIResp
}

func GetAnalyzeHandler(c echo.Context) error {
	var params GetAnalyzeHandlerReq
	err := c.Bind(&params)
	if err != nil {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v",
				params,
			),
		)
	}

	if params.AnalyzeText == "" {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v, analyze-text can't be empty;",
				params,
			),
		)
	}

	if params.Summarize != "on" &&
		params.Sentiment != "on" &&
		params.Keyword != "on" {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v, must select at least one analyze option;",
				params,
			),
		)
	}

	respChan := make(chan APIResponse, 3)
	wg := &sync.WaitGroup{}
	txt := url.QueryEscape(params.AnalyzeText)

	if params.Summarize == "on" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(
				fmt.Sprintf(
					"http://wordser:8080/api/v1/summary?txt=%s",
					txt,
				),
			)
			fmt.Printf("\n\n summarize resp: %v, err: %v \n\n", resp, err)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPISummary,
					err: err,
				}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				respChan <- APIResponse{
					api: supportedAPISummary,
					err: fmt.Errorf("failed to get summary; statusCode: %d", resp.StatusCode),
				}
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPISummary,
					err: err,
				}
				return
			}

			respChan <- APIResponse{
				api:  supportedAPISummary,
				data: data,
				err:  nil,
			}
		}()
	}

	if params.Sentiment == "on" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(
				fmt.Sprintf(
					"http://wordser:8080/api/v1/sentiment?txt=%s",
					txt,
				),
			)
			fmt.Printf("\n\n sentiment resp: %v, err: %v \n\n", resp, err)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPISentiment,
					err: err,
				}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				respChan <- APIResponse{
					api: supportedAPISentiment,
					err: fmt.Errorf("failed to get sentiment; statusCode: %d", resp.StatusCode),
				}
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPISentiment,
					err: err,
				}
				return
			}

			respChan <- APIResponse{
				api:  supportedAPISentiment,
				data: data,
				err:  nil,
			}
		}()
	}

	if params.Keyword == "on" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(
				fmt.Sprintf(
					"http://wordser:8080/api/v1/extract?txt=%s",
					txt,
				),
			)
			fmt.Printf("\n\n keyword resp: %v, err: %v \n\n", resp, err)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPIKeyword,
					err: err,
				}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				respChan <- APIResponse{
					api: supportedAPIKeyword,
					err: fmt.Errorf("failed to get extracted keywords; statusCode: %d", resp.StatusCode),
				}
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				respChan <- APIResponse{
					api: supportedAPIKeyword,
					err: err,
				}
				return
			}

			respChan <- APIResponse{
				api:  supportedAPIKeyword,
				data: data,
				err:  nil,
			}
		}()
	}

	wg.Wait()
	close(respChan)

	analysisData := AnalyzeData{
		OriginalText: params.AnalyzeText,
	}

	for resp := range respChan {
		fmt.Printf("\n\n %v \n\n", resp)
		if resp.err != nil {
			fmt.Printf(
				"failed get response from api: %v;",
				resp.api,
			)
			return c.String(
				http.StatusInternalServerError,
				"failed to unmarshal Summary API response",
			)
		}
		switch resp.api {
		case supportedAPISummary:
			analysisData.Summary = &SummaryAPIResp{}
			if err := json.Unmarshal(resp.data, analysisData.Summary); err != nil {
				fmt.Printf(
					"failed to unmarshal Summary API response: %v;",
					resp.data,
				)
				return c.String(
					http.StatusInternalServerError,
					"failed to unmarshal Summary API response",
				)
			}
		case supportedAPISentiment:
			analysisData.Sentiment = &SentimentAPIResp{}
			if err := json.Unmarshal(resp.data, analysisData.Sentiment); err != nil {
				fmt.Printf(
					"failed to unmarshal Sentiment API response: %v;",
					resp.data,
				)
				return c.String(
					http.StatusInternalServerError,
					"failed to unmarshal Sentiment API response",
				)
			}
		case supportedAPIKeyword:
			analysisData.Keywords = &KeywordAPIResp{}
			if err := json.Unmarshal(resp.data, analysisData.Keywords); err != nil {
				fmt.Printf(
					"failed to unmarshal Keyword API response: %v;",
					resp.data,
				)
				return c.String(
					http.StatusInternalServerError,
					"failed to unmarshal Keyword API response",
				)
			}
		}
	}

	fmt.Printf("\n\n analysisData: %v \n\n", analysisData)

	return c.Render(http.StatusOK, "analysis", analysisData)
}
