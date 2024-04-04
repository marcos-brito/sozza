Sozza is a CLI tool that helps populate databases using csv files. It uses yaml files and some especial symbols to map a csv to a certain database scheme.

# Mapping file

The mapping file uses the `yaml` format. We'r going to see some examples on how to write it.

Assume a csv with the followings fields:

- Email
- UserName
- Age
- HouseColor
- HouseSize

And a schema such as the following:

```sql
CREATE TABLE users (
    email TEXT,
    user_name TEXT,
    age int,
    house REFERENCES houses (house_id)
);

CREATE TABLE house (
    house_id PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    house_color TEXT,
    house_size TEXT,
);
```

The mapping file would look like this:

```yaml
user:
  - email: Email
    user_name: UserName
    age: Age
    house_id: __house__

house:
  - house_color: HouseColor
    house_size: HouseSize
```

If the field value has no special meaning as describe down below, then it must match a field from the csv otherwise a error will be returned.

## Table references

You can referece a table using `__table__`. If the table has [multiple insertions] you can referece a certain insertion using `__table__, X` where `X` is the insertion number. If no number is passed, it assumes it's `0`

```yaml
user:
  - email: Email
    user_name: UserName
    age: Age
    house_id: __house__, 1

house:
  - house_color: HouseColor
    house_size: HouseSize
  - house_color: HouseColor
    house_size: HouseSize
```

## Multiple insertions

A table can have multiple insertions going in:

```yaml
user:
  - email: Email
    user_name: UserName
    age: Age
    house_id: __house__
  - email: Email
    user_name: UserName
    age: Age
    house_id: __house__
```

There is no need for the fields to have the same values. As long as the database is happy about it, no error will occur.

## Formatted values

Fields values can also be formmated:

```yaml
user:
  - email: ./path/to/executable, UserName, HouseColor
    user_name: UserName
    age: Age
    house_id: __house__
  - email: Email
    user_name: UserName
    age: Age
    house_id: __house__
```

When evaluating the values, `UserName` and `HouseColor` will be passed as parameters to `./path/to/executable`. The `stdout` of the execution
will be inserted in the database. The parameters can be any thing or no thing at all.

# Building

First, make sure you installed:

- go
- some C compiler e.g clang, gcc
- A C standard library

> The C stuff is necessary because some of the database drivers uses `cgo`.

After that, just clone the repo and build:

```bash
git clone https://github.com/marcos-brito/sozza
cd sozza
go build -o sozza main.go
```

Now you can do whatever you want it the binary. One might want to run the `--help` command
