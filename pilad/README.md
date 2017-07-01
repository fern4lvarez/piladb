pilad
=====

**pilad** is the daemon that runs the **piladb** server, manages the main pila,
databases and stacks. It exposes a RESTful HTTP server that listens to requests
in order to interact with the engine.

> Note: pilad API does not come with a built-in `pretty` option. We encourage
> to use [`jq`](https://stedolan.github.io/jq/) to visualize JSON data on the terminal,
> or advanced CLI HTTP clients like [HTTPie](https://httpie.org/).

Endpoints
---------

### STATUS

#### GET `/`

Returns information about **piladb** and `200 OK`.

#### GET `/_ping`

Returns `pong` and `200 OK`.

#### HEAD `/_ping`

Returns  headers and `200 OK`.

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

### CONFIG

#### GET `/_config`

Returns `200 OK` and a representation of the configuration values
in key-value format.

```json
{
  "stacks": {
    "MAX_STACK_SIZE": 10
  }
}
```

Returns `400 BAD REQUEST` if there's an error serializing the config
response.

#### `GET /_config/$CONFIG_KEY`

> GET a Config value.

Returns `200 OK` and value associated to `$CONFIG_KEY`.

```json
{
  "element": 10
}
```

Returns `410 GONE` if configuration key does not exist.

Returns `400 BAD REQUEST` if there's an error serializing the config
response.

#### POST `/_config/$CONFIG_KEY` + `{"element":$CONFIG_VALUE}`

> SET a Config value.

Returns `200 OK` and the new value set to `$CONFIG_KEY`.

```json
{
  "element": 10
}
```

Returns `410 GONE` if configuration key does not exist.

Returns `400 BAD REQUEST` if `$CONFIG_VALUE` is not provided or there's
an error serializing the config response.

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
      "id": "91010edc-36f6-55cc-9b10-f2648eb2b322"
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
  "id": "91010edc-36f6-55cc-9b10-f2648eb2b322"
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
  "id": "91010edc-36f6-55cc-9b10-f2648eb2b322"
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
      "id":"ef7199db-821c-50df-8dc0-ddc77fc6397e",
      "name":"stack1",
      "peek":"foo",
      "size":1,
      "created_at":"2016-12-08T17:45:50.668575679+01:00",
      "updated_at":"2016-12-08T18:21:270.813642732+01:00",
      "read_at":"2016-12-08T18:21:270.813642732+01:00"
    },
    {
      "id":"15860a24-e97c-5a2a-be81-3d5066246cb6",
      "name":"stack2",
      "peek":8,
      "size":2,
      "created_at": "2016-12-08T17:48:65.122475579+01:00",
      "updated_at":"2016-12-08T18:16:120.4267723134+01:00",
      "read_at":"2016-12-08T18:17:32.456823273254+01:00"
    }
  ]
}
```

Returns `410 GONE` if the database does not exist.

Returns `400 BAD REQUEST` if there's an error serializing the stacks
response.

#### GET `/databases/$DATABASE_ID/stacks?kv`

Returns `200 OK` and a key-value representation of the stacks of
the database `$DATABASE_ID`, where key is the Name and value is the Peek.
You can use either the ID or the Name of the database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "stacks" : {
    "stack1":"foo",
    "stack2":8
  }
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
  "id": "91010edc-36f6-55cc-9b10-f2648eb2b322",
  "created_at": "2016-12-08T17:45:50.668575679+01:00",
  "updated_at": "2016-12-08T17:45:50.668575679+01:00",
  "read_at": "2016-12-08T17:45:50.668575679+01:00"
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
  "id": "91010edc-36f6-55cc-9b10-f2648eb2b322",
  "created_at": "2016-12-08T17:45:50.668575679+01:00",
  "updated_at": "2016-12-08T17:45:50.668575679+01:00",
  "read_at":"2016-12-08T18:17:32.456823273254+01:00"
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

#### GET `/databases/$DATABASE_ID/stacks/$STACK_ID?size`

> SIZE operation.

Returns the size of the `$STACK_ID` stack of database `$DATABASE_ID`, and
`200 OK`.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
6
```

Returns `410 GONE` if the database or stack do not exist.

#### GET `/databases/$DATABASE_ID/stacks/$STACK_ID?empty`

> EMPTY operation.

Returns true if the stack identify by `$STACK_ID` in database `$DATABASE_ID` is empty,
and `200 OK`.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
false
```

Returns `410 GONE` if the database or stack do not exist.

#### GET `/databases/$DATABASE_ID/stacks/$STACK_ID?full`

> FULL operation.

Returns true if the stack identify by `$STACK_ID` in database `$DATABASE_ID` is full,
and `200 OK`. Full means that the size of the Stack is equals or bigger to the `MAX_STACK_SIZE` config value, if set.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
false
```

Returns `410 GONE` if the database or stack do not exist.

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

Returns `406 NOT ACCEPTABLE` if the stack is full.

Returns `410 GONE` if the database or stack do not exist.

Returns `400 BAD REQUEST` if there's an error serializing the element.

#### POST `/databases/$DATABASE_ID/stacks/$STACK_ID?base` + `{"element":$ELEMENT}`

> BASE operation.

Puts an `ELEMENT` on the bottom of the `$STACK_ID` stack of database `$DATABASE_ID`, and
returns `200 OK`, and the based element.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "element": "this is an element"
}
```

Returns `406 NOT ACCEPTABLE` if the stack is full.

Returns `410 GONE` if the database or stack do not exist.

Returns `400 BAD REQUEST` if there's an error serializing the element.

#### POST `/databases/$DATABASE_ID/stacks/$STACK_ID?rotate`

> ROTATE operation.

Puts the bottommost element on the top of the `$STACK_ID` stack of database `$DATABASE_ID`,
and returns `200 OK`, and the rotated element. The element next to the former bottommost
element becomes the new one.
You can use either the ID or the Name of the stack and database, although the former
is used as default, the latter as fallback.

```json
200 OK
{
  "element": "this is an element"
}
```

Returns `204 NO CONTENT` if the stack is empty.

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
  "id": "91010edc-36f6-55cc-9b10-f2648eb2b322",
  "created_at": "2016-12-08T17:45:50.668575679+01:00",
  "updated_at": "2016-12-08T17:46:23.133256135+01:00",
  "read_at": "2016-12-08T17:46:23.133256135+01:00"
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
