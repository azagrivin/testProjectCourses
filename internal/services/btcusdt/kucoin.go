package btcusdt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/azagrivin/testProjectCourses/internal/models"

	"github.com/corpix/uarand"
)

const (
	statsRequest     = "https://api.kucoin.com/api/v1/market/stats"
	statsSuccessCode = "200000"
)

type kucoinStatsResponse struct {
	Code string          `json:"code"`
	Data json.RawMessage `json:"data"`
}
type stats struct {
	Time             int64  `json:"time"`
	Symbol           string `json:"symbol"`
	Buy              string `json:"buy"`
	Sell             string `json:"sell"`
	ChangeRate       string `json:"changeRate"`
	ChangePrice      string `json:"changePrice"`
	High             string `json:"high"`
	Low              string `json:"low"`
	Vol              string `json:"vol"`
	VolValue         string `json:"volValue"`
	Last             string `json:"last"`
	AveragePrice     string `json:"averagePrice"`
	TakerFeeRate     string `json:"takerFeeRate"`
	MakerFeeRate     string `json:"makerFeeRate"`
	TakerCoefficient string `json:"takerCoefficient"`
	MakerCoefficient string `json:"makerCoefficient"`
}

func kucoinStats24hr() (*models.Btcusdt, error) {
	req, err := http.NewRequest(http.MethodGet, statsRequest, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error, %v", err)
	}

	q := req.URL.Query()
	q.Set("symbol", "BTC-USDT")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request error, %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is NOT 200, kucoinStatsResponse=%s", msgError{res.Body})
	}

	var responseModel kucoinStatsResponse

	err = json.NewDecoder(res.Body).Decode(&responseModel)
	if err != nil {
		return nil, fmt.Errorf("unmarshal kucoinStatsResponse model error, %v", err)
	}
	if responseModel.Code != statsSuccessCode {
		return nil, fmt.Errorf("api code is NOT %s, kucoinStatsResponse=%s", statsSuccessCode, msgError{res.Body})
	}

	responseData := &stats{}
	err = json.Unmarshal(responseModel.Data, responseData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal kucoinStatsResponse data error, %v", err)
	}

	t := time.Unix(0, responseData.Time*int64(time.Millisecond))
	buy, err := strconv.ParseFloat(responseData.Buy, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "buy" error, %w`, err)
	}
	sell, err := strconv.ParseFloat(responseData.Sell, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "sell" error, %w`, err)
	}
	high, err := strconv.ParseFloat(responseData.High, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "high" error, %w`, err)
	}
	low, err := strconv.ParseFloat(responseData.Low, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "low" error, %w`, err)
	}
	last, err := strconv.ParseFloat(responseData.Last, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "last" error, %w`, err)
	}
	averagePrice, err := strconv.ParseFloat(responseData.AveragePrice, 64)
	if err != nil {
		return nil, fmt.Errorf(`parsing "averagePrice" error, %w`, err)
	}

	return &models.Btcusdt{
		Time:         &t,
		Buy:          buy,
		Sell:         sell,
		High:         high,
		Low:          low,
		Last:         last,
		AveragePrice: averagePrice,
	}, nil
}
