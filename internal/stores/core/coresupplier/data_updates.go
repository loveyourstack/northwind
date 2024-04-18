package coresupplier

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
)

type DuModel struct {
	Input
	Model
	lyspg.DataUpdateCols
}

var (
	gDuDbTags, gDuJsonTags []string
	gDuViewName            string
)

func init() {
	gDuDbTags, gDuJsonTags, _ = lysmeta.GetStructTags(reflect.ValueOf(&Input{}).Elem(), reflect.ValueOf(&Model{}).Elem())
	gDuDbTags = append(gDuDbTags, lyspg.DataUpdateColTags...)
	gDuJsonTags = append(gDuJsonTags, lyspg.DataUpdateColTags...)

	gDuViewName = viewName + lyspg.DataUpdateViewSuffix
}

type DuStore struct {
	Db *pgxpool.Pool
}

func (s DuStore) GetJsonFields() []string {
	return gDuJsonTags
}

func (s DuStore) Select(ctx context.Context, params lyspg.SelectParams) (items []DuModel, unpagedCount lyspg.TotalCount, stmt string, err error) {
	return lyspg.Select[DuModel](ctx, s.Db, schemaName, gDuViewName, gDuViewName, defaultOrderBy, gDuDbTags, params)
}

func (s DuStore) SelectById(ctx context.Context, fields []string, id int64) (item DuModel, stmt string, err error) {
	return lyspg.SelectUnique[DuModel](ctx, s.Db, schemaName, gDuViewName, pkColName, fields, gDuDbTags, id)
}
