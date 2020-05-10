package main

import (
	"fmt"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/providers"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)


func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		config := config.New(ctx, "")
		org := config.Require("org")
		clusterProject := config.Require("clusterProject")

		// Get stack reference
		slug := fmt.Sprintf("%s/%s/%v", org, clusterProject, ctx.Stack())
		stackRef, err := pulumi.NewStackReference(ctx, slug, nil)

		if err != nil {
			return fmt.Errorf("Error getting stack reference")
		}

		kubeConfig := stackRef.GetOutput(pulumi.String("kubeconfig"))

		// provider init
		provider, err := providers.NewProvider(ctx, "k8sprovider", &providers.ProviderArgs{
			Kubeconfig: pulumi.StringPtrOutput(kubeConfig),
		})
		if err != nil {
			return err
		}

		// Create new provider
		err = newMetallb(ctx, config, provider)

		if err != nil {
			return err
		}

		return nil
	})
}
