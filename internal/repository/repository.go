package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/azagrivin/testProjectCourses/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db    *sqlx.DB
	codes map[string]uint32
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:    db,
		codes: loadCodes(db),
	}
}

func loadCodes(db *sqlx.DB) map[string]uint32 {
	var codes []struct {
		Id    uint32 `db:"id"`
		Value string `db:"value"`
	}

	err := db.Select(&codes, `SELECT id, value FROM codes`)
	if err != nil {
		log.Fatalf("get codes from db error, %v", err)
	}

	result := make(map[string]uint32, len(codes))
	for _, c := range codes {
		result[c.Value] = c.Id
	}

	return result
}

func (r *Repository) StoreBtcusdt(ctx context.Context, in *models.Btcusdt) (uint32, error) {
	rows, err := r.db.NamedQueryContext(ctx, `
INSERT INTO btc_usdt(timestamp, buy, sell, high, low, last, average_price)
VALUES (:timestamp, :buy, :sell, :high, :low, :last, :average_price)
ON CONFLICT (timestamp)
    DO NOTHING
RETURNING id
`, in)
	if err != nil {
		return 0, err
	}

	var id uint32

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
func (r *Repository) HistoryBtcusdt(ctx context.Context, f *models.FilterBtcusdt) (result *models.HistoryBtcusdt, err error) {
	result = &models.HistoryBtcusdt{}

	if f != nil && !f.FirstOnly {
		err = r.db.Get(&result.Total, `SELECT COUNT(DISTINCT timestamp) FROM btc_usdt`)
		if err != nil {
			return nil, fmt.Errorf("get btcusdt total error, %w", err)
		}
	}

	args := map[string]interface{}{}
	query := strings.Builder{}
	{
		query.WriteString(`
SELECT timestamp, buy, sell, high, low, last, average_price
FROM btcusdt
WHERE true`)
		if f != nil {
			if f.TimeFrom != nil {
				query.WriteString(`
  AND timestamp >= :time_from`)
				args["time_from"] = f.TimeFrom
			}
			if f.TimeTo != nil {
				query.WriteString(`
  AND timestamp <= :time_to`)
				args["time_to"] = f.TimeTo
			}
		}

		query.WriteString(`
ORDER BY timestamp DESC
`)

		if f != nil && f.Limit != 0 {
			query.WriteString(fmt.Sprintf(`
OFFSET %d
LIMIT %d`, f.Page*f.Limit, f.Limit))
		}
	}

	rows, err := r.db.NamedQueryContext(ctx, query.String(), args)
	if err != nil {
		return nil, fmt.Errorf("get btcusdt history error, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h models.Btcusdt

		if err = rows.StructScan(&h); err != nil {
			return nil, fmt.Errorf("scan btsusdt error, %w", err)
		}

		result.History = append(result.History, &h)
	}

	return result, nil
}

func (r *Repository) StoreCurrencies(ctx context.Context, in []*models.Currency) error {

	for i := range in {
		codeID, ok := r.codes[in[i].CodeStr]
		if !ok {
			err := r.db.Get(&codeID, `
INSERT INTO codes(Value)
VALUES($1)
RETURNING Id`, in[i].CodeStr)
			if err != nil {
				return fmt.Errorf("insert code error, %v", err)
			}

			r.codes[in[i].CodeStr] = codeID
		}

		in[i].CodeID = codeID
	}

	_, err := r.db.NamedExecContext(ctx, `
INSERT INTO currencies(code_id, date, value)
VALUES (:code_id, date(:time), :value)
ON CONFLICT (code_id, date) 
    DO UPDATE 
       SET value = excluded.value`, in)

	return err
}
func (r *Repository) HistoryCurrency(ctx context.Context, f *models.FilterCurrencies) (result *models.HistoryCurrencies, err error) {

	result = &models.HistoryCurrencies{}

	if f != nil && !f.FirstOnly {
		err = r.db.Get(&result.Total, `SELECT COUNT(DISTINCT date) FROM currencies`)
		if err != nil {
			return nil, fmt.Errorf("get currencies total error, %w", err)
		}
	}

	args := map[string]interface{}{}
	query := strings.Builder{}
	{
		query.WriteString(`
WITH t AS (
    SELECT c.date, jsonb_object_agg(codes.value, c.value) objects
    FROM currencies c
             JOIN codes ON c.code_id = codes.id
    WHERE true`)

		if f != nil {
			if f.TimeFrom != nil {
				query.WriteString(`
      AND c.date >= DATE(:date_from)`)
				args["date_from"] = f.TimeFrom
			}
			if f.TimeTo != nil {
				query.WriteString(`
      AND c.date <= DATE(:date_to)`)
				args["date_to"] = f.TimeTo
			}

			if len(f.Currencies) != 0 {
				query.WriteString(`
      AND codes.value = ANY(:currencies)`)
				args["currencies"] = pq.Array(f.Currencies)
			}
		}

		query.WriteString(`
    GROUP BY c.date
    ORDER BY c.date DESC`)

		if f != nil && f.Limit != 0 {
			query.WriteString(fmt.Sprintf(`
    OFFSET %d
    LIMIT %d`, f.Page*f.Limit, f.Limit))
		}

		query.WriteString(`
)
SELECT ('{"date"::"' || TO_CHAR(t.date, 'YYYY-MM-DD') || '"}') :::: jsonb || t.objects :::: jsonb objects
FROM t`)
	}

	rows, err := r.db.NamedQueryContext(ctx, query.String(), args)
	if err != nil {
		return nil, fmt.Errorf("get currencies history error, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h models.History

		if err = rows.Scan(&h); err != nil {
			return nil, fmt.Errorf("scan currency error, %w", err)
		}

		result.History = append(result.History, h)
	}

	return result, nil
}

func (r *Repository) StoreBtc(ctx context.Context, in []*models.BtcDB) error {

	for i := range in {
		codeID, ok := r.codes[in[i].CodeStr]
		if !ok {
			err := r.db.Get(&codeID, `
INSERT INTO codes(Value)
VALUES($1)
RETURNING Id`, in[i].CodeStr)
			if err != nil {
				return fmt.Errorf("insert code error, %v", err)
			}

			r.codes[in[i].CodeStr] = codeID
		}

		in[i].CodeID = codeID
	}

	_, err := r.db.NamedExecContext(ctx, `
INSERT INTO btc(code_id, timestamp, buy, sell, high, low, last, average_price)
VALUES (:code_id, :timestamp, :buy, :sell, :high, :low, :last, :average_price)
`, in)
	if err != nil {
		return fmt.Errorf("insert btc error, %w", err)
	}

	return nil
}

func (r *Repository) HistoryBtc(ctx context.Context, f *models.FilterCurrencies) (result *models.HistoryBtc, err error) {
	result = &models.HistoryBtc{}

	if f != nil && !f.FirstOnly {
		err = r.db.Get(&result.Total, `SELECT COUNT(DISTINCT timestamp) FROM btc`)
		if err != nil {
			return nil, fmt.Errorf("get btc total error, %w", err)
		}
	}

	args := map[string]interface{}{}
	query := strings.Builder{}
	{

		query.WriteString(`
WITH t AS (
    SELECT timestamp
    FROM btc
    WHERE true`)

		if f != nil {
			if f.TimeFrom != nil {
				query.WriteString(`
  AND timestamp >= :time_from`)
				args["time_from"] = f.TimeFrom
			}
			if f.TimeTo != nil {
				query.WriteString(`
  AND timestamp <= :time_to`)
				args["time_to"] = f.TimeTo
			}
		}

		query.WriteString(`
    GROUP BY timestamp
    ORDER BY timestamp DESC
`)

		if f != nil && f.Limit != 0 {
			query.WriteString(fmt.Sprintf(`
OFFSET %d
LIMIT %d`, f.Page*f.Limit, f.Limit))
		}

		query.WriteString(`
)
SELECT timestamp, code_id, c.value code, buy, sell, high, low, last, average_price
FROM btc
    JOIN codes c ON btc.code_id = c.id
WHERE timestamp = ANY (SELECT timestamp FROM t)`)

		if len(f.Currencies) != 0 {
			query.WriteString(`
  AND c.value = ANY(:currencies)`)
			args["currencies"] = pq.Array(f.Currencies)
		}

		query.WriteString(`
ORDER BY timestamp DESC, c.value`)
	}

	fmt.Println(query.String())
	fmt.Println(args)

	rows, err := r.db.NamedQueryContext(ctx, query.String(), args)
	if err != nil {
		return nil, fmt.Errorf("get btc history error, %w", err)
	}
	defer rows.Close()

	result.History = make(map[string]map[string]*models.BtcDB)
	for rows.Next() {
		res := &models.BtcDB{}

		if err = rows.StructScan(res); err != nil {
			return nil, fmt.Errorf("scan btc error, %w", err)
		}

		t := res.Time.Format(models.BtcusdtTmeFormat)

		byTimestamp, ok := result.History[t]
		if !ok {
			result.History[t] = make(map[string]*models.BtcDB)
			byTimestamp = result.History[t]
		}

		byTimestamp[res.CodeStr] = res
	}

	return result, nil
}
