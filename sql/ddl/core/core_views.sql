
CREATE OR REPLACE VIEW core.v_category AS
	SELECT 
    cat.color_hex,
    cat.color_is_light,
    cat.created_at,
    cat.created_by,
    cat.description,
    cat.id,
    cat.last_user_update_by,
    cat.name,
    cat.updated_at,
    COALESCE(p.active_product_count,0) AS active_product_count
  FROM core.category cat
  LEFT JOIN (SELECT category_fk, count(*) FILTER (WHERE is_discontinued = false) AS active_product_count FROM core.product GROUP BY 1) p ON p.category_fk = cat.id;

CREATE OR REPLACE VIEW core.v_category_data_update AS
  SELECT
    cat.*,
    sys_du.id AS data_update_id,
    sys_du.affected_id,
    sys_du.affected_at,
    sys_du.affected_by,
    sys_du.affected_old_values,
    sys_du.affected_new_values
  FROM core.v_category cat
  JOIN system.data_update sys_du ON cat.id = sys_du.affected_id
  WHERE sys_du.affected_schema = 'core' AND sys_du.affected_table = 'category';


CREATE OR REPLACE VIEW core.v_product AS
	SELECT 
    p.category_fk,
    cat.name AS category,
    cat.color_hex AS category_color_hex,
    cat.color_is_light AS category_color_is_light,
    p.created_at,
    p.created_by,
    p.id,
    p.is_discontinued,
    p.last_user_update_by,
    p.name,
    p.quantity_per_unit,
    p.reorder_level,
    p.supplier_fk,
    s.company_name AS supplier_company_name,
    co.iso2 AS supplier_country_iso2,
    p.tags,
    p.unit_price,
    p.units_in_stock,
    p.units_on_order,
    p.updated_at
  FROM core.product p
  JOIN core.supplier s ON p.supplier_fk = s.id
  JOIN core.country co ON s.country_fk = co.id
  JOIN core.category cat ON p.category_fk = cat.id;


CREATE OR REPLACE VIEW core.v_supplier AS
	SELECT 
    s.address,
    s.city,
    s.company_name,
    s.contact_name,
    s.contact_title,
    s.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    s.created_at,
    s.created_by,
    s.id,
    s.last_user_update_by,
    s.phone,
    s.postal_code,
    s.state,
    s.company_name || ' (' || co.iso2 || ')' AS name, --unique
    s.updated_at,
    COALESCE(p.active_product_count,0) AS active_product_count
  FROM core.supplier s
  JOIN core.country co ON s.country_fk = co.id
  LEFT JOIN (SELECT supplier_fk, count(*) FILTER (WHERE is_discontinued = false) AS active_product_count FROM core.product GROUP BY 1) p ON p.supplier_fk = s.id;

CREATE OR REPLACE VIEW core.v_supplier_data_update AS
  SELECT
    s.*,
    sys_du.id AS data_update_id,
    sys_du.affected_id,
    sys_du.affected_at,
    sys_du.affected_by,
    sys_du.affected_old_values,
    sys_du.affected_new_values
  FROM core.v_supplier s
  JOIN system.data_update sys_du ON s.id = sys_du.affected_id
  WHERE sys_du.affected_schema = 'core' AND sys_du.affected_table = 'supplier';
