---
title: "Go libs I don't use but are popular"
publishdate: 2024-09-10
short: An overview of widely used Go libraries such as sqlmock, GORM, and gorilla/mux, with a discussion of their drawbacks in real-world applications. The post explains when these tools are helpful and when avoiding them leads to cleaner, more maintainable code.
categories: 
    - Golang
    - Programming
tags:
  - golang
---

## Introduction

In the Go programming ecosystem, developers have a plethora of libraries available for solving common problems. However, some libraries may not always be the best fit for every project or developer preference. This article highlights a few Go libraries that I personally avoid using, along with the reasons behind these choices. The intention is not to discourage the use of these libraries universally but to shed light on potential challenges that may arise when using them, especially in larger or more complex projects.

## Libraries I Avoid Using

### 1. **sqlmock**

[sqlmock](https://github.com/DATA-DOG/go-sqlmock) is a library that allows developers to mock SQL queries and test database interactions without connecting to a real database. It implements the [sql/driver](https://godoc.org/database/sql/driver) interface, making it easy to validate the execution of specific SQL queries, ensuring they run in the desired order with the expected parameters.

At first glance, this seems ideal for testing repositories that rely on SQL queries. It enables achieving high test coverage, and the tests are fast since no actual database interactions occur. However, there are several drawbacks.

Imagine a scenario where SQL queries are generated dynamically through an SQL builder. We want to test that the generated query executes correctly. With sqlmock, there are two common approaches:

- Provide the entire query to the `mock.ExpectExec` function.
- Provide a partial query to `mock.ExpectExec`.

Both approaches come with significant downsides. If you supply the entire query, your tests become highly brittle, as even minor changes (e.g., adding a new column) will require updating multiple tests. This becomes especially problematic in larger projects with many tables and tests. Conversely, using a partial query risks missing potential issues, such as typos or logic errors, in the untested portions of the query. These errors may only surface during manual testing, end-to-end tests, or, worst-case scenario, in production.

Additionally, sqlmock creates overly rigid tests. Even if changes to the code do not impact the actual query execution, the tests may still fail due to added queries or parameters.

In small projects, these issues might not pose a significant challenge, but in larger applications where tables are reused in various contexts, they can lead to headaches.

#### **Alternative Approach**

I prefer using separate integration tests to verify database interactions. Though slower, they provide greater confidence by running against a real database. For an example of this approach, you can refer to [Writing tests in Go (business apps)](https://developer20.com/testing-in-go/).

### 2. **GORM**

GORM is a popular ORM in the Go community, but I avoid using it unless necessary. Let's explore why many developers choose [GORM](https://gorm.io/), and why I personally refrain from relying on it.

#### **Migrations**

One of GORM’s key selling points is automatic schema generation and migrations. It allows you to define database schema through Go structs and automatically generate and apply migrations. This feature is perfect for quickly bootstrapping projects, as it saves time and effort.

```go
db.AutoMigrate(&User{}, &Product{}, &Order{})
```

However, this convenience diminishes as the project grows. In larger applications, more control is needed over how schema changes are handled. While GORM offers a [Migrator interface](https://gorm.io/docs/migration.html#Migrator-Interface) for manually managing migrations, at that point, the advantages of using an ORM diminish, as manual migrations become a necessity sooner than later.

#### **Fetching and Saving Data**

Fetching data with GORM is quite simple and convenient, especially when compared to Go’s standard library.

```go
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// Generates:
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;
```

However, a library like [sqlx](https://jmoiron.github.io/sqlx/) offers similar ease with additional features, such as named parameters, while staying closer to the standard `database/sql` package. This allows for more control over query execution, which can be crucial in performance-sensitive applications.

```go
people := []Person{}
db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
```

I prefer using `sqlx` over GORM, as it gives me greater visibility and control over how queries are executed, which is critical in more complex scenarios.

#### **Relations**

GORM simplifies defining relationships between structs and tables, automatically handling foreign keys and joins.

```go
type User struct {  
  gorm.Model  
  Name      string  
  CompanyID int  
  Company   Company  
}  
  
type Company struct {  
  ID   int  
  Name string  
}
```

While this sounds great, the reality in production is often more complex. Issues like the infamous "n+1 problem" frequently arise, leading to performance bottlenecks. GORM’s preloading feature addresses this issue partially, but often results in overly complex queries with excessive joins.

### 3. **Non-Standard Mux Server**

In the past, I frequently used [gorilla mux](https://github.com/gorilla/mux) for routing due to its named parameters and API compatibility with Go’s standard library. However, with the release of [Go 1.22](https://go.dev/blog/routing-enhancements), the standard library now includes a pattern matcher that handles routing effectively, rendering gorilla mux largely unnecessary.

Other alternatives, such as [fasthttp](https://github.com/valyala/fasthttp), prioritize speed over standard library compatibility. However, cases where the router becomes a bottleneck are rare, with database or I/O operations typically being the root cause of performance issues. As a result, I tend to stick with Go’s built-in `net/http` package for most projects.

## Summary

While libraries like sqlmock, GORM, and gorilla mux have their merits, they can introduce complexity or rigidity in larger, more dynamic projects. In my experience, relying on more straightforward or flexible alternatives, such as integration testing, `sqlx`, or Go’s standard `net/http` package, leads to more maintainable and scalable solutions. It’s essential to consider the trade-offs of using specific libraries and evaluate whether their features align with your project’s long-term needs.
