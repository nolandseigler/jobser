package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetAnalyzeHandlerReq struct {
	AnalyzeText string `query:"analyze-text"`
	Summarize   string `query:"summarize"`
	Sentiment   string `query:"sentiment"`
	Keyword     string `query:"keyword"`
	Unmask      string `query:"unmask"`
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
		params.Keyword != "on" &&
		params.Unmask != "on" {
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf(
				"invalid parameters: %v, must select at least one analyze option;",
				params,
			),
		)
	}

	// if unmask is set we must do unmask synchornous first then use that for the rest of the calls.

	// resp, err := http.Get(
	// 	fmt.Sprintf(
	// 		"http://backend:5000/translate?inputText=%s&sourceLanguage=%s&targetLanguage=%s",
	// 		url.QueryEscape(params.TranslateText),
	// 		url.QueryEscape(params.SourceLanguage.String()),
	// 		url.QueryEscape(params.TargetLanguage.String()),
	// 	),
	// )
	// fmt.Printf("\n\n resp: %v, err: %v \n\n", resp, err)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("failed to get translation from partner's backend; statusCode: %d", resp.StatusCode)
	// }

	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// transResp := &TranslateResp{}
	// json.Unmarshal(data, transResp)

	// fmt.Printf("\n translated text: %v \n", transResp)

	// return c.HTML(
	// 	http.StatusOK,
	// 	fmt.Sprintf(
	// 		`
	// 		<div class="card" style="width: 18rem;">
	// 			<div class="card-body">
	// 				<h5 class="card-title">%s -> %s</h5>
	// 				<h6 class="card-subtitle mb-2 text-muted">Original Text: %s</h6>
	// 				<p class="card-text font-weight-bold">Translated Text: %s</p>
	// 			</div>
	// 		</div>
	// 		`,
	// 		sourceLangName,
	// 		targLangName,
	// 		params.TranslateText,
	// 		transResp.TranslatedText,
	// 	),
	// )

	return c.Render(http.StatusOK, "dashboard", "")
}
