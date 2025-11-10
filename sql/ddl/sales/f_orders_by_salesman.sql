
CREATE OR REPLACE FUNCTION sales.orders_by_salesman (
	_from_date date,
  _until_date date
)
RETURNS TABLE (
  order_count int,
  salesman text_medium,
  shipped_count int,
  total_value numeric
) AS
$BODY$

SELECT 
  COUNT(*) AS order_count,
  emp.name AS salesman,
  COALESCE(COUNT(*) FILTER (WHERE o.is_shipped = true),0) AS shipped_count,
  SUM(COALESCE(oi.gross_value,0)) AS total_value
FROM sales.order o
JOIN hr.employee emp ON o.salesman_fk = emp.id
LEFT JOIN (SELECT order_fk, count(*), SUM(gross_value) AS gross_value FROM sales.order_item GROUP BY 1) oi ON oi.order_fk = o.id
WHERE order_date BETWEEN '2021-04-20' AND '2021-05-06'
GROUP BY emp.name;

$BODY$
  LANGUAGE sql VOLATILE
  COST 100;