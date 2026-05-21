package main

import (
	"context"
	"log"

	"github.com/cyc115/terraform-provider-elevenlabs/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/cyc115/elevenlabs",
	})
	if err != nil {
		log.Fatal(err)
	}
}
