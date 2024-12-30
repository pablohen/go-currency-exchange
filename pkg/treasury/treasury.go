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

var (
	ErrDateFuture         = fmt.Errorf("date cannot be in the future")
	ErrOlderThanSixMonths = fmt.Errorf("rate is older than 6 months")
)

func GetRatesByDate(date time.Time) (Rate, error) {
	if date.After(time.Now().UTC()) {
		return Rate{}, ErrDateFuture
	}

	formattedDate := date.Format("2006-01-02")

	url := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?fields=country_currency_desc,exchange_rate,record_date&filter=country_currency_desc:in:(Brazil-Real),record_date:lte:" + formattedDate + "&sort=-record_date"
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

	if len(rates.Data) == 0 {
		return Rate{}, fmt.Errorf("no rates found for date: %s", formattedDate)
	}

	firstRate := rates.Data[0]
	firstRateDate, err := time.Parse("2006-01-02", firstRate.RecordDate)
	if err != nil {
		return Rate{}, err
	}

	sixMonthsBeforeDate := date.AddDate(0, -6, -date.Day()+1)
	if firstRateDate.Before(sixMonthsBeforeDate) {
		return Rate{}, ErrOlderThanSixMonths
	}

	return firstRate, nil
}
