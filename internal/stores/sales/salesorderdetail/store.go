package salesorderdetail

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
)

const (
	schemaName     string = "sales"
	tableName      string = "order_detail"
	viewName       string = "v_order_detail"
	pkColName      string = "id"
	defaultOrderBy string = "id"
)

type Input struct {
	Discount       float32          `db:"discount" json:"discount,omitempty"`
	LastModifiedAt lystype.Datetime `db:"last_modified_at" json:"last_modified_at,omitempty"` // assigned in Update funcs
	OrderFk        int64            `db:"order_fk" json:"order_fk,omitempty" validate:"required"`
	ProductFk      int64            `db:"product_fk" json:"product_fk,omitempty" validate:"required"`
	Quantity       int32            `db:"quantity" json:"quantity,omitempty" validate:"required"`
	UnitPrice      float32          `db:"unit_price" json:"unit_price,omitempty" validate:"required"`
}

type Model struct {
	Id          int64            `db:"id" json:"id"`
	EntryAt     lystype.Datetime `db:"entry_at" json:"entry_at,omitempty"`
	OrderNumber int32            `db:"order_number" json:"order_number,omitempty"`
	ProductName string           `db:"product_name" json:"product_name,omitempty"`
	Input
}

var (
	gDbTags      []string
	gJsonTags    []string
	gInputDbTags []string
)

func init() {
	var err error
	gDbTags, gJsonTags, err = lysmeta.GetStructTags(reflect.ValueOf(&Input{}).Elem(), reflect.ValueOf(&Model{}).Elem())
	if err != nil {
		log.Fatalf("lysmeta.GetStructTags failed for %s.%s: %s", schemaName, tableName, err.Error())
	}
	gInputDbTags, _, _ = lysmeta.GetStructTags(reflect.ValueOf(&Input{}).Elem())
}

type Store struct {
	Db *pgxpool.Pool
}

func (s Store) Delete(ctx context.Context, id int64) (stmt string, err error) {
	return lyspg.DeleteUnique(ctx, s.Db, schemaName, tableName, pkColName, id)
}

func (s Store) GetJsonFields() []string {
	return gJsonTags
}

func (s Store) Insert(ctx context.Context, input Input) (newItem Model, stmt string, err error) {
	return lyspg.Insert[Input, Model](ctx, s.Db, schemaName, tableName, viewName, pkColName, gDbTags, input)
}

func (s Store) Restore(ctx context.Context, tx pgx.Tx, id int64) (stmt string, err error) {
	return lyspg.Restore(ctx, tx, schemaName, tableName, pkColName, id, false)
}

func (s Store) RestoreCascadedByOrder(ctx context.Context, tx pgx.Tx, orderId int64) (stmt string, err error) {
	return lyspg.Restore(ctx, tx, schemaName, tableName, "order_fk", orderId, true)
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, stmt string, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, gDbTags, params)
}

func (s Store) SelectById(ctx context.Context, fields []string, id int64) (item Model, stmt string, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, fields, gDbTags, id)
}

func (s Store) SoftDelete(ctx context.Context, tx pgx.Tx, id int64) (stmt string, err error) {
	return lyspg.SoftDelete(ctx, tx, schemaName, tableName, pkColName, id, false)
}

func (s Store) SoftDeleteCascadedByOrder(ctx context.Context, tx pgx.Tx, orderId int64) (stmt string, err error) {
	return lyspg.SoftDelete(ctx, tx, schemaName, tableName, "order_fk", orderId, true)
}

func (s Store) Update(ctx context.Context, input Input, id int64) (stmt string, err error) {
	input.LastModifiedAt = lystype.Datetime(time.Now())
	return lyspg.Update[Input](ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) (stmt string, err error) {
	assignmentsMap["last_modified_at"] = lystype.Datetime(time.Now())
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, gInputDbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
