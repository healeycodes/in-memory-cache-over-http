# ⛷️ in-memory-cache-over-http

> My blog post: [Cloning Memcached with Go](https://healeycodes.com/go/tutorial/beginners/showdev/2019/10/21/cloning-memcached-with-go.html)

<br>

[![](https://github.com/healeycodes/in-memory-cache-over-http/workflows/Go/badge.svg)](https://github.com/healeycodes/in-memory-cache-over-http/actions?query=workflow%3AGo)

<br>

An in-memory key/value cache server over HTTP with no dependencies.

Keys and values are strings. Integer math can be applied in some situations (like Memcached does).

The caching method is Least Recently Used (LRU).

<br>

### Install

`go get healeycodes/in-memory-cache-over-http`

<br>

### Setup

- Set your PORT environmental variable.
- Set APP_ENV to `production` to turn off logging.
- Set SIZE to limit the number of key/value pairs, (0 is default - unlimited).

```bash
# Linux/macOS
export PORT=8000
export APP_ENV=production

# Command Prompt
set PORT=8000
set APP_ENV=production

# PowerShell
$env:PORT = "8000"
$env:APP_ENV = "production"
```

- Run

`go run .\main.go`

- Build

`go build`

<br>

### Usage

Adding an expire parameter is always optional. Not setting it, or setting it to zero means that the key will not expire. It uses **Unix time**.

Example usage.

Set `name` to be `Andrew` with an expire time of `01/01/2030 @ 12:00am (UTC)`

GET `localhost:8000/set?key=name&value=Andrew&expire=1893456000` (204 status code)

Retrieve the value located at `name`.

GET `localhost:8000/get?key=name` (200 status code, body: `Andrew`)

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

The majority of tests are integration tests that test API routes while checking the underlying cache.

There are some unit tests for the cache.

<br>

Run tests recursively.

`go test ./...`

Example output.

```bash
?       healeycodes/in-memory-cache-over-http   [no test files]
ok      healeycodes/in-memory-cache-over-http/api       0.527s
ok      healeycodes/in-memory-cache-over-http/cache     0.340s
```

<br>


### Contributing

Feel free to raise any issues and pull requests 👍

There is no road map for this project. My main motivations were to learn more about Go!
