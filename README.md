# Northwind

Sample application showing the implementation of [LoveYourStack (lys)](https://github.com/loveyourstack/lys) packages.

Adapts the Northwind sample database from [here](https://github.com/pthom/northwind_psql/blob/master/northwind.sql).

Is a work in progress - some functionality is missing, especially in the UI, and no tests are written yet.

## Setup

* Follow the instructions in CONTRIBUTING.md for creating a local database
* Fetch API packages: from root directory, `go get -d ./...`
* Run the API server: from root directory, `make srv`
* Fetch UI packages: from frontend/northwind-ui directory, `yarn`
* Run the UI: from frontend/northwind-ui directory, `yarn dev`

## Testing

See CONTRIBUTING.md for setup instructions.

## Supported Go and PostgreSQL Versions

Go 1.16+ (due to embed.FS)

PostgreSQL TODO