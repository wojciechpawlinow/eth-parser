# eth-parser

Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Description 

Users not able to receive push notifications for incoming/outgoing transactions. By Implementing Parser interface we would be able to hook this up to notifications service to notify about any incoming/outgoing transactions.


## ETH Info 

Transactions are being fetched from block "0x0" to block "latest".

## Usage

```bash
git clone https://github.com/wojciechpawlinow/eth-parser.git
cd eth-parser
make run # or `go run cmd/server/main.go `

# listening at 0.0.0.0:8080
```

API:

```
GET /current-block
POST /subscribe
GET /address/{address}/transactions
```

### Sample output
```
❯ curl -X GET localhost:8080/current-block
{"current_block":20082405}

❯ curl -X POST localhost:8080/subscribe -i -d '{"address": "0xCcE1351B6553040894fAf0490d8B7879B035DeF9"}'                                                                                                                                                                                                                            
HTTP/1.1 200 OK
Date: Thu, 13 Jun 2024 11:00:43 GMT
Content-Length: 0

❯ curl -X GET localhost:8080/address/0xCcE1351B6553040894fAf0490d8B7879B035DeF9/transactions

[
  {
    "transactionIndex": "0x74",
    "hash": "0xd742eacd4786c5379a9109ae95daabe54673e7e89bbc42cba335975f841c5a53",
    "from": "0xcce1351b6553040894faf0490d8b7879b035def9",
    "to": "0x11111254369792b2ca5d084ab5eea397ca8fa48b",
    "gas": "0x74cee",
    "gasPrice": "0x7b2879100",
    "value": "0x0"
  },
  {
    "transactionIndex": "0x39",
    "hash": "0x9efb276666a9f97c21a2b37b5e73dde79f94def8aab5618fa4eb3295c453faa3",
    "from": "0xcce1351b6553040894faf0490d8b7879b035def9",
    "to": "0x3efa30704d2b8bbac821307230376556cf8cc39e",
    "gas": "0x114d8",
    "gasPrice": "0xf50a1254b",
    "value": "0x0"
  }
]
```

## Unit tests

```bash
make test
```

## Lint

```bash
make format