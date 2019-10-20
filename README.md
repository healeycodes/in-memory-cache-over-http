# ⛷️

An in memory key/value cache over HTTP with no dependencies. Kinda like memcached lite.

The caching method is Least Recently Used (LRU).

[Install](#Install)
[Setup](#Setup)
[Usage](#Usage)
[Methods](#Methods)
[Tests](#Tests)

<br>

### Install

`go get healeycodes/in-memory-cache-over-http`

<br>

### Setup

Set your PORT environmental variable.

```bash
# Linux/macOS
export PORT=8000

# Command Prompt
set PORT=8000

# PowerShell
$env:PORT = "8000"
```

Run.

`go run .\main.go`

<br>

### Usage

The API is accessed via GET request.

Keys and values are strings. Integer math can be applied in some situations (like memcached).

Adding an expire parameter is always optional. Not setting it, or setting it to zero means that the key will not expire. It uses **Unix time**.

Example usage.

Set `name` to be `Andrew` with an expire time of `10/20/2019 @ 1:23pm (UTC)`

GET `localhost:8000/set?key=name&value=Andrew&expire=1571577784` (204 status code)

Retrieve the value located at `name`.

GET `localhost:8000/get?key=Andrew` (200 status code, body: `Andrew`)

<br>

### Methods

#### Set (params: key, value, expire) `/set`

Set a key/value. Existing will be overwritten.

<br>

#### Get `/get`

Get a value from a key.

<br>

#### Delete (params: key) `/delete`

Delete a key.

<br>

#### CheckAndSet (params: key, value, expire, compare) `/checkandset`

Set a key/value if the current value at that key matches the compare.

If no existing key, set the key/value.

<br>

#### Increment (params: key, value, expire) `/increment`

Increment a value. Both the existing value and the new value amount should be integers.

If no existing key, set the key/value.

<br>

#### Decrement (params: key, value, expire) `/decrement`

Decrement a value. Both the existing value and the new value amount should be integers.

If no existing key, set the key/value.

<br>

#### Append (params: key, value, expire) `/append`

Concatenates the new value onto the existing.

If no existing key, set the key/value.

<br>

#### Prepend (params: key, value, expire) `/prepend`

Concatenates the existing value onto the new value.

If no existing key, set the key/value.

<br>

#### Flush `/flush`

Clear the cache. Delete all keys and values.

<br>

#### Stats `/stats`

Return statistics about the cache.

Example.

```json
{
    "keyCount": 1,
    "maxSize": 0
}
```

<br>

### Tests

**TODO**
