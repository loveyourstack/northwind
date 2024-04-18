
CREATE OR REPLACE VIEW sales.v_customer AS
	SELECT 
    s_c.id,
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
    s_c.last_modified_at,
    s_c.last_modified_by,
    s_c.code || ' (' || s_c.company_name || ')' AS name,
    s_c.phone,
    s_c.postal_code,
    s_c.state,
    COALESCE(s_o.order_count,0) AS order_count
  FROM sales.customer s_c
  JOIN common.country co ON s_c.country_fk = co.id
  LEFT JOIN (SELECT customer_fk, count(*) AS order_count FROM sales.order GROUP BY 1) s_o ON s_o.customer_fk = s_c.id;


CREATE OR REPLACE VIEW sales.v_order AS
	SELECT 
    s_o.id,
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
    COALESCE(s_od.count,0) AS order_detail_count,
    COALESCE(s_od.value,0) AS order_value
  FROM sales.order s_o
  JOIN sales.customer s_c ON s_o.customer_fk = s_c.id
  JOIN common.country co ON s_o.dest_country_fk = co.id
  JOIN hr.employee hr_e ON s_o.salesman_fk = hr_e.id
  JOIN sales.shipper s_s ON s_o.shipper_fk = s_s.id
  LEFT JOIN (SELECT order_fk, count(*), SUM(quantity * unit_price * (1 - discount)) AS value FROM sales.order_detail GROUP BY 1) s_od ON s_od.order_fk = s_o.id;;


CREATE OR REPLACE VIEW sales.v_order_detail AS
	SELECT 
    s_od.id,
    s_od.discount,
    s_od.entry_at,
    s_od.entry_by,
    s_od.last_modified_at,
    s_od.last_modified_by,
    s_od.order_fk,
    s_o.order_number,
    s_od.product_fk,
    c_p.name AS product_name,
    s_od.quantity,
    s_od.unit_price
  FROM sales.order_detail s_od
  JOIN sales.order s_o ON s_od.order_fk = s_o.id
  JOIN core.product c_p ON s_od.product_fk = c_p.id;


CREATE OR REPLACE VIEW sales.v_employee_territory AS
	SELECT 
    s_et.id,
    s_et.employee_fk,
    hr_e.name AS employee_name,
    s_et.entry_at,
    s_et.entry_by,
    s_et.last_modified_at,
    s_et.last_modified_by,
    s_et.territory_fk,
    s_t.name AS territory_name
  FROM sales.employee_territory s_et
  JOIN hr.employee hr_e ON s_et.employee_fk = hr_e.id
  JOIN sales.territory s_t ON s_et.territory_fk = s_t.id;