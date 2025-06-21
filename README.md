# Tradovate data importer and visualizer

## Prereqs

- Bash
- Docker
- Go

## Usage

1. Download CSV data files from Tradovate and place them in a folder named ".data" at the root of this repo.
2. Run the script `./import_data.sh` to create a SQLite database from that data.  This database file will be stored in a folder named ".db" at the root of this repo.
3. Run `./upgrade_metabase.sh` to run Metabase in a new Docker container, or upgrade an already-running Metabase Docker container.

The Docker container will run Metabase on <http://localhost:3000>.

You can run `./upgrade_metabase.sh` to take a backup of the Metabase state, i.e. dashboards, models, questions, etc.

## Development

To add support for more data (different CSV files), implement a new function in the `pkg/table` package.
