package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type TranslateLanguage string

const (
	Arabic    TranslateLanguage = "ar"
	Chinese   TranslateLanguage = "zh"
	English   TranslateLanguage = "en"
	French    TranslateLanguage = "fr"
	German    TranslateLanguage = "de"
	Greek     TranslateLanguage = "el"
	Italian   TranslateLanguage = "it"
	Portugese TranslateLanguage = "pt"
	Spanish   TranslateLanguage = "es"
	Russian   TranslateLanguage = "ru"
)

func (t TranslateLanguage) String() string {
	return string(t)
}

func (t TranslateLanguage) ToLanguageName() (string, error) {
	switch t {
	case Arabic:
		return "Arabic", nil
	case Chinese:
		return "Chinese", nil
	case English:
		return "English", nil
	case French:
		return "French", nil
	case German:
		return "German", nil
	case Greek:
		return "Greek", nil
	case Italian:
		return "Italian", nil
	case Portugese:
		return "Portugese", nil
	case Spanish:
		return "Spanish", nil
	case Russian:
		return "Russian", nil
	default:
		return "", fmt.Errorf(
			"invalid parameters: %v; source-language not supported;",
			t,
		)
	}
}

type TranslateResp struct {
	TranslatedText string `json:"translated_text"`
}

func GetDashboardHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", "")
}

type GetTranslateHandlerReq struct {
	TranslateText  string            `query:"translate-text"`
	SourceLanguage TranslateLanguage `query:"source-language"`
	TargetLanguage TranslateLanguage `query:"target-language"`
}

func GetTranslateHandler(c echo.Context) error {
	var params GetTranslateHandlerReq
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

	if params.TranslateText == "" {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v, translate-text can't be empty;",
				params,
			),
		)
	}

	sourceLangName, err := params.SourceLanguage.ToLanguageName()
	if err != nil {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v; source-language not supported;",
				params,
			),
		)
	}

	targLangName, err := params.TargetLanguage.ToLanguageName()
	if err != nil {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v; target-language not supported;",
				params,
			),
		)
	}

	if params.SourceLanguage == params.TargetLanguage {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v; support-language cannot equal target-language;",
				params,
			),
		)
	}

	// TODO: remove me. this is partner's micorservice.
	// after course just use wordser translate
	resp, err := http.Get(
		fmt.Sprintf(
			"http://backend:5000/translate?inputText=%s&sourceLanguage=%s&targetLanguage=%s",
			url.QueryEscape(params.TranslateText),
			url.QueryEscape(params.SourceLanguage.String()),
			url.QueryEscape(params.TargetLanguage.String()),
		),
	)
	fmt.Printf("\n\n resp: %v, err: %v \n\n", resp, err)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get translation from partner's backend; statusCode: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	transResp := &TranslateResp{}
	json.Unmarshal(data, transResp)

	fmt.Printf("\n translated text: %v \n", transResp)

	return c.HTML(
		http.StatusOK,
		fmt.Sprintf(
			`
			<div class="card" style="width: 18rem;">
				<div class="card-body">
					<h5 class="card-title">%s -> %s</h5>
					<h6 class="card-subtitle mb-2 text-muted">Original Text: %s</h6>
					<p class="card-text font-weight-bold">Translated Text: %s</p>
				</div>
			</div>
			`,
			sourceLangName,
			targLangName,
			params.TranslateText,
			transResp.TranslatedText,
		),
	)
}
