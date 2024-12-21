package main

import (
	"fmt"

	"github.com/dipghoshraj/media-service/node-manager/internal"
	"github.com/dipghoshraj/media-service/replication-manager/pkg"
)

func main() {
	fmt.Println("Hello from Service A!")
	fmt.Println(internal.Greet("User from A"))
	fmt.Println(pkg.HelperFunction())
}
