
CREATE OR REPLACE FUNCTION system.remove_jsonb_fields(p_row jsonb)
  RETURNS jsonb AS
$BODY$
DECLARE
BEGIN

-- remove these cols in all tables
p_row = p_row - 'id';
p_row = p_row - 'entry_at';
p_row = p_row - 'entry_by';
p_row = p_row - 'last_modified_at';
p_row = p_row - 'last_modified_by';

RETURN p_row;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
