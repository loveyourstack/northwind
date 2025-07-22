package coresupplier

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
	name           string = "Suppliers"
	schemaName     string = "core"
	tableName      string = "supplier"
	viewName       string = "v_supplier"
	pkColName      string = "id"
	defaultOrderBy string = "name"
)

type Input struct {
	Address        string           `db:"address" json:"address,omitempty"`
	City           string           `db:"city" json:"city,omitempty"`
	CompanyName    string           `db:"company_name" json:"company_name,omitempty" validate:"required"`
	ContactName    string           `db:"contact_name" json:"contact_name,omitempty" validate:"required"`
	ContactTitle   string           `db:"contact_title" json:"contact_title,omitempty"`
	CountryFk      int64            `db:"country_fk" json:"country_fk,omitempty" validate:"required"`
	LastModifiedAt lystype.Datetime `db:"last_modified_at" json:"last_modified_at,omitzero"` // assigned in Update funcs
	Phone          string           `db:"phone" json:"phone,omitempty"`
	PostalCode     string           `db:"postal_code" json:"postal_code,omitempty"`
	State          string           `db:"state" json:"state,omitempty"`
}

type Model struct {
	Id                 int64            `db:"id" json:"id"`
	ActiveProductCount int              `db:"active_product_count" json:"active_product_count"`
	Country            string           `db:"country" json:"country,omitempty"`
	CountryIso2        string           `db:"country_iso2" json:"country_iso2,omitempty"`
	EntryAt            lystype.Datetime `db:"entry_at" json:"entry_at,omitzero"`
	EntryBy            string           `db:"entry_by" json:"entry_by,omitempty"`
	LastModifiedBy     string           `db:"last_modified_by" json:"last_modified_by,omitempty"`
	Name               string           `db:"name" json:"name,omitempty"`
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
	return lyspg.Insert[Input, int64](ctx, s.Db, schemaName, tableName, pkColName, input)
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, id int64) (item Model, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, id)
}

func (s Store) Update(ctx context.Context, input Input, id int64) error {
	input.LastModifiedAt = lystype.Datetime(time.Now())
	return lyspg.Update(ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	assignmentsMap["last_modified_at"] = lystype.Datetime(time.Now())
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
