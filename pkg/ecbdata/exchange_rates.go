package ecbdata

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ExchangeRate struct {
	FromCurr  string // base currency
	ToCurr    string
	Freq      Frequency
	PeriodStr string // daily: YYYY-MM-DD, monthly: YYYY-MM
	Rate      float32
}

// GetExchangeRates returns average daily or monthly exchange rates from baseCurr to all other available currencies
func GetExchangeRates(baseCurr string, freq Frequency, startDate, endDate time.Time) (exRates []ExchangeRate, err error) {

	// validate dates
	if startDate.After(time.Now()) {
		return nil, fmt.Errorf("startDate must be before now")
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("startDate must be before endDate")
	}
	if endDate.After(time.Now()) {
		return nil, fmt.Errorf("endDate must be before now")
	}

	// set vars depending on freq
	var dateFormat string
	switch freq {
	case Daily:
		dateFormat = "2006-01-02"
	case Monthly:
		dateFormat = "2006-01"
	default:
		return nil, fmt.Errorf("invalid freq '%s'", freq)
	}

	// build URL
	exrBaseUrl := "https://data-api.ecb.europa.eu/service/data/EXR"
	path := fmt.Sprintf("/%s..%s.SP00.A", freq, baseCurr)
	params := url.Values{}
	params.Add("detail", "dataonly")
	params.Add("format", "csvdata")
	params.Add("startPeriod", startDate.Format(dateFormat))
	params.Add("endPeriod", endDate.Format(dateFormat))
	exrUrl := exrBaseUrl + path + "?" + params.Encode()

	// get rates
	resp, err := http.Get(exrUrl)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed: %w", err)
	}
	defer resp.Body.Close()

	// read csv content
	csvContent, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv.NewReader().ReadAll failed: %w", err)
	}

	if len(csvContent) < 2 {
		return nil, fmt.Errorf("no rates found for these params")
	}

	/* csvContent looks like this:
	KEY,FREQ,CURRENCY,CURRENCY_DENOM,EXR_TYPE,EXR_SUFFIX,TIME_PERIOD,OBS_VALUE
	EXR.D.AUD.EUR.SP00.A,D,AUD,EUR,SP00,A,2024-09-02,1.6322
	EXR.D.AUD.EUR.SP00.A,D,AUD,EUR,SP00,A,2024-09-03,1.6394
	*/

	// for each line
	for i, lineA := range csvContent {

		// skip header
		if i == 0 {
			continue
		}

		// parse out the values
		exRate := ExchangeRate{
			FromCurr:  baseCurr,
			ToCurr:    lineA[2],
			Freq:      freq,
			PeriodStr: lineA[6],
		}

		rateFl64, err := strconv.ParseFloat(lineA[7], 32)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseFloat failed for rate '%s': %w", lineA[7], err)
		}
		exRate.Rate = float32(rateFl64)

		exRates = append(exRates, exRate)
	}

	return exRates, nil
}
