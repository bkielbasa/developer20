---
layout: post
title:  Indexing in MySQL
mainPhoto: mysql.jpg
---

Why do we use indexes? Searching through a row in a sorted file with N length takes O(log2N) comparisons and the same number of reads from a filesystem which is heavy itself. However, tables in databases are not sorted which complicates the operation, Especially, if you have a lot of reads, updates and deletions on them. Writing the sorted version of the file (table) would dramatically slow the database down. There is one more thing which makes it even more complicated: every table may be sorted in more than one order. That’s why we use indexes which hold only attributes used in sorting and a reference to the place where the data are kept.

## Types of indexes in MySQL

In MySQL, we have few indexes available which can be used. Each of them has a different application and it’s a good idea to know a bit about them.
### BTree+ index

B-Tree+ is the most commonly used index in MySQL. It is useful when you have many reads and few writes in the table because reindexing of the table is very costly. There are few key things about B-Tree+:

* it’s a rooted tree
* every node does not have more than mm children
* every node (excluding root and leaves) has not less than m/2 children. m/2 is rounded down to a natural digit and called a factor of minimization.
* root, if it’s not a leaf, has at least two children.
* every leaf is always on the same level
* nodes, which are not a leaf heaving kk children has k−1 keys.
* node does not have more than m−1 keys. Leaf, which is not a root, at least has m/2 some keys.

Where mm is a number of nodes in one level. In practical usage, the mm value is very high which makes the tree very weight but very short. It helps to reduce the number of reads from a hard drive which is very costly and has a huge impact on a performance.

An example of the tree is shown in the picture below.

![BTree+](/assets/posts/BPlusTree.png)

That’s a theory. In practice you can create the index using a query

    CREATE INDEX test_indexing USING BTREE ON bkielbasa_index.indexing(id);


or in a less verbose version

    CREATE INDEX test_indexing ON bkielbasa_index.indexing(id);

That is because the B-tree+ index is a default one in MySQL. In order to create an index on more than one column, all of them must be NOT NULL. Remember that the primary key is always an index.

### RTREE index

R-Tree index is used with a purpose of indexing spatial data like geographical coordinates, rectangles or polygons. You can see how it looks like in the 2d picture shown below.


![RTree](/assets/posts/R-tree.png)

To create T-Tree index you can use CREATE SPATIAL INDEX like shown below:

    CREATE SPATIAL INDEX some_index ON some_table_with_geometry (geometry_coords);

### HASH index

Another index type available in MySQL database is hash index. The key thing you should know about it is that the index is kept in the memory. What’s important to stress is that HASH index has O(1)O(1) average search and the worst case search. Hash Indexes are the general recommendation for equality based lookups. The optimizer cannot use the hash index to speed up `ORDER BY` operations. (This type of index cannot be used to search for the next entry in order.) Additionally, it is the only whole index which can be used in the search. The InnoDB indexes can be used in open and range lookups which will be described further in the article.

Creating indexes is very similar to the above.

### FULLTEXT index

Fulltext index was introduced to InnoDB in MySQL 5.6 and helps in searching a text (`CHAR`, `VARCHAR`, or `TEXT` columns). How does it works and why it’s faster than a regular search? In Fulltext index, they split the text into words and make an index of the words and not of the whole text. This works a lot faster for text searches when looking for specific words.

    CREATE TABLE opening_lines (
           id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
           opening_line TEXT(500),
           author VARCHAR(200),
           title VARCHAR(200),
           FULLTEXT idx (opening_line)
           ) ENGINE=InnoDB;

    SELECT id, author, title FROM opening_lines WHERE MATCH(opening_line) AGAINST('Ishmael');

Remember – it works only on columns where you have already added the index. The fulltext index is not used very often because there are better alternatives, such as Lucene, Lucene with Compass or Solr, which have much more features and are faster than fulltext index.

To fully understand next chapters you should be familiar with command `EXPLAIN` in MySQL. EXPLAIN provides information on how MySQL will execute your code statements. It works with SELECT, DELETE, INSERT, REPLACE, and UPDATE statements. Here are columns you will see in the output of the command.

<table class="table table-striped">
<thead>
<tr>
<th>COLUMN</th>
<th>JSON NAME</th>
<th>MEANING</th>
</tr>
</thead>
<tbody>
<tr>
<td>id</td>
<td>select_id</td>
<td>The SELECT identifier</td>
</tr>
<tr>
<td>select_type</td>
<td>None</td>
<td>The SELECT type</td>
</tr>
<tr>
<td>table</td>
<td>table_name</td>
<td>The table for the output row</td>
</tr>
<tr>
<td>partitions</td>
<td>partitions</td>
<td>The matching partitions</td>
</tr>
<tr>
<td>type</td>
<td>access_type</td>
<td>The join type</td>
</tr>
<tr>
<td>possible_keys</td>
<td>possible_keys</td>
<td>The possible indexes to choose</td>
</tr>
<tr>
<td>key</td>
<td>key</td>
<td>The index actually chosen</td>
</tr>
<tr>
<td>key_len</td>
<td>key_length</td>
<td>The length of the chosen key</td>
</tr>
<tr>
<td>ref</td>
<td>ref</td>
<td>The columns compared to the index</td>
</tr>
<tr>
<td>rows</td>
<td>rows</td>
<td>Estimate of rows to be examined</td>
</tr>
<tr>
<td>filtered</td>
<td>filtered</td>
<td>Percentage of rows filtered by table condition</td>
</tr>
<tr>
<td>Extra</td>
<td>None</td>
<td>Additional information</td>
</tr>
</tbody>
</table>


The columns are more detailed described [in the docs](https://dev.mysql.com/doc/refman/5.7/en/explain-output.html#explain-output-columns). We will be focused on a few of them: Extra, key, key_len and rows. Please notice that multiple records may be returned for a single query. A simple example is shown below.

    EXPLAIN SELECT * FROM
    orderdetails d
    INNER JOIN orders o ON d.orderNumber = o.orderNumber
    INNER JOIN products p ON p.productCode = d.productCode
    INNER JOIN productlines l ON p.productLine = l.productLine
    INNER JOIN customers c on c.customerNumber = o.customerNumber
    WHERE o.orderNumber = 10101

where the output is:

    ********************** 1. row **********************
               id: 1
      select_type: SIMPLE
            table: l
             type: ALL
    possible_keys: NULL
              key: NULL
          key_len: NULL
              ref: NULL
             rows: 7
            Extra:
    ********************** 2. row **********************
               id: 1
      select_type: SIMPLE
            table: p
             type: ALL
    possible_keys: NULL
              key: NULL
          key_len: NULL
              ref: NULL
             rows: 110
            Extra: Using where; Using join buffer
    ********************** 3. row **********************
               id: 1
      select_type: SIMPLE
            table: c
             type: ALL
    possible_keys: NULL
              key: NULL
          key_len: NULL
              ref: NULL
             rows: 122
            Extra: Using join buffer
    ********************** 4. row **********************
               id: 1
      select_type: SIMPLE
            table: o
             type: ALL
    possible_keys: NULL
              key: NULL
          key_len: NULL
              ref: NULL
             rows: 326
            Extra: Using where; Using join buffer
    ********************** 5. row **********************
               id: 1
      select_type: SIMPLE
            table: d
             type: ALL
    possible_keys: NULL
              key: NULL
          key_len: NULL
              ref: NULL
             rows: 2996
            Extra: Using where; Using join buffer
    5 rows in set (0.00 sec)

While studying the `EXPLAIN` output for performance it is important to fetch as small columns as it’s possible (the rows column) and always use an index. Using indexes indicates a small number of reads from the hard drive what’s less expensive. Using a low number of rows decreases the complexity of searching. What’s interesting, you may get the information in JSON format. You can achieve that by passing `format=json` to the query.

```sql
    explain format=json SELECT * FROM film WHERE (film_id BETWEEN 1 and 10) or (film_id BETWEEN 911 and 920)
```

    ********* 1. row *********
    EXPLAIN: {
      "query_block": {
        "select_id": 1,
        "cost_info": {
          "query_cost": "10.04"
        },
        "table": {
          "table_name": "film",
          "access_type": "range",
          "possible_keys": [
            "PRIMARY"
          ],
          "key": "PRIMARY",
          "used_key_parts": [
            "film_id"
          ],
          "key_length": "2",
          "rows_examined_per_scan": 20,
          "rows_produced_per_join": 20,
          "filtered": "100.00",
          "cost_info": {
            "read_cost": "6.04",
            "eval_cost": "4.00",
            "prefix_cost": "10.04",
            "data_read_per_join": "15K"
          },
          "used_columns": [
            "film_id",
            "title",
            "description",
            "release_year",
            "language_id",
            "original_language_id",
            "rental_duration",
            "rental_rate",
            "length",
            "replacement_cost",
            "rating",
            "special_features",
            "last_update"
          ],
          "attached_condition": "((`film`.`film_id` between 1 and 10) or (`film`.`film_id` between 911 and 920))"
        }
      }
    }

In forks of MySQL like MariaDB or Percona the feature is even more useful because it shows even more information.

## Example

When you have a general overview on the topic, let’s check how it works in practice. I’ve created two tables: `clients`

```sql
    CREATE TABLE `clients` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `username` varchar(255) NOT NULL,
      `password` varchar(255) NOT NULL,
      `email` varchar(255) NOT NULL,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

and `clients_indexed`

```sql
    CREATE TABLE `clients_indexed` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `username` varchar(255) NOT NULL,
      `password` varchar(255) NOT NULL,
      `email` varchar(255) NOT NULL,
      PRIMARY KEY (`id`),
      KEY `username_idx` (`username`,`email`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8
```



and added to each of them exactly 200002 identical records. The only difference in the schema is that the second table has an index called username_idx. Maybe let’s start with a trivial example. Our task is to find a record with` ID = 1053`.

and added exactly 200002 identical records to each of them. The only difference in the schema is that the second table has an index called username_idx. Maybe let’s start with a trivial example. Our task is to find a record with `ID = 1053`.

<table class="table table-striped">
<thead>
<tr>
<th>QUERY</th>
<th>RESULT</th>
</tr>
</thead>
<tbody>
<tr>
<td>SELECT * FROM clients where id = 1053;</td>
<td>1 row in set (0.00 sec)</td>
</tr>
<tr>
<td>SELECT * FROM clients_indexed where id = 1053;</td>
<td>1 row in set (0.00 sec)</td>
</tr>
</tbody>
</table>

We could not expect a different result. In Both queries, we used the primary key. More interesting results we’ll get when we try to find it in a column without any index.

<table class="table table-striped">
<thead>
<tr>
<th>QUERY</th>
<th>RESULT</th>
</tr>
</thead>
<tbody>
<tr>
<td>select * from clients where username = 'username_0.12939831440280966';</td>
<td>1 row in set (0.35 sec)</td>
</tr>
<tr>
<td>select * from clients_indexed where username = 'username_0.12939831440280966';</td>
<td>1 row in set (0.00 sec)</td>
</tr>
</tbody>
</table>

In the table without the index in username column we had to wait much longer than in the second table. That’s because MySQL had to read all the columns from the hard drive and try to match our criteria. Here is an explanation of both queries:

Table without an index:

<table class="table table-striped">
<thead>
<tr>
<th>ID</th>
<th>SELECT_TYPE</th>
<th>TABLE</th>
<th>PARTITIONS</th>
<th>TYPE</th>
<th>POSSIBLE_KEYS</th>
<th>KEY</th>
<th>KEY_LEN</th>
<th>REF</th>
<th>ROWS</th>
<th>FILTERED</th>
<th>EXTRA</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>SIMPLE</td>
<td>clients</td>
<td>NULL</td>
<td>ALL</td>
<td>NULL</td>
<td>NULL</td>
<td>NULL</td>
<td>NULL</td>
<td>776890</td>
<td>10.00</td>
<td>Using where</td>
</tr>
</tbody>
</table>

Table with index:

<table class="table table-striped">
<thead>
<tr>
<th>ID</th>
<th>SELECT_TYPE</th>
<th>TABLE</th>
<th>PARTITIONS</th>
<th>TYPE</th>
<th>POSSIBLE_KEYS</th>
<th>KEY</th>
<th>KEY_LEN</th>
<th>REF</th>
<th>ROWS</th>
<th>FILTERED</th>
<th>EXTRA</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>SIMPLE</td>
<td>clients_indexed</td>
<td>NULL</td>
<td>ref</td>
<td>username_idx</td>
<td>username_idx</td>
<td>767</td>
<td>const</td>
<td>1</td>
<td>100.00</td>
<td>NULL</td>
</tr>
</tbody>
</table>

It exactly shows why the second search was much faster – it filtered 100% of the rows! In the table without index before analyzing the data MySQL could filter only 10% of all records. It means that MySQL had to read from hard drive 90% of the records to try to match the’ where’ condition. That’s a huge difference. Please note that in the indexed table the engine did not even need to use ‘where’ condition to compare columns with the condition – it just already got what was needed!

What’s important to be mention is that you can use index to find values from the beginning. It means that in the query above the index was used but in the queries below it will not:

```sql
     --- won't use index
    SELECT * from clients where username LIKE '%a%';
    SELECT * from clients where username LIKE '%a';
    SELECT * from clients where username LIKE 'a%a';
```



BTree+ Indexes are useful in three kinds of look-ups:

* point lookup – an example of it you can find in queries above. Point lookups are look-ups where you use equal sign for example: where username = 'admin'.
* open range lookup – it’s a search where you specify start or end of the indexed value. Example where id > 10 or where price < 300.
* close range lookup – the same as open range lookup but you define both beginning and end of the indexed values. Example where id > 10 and id < 100

## Cons of using indexes

Using indexes has not only advantages but disadvantages, too. Below you have list of a few of them.

### Indexes take additional space on the hard drive
On huge tables or in case of badly designed indexes it may be a big deal. I think it’s visible very well in tables from this article.

    -rw-r-----   1 bartlomiejkielbasa  admin       8664 Oct  3 19:09 clients.frm
    -rw-r-----   1 bartlomiejkielbasa  admin  100663296 Oct  3 19:23 clients.ibd
    -rw-r-----   1 bartlomiejkielbasa  admin       8664 Oct  3 19:10 clients_indexed.frm
    -rw-r-----   1 bartlomiejkielbasa  admin  192937984 Oct  3 19:25 clients_indexed.ibd
    -rw-r-----   1 bartlomiejkielbasa  admin         61 Oct  3 19:09 db.opt

What is idb file?

> By default, all InnoDB tables and indexes are stored in the system tablespace. As an alternative, you can store each InnoDB table and associated indexes in its own data file. This feature is called “file-per-table tablespaces” because each table has its own tablespace, and each tablespace has its own .ibd data file.

Read more in the official docs. In other words, ibd files contain indexes. The file with indexes for `clients_indexed` table takes almost twice more size than indexes for the first table!

### Write operations take more time

The indexes must be updated in this situation which is an additional operation on your hard drive. You will have more performance improvement on tables with a huge number of reads and a low number of updates, deletes and so on.

## Summary

Indexes can save our time but can be problematic in cases of very huge tables or if they are not well designed. And remember – Use indexes find specific records – not to find all of them.
