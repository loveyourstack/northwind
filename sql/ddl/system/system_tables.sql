
CREATE TABLE system.data_update
(
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  affected_id bigint NOT NULL, -- column names are intended to be unambiguous when this table is joined to others
  affected_at tracking_at,
  affected_by tracking_by,
  affected_schema text NOT NULL,
  affected_table text NOT NULL,
  affected_old_values jsonb NOT NULL,
  affected_new_values jsonb NOT NULL
);
COMMENT ON TABLE system.data_update IS 'shortname: sys_du';

GRANT SELECT,INSERT ON system.data_update TO northwind_server;
GRANT SELECT,INSERT ON system.data_update TO northwind_cli;
