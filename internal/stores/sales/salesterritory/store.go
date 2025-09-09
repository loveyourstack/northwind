package salesterritory

import (
	"context"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
	"github.com/loveyourstack/northwind/internal/enums/salesregion"
)

const (
	name           string = "Territories"
	schemaName     string = "sales"
	tableName      string = "territory"
	viewName       string = "v_territory"
	pkColName      string = "id"
	defaultOrderBy string = "name"
)

type Input struct {
	Code       string           `db:"code" json:"code,omitempty" validate:"required"`
	CreatedBy  string           `db:"created_by" json:"created_by,omitempty"` // omitted from Update, assigned in Insert func
	SalesmanFk int64            `db:"salesman_fk" json:"salesman_fk,omitempty" validate:"required"`
	Name       string           `db:"name" json:"name,omitempty" validate:"required"`
	Region     salesregion.Enum `db:"region" json:"region,omitempty" validate:"required"`
	UpdatedBy  string           `db:"updated_by" json:"updated_by,omitempty"` // assigned in Update funcs
}

type Model struct {
	Id        int64            `db:"id" json:"id"`
	CreatedAt lystype.Datetime `db:"created_at" json:"created_at,omitzero"`
	Salesman  string           `db:"salesman" json:"salesman"`
	UpdatedAt lystype.Datetime `db:"updated_at" json:"updated_at,omitzero"` // assigned by trigger
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

func (s Store) Delete(ctx context.Context, id int64) error {
	return lyspg.DeleteUnique(ctx, s.Db, schemaName, tableName, pkColName, id)
}

func (s Store) GetMeta() lysmeta.Result {
	return meta
}
func (s Store) GetName() string {
	return name
}

func (s Store) Insert(ctx context.Context, input Input) (newId int64, err error) {
	input.CreatedBy = lys.GetUserNameFromCtx(ctx, "Unknown")
	return lyspg.Insert[Input, int64](ctx, s.Db, schemaName, tableName, pkColName, input)
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, id int64) (item Model, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, id)
}

func (s Store) Update(ctx context.Context, input Input, id int64) error {
	input.UpdatedBy = lys.GetUserNameFromCtx(ctx, "Unknown")
	return lyspg.Update(ctx, s.Db, schemaName, tableName, pkColName, input, id, lyspg.UpdateOption{OmitFields: []string{"entry_by"}})
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	assignmentsMap["updated_by"] = lys.GetUserNameFromCtx(ctx, "Unknown")
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
