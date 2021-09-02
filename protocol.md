# Pglet WebSockets protocol

## Register Web Client

```json
"id": "",
"action": "registerWebClient",
"payload": {
    "pageName": "",
    "pageHash": "",
    "sessionID": ""
}
```

## Host client

### Register Host Client

Request:

```json
"id": "{correlation-id}", // request/response correlation ID
"action": "registerHostClient",
"payload": {
    "hostClientID": "", // set if re-connecting to the same host client
    "pageName": "",
    "isApp": true,
    "authToken": "",
    "permissions": ""
}
```

Response:

```json
"id": "{correlation-id}",
"action": "",
"payload": {
    "hostClientID": "", // generated on first connect if not set
    "sessionID": "0",   // always "0"
    "pageName": "",     // parsed/normalized page full name
    "error": ""         // set if there was an error registering host agent
}
```

### New session started

The message is sent to a host client when a new app session is started.

```json
"id": "",
"action": "sessionCreated",
"payload": {
    "pageName": "",
    "sessionID": ""
}
```

### Page modification command

Request from a host client:

```json
"id": "{correlation-id}",
"action": "pageCommandFromHost",
"payload": {
    "pageName": "",
    "sessionID": "",
    "command": {
        "Indent": 0,
        "Name": "add" // mandatory command name
        "Values": ["value_1", "value_2", ...],
        "Attrs": {
            "attr_1": "value_1",
            "attr_2": "value_2",
            ...
        }
        "Lines": ["line1", ...]
    }
}
```

Response:

```json
"id": "{correlation-id}",
"action": "",
"payload": {
    "result": "", // control property value or space-delimited IDs of newly added controls
    "error": ""
}
```

### Page modification commands batch

Request from a host client:

```json
"id": "{correlation-id}",
"action": "pageCommandsBatchFromHost",
"payload": {
    "pageName": "",
    "sessionID": "",
    "commands": [
        {
            "Indent": 0,
            "Name": "add"
            "Values": ["value_1", "value_2", ...],
            "Attrs": {
                "attr_1": "value_1",
                "attr_2": "value_2",
                ...
            }
            "Lines": ["line1", ...]
        },
        ...
    ]
}
```

Response:

```json
"id": "{correlation-id}",
"action": "",
"payload": {
    "results": ["", "", ...],
    "error": ""
}
```

### Notify server about inactive app

One-way message from a host to Pglet server:

```json
"id": "",
"action": "inactiveAppFromHost",
"payload": {
    "pageName": ""
}
```

### Page event

One-way message from Pglet server to a host:

```json
"id": "",
"action": "pageEventToHost",
"payload": {
    "pageName": "",
    "sessionID": "",
    "eventTarget": "",
    "eventName": "",
    "eventData": ""    
}
```