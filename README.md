# StatusThing
StatusThing is an open source status page application

It originally started life as a quick thought exercise on a basic status page service for internal status pages.

## WARNING
This application is in active development. At the time of this readme update, the only implemented store is in-memory. Data is lost at shutdown.
Documentation below this section may be out of state but I'll try and not let that happen.

This section will go away when that situation changes at a minimum requiring:

- [X] a non-in-memory store implementation 
- [ ] simple ability to create your own binary linking to your own store implementation
- [ ] documentation updates
- [ ] ability to support bug requests

At that point, I will feel comfortable with folks using it and able to support it properly

## QuickStart
`go run cmd/statusthing-api/main.go`

Everything is stored in memory right now and is lost at shutdown. 

## Concepts
All core concepts are represented as protobuf types in the file `proto/statusthing/v1/types.proto`. If you've never worked with protobuf before, that's fine as you never need to deal with anything protobuf-specific to use the service.

The following classifications will map to a protobuf type under the covers.

### Items
Items are "things" that have a `Status`. Some systems call these "services" or "components". These are the things you want to communicate information about to others.

Items have, at a minimum, a name and a `Status`. They can optionally have:
- a description
- Notes

### Status
A status is a user-driven concept of the state of something. It will have at a minimum a name and a "kind".

Currently there are the following kinds of statuses:
- Up
- Down
- Warning
- Unavailable 
- Available
- Investigating
- Observing

### Notes
Note are updates about an `Item`. These mostly align with the concept of a status update.

## API
The API can be interacted with in multiple ways:

- a gRPC client
- HTTP Post + JSON

Both types of interaction happen over the same url.

### gRPC API
gRPC services are defined in `proto/statusthing/v1/services.proto`.

Currently I only generate go gRPC clients

### HTTP API
The HTTP API is provided automatically by buf-connect. It maps as follows:

- HTTP Path: `/statusthing.v1.<ServiceName>/<RpcName>`
- HTTP Verb: `POST`
- Content-Type: `application/json`

Every request requires a body and at a minimum it must be an empty JSON object `{}` and API response will always return at a minimum an empty JSON object `{}`

Mapping of protobuf messages follows a predictable pattern and any empty fields are not returned in the JSON object

### Example 1: `Item`
The protobuf message for `Item` looks like so:

```proto3
// Item represents a status page entry
message Item {
    // the unique identifier
    string id = 1;
    // the name
    string name = 2;
    // the description
    string description = 3;
    // the status
    Status status = 4;
    // any associated notes
    repeated Note notes = 5;

    Timestamps timestamps = 15;
}
```

and will be represented as

```json
{
    "id":"my-id",
    "name":"my-name",
    "description":"my-desc",
    "status":{"id","my-status-id","name":"my-status-name","kind":1},
    "notes":[{"id":"my-note-id","text":"this thing is broken"}],
    "timestamps:"{"created":"<some timestamp>", "updated":"<some timestamp>"}
}
```

### Example 2: `FindItemsRequest`
The protobuf message for `FindItemsRequest` looks like so

```proto3
message FindItemsRequest {
    // return results having a status with any of the provided ids
    repeated string status_ids = 1;
    // return results having any of the provided [StatusKind] 
    repeated statusthing.v1.StatusKind kinds = 2;
    // by default notes are note returned. extended will include notes
    bool extended = 14;
}
```

This will translate to the follow request body (JSON):

```json
{
    "status_ids":["id1","id2"],
    "extended":true,
    "kinds":["STATUS_KIND_UP","STATUS_KIND_DOWN"]
}
```
Alternately, enums (the `kinds` field above) can also be provided as the integer value for a `StatusKind`: 

```json
{"kinds":[1,2]}
```

The choice is personal preference.
