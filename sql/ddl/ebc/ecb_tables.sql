
CREATE TABLE ecb.currency
(	
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  entry_at tracking_at,
  last_modified_at tracking_at,
  code text NOT NULL UNIQUE, -- natural key
  name text NOT NULL
);
COMMENT ON TABLE ecb.currency IS 'shortname: curr';


CREATE TABLE ecb.exchange_rate
(	
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  frequency ecb.frequency NOT NULL,
  from_currency_fk bigint NOT NULL REFERENCES ecb.currency(id),
  to_currency_fk bigint NOT NULL REFERENCES ecb.currency(id),
  rate numeric(12,4) NOT NULL,
  day date NOT NULL,
  entry_at tracking_at,
  last_modified_at tracking_at,
  UNIQUE (frequency, day, from_currency_fk, to_currency_fk)
);
COMMENT ON TABLE ecb.exchange_rate IS 'shortname: xr';
