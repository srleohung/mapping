# mapping (under development)

# How To Use

## Installation

1. Add the client package your to your project dependencies (go.mod).
```bash
go get github.com/srleohung/mapping/structure
```
2. Add import `github.com/srleohung/mapping/structure` to your source code.

## Basic Example
The following example demonstrates how to use `github.com/srleohung/mapping/structure` to transform one structure into another.
- This is the source structure.
```go
type Order struct {
	ID          int
	Store       StoreInformation
	Product     []ProductInformation
	Price       float64
	CreatedTime time.Time
}

type StoreInformation struct {
	ID   int
	Name string
}

type ProductInformation struct {
	ID   int
	Name string
}
```
- This is the destination structure.
```go
type TransactionRecord struct {
	OrderID         int       `struct:"ID"`          // Mapping from Order.ID to OrderID
	StoreID         int       `struct:"Store.ID"`    // Mapping from Order.Store.ID to StoreID
	Store           Store     `struct:"Store"`       // Mapping from Order.Store(StoreInformation) to Store(Store)
	Product         []Product `struct:"Product"`     // Mapping from Order.Product(array of ProductInformation) to Product(array of Product)
	TransactionTime time.Time `struct:"CreatedTime"` // Mapping from Order.CreatedTime(time.Time) to TransactionTime(time.Time)
}

type Store struct {
	ID   string `struct:"Store.ID"`   // Mapping from Order.Store.ID(int) to ID(string)
	Name string `struct:"Store.Name"` // Mapping from Order.Store.Name to Name
}

type Product struct {
	ID    int     `struct:"Product.ID"`   // Mapping from Order.Product.ID to ID
	Name  string  `struct:"Product.Name"` // Mapping from Order.Product.Name to Name
	Price float64 `struct:"Price"`        // Mapping from Order.Price to Price
}
```
- This is the transformation structure function.
```go
	var source Order = Order{
		ID:          1,
		Store:       StoreInformation{ID: 1, Name: "STORE_NAME_1"},
		Product:     []ProductInformation{{ID: 1, Name: "PRODUCT_NAME_1"}, {ID: 2, Name: "PRODUCT_NAME_2"}, {ID: 3, Name: "PRODUCT_NAME_3"}},
		Price:       9.9,
		CreatedTime: time.Now(),
	}
	fmt.Printf("source %+v\n", source)
	/*
	source {ID:1 Store:{ID:1 Name:STORE_NAME_1} Product:[{ID:1 Name:PRODUCT_NAME_1} {ID:2 Name:PRODUCT_NAME_2} {ID:3 Name:PRODUCT_NAME_3}] Price:9.9 CreatedTime:2021-01-12 10:48:29.649173 +0800 HKT m=+0.000213016}
	*/
	var destination TransactionRecord
	structure.StructToStruct(source, &destination)
	fmt.Printf("destination %+v\n", destination)
	/*
	destination {OrderID:1 StoreID:1 Store:{ID:1 Name:STORE_NAME_1} Product:[{ID:1 Name:PRODUCT_NAME_1 Price:9.9} {ID:2 Name:PRODUCT_NAME_2 Price:9.9} {ID:3 Name:PRODUCT_NAME_3 Price:9.9}] TransactionTime:2021-01-12 10:48:29.649173 +0800 HKT m=+0.000213016}
	*/
```
