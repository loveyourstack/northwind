
CREATE OR REPLACE VIEW hr.v_employee AS
	SELECT 
    hr_e.id,
    hr_e.address,
    EXTRACT('YEAR' FROM AGE(current_date, hr_e.date_of_birth))::int AS age,
    hr_e.city,
    hr_e.country_fk,
    co.name AS country,
    co.iso2 AS country_iso2,
    hr_e.date_of_birth,
    hr_e.entry_at,
    hr_e.entry_by,
    hr_e.first_name,
    hr_e.hire_date,
    hr_e.home_phone,
    hr_e.job_title,
    hr_e.last_modified_at,
    hr_e.last_modified_by,
    hr_e.last_name,
    hr_e.name,
    hr_e.notes,
    hr_e.postal_code,
    hr_e.reports_to_fk,
    CASE WHEN hr_e.reports_to_fk = hr_e.id THEN '' ELSE report.name END AS reports_to,
    hr_e.state,
    hr_e.title
  FROM hr.employee hr_e
  JOIN common.country co ON hr_e.country_fk = co.id
  JOIN hr.employee report ON hr_e.reports_to_fk = report.id;
