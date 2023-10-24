package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func EvalStock(key string) string {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", url.QueryEscape(key))
	response, err := http.Get(url)
	
	if err != nil {
		log.Println("error :", err)
		return "Stock service is not available"
	}

	if response.StatusCode == http.StatusOK {
		content, err := csv.NewReader(response.Body).ReadAll()

		if err != nil {
			log.Println("error :", err)
			return "Stock service CSV error"
		}

		token := content[1][0]
		close  := content[1][6]

		if close == "N/D" {
			return fmt.Sprintf("%s quote is not available", strings.ToUpper(token))
		}

		return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(token), close)
	}

	log.Println("error : response.StatusCode is ", response.StatusCode)

	return "Stock service is not available"
}