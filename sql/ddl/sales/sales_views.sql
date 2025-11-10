
CREATE OR REPLACE VIEW sales.v_customer AS
	SELECT 
    c.address,
    c.city,
    c.code,
    c.company_name,
    c.contact_name,
    c.contact_title,
    c.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    c.created_at,
    c.created_by,
    c.id,
    c.code || ' (' || c.company_name || ')' AS name,
    c.last_user_update_by,
    c.phone,
    c.postal_code,
    c.state,
    c.updated_at,
    COALESCE(o.order_count,0) AS order_count
  FROM sales.customer c
  JOIN core.country co ON c.country_fk = co.id
  LEFT JOIN (SELECT customer_fk, count(*) AS order_count FROM sales.order GROUP BY 1) o ON o.customer_fk = c.id;


CREATE OR REPLACE VIEW sales.v_order AS
	SELECT 
    o.customer_fk,
    c.code AS customer_code,
    c.company_name AS customer_company_name,
    o.dest_address,
    o.dest_city,
    o.dest_company_name,
    o.dest_country_fk,
    co.iso2 AS dest_country_iso2,
    o.dest_postal_code,
    o.dest_state,
    o.created_at,
    o.created_by,
    o.freight_cost,
    o.id,
    o.is_shipped,
    o.last_user_update_by,
    o.order_date,
    o.order_number,
    o.required_date,
    o.salesman_fk,
    emp.name AS salesman,
    o.shipped_date,
    o.shipper_fk,
    ship.company_name AS shipper_company_name,
    o.updated_at,
    COALESCE(oi.count,0) AS order_item_count,
    COALESCE(oi.gross_value,0) AS order_value
  FROM sales.order o
  JOIN sales.customer c ON o.customer_fk = c.id
  JOIN core.country co ON o.dest_country_fk = co.id
  JOIN hr.employee emp ON o.salesman_fk = emp.id
  JOIN sales.shipper ship ON o.shipper_fk = ship.id
  LEFT JOIN (SELECT order_fk, count(*), SUM(gross_value) AS gross_value FROM sales.order_item GROUP BY 1) oi ON oi.order_fk = o.id;


CREATE OR REPLACE VIEW sales.v_order_item AS
	SELECT 
    oi.discount,
    oi.created_at,
    oi.created_by,
    oi.gross_value,
    oi.id,
    oi.last_user_update_by,
    oi.order_fk,
    o.order_number,
    oi.product_fk,
    p.name AS product_name,
    oi.quantity,
    oi.unit_price,
    oi.updated_at
  FROM sales.order_item oi
  JOIN sales.order o ON oi.order_fk = o.id
  JOIN core.product p ON oi.product_fk = p.id;


CREATE OR REPLACE VIEW sales.v_order_value_latest_weeks AS
	WITH latest_weeks AS (
	  SELECT EXTRACT('week' FROM order_date)::int AS order_week, ROUND((SUM(order_value)/1000),1) AS total_value 
	  FROM sales.v_order WHERE order_date BETWEEN '2021-01-03' AND '2021-12-31'
	  GROUP BY 1 ORDER BY 1 DESC LIMIT 10
	)
	SELECT * FROM latest_weeks ORDER BY order_week;


CREATE OR REPLACE VIEW sales.v_territory AS
	SELECT 
    terr.code,
    salesman_fk,
    emp.name AS salesman,
    terr.created_at,
    terr.created_by,
    terr.id,
    terr.last_user_update_by,
    terr.name,
    terr.region::text, -- cast so that sorting works as expected
    terr.updated_at
  FROM sales.territory terr
  LEFT JOIN hr.employee emp ON terr.salesman_fk = emp.id;

