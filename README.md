## STRATA

A helper library for intuitively creating SQL strings. This was created to simplify the process of writing SQL query strings in codebases where that was already the norm
If one would like to read the data, they would need to interface with third-party drivers. A potential option (as a personal project, could include fleshing this out further to be a small ORM framework - but much is to be done, as this repo is a shallow hack of about 12 hours' worth of work) 

#### Requirements
* Go v1.15 and later
* A PostgreSQL database - as the SQL syntax is directly optimized for Postgres (In the case of other databases like Microsoft SQL Server, Entity Framework is mature enough to be the de-facto ORM)


#### Generating a query

The idea would be that the user should create a top-level Query object
```go
q := strata.Query{}
```
After furnishing the Query object with what it needs, a set of tables (base table and join tables) and a few other options. 
Should an object be required, it will fail to create a query string and a native error will be returned.
and from there add the tables that are required
```go
	q := strata.Query{}

	ts := strata.Table{
		Name: "township", Schema: "cadastral",
	}
	whereConditionName := "arbitrary_name"
	ts.AddFields(
		// Ordered set of fields
		strata.NumberField("_id"),
		strata.StringField(whereConditionName),
		strata.StringField("tsg_id"),
	)

	if err := ts.SetWhereConditions(
		whereConditionName, 
		strata.ILike, 
		"search r-value literal",
	); err != nil {
		return nil, err
	}

	q.SetBaseTable(ts)
	q.Limit = 25
	returnedSQLString, err := q.SQL()
```