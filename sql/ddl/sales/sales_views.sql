
CREATE OR REPLACE VIEW sales.v_customer AS
	SELECT 
    s_c.address,
    s_c.city,
    s_c.code,
    s_c.company_name,
    s_c.contact_name,
    s_c.contact_title,
    s_c.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    s_c.entry_at,
    s_c.entry_by,
    s_c.id,
    s_c.last_modified_at,
    s_c.last_modified_by,
    s_c.code || ' (' || s_c.company_name || ')' AS name,
    s_c.phone,
    s_c.postal_code,
    s_c.state,
    COALESCE(s_o.order_count,0) AS order_count
  FROM sales.customer s_c
  JOIN core.country co ON s_c.country_fk = co.id
  LEFT JOIN (SELECT customer_fk, count(*) AS order_count FROM sales.order GROUP BY 1) s_o ON s_o.customer_fk = s_c.id;


CREATE OR REPLACE VIEW sales.v_order AS
	SELECT 
    s_o.customer_fk,
    s_c.code AS customer_code,
    s_c.company_name AS customer_company_name,
    s_o.dest_address,
    s_o.dest_city,
    s_o.dest_company_name,
    s_o.dest_country_fk,
    co.iso2 AS dest_country_iso2,
    s_o.dest_postal_code,
    s_o.dest_state,
    s_o.entry_at,
    s_o.entry_by,
    s_o.freight_cost,
    s_o.id,
    s_o.is_shipped,
    s_o.last_modified_at,
    s_o.last_modified_by,
    s_o.order_date,
    s_o.order_number,
    s_o.required_date,
    s_o.salesman_fk,
    hr_e.name AS salesman,
    s_o.shipped_date,
    s_o.shipper_fk,
    s_s.company_name AS shipper_company_name,
    COALESCE(s_oi.count,0) AS order_item_count,
    COALESCE(s_oi.value,0) AS order_value
  FROM sales.order s_o
  JOIN sales.customer s_c ON s_o.customer_fk = s_c.id
  JOIN core.country co ON s_o.dest_country_fk = co.id
  JOIN hr.employee hr_e ON s_o.salesman_fk = hr_e.id
  JOIN sales.shipper s_s ON s_o.shipper_fk = s_s.id
  LEFT JOIN (SELECT order_fk, count(*), SUM(quantity * unit_price * (1 - discount)) AS value FROM sales.order_item GROUP BY 1) s_oi ON s_oi.order_fk = s_o.id;


CREATE OR REPLACE VIEW sales.v_order_item AS
	SELECT 
    s_oi.discount,
    s_oi.entry_at,
    s_oi.entry_by,
    s_oi.id,
    s_oi.last_modified_at,
    s_oi.last_modified_by,
    s_oi.order_fk,
    s_o.order_number,
    s_oi.product_fk,
    c_p.name AS product_name,
    s_oi.quantity,
    s_oi.unit_price
  FROM sales.order_item s_oi
  JOIN sales.order s_o ON s_oi.order_fk = s_o.id
  JOIN core.product c_p ON s_oi.product_fk = c_p.id;


CREATE OR REPLACE VIEW sales.v_order_value_latest_weeks AS
	WITH latest_weeks AS (
	  SELECT EXTRACT('week' FROM order_date)::int AS order_week, ROUND((SUM(order_value)/1000),1) AS total_value 
	  FROM sales.v_order WHERE order_date BETWEEN '2021-01-03' AND '2021-12-31'
	  GROUP BY 1 ORDER BY 1 DESC LIMIT 10
	)
	SELECT * FROM latest_weeks ORDER BY order_week;


CREATE OR REPLACE VIEW sales.v_territory AS
	SELECT 
    s_t.code,
    salesman_fk,
    hr_e.name AS salesman,
    s_t.entry_at,
    s_t.entry_by,
    s_t.id,
    s_t.last_modified_at,
    s_t.last_modified_by,
    s_t.name,
    s_t.region::text -- cast so that sorting works as expected
  FROM sales.territory s_t
  LEFT JOIN hr.employee hr_e ON s_t.salesman_fk = hr_e.id;

