package http

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	er "github.com/azagrivin/testProjectCourses/internal/delivery/http/error"
	"github.com/azagrivin/testProjectCourses/internal/logger"
	"github.com/azagrivin/testProjectCourses/internal/models"
	"github.com/azagrivin/testProjectCourses/internal/services/btcusdt"
)

const (
	keyTimeFrom = "time_from"
	keyTimeTo   = "time_to"
)

func NewHandlerGetBtcusdt(svc *btcusdt.Service, log logger.HttpLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := log.WithRequest(r)

		now := time.Now().UTC()

		f := &models.FilterBtcusdt{
			TimeFrom:  &now,
			FirstOnly: true,
			Limit:     1,
		}

		log.Infof("1")
		response, err := svc.HistoryBtcusdt(ctx, f)
		if err != nil {
			log.Errorf("get last btcusdt error, %v", err)
			er.ErrInternalError.Handle(w)
			return
		}

		var resp interface{}
		{
			if response != nil && len(response.History) != 0 {
				resp = response.History[0]
			}
		}

		writeResponse(w, resp, log)
	}
}

func NewHandlerPostBtcusdt(svc *btcusdt.Service, log logger.HttpLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := log.WithRequest(r)

		var f models.FilterBtcusdt
		{
			queryValues, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				log.Errorf("parse query error, %v", err)
				er.ErrIncorrectInput.Handle(w)
				return
			}

			if values, ok := queryValues[keyTimeFrom]; ok {
				date, err := time.Parse(models.BtcusdtTmeFormat, values[0])
				if err != nil {
					log.Errorf("parse %q from query error, %v", keyTimeFrom, err)
					er.ErrIncorrectInput.Handle(w)
					return
				}

				f.TimeFrom = &date
			}

			if values, ok := queryValues[keyTimeTo]; ok {
				date, err := time.Parse(models.BtcusdtTmeFormat, values[0])
				if err != nil {
					log.Errorf("parse %q from query error, %v", keyTimeTo, err)
					er.ErrIncorrectInput.Handle(w)
					return
				}

				f.TimeTo = &date
			}

			f.Limit = defaultLimit
			if values, ok := queryValues[keyLimit]; ok {
				f.Limit, err = strconv.ParseUint(values[0], 10, 64)
				if err != nil {
					log.Errorf("parse %q from query error, %v", keyLimit, err)
					er.ErrIncorrectInput.Handle(w)
					return
				}
			}

			if values, ok := queryValues[keyPage]; ok {
				f.Page, err = strconv.ParseUint(values[0], 10, 64)
				if err != nil {
					log.Errorf("parse %q from query error, %v", keyPage, err)
					er.ErrIncorrectInput.Handle(w)
					return
				}

				if f.Page > 0 {
					f.Page--
				}
			}
		}

		response, err := svc.HistoryBtcusdt(ctx, &f)
		if err != nil {
			log.Error("get btcusdt history error, %v", err)
			er.ErrInternalError.Handle(w)
			return
		}

		writeResponse(w, response, log)
	}
}
