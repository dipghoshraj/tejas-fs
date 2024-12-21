package main

import (
	"fmt"

	"github.com/dipghoshraj/media-service/replication-manager/pkg"
)

func main() {
	fmt.Println("Hello from Service B!")
	fmt.Println(pkg.HelperFunction())
}
