# mini-chain


$ go run . 3000
$ go run . 3001

$ curl -X POST http://localhost:43000/tx \
    -H "Content-Type: application/json" \
    -d '{
    "from": "miner",
    "to": "alice",
    "amount": 10
    }'


$ go run . 3000 getchain

$ go run . 3001 getchain