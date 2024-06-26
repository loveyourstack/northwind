package salesemployeeterritory

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
)

const (
	name           string = "Employee territories"
	schemaName     string = "sales"
	tableName      string = "employee_territory"
	viewName       string = "v_employee_territory"
	pkColName      string = "id"
	defaultOrderBy string = "employee_name, territory_name"
)

type Input struct {
	EmployeeFk     int64            `db:"employee_fk" json:"employee_fk,omitempty" validate:"required"`
	LastModifiedAt lystype.Datetime `db:"last_modified_at" json:"last_modified_at,omitempty"` // assigned in Update funcs
	TerritoryFk    int64            `db:"territory_fk" json:"territory_fk,omitempty" validate:"required"`
}

type Model struct {
	Id            int64            `db:"id" json:"id"`
	EmployeeName  string           `db:"employee_name" json:"employee_name,omitempty"`
	EntryAt       lystype.Datetime `db:"entry_at" json:"entry_at,omitempty"`
	TerritoryName string           `db:"territory_name" json:"territory_name,omitempty"`
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

func (s Store) Delete(ctx context.Context, id int64) (stmt string, err error) {
	return lyspg.DeleteUnique(ctx, s.Db, schemaName, tableName, pkColName, id)
}

func (s Store) GetMeta() lysmeta.Result {
	return meta
}
func (s Store) GetName() string {
	return name
}

func (s Store) Insert(ctx context.Context, input Input) (newItem Model, stmt string, err error) {
	return lyspg.Insert[Input, Model](ctx, s.Db, schemaName, tableName, viewName, pkColName, meta.DbTags, input)
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, stmt string, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, fields []string, id int64) (item Model, stmt string, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, fields, meta.DbTags, id)
}

func (s Store) Update(ctx context.Context, input Input, id int64) (stmt string, err error) {
	input.LastModifiedAt = lystype.Datetime(time.Now())
	return lyspg.Update[Input](ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) (stmt string, err error) {
	assignmentsMap["last_modified_at"] = lystype.Datetime(time.Now())
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
