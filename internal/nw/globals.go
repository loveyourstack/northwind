package nw

import (
	"log"
	"time"

	"github.com/loveyourstack/lys/lystype"
)

const CommencementOfTradingDateStr string = "2010-01-01"

var CommencementOfTradingDate time.Time

func init() {
	var err error
	CommencementOfTradingDate, err = time.Parse(lystype.DateFormat, CommencementOfTradingDateStr)
	if err != nil {
		log.Fatalf("initialization: time.Parse of %s failed: "+err.Error(), CommencementOfTradingDateStr)
	}
}
