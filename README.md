# Just Another Go(JAG) SQL Builder

This library aims to provide an easy ans simple way to build SQL queries.


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
`GROUP BY` is also not supported in this version. Additionally, type casting using `::` is also unsupported. 
These maybe added in the future.

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

This will result in a `queryStr` value of `SELECT "customers".*, "orders".* FROM "customers", "orders";`

Aliasing table specifiers and column specifiers are also possible as well:

```go
queryStr, queryParams, err := sqlBuilder.Select(
  "persons AS p", "given_name AS first_name", "family_name AS last_name",
).Table("contact_info AS ci", "*").Build()
```
This will result in a `queryStr` value of `SELECT "p"."given_name" AS "first_name", "p"."family_name" AS "last_name", "ci".* FROM "persons" AS "p", "contact_info" AS "ci";` By itself like this, this isn't too useful, but comes in rather handy when dealing with Joins.

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
  condition.GroupedOr(condition.Equals("family_name", "Smith"), condition.Equals("family_name", "Lee"))
).Or(condition.LessThan("dob", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC))).Build()
```

This code will produce the following for `queryStr`(minus the new line and tab characters): 
```sql
SELECT * FROM person 
WHERE "dob" BETWEEN $1 AND $2 AND (
  "family_name" = $3 OR "family_name" = $4
) OR "dob" < $5;
```

This will also populate the `queryParams` return variable with the values provided to each condition like so:
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
SELECT "c".*, "s"."ends" AS "payment_due" FROM "customer" AS "c" INNER JOIN "subscriptions" AS "s" USING("customer_id")
WHERE "s"."ends" < $1 AND "c"."last_active" > "s"."ends"
```

With `queryParams` looking like this:
```
[]any{
  time.Now().Truncate(24*time.Hour) // The evaluation of this expression, of course.
}
```

As you can see, if you want to compare two columns, then you will need to use `condtion.ColumnValue`. Otherwise, it will
get parameterized as a string value, which will cause erronious behavior.
