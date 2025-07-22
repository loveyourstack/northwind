
CREATE OR REPLACE VIEW core.v_category AS
	SELECT 
    c_c.color_hex,
    c_c.color_is_light,
    c_c.description,
    c_c.entry_at,
    c_c.entry_by,
    c_c.id,
    c_c.last_modified_at,
    c_c.last_modified_by,
    c_c.name,
    COALESCE(c_p.active_product_count,0) AS active_product_count
  FROM core.category c_c
  LEFT JOIN (SELECT category_fk, count(*) FILTER (WHERE is_discontinued = false) AS active_product_count FROM core.product GROUP BY 1) c_p ON c_p.category_fk = c_c.id;

CREATE OR REPLACE VIEW core.v_category_data_update AS
  SELECT
    c_c.*,
    sys_du.id AS data_update_id,
    sys_du.affected_id,
    sys_du.affected_at,
    sys_du.affected_by,
    sys_du.affected_old_values,
    sys_du.affected_new_values
  FROM core.v_category c_c
  JOIN system.data_update sys_du ON c_c.id = sys_du.affected_id
  WHERE sys_du.affected_schema = 'core' AND sys_du.affected_table = 'category';


CREATE OR REPLACE VIEW core.v_product AS
	SELECT 
    c_p.category_fk,
    c_c.name AS category,
    c_c.color_hex AS category_color_hex,
    c_c.color_is_light AS category_color_is_light,
    c_p.entry_at,
    c_p.entry_by,
    c_p.id,
    c_p.is_discontinued,
    c_p.last_modified_at,
    c_p.last_modified_by,
    c_p.name,
    c_p.quantity_per_unit,
    c_p.reorder_level,
    c_p.supplier_fk,
    c_s.company_name AS supplier_company_name,
    co.iso2 AS supplier_country_iso2,
    c_p.tags,
    c_p.unit_price,
    c_p.units_in_stock,
    c_p.units_on_order
  FROM core.product c_p
  JOIN core.supplier c_s ON c_p.supplier_fk = c_s.id
  JOIN core.country co ON c_s.country_fk = co.id
  JOIN core.category c_c ON c_p.category_fk = c_c.id;


CREATE OR REPLACE VIEW core.v_supplier AS
	SELECT 
    c_s.address,
    c_s.city,
    c_s.company_name,
    c_s.contact_name,
    c_s.contact_title,
    c_s.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    c_s.entry_at,
    c_s.entry_by,
    c_s.id,
    c_s.last_modified_at,
    c_s.last_modified_by,
    c_s.phone,
    c_s.postal_code,
    c_s.state,
    c_s.company_name || ' (' || co.iso2 || ')' AS name, --unique
    COALESCE(c_p.active_product_count,0) AS active_product_count
  FROM core.supplier c_s
  JOIN core.country co ON c_s.country_fk = co.id
  LEFT JOIN (SELECT supplier_fk, count(*) FILTER (WHERE is_discontinued = false) AS active_product_count FROM core.product GROUP BY 1) c_p ON c_p.supplier_fk = c_s.id;

CREATE OR REPLACE VIEW core.v_supplier_data_update AS
  SELECT
    c_s.*,
    sys_du.id AS data_update_id,
    sys_du.affected_id,
    sys_du.affected_at,
    sys_du.affected_by,
    sys_du.affected_old_values,
    sys_du.affected_new_values
  FROM core.v_supplier c_s
  JOIN system.data_update sys_du ON c_s.id = sys_du.affected_id
  WHERE sys_du.affected_schema = 'core' AND sys_du.affected_table = 'supplier';
