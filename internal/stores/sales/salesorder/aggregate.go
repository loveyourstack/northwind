package salesorder

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	orderValueLatestWeeksViewName string = "v_order_value_latest_weeks"
)

type OrderValueLatestWeeksModel struct {
	OrderWeek  int     `db:"order_week" json:"order_week,omitempty"`
	TotalValue float32 `db:"total_value" json:"total_value"`
}

func (s Store) SelectOrderValueLatestWeeks(ctx context.Context) (items []OrderValueLatestWeeksModel, stmt string, err error) {

	stmt = fmt.Sprintf("SELECT * FROM %s.%s;", schemaName, orderValueLatestWeeksViewName)
	rows, _ := s.Db.Query(ctx, stmt)
	items, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[OrderValueLatestWeeksModel])
	if err != nil {
		return nil, stmt, fmt.Errorf("pgx.CollectRows failed: %w", err)
	}

	return items, "", nil
}
