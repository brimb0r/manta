package main

import (
	"context"
	"log"

	"github.com/gagno/terraform-provider-manta/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/gagno/manta",
	})
	if err != nil {
		log.Fatal(err)
	}
}
