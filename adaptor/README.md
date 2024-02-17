## Note:

- When modifying the resp header in the server handler, please use `write` function to write the header. Otherwise, the modification will not take effect.

```go
func handler(resp http.ResponseWriter, req *http.Request) {
    resp.Header().Add("Content-Encoding", "test")
    _, err := resp.Write([]byte("Content-Encoding: test\n"))
    if err != nil {
        panic(err)
    }
}
```