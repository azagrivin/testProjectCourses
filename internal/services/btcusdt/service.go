package btcusdt

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/azagrivin/testProjectCourses/internal/models"
	"github.com/azagrivin/testProjectCourses/internal/repository"
)

const (
	codeStrUSD = "USD"
	codeStrRUB = "RUB"
)

type msgError struct {
	r io.Reader
}

func (e msgError) String() string {
	body, err := ioutil.ReadAll(e.r)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

type Service struct {
	repo *repository.Repository

	lastCurrencies  []*models.Currency
	lastCurrencyUSD *models.Currency
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetBtcusdt(ctx context.Context) error {
	btcsudtDB, err := kucoinStats24hr()
	if err != nil {
		return fmt.Errorf("get btcusdt error, %w", err)
	}

	id, err := s.repo.StoreBtcusdt(ctx, btcsudtDB)
	if err != nil {
		return fmt.Errorf("store btcusdt error, %w", err)
	}

	if id != 0 {
		for s.lastCurrencies == nil {
		}

		res := make([]*models.BtcDB, 0, len(s.lastCurrencies)+1)
		btcInUSD := &models.BtcDB{
			CodeStr:      codeStrUSD,
			Time:         btcsudtDB.Time,
			Buy:          btcsudtDB.Buy,
			Sell:         btcsudtDB.Sell,
			High:         btcsudtDB.High,
			Low:          btcsudtDB.Low,
			Last:         btcsudtDB.Last,
			AveragePrice: btcsudtDB.AveragePrice,
		}
		btcInRUB := &models.BtcDB{
			CodeStr:      codeStrRUB,
			Time:         btcsudtDB.Time,
			Buy:          s.lastCurrencyUSD.Value * btcsudtDB.Buy,
			Sell:         s.lastCurrencyUSD.Value * btcsudtDB.Sell,
			High:         s.lastCurrencyUSD.Value * btcsudtDB.High,
			Low:          s.lastCurrencyUSD.Value * btcsudtDB.Low,
			Last:         s.lastCurrencyUSD.Value * btcsudtDB.Last,
			AveragePrice: s.lastCurrencyUSD.Value * btcsudtDB.AveragePrice,
		}

		res = append(res, btcInRUB, btcInUSD)

		for _, c := range s.lastCurrencies {
			if c.CodeStr == codeStrUSD {
				continue
			}

			res = append(res, &models.BtcDB{
				CodeStr:      c.CodeStr,
				Time:         btcsudtDB.Time,
				Buy:          btcInRUB.Buy / c.Value,
				Sell:         btcInRUB.Sell / c.Value,
				High:         btcInRUB.High / c.Value,
				Low:          btcInRUB.Low / c.Value,
				Last:         btcInRUB.Last / c.Value,
				AveragePrice: btcInRUB.AveragePrice / c.Value,
			})
		}

		if err = s.repo.StoreBtc(ctx, res); err != nil {
			return fmt.Errorf("store btc error, %w", err)
		}
	}

	return nil
}
func (s *Service) HistoryBtcusdt(ctx context.Context, f *models.FilterBtcusdt) (*models.HistoryBtcusdt, error) {
	return s.repo.HistoryBtcusdt(ctx, f)
}

func (s *Service) GetCurrencies(ctx context.Context, t time.Time) error {
	cbrCurrency, err := cbrDaily(ctx, t)
	if err != nil {
		return fmt.Errorf("get currenies error, %w", err)
	}

	cur := make([]*models.Currency, len(cbrCurrency))
	for i := range cbrCurrency {
		cbrCurrency[i].Value = strings.Replace(cbrCurrency[i].Value, ",", ".", 1)
		value, err := strconv.ParseFloat(cbrCurrency[i].Value, 64)
		if err != nil {
			return fmt.Errorf("convert value of cyrrency error, %w", err)
		}

		cur[i] = &models.Currency{
			CodeStr: cbrCurrency[i].CharCode,
			Value:   value / float64(cbrCurrency[i].Nom),
			Time:    t,
		}

		if cur[i].CodeStr == codeStrUSD {
			s.lastCurrencyUSD = cur[i]
		}
	}

	s.lastCurrencies = cur

	if err = s.repo.StoreCurrencies(ctx, cur); err != nil {
		return fmt.Errorf("store cyrrencies error, %w", err)
	}

	return nil
}
func (s *Service) HistoryCurrencies(ctx context.Context, f *models.FilterCurrencies) (*models.HistoryCurrencies, error) {
	return s.repo.HistoryCurrency(ctx, f)
}

func (s *Service) HistoryBtc(ctx context.Context, f *models.FilterCurrencies) (*models.HistoryBtc, error) {
	return s.repo.HistoryBtc(ctx, f)
}
