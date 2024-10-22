
CREATE OR REPLACE VIEW ecb.v_exchange_rate AS
  SELECT
    xr.day,
    xr.frequency,
    xr.from_currency_fk,
    from_curr.code AS from_currency,
    xr.entry_at,
    xr.last_modified_at,
    xr.id,
    xr.rate,
    xr.to_currency_fk,
    to_curr.code AS to_currency
  FROM ecb.exchange_rate xr
  JOIN ecb.currency from_curr ON xr.from_currency_fk = from_curr.id
  JOIN ecb.currency to_curr ON xr.to_currency_fk = to_curr.id;
