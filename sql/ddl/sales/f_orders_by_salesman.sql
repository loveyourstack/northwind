
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
  salesman,
  COALESCE(COUNT(*) FILTER (WHERE is_shipped = true),0) AS shipped_count,
  SUM(order_value) AS total_value
FROM sales.v_order
WHERE order_date BETWEEN _from_date AND _until_date
GROUP BY salesman;

$BODY$
  LANGUAGE sql VOLATILE
  COST 100;