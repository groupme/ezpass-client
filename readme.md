# EZPass client

A pure Go client for interacting with GroupMe's authentication service.

## Development

    go get ./...
    go test ./...

## Usage

    import ezpass "github.com/groupme/ezpass-client"

### As a library
    // get token
    token := ezpass.Token(r) // r is *http.Request

    // to authenticate a user
    pass, err := ezpass.Get(token)
    fmt.Printf("Hi my name is %s", pass.User.Name)

    // to also check membership
    pass, err := ezpass.GetMembership(token, "group-id")
    fmt.Printf("You can call me %", pass.Membership.Nickname)

### As a Handler

You can also implement an `ezpassHandler`. In this example, appHandler will be
called with a `*ezpass.PAss` if successful and will not be called if
unsuccessful.

    type ezpassHandler func(http.ResponseWriter, *http.Request, *ezpass.Pass)

    func appHandler(w http.ResponseWriter, r *http.Request, p *ezpass.Pass) {
        // ...
    }

    http.Handle("/foo", ezpass.AuthHandler(appHandler))

## Error Handling

* `ezpass.ErrTimeout` - request timeout
* `ezpass.ErrUnauthorized` - token is invalid
* `ezpass.ErrNotFound` - token is valid, but user is not in group
* `ezpass.ErrUnknown` - an unknown error occurred

## Testing

Use the `net/http/httptest` package to simulate a real service:

    import (
      ezpass "github.com/groupme/ezpass-client"
      "github.com/groupme/ezpass-client/ezpasstest"
    )

    // start a test server and point ezpass at it
    ezpass.URL = ezpasstest.NewServer().URL

    // ... your tests here ...

