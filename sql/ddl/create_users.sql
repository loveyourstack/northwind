-- run as superuser

CREATE USER northwind_owner WITH PASSWORD '123' CREATEROLE;
CREATE USER northwind_server WITH PASSWORD '456';
CREATE USER northwind_cli WITH PASSWORD '789';