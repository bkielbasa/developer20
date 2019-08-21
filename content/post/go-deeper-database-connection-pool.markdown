---
title: Go deeper – Database connection pool
publishdate: 2019-01-25
categories: [Golang, Programming]
tags:
  - golang
  - database
  - connection pool
---
Golang uses a connection pool to manage opened connections for us. As a result, new connections are used when no free connection left and reuses them when golang finds an idle connection. The most important thing is that when two queries are called one by one it does not mean that the queries will use the same connection. It may be true if not in every case.

In the example below, you can find two queries which may seem to be executed in one connection. The problem is that the first query may use a different connection than the second one.

{{< highlight sql >}}
db.Exec('LOCK TABLE table1 WRITE;');
db.Exec("INSERT INTO table1 VALUES ($1, $2)", val1, val2);
{{< / highlight >}}

This may cause a bug which is hard to find and reproduce. If the INSERT statement uses not the same connection as the first one, it will produce an error. It may happen because the table1 table is blocked for writing. In this scenario, the second query will not succeed to execute. If the error will not be handled correctly it may lead lock the table until the connection is closed.

To avoid this we have to make sure that queries use the same connection from the connection pool. To achieve that a transaction should be used.

{{<  highlight go >}}
tx, err := db.BeginTx(ctx, nil)
tx.Exec("LOCK TABLE table1 WRITE;")
tx.Exec("INSERT INTO table1 VALUES ($1, $2)", val1, val2)
err = tx.Commit()
{{< / highlight >}}

When the transaction is created, a connection from the connection pool is related to it and only this one is used to execute the queries. Thanks to that we can be sure the same connection is used when we really want it.

## Managing connections

There is a way that the number of opened connections or idle connection is configured. To set the maximum number of open connections `db.SetMaxOpenConns(N)` should be used. By default, this value is set to 0 which means that there’s no limit on the client site. However, when the default value is used it does not mean that the limit cannot be reached on the database server site. It is an important thing worth remembering.

When the maximum open connection is reached the goroutinge that want to execute SQL statement will be waiting for a connection to back to the connection pool. This may take a long time or never happen.

Queries to the databases are sent using TCP protocol. However, opening a new TCP connection is costly. To minimize the cost, golang keeps some idle connections opened and ready for reuse. On the other hand, the max idle connection should not be too huge because, as we talked earlier, there is a limit of available open connections on the server side. Furthermore, other services can share the limit so the limit can be reached even faster than we assumed.
The maximum idle connections can be configured using `db.SetMaxIdleConns(N)` function. This value can help to optimize the number of opened connections. When the number of idle connections reaches the limit other connections are closed and resources come back to be reused. It is important when you have a quite low number of available connections on the server side.

## Summary

The `database/sql` library helps developers by managing the connection pool for them. On the other hand, developers have to remember how the managing connection works to prevent non-trivial bugs to find. Keeping queries which rely on each other in a single transaction can solve many problems.
