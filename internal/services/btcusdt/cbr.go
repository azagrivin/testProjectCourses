package btcusdt

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/corpix/uarand"
	"golang.org/x/text/encoding/charmap"
)

const (
	dateFormat = "02/01/2006"
	cbrRequest = "http://www.cbr.ru/scripts/XML_daily.asp"
)

type currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  uint64 `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nom      uint64 `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type cbrDailyResponse struct {
	XMLName    xml.Name    `xml:"ValCurs"`
	Date       string      `xml:"Date,attr"`
	Currencies []*currency `xml:"Valute"`
}

func cbrDaily(ctx context.Context, t time.Time) ([]*currency, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cbrRequest, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error, %w", err)
	}

	q := req.URL.Query()
	q.Set("date_req", t.Format(dateFormat))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request error, %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is NOT 200, response=%s", msgError{res.Body})
	}

	var responseModel cbrDailyResponse

	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = decoder.Decode(&responseModel)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response model error, %w", err)
	}

	return responseModel.Currencies, nil
}
