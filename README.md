## Usage

```shell
Run the following commands:
git clone https://github.com/ashtanko/octo-server

cd octo-server

make build
 
make up

make migrate-up

make run
```

At this time the server running at `http://127.0.0.1:8000`. It provides the following endpoints:

* `POST /transaction`: creates a new transaction, example request:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"state": "win", "amount": 10.15, "transactionId": "2111038f-dcd9-40b4-8ae5-18bfec41a66b"}' http://localhost:8000/transaction
```

* `GET /account/balance`: returns the account balance, example request:

```shell
curl -X GET  http://localhost:8000/account/balance
```
