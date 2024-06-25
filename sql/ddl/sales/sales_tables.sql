
CREATE TABLE sales.customer
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  country_fk bigint NOT NULL REFERENCES common.country (id),
  entry_at tracking_at,
  last_modified_at tracking_at,
  code text NOT NULL UNIQUE,
  address text NOT NULL,
  city text NOT NULL,
  company_name text NOT NULL,
  contact_name text NOT NULL,
  contact_title text NOT NULL,
  entry_by tracking_by,
  last_modified_by tracking_by,
  phone text NOT NULL,
  postal_code text NOT NULL,
  state text NOT NULL,
  CONSTRAINT s_c_code_len CHECK (char_length(code) = 5)
);
COMMENT ON TABLE sales.customer IS 'shortname: s_c';


CREATE TABLE sales.shipper
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  entry_at tracking_at,
  last_modified_at tracking_at,
  company_name text NOT NULL UNIQUE,
  entry_by tracking_by,
  last_modified_by tracking_by,
  phone text NOT NULL
);
COMMENT ON TABLE sales.shipper IS 'shortname: s_s';


CREATE TABLE sales.order
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  is_shipped boolean NOT NULL DEFAULT false,
  order_number int NOT NULL UNIQUE CHECK (order_number > 0),
  customer_fk bigint NOT NULL REFERENCES sales.customer(id),
  dest_country_fk bigint NOT NULL REFERENCES common.country(id),
  salesman_fk bigint NOT NULL REFERENCES hr.employee(id),
  shipper_fk bigint NOT NULL REFERENCES sales.shipper(id),
  freight_cost numeric(12,2) NOT NULL CHECK (freight_cost >= 0.0),
  order_date date NOT NULL CHECK (order_date >= '2010-01-01'),
  required_date date NOT NULL CHECK (required_date >= '2010-01-01'),
  shipped_date date NOT NULL DEFAULT '0001-01-01'::date,
  entry_at tracking_at,
  last_modified_at tracking_at,
  dest_address text NOT NULL,
  dest_city text NOT NULL,
  dest_company_name text NOT NULL,
  dest_postal_code text NOT NULL,
  dest_state text NOT NULL,
  entry_by tracking_by,
  last_modified_by tracking_by
);
COMMENT ON TABLE sales.order IS 'shortname: s_o';
--- change columns with ---
CREATE TABLE sales.order_deleted
(
  id bigint NOT NULL PRIMARY KEY,
  is_shipped boolean NOT NULL,
  order_number int NOT NULL,
  customer_fk bigint NOT NULL,
  dest_country_fk bigint NOT NULL,
  salesman_fk bigint NOT NULL,
  shipper_fk bigint NOT NULL,
  freight_cost numeric(12,2) NOT NULL,
  order_date date NOT NULL,
  required_date date NOT NULL,
  shipped_date date NOT NULL,
  entry_at tracking_at,
  last_modified_at tracking_at,
  dest_address text NOT NULL,
  dest_city text NOT NULL,
  dest_company_name text NOT NULL,
  dest_postal_code text NOT NULL,
  dest_state text NOT NULL,
  entry_by tracking_by,
  last_modified_by tracking_by,
  deleted_at tracking_at,
  deleted_by tracking_by,
  deleted_by_cascade bool NOT NULL
);
COMMENT ON TABLE sales.order_deleted IS 'shortname: s_o_del';


CREATE TABLE sales.order_item
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  quantity int NOT NULL CHECK (quantity >= 0),
  order_fk bigint NOT NULL REFERENCES sales.order(id),
  product_fk bigint NOT NULL REFERENCES core.product(id),
  unit_price numeric(12,2) NOT NULL CHECK (unit_price >= 0.0),
  discount numeric(5,4) NOT NULL CHECK (discount BETWEEN 0.0 AND 1.0),
  entry_at tracking_at,
  last_modified_at tracking_at,
  entry_by tracking_by,
  last_modified_by tracking_by,
  UNIQUE(order_fk, product_fk)
);
COMMENT ON TABLE sales.order_item IS 'shortname: s_oi';
--- change columns with ---
CREATE TABLE sales.order_item_deleted
(
  id bigint NOT NULL PRIMARY KEY,
  quantity int NOT NULL,
  order_fk bigint NOT NULL,
  product_fk bigint NOT NULL,
  unit_price numeric(12,2) NOT NULL,
  discount numeric(5,4) NOT NULL,
  entry_at tracking_at,
  last_modified_at tracking_at,
  entry_by tracking_by,
  last_modified_by tracking_by,
  deleted_at tracking_at,
  deleted_by tracking_by,
  deleted_by_cascade bool NOT NULL
);
COMMENT ON TABLE sales.order_item_deleted IS 'shortname: s_oi_del';

CREATE INDEX order_item_order_fk_idx ON sales.order_item USING btree(order_fk);
CLUSTER sales.order_item USING order_item_order_fk_idx;


CREATE TABLE sales.territory
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  region sales.region NOT NULL,
  code text NOT NULL UNIQUE,
  entry_at tracking_at,
  last_modified_at tracking_at,
  entry_by tracking_by,
  last_modified_by tracking_by,
  name text NOT NULL UNIQUE,
  CONSTRAINT s_t_code_len CHECK (char_length(code) = 5)
);
COMMENT ON TABLE sales.territory IS 'shortname: s_t';


CREATE TABLE sales.employee_territory
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  employee_fk bigint NOT NULL REFERENCES hr.employee(id),
  territory_fk bigint NOT NULL REFERENCES sales.territory(id),
  entry_at tracking_at,
  last_modified_at tracking_at,
  entry_by tracking_by,
  last_modified_by tracking_by,
  UNIQUE (employee_fk, territory_fk)
);
COMMENT ON TABLE sales.employee_territory IS 'shortname: s_et';
