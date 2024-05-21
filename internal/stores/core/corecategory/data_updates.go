package corecategory

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
	duMeta      lysmeta.Result
	gDuViewName string
)

func init() {
	duMeta, _ = lysmeta.AnalyzeStructs(reflect.ValueOf(&Input{}).Elem(), reflect.ValueOf(&Model{}).Elem())
	duMeta.DbTags = append(duMeta.DbTags, lyspg.DataUpdateColTags...)
	duMeta.JsonTags = append(duMeta.JsonTags, lyspg.DataUpdateColTags...)

	gDuViewName = viewName + lyspg.DataUpdateViewSuffix
}

type DuStore struct {
	Db *pgxpool.Pool
}

func (s DuStore) GetJsonFields() []string {
	return duMeta.JsonTags
}

func (s DuStore) Select(ctx context.Context, params lyspg.SelectParams) (items []DuModel, unpagedCount lyspg.TotalCount, stmt string, err error) {
	return lyspg.Select[DuModel](ctx, s.Db, schemaName, gDuViewName, gDuViewName, defaultOrderBy, duMeta.DbTags, params)
}

func (s DuStore) SelectById(ctx context.Context, fields []string, id int64) (item DuModel, stmt string, err error) {
	return lyspg.SelectUnique[DuModel](ctx, s.Db, schemaName, gDuViewName, pkColName, fields, duMeta.DbTags, id)
}
