
-- example
--CREATE TRIGGER t_set_updated_at BEFORE UPDATE ON core.category FOR EACH ROW EXECUTE PROCEDURE system.set_updated_at();

-- this trigger is added automatically where needed by lyspgmon.CheckDDL