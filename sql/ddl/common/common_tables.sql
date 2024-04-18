
CREATE TABLE common.country
(	
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  is_active bool NOT NULL DEFAULT false,
  iso2 text NOT NULL UNIQUE,
  entry_at tracking_at,
  last_modified_at tracking_at,
  entry_by tracking_by,
  last_modified_by tracking_by,
  name text NOT NULL UNIQUE,
  CONSTRAINT co_iso2_len CHECK (char_length(iso2) = 2)
);
COMMENT ON TABLE common.country IS 'shortname: co';
