package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/providers"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/imdario/mergo"
)

type LabelledNamespace struct {
	pulumi.ResourceState
}

func NewLabelledNamespace(ctx *pulumi.Context, name string, l pulumi.StringMap, opts ...pulumi.ResourceOption) (*LabelledNamespace, error) {
	ns := &LabelledNamespace{}
	// Create a new namespace

	labels := pulumi.StringMap{
		"app": pulumi.String(name),
	}

	if err := mergo.Merge(&labels, l, mergo.WithOverride); err != nil {
		if err != nil {
			return nil, err
		}
	}

	_, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
			Labels: labels,
		},
	})
	if err != nil {
		return nil, err
	}
	err = ctx.RegisterComponentResource("lbrlabs:kubernetes:LabelledNamespace", name, ns, opts...)
	if err != nil {
		return nil, err
	}
	return ns, nil
}
