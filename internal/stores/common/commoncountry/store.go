package commoncountry

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
	"golang.org/x/exp/maps"
)

const (
	name           string = "Countries"
	schemaName     string = "common"
	tableName      string = "country"
	viewName       string = "country"
	pkColName      string = "id"
	defaultOrderBy string = "name"
)

type Input struct {
	IsActive       bool             `db:"is_active" json:"is_active"`
	Iso2           string           `db:"iso2" json:"iso2,omitempty" validate:"required,len=2"`
	LastModifiedAt lystype.Datetime `db:"last_modified_at" json:"last_modified_at,omitzero"` // assigned in Update funcs
	Name           string           `db:"name" json:"name,omitempty" validate:"required"`
}

type Model struct {
	Id      int64            `db:"id" json:"id"`
	EntryAt lystype.Datetime `db:"entry_at" json:"entry_at,omitzero"`
	Input
}

var (
	meta, inputMeta lysmeta.Result
)

func init() {
	var err error
	meta, err = lysmeta.AnalyzeStructs(reflect.ValueOf(&Input{}).Elem(), reflect.ValueOf(&Model{}).Elem())
	if err != nil {
		log.Fatalf("lysmeta.AnalyzeStructs failed for %s.%s: %s", schemaName, tableName, err.Error())
	}
	inputMeta, _ = lysmeta.AnalyzeStructs(reflect.ValueOf(&Input{}).Elem())
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

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, id int64) (item Model, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {

	// only allowed to update is_active

	if len(assignmentsMap) != 1 {
		return fmt.Errorf("assignmentsMap length must be 1")
	}

	if maps.Keys(assignmentsMap)[0] != "is_active" {
		return fmt.Errorf(`assignmentsMap key must be "is_active"`)
	}

	assignmentsMap["last_modified_at"] = lystype.Datetime(time.Now())

	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}
