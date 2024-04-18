
CREATE TRIGGER t_audit_update AFTER UPDATE ON core.category FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON core.product FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON core.supplier FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();

CREATE TRIGGER t_audit_update AFTER UPDATE ON hr.employee FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();

CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.customer FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.employee_territory FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.order FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.order_detail FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.shipper FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
CREATE TRIGGER t_audit_update AFTER UPDATE ON sales.territory FOR EACH ROW EXECUTE PROCEDURE system.update_trigger();
