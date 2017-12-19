package main


import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/vitess"

	"fmt"
	"github.com/youtube/vitess/go/vt/vitessdriver"
	"github.com/youtube/vitess/go/vt/vtgate/vtgateconn"
	"time"
)

type Product struct {
	gorm.Model
	Code string
	Discount int8
	Price uint
}

func main() {

	keyspace   := "metering_keyspace"
	timeout    := 10
	vtgatePort := 15991

	addr := fmt.Sprintf("35.189.158.62:%d", vtgatePort)

	//Connect to vitess db
	c := vitessdriver.Configuration{
		Protocol:  *vtgateconn.VtgateProtocol,
		Address:   addr,
		Target:    keyspace + "@master", //"keyspace:shard@tablet_type"
		Timeout:   time.Duration(timeout) * time.Second,
		Streaming: false,
	}

	db, err := gorm.Open("vitess", c)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Discount: 0, Price: 1000})

	// Read
	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}