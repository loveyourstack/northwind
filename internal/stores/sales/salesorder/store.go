package salesorder

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys/lyserr"
	"github.com/loveyourstack/lys/lysmeta"
	"github.com/loveyourstack/lys/lyspg"
	"github.com/loveyourstack/lys/lystype"
	"github.com/loveyourstack/northwind/internal/nw"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesorderitem"
	"github.com/loveyourstack/northwind/pkg/lyspmath"
)

const (
	name           string = "Orders"
	schemaName     string = "sales"
	tableName      string = "order"
	viewName       string = "v_order"
	pkColName      string = "id"
	defaultOrderBy string = "order_number DESC"
)

type Input struct {
	CustomerFk      int64            `db:"customer_fk" json:"customer_fk,omitempty" validate:"required"`
	DestAddress     string           `db:"dest_address" json:"dest_address,omitempty" validate:"required"`
	DestCity        string           `db:"dest_city" json:"dest_city,omitempty" validate:"required"`
	DestCompanyName string           `db:"dest_company_name" json:"dest_company_name,omitempty" validate:"required"`
	DestCountryFk   int64            `db:"dest_country_fk" json:"dest_country_fk,omitempty" validate:"required"`
	DestPostalCode  string           `db:"dest_postal_code" json:"dest_postal_code,omitempty" validate:"required"`
	DestState       string           `db:"dest_state" json:"dest_state,omitempty"`
	FreightCost     float32          `db:"freight_cost" json:"freight_cost,omitempty" validate:"number,gte=0"`
	IsShipped       bool             `db:"is_shipped" json:"is_shipped,omitempty"`
	LastModifiedAt  lystype.Datetime `db:"last_modified_at" json:"last_modified_at,omitzero"` // assigned in Update funcs
	OrderDate       lystype.Date     `db:"order_date" json:"order_date,omitzero" validate:"required"`
	OrderNumber     int32            `db:"order_number" json:"order_number,omitempty"` // assigned in Insert
	RequiredDate    lystype.Date     `db:"required_date" json:"required_date,omitzero" validate:"required"`
	SalesmanFk      int64            `db:"salesman_fk" json:"salesman_fk,omitempty" validate:"required"`
	ShippedDate     lystype.Date     `db:"shipped_date" json:"shipped_date,omitzero"`
	ShipperFk       int64            `db:"shipper_fk" json:"shipper_fk,omitempty" validate:"required"`
}

type Model struct {
	Id                  int64            `db:"id" json:"id"`
	CustomerCode        string           `db:"customer_code" json:"customer_code,omitempty"`
	CustomerCompanyName string           `db:"customer_company_name" json:"customer_company_name,omitempty"`
	DestCountryIso2     string           `db:"dest_country_iso2" json:"dest_country_iso2,omitempty"`
	EntryAt             lystype.Datetime `db:"entry_at" json:"entry_at,omitzero"`
	EntryBy             string           `db:"entry_by" json:"entry_by,omitempty"`
	LastModifiedBy      string           `db:"last_modified_by" json:"last_modified_by,omitempty"`
	Salesman            string           `db:"salesman" json:"salesman,omitempty"`
	ShipperCompanyName  string           `db:"shipper_company_name" json:"shipper_company_name,omitempty"`
	OrderItemCount      int              `db:"order_item_count" json:"order_item_count"`
	OrderValue          float32          `db:"order_value" json:"order_value"`
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

	// cascade to order items
	orderItemStore := salesorderitem.Store{Db: s.Db}
	err := orderItemStore.ArchiveCascadedByOrder(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("orderItemStore.ArchiveCascadedByOrder failed: %w", err)
	}

	return lyspg.Archive(ctx, tx, schemaName, tableName, pkColName, id, false)
}

func (s Store) BulkInsert(ctx context.Context, inputs []Input) (rowsAffected int64, err error) {
	return lyspg.BulkInsert[Input](ctx, s.Db, schemaName, tableName, inputs)
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

	// get next order number and add it to input
	stmt := fmt.Sprintf("SELECT max(order_number)+1 FROM %s.%s;", schemaName, tableName)
	rows, _ := s.Db.Query(ctx, stmt)
	nextOrderNum, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int32])
	if err != nil {
		return 0, lyserr.Db{Err: fmt.Errorf("pgx.CollectExactlyOneRow failed: %w", err), Stmt: stmt}
	}
	input.OrderNumber = nextOrderNum

	return lyspg.Insert[Input, int64](ctx, s.Db, schemaName, tableName, pkColName, input)
}

func (s Store) RestoreById(ctx context.Context, tx pgx.Tx, id int64) error {

	err := lyspg.Restore(ctx, tx, schemaName, tableName, pkColName, id, false)
	if err != nil {
		return fmt.Errorf("lyspg.Restore failed: %w", err)
	}

	// cascade to order items
	orderItemStore := salesorderitem.Store{Db: s.Db}
	err = orderItemStore.RestoreCascadedByOrder(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("orderItemStore.RestoreCascadedByOrder failed: %w", err)
	}

	return nil
}

func (s Store) Select(ctx context.Context, params lyspg.SelectParams) (items []Model, unpagedCount lyspg.TotalCount, err error) {
	return lyspg.Select[Model](ctx, s.Db, schemaName, tableName, viewName, defaultOrderBy, meta.DbTags, params)
}

func (s Store) SelectById(ctx context.Context, id int64) (item Model, err error) {
	return lyspg.SelectUnique[Model](ctx, s.Db, schemaName, viewName, pkColName, id)
}

func (s Store) Update(ctx context.Context, input Input, id int64) error {
	input.LastModifiedAt = lystype.Datetime(time.Now())
	return lyspg.Update[Input](ctx, s.Db, schemaName, tableName, pkColName, input, id)
}

func (s Store) UpdatePartial(ctx context.Context, assignmentsMap map[string]any, id int64) error {
	assignmentsMap["last_modified_at"] = lystype.Datetime(time.Now())
	return lyspg.UpdatePartial(ctx, s.Db, schemaName, tableName, pkColName, inputMeta.DbTags, assignmentsMap, id)
}

func (s Store) Validate(validate *validator.Validate, input Input) (err error) {

	err = lysmeta.Validate(validate, input)
	if err != nil {
		return err
	}

	// further business rules

	if input.FreightCost != lyspmath.RoundFloat32(input.FreightCost, 2) {
		return fmt.Errorf("freight_cost may have max 2 decimal places")
	}

	if time.Time(input.OrderDate).Before(nw.CommencementOfTradingDate) {
		return fmt.Errorf("order_date must be from %s", nw.CommencementOfTradingDateStr)
	}

	if time.Time(input.RequiredDate).Before(nw.CommencementOfTradingDate) {
		return fmt.Errorf("required_date must be from %s", nw.CommencementOfTradingDateStr)
	}

	return nil
}
