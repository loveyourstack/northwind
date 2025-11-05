package salesordersbysalesman

import (
	"context"
	"log"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
)

const (
	name           string = "Orders by salesman"
	schemaName     string = "sales"
	tableName      string = "" // will result in no total_count in result
	setFuncName    string = "orders_by_salesman"
	defaultOrderBy string = "total_value DESC"
)

type Model struct {
	OrderCount   int     `db:"order_count" json:"order_count"`
	Salesman     string  `db:"salesman" json:"salesman"`
	ShippedCount int     `db:"shipped_count" json:"shipped_count"`
	TotalValue   float32 `db:"total_value" json:"total_value"`
}

var (
	meta lysmeta.Result
)

func init() {
	var err error
	meta, err = lysmeta.AnalyzeStructs(reflect.ValueOf(&Model{}).Elem())
	if err != nil {
		log.Fatalf("lysmeta.AnalyzeStructs failed for %s.%s: %s", schemaName, tableName, err.Error())
	}
}

type Store struct {
	Db *pgxpool.Pool
}

func (s Store) GetMeta() lysmeta.Result {
	return meta
}
func (s Store) GetName() string {
	return name
}
func (s Store) GetSetFuncUrlParamNames() []string {
	return []string{"from_date", "until_date"}
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, setFuncName, defaultOrderBy, meta.DbTags, params)
}
