pilad
=====

**pilad** is the daemon that runs the **piladb** server, manages the main pila,
databases and stacks. It exposes a RESTful HTTP server that listens to requests
in order to interact with the engine.

> Note: pilad API does not come with a built-in `pretty` option. We encourage
  to use [`jq`](https://stedolan.github.io/jq/) to visualize JSON data on the terminal.

Endpoints
---------

### STATUS

#### GET `/_status`

Returns `200 OK` and a JSON document with the current piladb status.

```json
200 OK
{
  "status": "OK",
  "version": "511016882554615139ba590753af00519513f765",
  "pid": 26345,
  "host": "linux_amd64",
  "started_at": "2015-09-25T23:01:04.181146284+02:00",
  "running_for": 12.215756477,
  "memory_alloc": "1.28MiB",
  "number_goroutines": 3
}
```

### `DATABASES`

#### `GET /databases`

Returns `200 OK` and the status of the currently running databases.

```json
200 OK
{
  "number_of_databases": 3,
  "databases": [
    {
      "number_of_stacks": 0,
      "name": "db0",
      "id": "714e49277eb730717e413b167b76ef78"
    },
    {
      "number_of_stacks": 0,
      "name": "db1",
      "id": "93c6f621b761cd88017846beae63f4be"
    },
    {
      "number_of_stacks": 0,
      "name": "db2",
      "id": "5d02dd2c3917fdd29abe20a2c1b5ea1c"
    }
  ]
}

```

#### `GET /databases/$DATABASE_ID`

Returns `200 OK` and the status of database `$DATABASE_ID`.
You can use either the ID or the name of the database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "number_of_stacks": 0,
  "name": "db0",
  "id": "714e49277eb730717e413b167b76ef78"
}
```

Returns `410 GONE` if database does not exist.

#### `DELETE /databases/$DATABASE_ID`

Returns `204 NO CONTENT` and deletes database `$DATABASE_ID`.
You can use either the ID or the name of the database, although
the former is used as default, the latter as fallback.

Returns `410 GONE` if database does not exist.

#### `PUT /databases?name=$DATABASE_NAME`

Returns `201 CREATED` and creates a new $DATABASE_NAME database.

```json
201 CREATED
{
  "number_of_stacks": 0,
  "name": "db0",
  "id": "714e49277eb730717e413b167b76ef78"
}
```

Returns `400 BAD REQUEST` if `name` is not provided

Returns `409 CONFLICT` if `$DATABASE_NAME` already exists.

### STACKS

#### GET `/databases/$DATABASE_ID/stacks`

Returns `200 OK` and the status of the stacks of the database `$DATABASE_ID`.
You can use either the ID or the Name of the database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "stacks" : [
    {
      "id":"f0306fec639bd57fc2929c8b897b9b37",
      "name":"stack1",
      "peek":"foo",
      "size":1
    },
    {
      "id":"dde8f895aea2ffa5546336146b9384e7",
      "name":"stack2",
      "peek":8,
      "size":2
    }
  ]
}
```

Returns `410 GONE` if the database does not exist.

Returns `400 BAD REQUEST` if there's an error serializing the stacks
response.

#### PUT `/databases/$DATABASE_ID/stacks?name=$STACK_NAME`

Creates a new $STACK_NAME stack belonging to database $DATABASE_ID.

```json
201 CREATED
{
  "size": 0,
  "peek": null,
  "name": "stack",
  "id": "714e49277eb730717e413b167b76ef78"
}
```

Returns `410 GONE` if the database does not exist.

Returns `400 BAD REQUEST` if `name` is not provided.

Returns `409 CONFLICT` if `$STACK_NAME` already exists.

#### GET `/databases/$DATABASE_ID/stacks/$STACK_ID`

Returns the status of the `$STACK_ID` stack of database `$DATABASE_ID`, and `200 OK`.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "size": 0,
  "peek": null,
  "name": "stack",
  "id": "714e49277eb730717e413b167b76ef78"
}
```

Returns `410 GONE` if the database or stack do not exist.

#### GET `/databases/$DATABASE_ID/stacks/$STACK_ID?peek`

> PEEK operation.

Returns the peek of the `$STACK_ID` stack of database `$DATABASE_ID`, and
`200 OK`.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "element": "this is an element"
}
```

#### POST `/databases/$DATABASE_ID/stacks/$STACK_ID` + `{"element":$ELEMENT}`

> PUSH operation.

Pushes `ELEMENT` on top of the `$STACK_ID` stack of database `$DATABASE_ID`, and
returns `200 OK`, and the pushed element.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "element": "this is an element"
}
```

Returns `410 GONE` if the database or stack do not exist.

Returns `400 BAD REQUEST` if there's an error serializing the element.

#### DELETE `/databases/$DATABASE_ID/stacks/$STACK_ID`

> POP operation.

Pops the element on top of the `$STACK_ID` stack of database `$DATABASE_ID`, and
returns `200 OK`, and the popped element.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "element": "this is an element"
}
```

Returns `204 NO CONTENT` if the stack is empty and no element was popped.

Returns `410 GONE` if the database or stack do not exist.

#### DELETE `/databases/$DATABASE_ID/stacks/$STACK_ID?flush`

> FLUSH operation.

Flushes the content of the `$STACK_ID` stack of database `$DATABASE_ID`,
and returns `200 OK`, and the stack status.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "size": 0,
  "peek": null,
  "name": "stack",
  "id": "714e49277eb730717e413b167b76ef78"
}
```

Returns `410 GONE` if the database or stack do not exist.

#### DELETE `/databases/$DATABASE_ID/stacks/$STACK_ID?full`

> DELETE stack operation.

Deletes `$STACK_ID` stack from database `$DATABASE_ID`,
and returns `204 No Content`.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

Returns `410 GONE` if the database or stack do not exist.
