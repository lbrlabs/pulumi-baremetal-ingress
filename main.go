package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/providers"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
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

		// Create metallb namespace
		metallbNamespace := config.Require("metallbNamespace")
		metallbChartVersion := config.Get("metallbChartVersion")
		metallbAddress := config.Require("metallbAddress")
		err = createNamespace(ctx, metallbNamespace, provider)
		if err != nil {
			return fmt.Errorf("Error creating namespace %w", err)
		}

		rawMetallbConfig := &MetallbConfig{
			AddressPools: []AddressPool{
				{
					Name: "default",
					Protocol: "layer2",
					Addresses: []string{metallbAddress},
				},
			},
		}

		metallbConfig, err := json.Marshal(rawMetallbConfig)

		if err != nil {
			panic("Error formatting metallbconfig")
		}



		// Set up the configmap
		cm, err := corev1.NewConfigMap(ctx, "metallb-config", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(metallbNamespace),
			},
			Data: pulumi.StringMap{
				"config": pulumi.String(string(metallbConfig)),
			},
		})

		// Install metallb
		_, err = helm.NewChart(ctx, "metallb", helm.ChartArgs{
			Chart:   pulumi.String("metallb"),
			Version: pulumi.String(metallbChartVersion),
			FetchArgs: &helm.FetchArgs{
				Repo: pulumi.String("https://kubernetes-charts.storage.googleapis.com/"),
			},
			Values: pulumi.Map{
				"existingConfigMap": cm.Metadata.Name(),
			},
		}, pulumi.Provider(provider))

		if err != nil {
			fmt.Errorf("Error creating chart: %w", err)
		}

		return nil
	})
}
