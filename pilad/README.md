pilad
=====

**pilad** is the daemon that runs the **piladb** server, manages the main pila,
databases and stacks. It exposes a RESTful HTTP server that listens to requests
in order to interact with the system.

Endpoints
---------

### `/_status`

`GET` Returns a JSON document with the current piladb status.

```
{
  "status": "OK",
  "version": "511016882554615139ba590753af00519513f765",
  "host": "linux_amd64",
  "started_at": "2015-09-25T23:01:04.181146284+02:00",
  "running_for": 12.215756477
}
```

### `/databases`

`GET /databases` Returns the status of the currently running databases.

```
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

`GET /databases/$DATABASE_ID` Returns the status of database $DATABASE_ID.

`PUT /databases/$DATABASE_NAME` Creates a new $DATABASE_NAME database.

### `/databases/$DATABASE_NAME/stacks`

### `/databases/$DATABASE_NAME/stacks/$STACK_NAME`
