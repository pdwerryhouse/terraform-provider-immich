package main

import (
	"context"
	"log"

	"terraform-provider-immich/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	opts := providerserver.ServeOpts{
		Address: "dwerryhouse.com/something/immich",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
