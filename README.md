# How to run


## Directly with golang

If you have go installed on your machine (this app is tested with the following config go version go1.24.0 darwin/arm64) you can do the following in the ```src``` subdir

```
go run .
```

## With docker 

Just do

```
docker compose up
```

# Assumptions

I make the following assumptions:
1. I check for price and total to be > 0.
2. For error cases, I just return the response as plain text instead of structured json.