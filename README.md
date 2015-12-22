# yodlee-client

Yodlee client for Go.

View the docs.

## Installation

```
go get github.com/travisjeffery/yodlee-client
```

## Example

``` go
c := &yodlee.Client{
  Login:    os.Getenv("YODLEE_COB_LOGIN"),
  Password: os.Getenv("YODLEE_COB_PASSWORD"),
}
c.Authenticate()
token, err := c.GetUserSessionToken(login, pass)
if err != nil {
   panic(err);
}
transactions, errs := c.GetTransactions(token, yodlee.NewGetTransactionInput())
// ...
```

## License

MIT
