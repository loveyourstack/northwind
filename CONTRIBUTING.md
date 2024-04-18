# Contributing

## Set up development environment

Create the Northwind PostgreSQL database on a Linux system as follows:
* find the file "sql/ddl/create_users.sql" and run the contents as superuser on your local cluster
* copy nw_config_sample.toml to nw_config.toml
* in nw_config.toml, enter your db superuser password
* ensure you don't currently have a database with the database name shown (default: "northwind")
* enter `make db`. This will (re-)create the database and populate it with tables and data

## Run tests

TODO