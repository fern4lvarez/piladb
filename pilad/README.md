pilad
=====

**pilad** is the daemon that runs the **piladb** server, manages the main pila,
databases and stacks. It exposes a RESTful HTTP server that listens to requests
in order to interact with the system.

Endpoints
---------

### `/_status`

`GET` Returns a JSON document with the current piladb status.

### `/databases`

`GET /databases` Returns the status of the currently running databases.

`GET /databases/$DATABASE_ID` Returns the status of database $DATABASE_ID.

`PUT /databases/$DATABASE_NAME` Creates a new $DATABASE_NAME database.

### `/databases/$DATABASE_NAME/stacks`

### `/databases/$DATABASE_NAME/stacks/$STACK_NAME`
