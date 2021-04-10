# Ticker Wall

For linux you will need: `libgl1-mesa-dev` and `xorg-dev` packages.

To run:
`go run cmd/client/*.go`

To run second screen(or more):
`go run cmd/client/*.go --screenindex=1`

### APIs

There is a RESTful HTTP API which you can interact with to update the ticker wall in real-time.

Updating presentation data:

```
POST /v1/presentation
{
    "tickerBoxWidth": 1400,
    "scrollSpeed": 20
}
```

This will slow down the scroll speed, and increase the width of each ticker symbol.
