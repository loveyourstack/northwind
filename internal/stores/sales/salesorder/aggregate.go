package salesorder

import (
	"context"
	"fmt"

	"github.com/loveyourstack/lys/lyspg"
)

const (
	orderValueLatestWeeksViewName string = "v_order_value_latest_weeks"
)

type OrderValueLatestWeeksModel struct {
	OrderWeek  int     `db:"order_week" json:"order_week"`
	TotalValue float32 `db:"total_value" json:"total_value"`
}

func (s Store) SelectOrderValueLatestWeeks(ctx context.Context) (items []OrderValueLatestWeeksModel, stmt string, err error) {
	return lyspg.SelectT[OrderValueLatestWeeksModel](ctx, s.Db, fmt.Sprintf("SELECT * FROM %s.%s;", schemaName, orderValueLatestWeeksViewName))
}
