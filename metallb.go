package main

import (
	"encoding/json"
	"fmt"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/helm/v2"

	//"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/helm/v2"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/providers"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

func newMetallb(ctx *pulumi.Context, config *config.Config, provider *providers.Provider) error {

	// Create metallb namespace
	metallbNamespace := config.Require("metallbNamespace")
	metallbChartVersion := config.Get("metallbChartVersion")
	metallbAddress := config.Require("metallbAddress")

	_, err := NewLabelledNamespace(ctx, metallbNamespace, nil, pulumi.Provider(provider))
	if err != nil {
		return err
	}

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
		return fmt.Errorf("Error unmarshalling metallb config: %w", err)
	}

	// Set up the configmap
	cm, err := corev1.NewConfigMap(ctx, "metallb-config", &corev1.ConfigMapArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(metallbNamespace),
		},
		Data: pulumi.StringMap{
			"config": pulumi.String(string(metallbConfig)),
		},
	}, pulumi.Provider(provider))

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
		Namespace: pulumi.String(metallbNamespace),
	})

	if err != nil {
		return fmt.Errorf("Error creating chart: %w", err)
	}

	return nil
}
