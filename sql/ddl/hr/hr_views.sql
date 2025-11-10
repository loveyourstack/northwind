
CREATE OR REPLACE VIEW hr.v_employee AS
	SELECT 
    emp.address,
    EXTRACT('YEAR' FROM AGE(current_date, emp.date_of_birth))::int AS age,
    emp.city,
    emp.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    emp.created_at,
    emp.created_by,
    emp.date_of_birth,
    emp.family_name,
    emp.given_name,
    emp.hire_date,
    emp.home_phone,
    emp.job_title,
    emp.id,
    emp.last_user_update_by,
    emp.name,
    emp.notes,
    emp.postal_code,
    emp.reports_to_fk,
    CASE WHEN emp.reports_to_fk = emp.id THEN '' ELSE report.name END AS reports_to,
    emp.state,
    emp.title,
    emp.updated_at,
    COALESCE(terr.cnt,0) AS territory_count
  FROM hr.employee emp
  JOIN core.country co ON emp.country_fk = co.id
  JOIN hr.employee report ON emp.reports_to_fk = report.id
  LEFT JOIN (SELECT salesman_fk, count(*) AS cnt FROM sales.territory GROUP BY 1) terr ON terr.salesman_fk = emp.id;
