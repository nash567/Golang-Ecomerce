# Golang-Ecomerce
An Ecomerse golang Backend with GRPC  

# Folder Structure --

```
.
├── api
│   ├── http
│   ├── proto
│   │   └── gocomerse
│   │       ├── buf.lock
│   │       ├── buf.yaml
│   │       ├── product
│   │       │   └── product.proto
│   │       └── user
│   │           └── user.proto
│   └── rpc
│       └── gocomerse
│           └── user
│               ├── error.go
│               └── rpc.go
├── buf.gen.yaml
├── buf.work.yaml
├── build
│   └── db
│       └── migrations
│           ├── 1_UserTable_sql.sql
│           ├── 2_ProductTable.sql
│           ├── 3_OrderTable.sql
│           ├── 5_Address.sql
│           ├── 6_Payment.sql
│           ├── 7_User_Order_Table.sql
│           └── Readme.md
├── cmd
│   └── goComerse
│       ├── app
│       │   ├── app.go
│       │   └── svc.go
│       └── main.go
├── config
│   ├── config.go
│   └── Readme.md
├── config.json
├── config.yml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   ├── config.go
│   │   ├── db.go
│   │   ├── helper
│   │   └── migrations
│   ├── logger
│   │   ├── logger.go
│   │   └── model
│   │       ├── config.go
│   │       ├── interface.go
│   │       ├── level.go
│   │       └── types.go
│   ├── pb
│   │   └── gocomerse
│   │       ├── product
│   │       │   ├── product_grpc.pb.go
│   │       │   ├── product.pb.go
│   │       │   └── product.pb.gw.go
│   │       └── user
│   │           ├── user_grpc.pb.go
│   │           ├── user.pb.go
│   │           └── user.pb.gw.go
│   └── service
│       ├── auth
│       │   ├── error.go
│       │   ├── interceptor.go
│       │   └── service.go
│       ├── product
│       │   ├── model
│       │   │   ├── error.go
│       │   │   ├── interface.go
│       │   │   └── model.go
│       │   ├── repo
│       │   │   ├── query.go
│       │   │   └── repository.go
│       │   └── service.go
│       └── user
│           ├── helper
│           │   └── helper.go
│           ├── model
│           │   ├── error.go
│           │   ├── interface.go
│           │   └── model.go
│           ├── repo
│           │   ├── filter.go
│           │   ├── query.go
│           │   └── repository.go
│           └── service.go
├── MakeFile
├── pkg
│   └── error.go
└── README.md

```
