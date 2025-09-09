package coreproduct

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
	"github.com/loveyourstack/northwind/internal/stores/core/corecountry"
)

const (
	name           string = "Products"
	schemaName     string = "core"
	tableName      string = "product"
	viewName       string = "v_product"
	pkColName      string = "id"
	defaultOrderBy string = "name"
)

type Input struct {
	CategoryFk      int64    `db:"category_fk" json:"category_fk" validate:"required"`
	IsDiscontinued  bool     `db:"is_discontinued" json:"is_discontinued"`
	Name            string   `db:"name" json:"name,omitempty" validate:"required"`
	QuantityPerUnit string   `db:"quantity_per_unit" json:"quantity_per_unit" validate:"required"`
	ReorderLevel    int      `db:"reorder_level" json:"reorder_level"`
	SupplierFk      int64    `db:"supplier_fk" json:"supplier_fk" validate:"required"`
	Tags            []string `db:"tags" json:"tags,omitempty"`
	UnitPrice       float32  `db:"unit_price" json:"unit_price" validate:"required"`
	UnitsInStock    int      `db:"units_in_stock" json:"units_in_stock"`
	UnitsOnOrder    int      `db:"units_on_order" json:"units_on_order"`
}

type Model struct {
	Id                   int64            `db:"id" json:"id"`
	Category             string           `db:"category" json:"category,omitempty"`
	CategoryColorHex     string           `db:"category_color_hex" json:"category_color_hex,omitempty"`
	CategoryColorIsLight bool             `db:"category_color_is_light" json:"category_color_is_light"`
	CreatedAt            lystype.Datetime `db:"created_at" json:"created_at,omitzero"`
	CreatedBy            string           `db:"created_by" json:"created_by,omitempty"`
	SupplierCompanyName  string           `db:"supplier_company_name" json:"supplier_company_name,omitempty"`
	SupplierCountryIso2  string           `db:"supplier_country_iso2" json:"supplier_country_iso2,omitempty"`
	UpdatedAt            lystype.Datetime `db:"updated_at" json:"updated_at,omitzero"` // assigned by trigger
	UpdatedBy            string           `db:"updated_by" json:"updated_by,omitempty"`
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

func (s Store) DistinctSupplierCountries(ctx context.Context) (countries []corecountry.Model, err error) {

	stmt := fmt.Sprintf(`
		SELECT co.id, co.name FROM core.country co 
		JOIN (
			SELECT DISTINCT country_fk 
			FROM core.supplier s
			JOIN %s.%s p ON p.supplier_fk = s.id
		) t1 ON co.id = t1.country_fk 
		ORDER BY co.name;`,
		schemaName, tableName)

	return lyspg.SelectT[corecountry.Model](ctx, s.Db, stmt)
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
	return lyspg.Update(ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
