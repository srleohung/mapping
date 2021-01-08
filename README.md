# mapping (under development)

# Examples
- Case 1
```go
package main

import (
	"fmt"

	"github.com/srleohung/mapping/structure"
)

type Order struct {
	ID      int
	Store   StoreInformation
	Product ProductInformation
	Price   float64
}

type StoreInformation struct {
	ID   int
	Name string
}

type ProductInformation struct {
	ID   int
	Name string
}

type TransactionRecord struct {
	OrderID     int     `struct:"ID"`
	StoreID     int     `struct:"Store.ID"`
	StoreName   string  `struct:"Store.Name"`
	ProductID   int     `struct:"Product.ID"`
	ProductName string  `struct:"Product.Name"`
	Price       float64 `struct:"Price"`
}

func main() {
	var order Order = Order{
		ID:      1,
		Store:   StoreInformation{ID: 1, Name: "STORE_NAME_1"},
		Product: ProductInformation{ID: 1, Name: "PRODUCT_NAME_1"},
		Price:   9.99,
	}
	var record TransactionRecord
	structure.StructToStruct(order, &record)
    fmt.Printf("%+v\n", record)
    // output: {OrderID:1 StoreID:1 StoreName:STORE_NAME_1 ProductID:1 ProductName:PRODUCT_NAME_1 Price:9.99}
}
```
- Case 2
```go
package main

import (
	"fmt"

	"github.com/srleohung/mapping/structure"
)

type Order struct {
	ID      int `struct:"OrderID"`
	Store   StoreInformation
	Product ProductInformation
	Price   float64 `struct:"Price"`
}

type StoreInformation struct {
	ID   int    `struct:"StoreID"`
	Name string `struct:"StoreName"`
}

type ProductInformation struct {
	ID   int    `struct:"ProductID"`
	Name string `struct:"ProductName"`
}

type TransactionRecord struct {
	OrderID     int
	StoreID     int
	StoreName   string
	ProductID   int
	ProductName string
	Price       float64
}

func main() {
	var record TransactionRecord = TransactionRecord{
		OrderID:     1,
		StoreID:     1,
		StoreName:   "STORE_NAME_1",
		ProductID:   1,
		ProductName: "PRODUCT_NAME_1",
		Price:       9.99,
	}
	var order Order
	structure.StructToStruct(record, &order)
    fmt.Printf("%+v\n", order)
    // output: {ID:1 Store:{ID:1 Name:STORE_NAME_1} Product:{ID:1 Name:PRODUCT_NAME_1} Price:9.99}
}
```