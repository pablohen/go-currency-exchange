package treasury

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RateResponse struct {
	Data []Rate `json:"data"`
}

type Rate struct {
	CountryCurrencyDesc string `json:"country_currency_desc"`
	ExchangeRate        string `json:"exchange_rate"`
	RecordDate          string `json:"record_date"`
}

var ErrDateFuture = fmt.Errorf("date cannot be in the future")

func GetRatesByDate(date time.Time) (Rate, error) {
	if date.After(time.Now()) {
		return Rate{}, ErrDateFuture
	}

	dateStart := date.AddDate(0, -6, -date.Day()+1)
	formattedDateStart := dateStart.Format("2006-01-02")

	url := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?fields=country_currency_desc,exchange_rate,record_date&filter=country_currency_desc:in:(Brazil-Real),record_date:gte:" + formattedDateStart + "&sort=-record_date"
	resp, err := http.Get(url)
	if err != nil {
		return Rate{}, err
	}
	defer resp.Body.Close()

	var rates RateResponse
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		return Rate{}, err
	}

	return rates.Data[0], nil
}
