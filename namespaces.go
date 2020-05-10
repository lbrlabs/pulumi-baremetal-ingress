package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/providers"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
)

// createNamespace creates a new namespace
func createNamespace(ctx *pulumi.Context, name string, provider *providers.Provider) error {

	// Create a new namespace
	_, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
		},
	}, pulumi.Provider(provider))

	if err != nil {
		return err
	}

	return nil

}
