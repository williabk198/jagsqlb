[![go.mod](https://img.shields.io/github/go-mod/go-version/williabk198/jagsqlb)](go.mod)
[![go report](https://goreportcard.com/badge/github.com/williabk198/jagsqlb)](https://goreportcard.com/report/github.com/williabk198/jagsqlb)
[![test status](https://github.com/williabk198/jagsqlb/workflows/test/badge.svg)](https://github.com/williabk198/jagsqlb/actions/workflows/test.yaml)
[![LICENSE](https://img.shields.io/github/license/williabk198/jagsqlb)](LICENSE) 

# Just Another Go(JAG) SQL Builder

This library aims to provide an easy and simple way to build SQL queries that can be used with `database/sql`, or similar packages.


## Table of Contents

* [Usage](#usage)
  * [Select Builder](#select-builder)
  * [Insert Builder](#insert-builder)
  * [Update Builder](#update-builder)
  * [Delete Builder](#delete-builder)
* [Struct Tags](#struct-tags)
* [Marshalling Values](#marshalling-values)

## Usage

To create a new SQL builder it's as simple as this:

```go
import "github.com/williabk198/jagsqlb"

// ...

sqlBuilder := jagsqlb.NewSqlBuilder()
```

The value returned by `jagsqlb.NewSqlBuilder` can be reused as many times as you would like 
if multiple queries are required to be built.

### Select Builder

*__IMPORTANT:__* Wrapping columns in functions (e.g. `SUM(col1)`) is not supported. Which also means,
`GROUP BY` is also not supported in this version either. Additionally, type casting using `::` is also unsupported in this version. 

To create a simple `SELECT` statement like this: `SELECT * FROM "customers";` All you would need to write is this:

```go
queryStr, queryParams, err := sqlBuilder.Select("customers", "*").Build()
```

In this example, the returned values for `Build` would be `SELECT * FROM "customers";`, `nil` and `nil` respectively.

If you are using an SQL dialect that supports providing schema names in the table specifier (e.g. PostgreSQL),
that is also supported here:

```go
sqlBuilder.Select("public.customers", "*").Build()
```

Also, if you want to select a subset of columns, we can handle that too:

```go
sqlBuilder.Select("customers", "name", "street", "city", "state").Build()
```

On top of that, you can also query from multiple tables like so:

```go
queryStr, queryParams, err := sqlBuilder.Select("customers", "*").Table("orders", "*").Build()
```

This will result in a `queryStr` value of: 
```sql
SELECT "customers".*, "orders".* FROM "customers", "orders";
```

Aliasing table specifiers and column specifiers are also possible as well:

```go
queryStr, queryParams, err := sqlBuilder.Select(
  "persons AS p",
  "given_name AS first_name",
  "family_name AS last_name",
).Table("contact_info AS ci", "*").Build()
```

This will result in the following `queryStr` value:

```sql
SELECT "p"."given_name" AS "first_name", "p"."family_name" AS "last_name", "ci".* FROM "persons" AS "p", "contact_info" AS "ci";
``` 

By itself like this, this isn't too useful, but comes in rather handy when dealing with Joins.

*__IMPORTANT:__* `AS` _MUST_ be all uppercase. Otherwise, an error will be returned. This may change in future versions.

*__IMPORTANT:__* Alias names must not have any spaces in them. Even if quoted. An error will be returned otherwise.
This may change in future versions.

#### Join Clause

The `SELECT` statement builder can also handle creating joins. So, if you wanted to create a query like this:
```sql
SELECT "i".*, "s"."date" AS "sales_date" 
FROM "inventory" AS "i" 
LEFT JOIN "sales" AS "s" ON "i"."id" = "s"."inventory_id";
```

Then you would write the following code:
```go
// package & other imports...
import (
  "github.com/williabk198/jagsqlb/condition"
  "github.com/williabk198/jagsqlb/join"
)

// Other code...

queryStr, queryParams, err := sqlBuilder.Select("inventory AS i", "*").Join(join.TypeLeft, "sales AS s", join.On(
  condition.Equals("i.id", condition.ColumnValue("s.inventory_id")),
), "date AS sales_date").Build()
```

Hopefully it's easy to tell what's going on with the `Join` function. In the case that it isn't, here is a break down:

* The first parameter represents the type of join to perform, in which there are constants in the `join` package that
represent the most of the available join types. 

* The second parameter is the table to join. An alias can be provided like shown above and follows the same rules
  as aliases mentioned in the previous section.

* The third parameter represent how to join the tables. There are two functions in the `join` that will handle
  populating this parameter for you: `On` and `Using`.

* The forth parameter and beyond represent the columns from the table provided in parameter #2 
  to include in the result set of the query. Aliases can be provided like shown above and follows the same rules
  as aliases mentioned in the previous section.

Here's another example of building a query with a join. This time utilizing the `join.Using` function

```go
queryStr, queryParams, err := sqlBuilder.Select("inventory AS i", "*").Join(join.TypeLeft, "sales AS s", join.Using(
  "inventory_id",
), "date AS sales_date").Build()
```

#### Where Clause

You can also add a `WHERE` clause using the `SELECT` builder as well. That can be done like so:

```go
// package & other imports...
import "github.com/williabk198/jagsqlb/condition"

// Other code...
queryStr, queryParams, err := sqlBuilder.Select("person", "*").Where(
  condition.Between("dob", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
  condition.GroupedOr(
    condition.Equals("family_name", "Smith"),
    condition.Equals("family_name", "Lee"),
  )
).Or(condition.LessThan("dob", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC))).Build()
```

This code will produce the following for `queryStr` and `queryParams` values respectively: 
```sql
SELECT * FROM person WHERE "dob" BETWEEN $1 AND $2 AND ("family_name" = $3 OR "family_name" = $4) OR "dob" < $5;
```

```
[]any{
  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
  "Smith",
  "Lee",
  time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
}
```

You can also define a `WHERE` condition after a `JOIN` as well!
```go
queryStr, queryParams, err := sqlBuilder.Select("customers AS c", "*").Join(
  join.TypeInner,
  "subscriptions AS s",
  join.Using("customer_id"),
  "ends AS payment_due"
).Where(
  condition.LessThan("s.ends", time.Now().Truncate(24*time.Hour)),
  condition.GreaterThan("c.last_active", condition.ColumnValue("s.ends"))
).Build()
```

This code will have the following as the value for `queryStr` after execution(minus any new-line characters):
```sql
SELECT "c".*, "s"."ends" AS "payment_due" FROM "customer" AS "c"
INNER JOIN "subscriptions" AS "s" USING("customer_id")
WHERE "s"."ends" < $1 AND "c"."last_active" > "s"."ends"
```

With `queryParams` looking like this:
```
[]any{
  time.Now().Truncate(24*time.Hour) // The evaluation of this expression, of course.
}
```

As you can see, if you want to compare two columns, then you will need to use `condition.ColumnValue`. Otherwise, it will
get parameterized as a string value, which will cause erroneous behavior.

### Insert Builder

*__IMPORTANT:__* Using Select statements inside of the Insert Builder is not supported in this version.

There are a couple of ways that you can build and insert statement with this library.

Like this:

```go
queryStr, queryParams, err :=  sqlBuilder.Insert("inventory").Columns("name", "price").Values(
  []any{"Car", 18365.0},
).Returning("id").Build()
```

Or like this:

```go
type Inventory struct {
  ID int `jagsqlb:"id;omit"`
  ProductName string `jagsqlb:"name"`
  Price float64 `jagsqlb:"price"`
}

car := Inventory{
  ProductName: "Car"
  Price: 18365.0
}

queryStr, queryParams, err :=  sqlBuilder.Insert("inventory").Data(car).Returning("id").Build()
```

Both of the above code snippets will result in the same result in the same `queryStr` and `queryParams` result,
respectively:

```sql
INSERT INTO "inventory" ("name", "price") VALUES ($1, $2) RETURNING "id";
```

```go
[]any{"Car", 18365.0}
```

### Update Builder

Like the Insert Builder, the Update Builder also has two ways that to build out the query.

```go
queryStr, queryParams, err :=  sqlBuilder.Update("inventory").SetMap(map[string]any{
  "price": 19.99
}).Where(
  condition.Between("price", 19.75, 20.25),
).Build()
```

```go
type UpdateInvPrice struct {
  Price float64 `jagsqlb:"price"`
}

newInvPrice := UpdateInvPrice{
  Price: 19.99
}

queryStr, queryParams, err :=  sqlBuilder.Update("inventory").SetStruct(newInvPrice).Where(
  condition.Between("price", 19.75, 20.25),
).Build()
```

Both of the above code snippets will result in the same result in the same `queryStr` and `queryParams` result,
respectively:

```sql
UPDATE "inventory" SET "price"=$1 WHERE "price" BETWEEN $2 AND $3
```

```go
[]any{19.99, 19.75, 20.25}
```

### Delete Builder

To create a `DELETE` simple delete statement, all you'll need is this:

```go
queryStr, queryParams, err := sqlBuilder.Delete("customers").Build()
```

A `WHERE` clause can also be added like so:

```go
queryStr, queryParams, err := sqlBuilder.Delete("customers").Where(condition.Equals("id", customerID)).Build()
```

Additionally, a `USING` clause can be used as will which acts just like using `FROM` in a `SELECT` statement.

```go
queryStr, queryParams, err := sqlBuilder.Delete("customers AS c").Using("metadata.customers AS mc").Where(
  condition.Equals("c.id", condition.ColumnValue("mc.customerID")),
  condition.LessThan("mc.last_login", twoYearsAgo),
).Build()
```

## Struct Tags

As a part of this package, struct tags were added to make things easier to build `INSERT` and `UPDATE` queries.

To correlate a field to column name with a struct tag works very similarly to other ORM/ORM-like packages:

```go
type Person struct {
  DateOfBirth time.Time `jagsqlb:"dob"`
}
```

*__NOTE:__* If a field does not have a column name defined with the `jagsqlb` tag, then the name of field
will be used as the column name by default.

### Omitting Fields

If there are fields that you don't want to be included during the execution of the insert or update process,
then you annotate the field(s) in question by adding `;omit` after the column name like so:

```go
type TestData struct {
  ID uuid.UUID `;omit`
  DateOfBirth time.Time `jagsqlb:"dob"`
}
```

*__NOTE:__* As you can see with the above example, you do not have to provide a column name if so desired.

### Inlining Nested Structs

If the struct that you are using to insert or update entries in the database has a nested struct within it that
represents other columns within the table, you can use the `;inline` tag to ensure that those values are
properly included in the query.

For example:

```go
type NameData struct {
  GivenName  string `jagsqlb:"given_name"`
  FamilyName string `jagsqlb:"family_name"`
}

type Person struct {
  ID          uuid.UUID `jagsqlb:";omit"`
  Name        NameData  `jagsqlb:";inline"`
  DateOfBirth time.Time `jagsqlb:"dob"`
}

person := Person{
  Name: NameData{
    GivenName: "Some",
    FamilyName: "Guy"
  },
  DateOfBirth: time.Unix(0, 0),
}

queryStr, queryParams, err := sqlBuilder.Insert("person").Data(person).Build()
```

This will result in the following value for `queryStr`

```sql
INSERT INTO "person" ("given_name", "family_name", "dob") VALUES ($1, $2, $3);
```

## Marshalling Values

When using the Insert or Update Builders, if the value you want to store in the database is slightly different
from what is stored in the provided struct field, then you will need to implement the `QueryMarshaler` interface
for that value.

The `QueryMarshaler` interface that specifies that a `MarshalQuery` function to be implemented which takes in no arguments
and returns a `string` and an `error`.

Here's an example:

```go
type PronounData struct {
  Subject string
  Object  string
}

func (pd PronounData) MarshalQuery() (string, error) {
  return fmt.Sprintf("%s/%s", pd.Subject, pd.Object), nil
}

type NameData struct {
  GivenName  string `jagsqlb:"given_name"`
  FamilyName string `jagsqlb:"family_name"`
}

type Person struct {
  ID          uuid.UUID   `jagsqlb:";omit"`
  Name        NameData    `jagsqlb:";inline"`
  DateOfBirth time.Time   `jagsqlb:"dob"`
  Pronouns    PronounData `jagsqlb:"pronouns"`
}

examplePerson := Person{
  Name: NameData{
    GivenName:  "Testy",
    FamilyName: "McTesterson",
  }
  DateOfBirth: time.Unix(0, 0),
  Pronouns: PronounData{
    Subject: "they"
    Object:  "them"
  }
}

queryStr, queryParams, err := sqlBuilder.Insert("persons").Data(examplePerson).Build()
```

The `queryStr` and `queryParams` values will be the following, respectively:

```sql
INSERT INTO "persons" ("given_name", "family_name", "dob", "pronouns") VALUES ($1, $2, $3, $4);
```

```go
[]any{
  "Testy",
  "McTesterson",
  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local),
  "they/them",
}
```

*__IMPORTANT:__* `QueryMarshaler` can only convert a type to a `string` value. If you wish to convert something
into a non-string type, then you will need to do the conversion yourself and store the value as the appropriate
type within the struct.
