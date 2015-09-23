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

`GET` Returns the status of the currently running databases.

### `/databases/$DATABASE_ID`

### `/databases/$DATABASE_ID/stacks`

### `/databases/$DATABASE_ID/stacks/$STACK_ID`
