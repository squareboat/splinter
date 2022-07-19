# Splinter

Splinter is a platform agnostic migration tool written in go

## Build Instructions (from source)

You need to have installed the following packages : `go,make`

Step 1. Clone this repository :- `git clone https://github.com/squareboat/splinter.git` and `cd` into it `cd splinter`

Step 2. Run `make clean` to clean the build binaries (if built previously ).

Step 3. Run `make build` to build for your system.

Step 4. Look for splinter binary in `./bin/splinter`. Copy this to path or use this as you want.

Step 5. To add to the path run `mv .bin/splinter /usr/local/bin` (For MacOS).

### Config

Splinter supports a wide variety of config file. Some of the formats that splinter supports are `.env,json,yaml,toml`.
You can override the default config by using a flag with splinter as following:-

###### Syntax :- `splinter [command] --config <path to config>`

The Path can be relative or absolute based on your preference.

Here is an example on how to provide the config file

`splinter migrate --config configs/splinter.json`

Note: Default location for config is `.env` in the folder from where you invoked splinter `splinter`

For example if you are in `/home/user/projects/test` and you call `splinter migrate`
spinter will look for `/home/user/projects/test/.env` file.

## Config Reference

| Key                        |                            Description                            |        Default |
| -------------------------- | :---------------------------------------------------------------: | -------------: |
| `SPLINTER_MIGRATIONS_PATH` |      location where splinter will look for migrations file.       | `./migrations` |
| `SPLINTER_DRIVER`          |                 Which sql database are you using.                 |     `postgres` |
| `SPLINTER_CONN_URI`        | A connection url for your databse. Should be based on your driver |         `none` |

## Flags Reference

Usage :-
`splinter [command] --[flag] <value>`

| Flag              |             Description              | Default |
| ----------------- | :----------------------------------: | ------: |
| `driver`          |            same as config            |  `none` |
| `uri`             |            same as config            |  `none` |
| `migrations-path` |            same as config            |  `none` |
| `config`          |   Path to the config for splitner    |  `.env` |
| `help`            | displays help for particular command |  `none` |

## Commands Reference

#### Create

Create a new migration file.

Usage:
`splinter create [flags]`

Examples:
`splinter create <filename1> <filename2>`
`splinter create create_user_table`
This will create `{SPLINTER_MIGRATIONS_PATH}/{timestamp}_{filename}.sql`
Write your SQL in the files created.

```
[up]

BEGIN;
CREATE TABLE IF NOT EXISTS users (id int);
COMMIT;

[down]
DROP TABLE IF EXISTS users;

```

Note: Transactions are not needed here because all the migrations are ran in a single transaction.

#### Migrate

Run all the migration that are pending in the system to database.

Usage:
`splinter migrate [flags]`

Aliases:
`migrate, up`

#### Rollback

Rollback the last migration that was applied to the database.

Usage:
`splinter rollback [flags]`

Aliases:
`rollback, down`
