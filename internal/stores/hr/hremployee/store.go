package hremployee

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
	name           string = "Employees"
	schemaName     string = "hr"
	tableName      string = "employee"
	viewName       string = "v_employee"
	pkColName      string = "id"
	defaultOrderBy string = "last_name"
)

type Input struct {
	Address     string           `db:"address" json:"address,omitempty" validate:"required"`
	City        string           `db:"city" json:"city,omitempty" validate:"required"`
	CountryFk   int64            `db:"country_fk" json:"country_fk,omitempty" validate:"required"`
	DateOfBirth lystype.Date     `db:"date_of_birth" json:"date_of_birth,omitzero" validate:"required"`
	FirstName   string           `db:"first_name" json:"first_name,omitempty" validate:"required"`
	HireDate    lystype.Date     `db:"hire_date" json:"hire_date,omitzero" validate:"required"`
	HomePhone   string           `db:"home_phone" json:"home_phone,omitempty"`
	JobTitle    string           `db:"job_title" json:"job_title,omitempty" validate:"required"`
	LastName    string           `db:"last_name" json:"last_name,omitempty" validate:"required"`
	Name        string           `db:"name" json:"name,omitempty" validate:"required"`
	Notes       string           `db:"notes" json:"notes,omitempty"`
	PostalCode  string           `db:"postal_code" json:"postal_code,omitempty" validate:"required"`
	ReportsToFk int64            `db:"reports_to_fk" json:"reports_to_fk,omitempty" validate:"required"`
	State       string           `db:"state" json:"state,omitempty"`
	Title       string           `db:"title" json:"title,omitempty" validate:"required"`
	UpdatedAt   lystype.Datetime `db:"updated_at" json:"updated_at,omitzero"` // assigned in Update funcs
}

type Model struct {
	Id          int64            `db:"id" json:"id"`
	Age         int              `db:"age" json:"age,omitempty"`
	Country     string           `db:"country" json:"country,omitempty"`
	CountryIso2 string           `db:"country_iso2" json:"country_iso2,omitempty"`
	CreatedAt   lystype.Datetime `db:"created_at" json:"created_at,omitzero"`
	CreatedBy   string           `db:"created_by" json:"created_by,omitempty"`
	ReportsTo   string           `db:"reports_to" json:"reports_to,omitempty"`
	UpdatedBy   string           `db:"updated_by" json:"updated_by,omitempty"`
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
	input.UpdatedAt = lystype.Datetime(time.Now())
	return lyspg.Update(ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	assignmentsMap["updated_at"] = lystype.Datetime(time.Now())
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
