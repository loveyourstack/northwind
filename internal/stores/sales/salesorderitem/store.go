package salesorderitem

import (
	"context"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
)

const (
	name           string = "Order items"
	schemaName     string = "sales"
	tableName      string = "order_item"
	viewName       string = "v_order_item"
	pkColName      string = "id"
	defaultOrderBy string = "id"
)

type Input struct {
	Discount         float32 `db:"discount" json:"discount,omitempty"`
	LastUserUpdateBy string  `db:"last_user_update_by" json:"last_user_update_by,omitempty"` // assigned in Update funcs
	OrderFk          int64   `db:"order_fk" json:"order_fk,omitempty" validate:"required"`
	ProductFk        int64   `db:"product_fk" json:"product_fk,omitempty" validate:"required"`
	Quantity         int32   `db:"quantity" json:"quantity,omitempty" validate:"required"`
	UnitPrice        float32 `db:"unit_price" json:"unit_price,omitempty" validate:"required"`
}

type Model struct {
	Id          int64            `db:"id" json:"id"`
	CreatedAt   lystype.Datetime `db:"created_at" json:"created_at,omitzero"`
	CreatedBy   string           `db:"created_by" json:"created_by,omitempty"`
	OrderNumber int32            `db:"order_number" json:"order_number,omitempty"`
	ProductName string           `db:"product_name" json:"product_name,omitempty"`
	UpdatedAt   lystype.Datetime `db:"updated_at" json:"updated_at,omitzero"` // assigned by trigger
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

func (s Store) ArchiveById(ctx context.Context, tx pgx.Tx, id int64) error {
	return lyspg.Archive(ctx, tx, schemaName, tableName, pkColName, id, false)
}

func (s Store) ArchiveCascadedByOrder(ctx context.Context, tx pgx.Tx, orderId int64) error {
	return lyspg.Archive(ctx, tx, schemaName, tableName, "order_fk", orderId, true)
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

func (s Store) RestoreById(ctx context.Context, tx pgx.Tx, id int64) error {
	return lyspg.Restore(ctx, tx, schemaName, tableName, pkColName, id, false)
}

func (s Store) RestoreCascadedByOrder(ctx context.Context, tx pgx.Tx, orderId int64) error {
	return lyspg.Restore(ctx, tx, schemaName, tableName, "order_fk", orderId, true)
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, id int64) (item Model, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, id)
}

func (s Store) Update(ctx context.Context, input Input, id int64) error {
	input.LastUserUpdateBy = lys.GetUserNameFromCtx(ctx, "Unknown")
	return lyspg.Update(ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	assignmentsMap["last_user_update_by"] = lys.GetUserNameFromCtx(ctx, "Unknown")
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) error {
	return lysmeta.Validate(validate, input)
}
