package actions_test

import (
	"context"
	"testing"

	"github.com/onsi/gomega/gstruct"
	"github.com/rs/xid"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dscv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/datasciencecluster/v1"
	dsciv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/dscinitialization/v1"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster/gvk"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/actions"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/metadata/labels"

	. "github.com/onsi/gomega"
)

func TestDeleteResourcesAction(t *testing.T) {
	g := NewWithT(t)

	ctx := context.Background()
	ns := xid.New().String()

	client := NewFakeClient(
		&appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gvk.Deployment.GroupVersion().String(),
				Kind:       gvk.Deployment.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-deployment",
				Namespace: ns,
				Labels: map[string]string{
					labels.K8SCommon.PartOf: "foo",
				},
			},
		},
		&appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gvk.Deployment.GroupVersion().String(),
				Kind:       gvk.Deployment.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-deployment-2",
				Namespace: ns,
				Labels: map[string]string{
					labels.K8SCommon.PartOf: "baz",
				},
			},
		},
	)

	action := actions.NewDeleteResourcesAction(
		ctx,
		actions.WithDeleteResourcesTypes(&appsv1.Deployment{}),
		actions.WithDeleteResourcesLabel(labels.K8SCommon.PartOf, "foo"))

	err := action.Execute(ctx, &types.ReconciliationRequest{
		Client:   client,
		Instance: nil,
		DSCI:     &dsciv1.DSCInitialization{Spec: dsciv1.DSCInitializationSpec{ApplicationsNamespace: ns}},
		DSC:      &dscv1.DataScienceCluster{},
		Platform: cluster.OpenDataHub,
	})

	g.Expect(err).ShouldNot(HaveOccurred())

	deployments := appsv1.DeploymentList{}
	err = client.List(ctx, &deployments)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(deployments.Items).Should(HaveLen(1))
	g.Expect(deployments.Items[0]).To(
		gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"ObjectMeta": gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Name": Equal("my-deployment-2"),
			}),
		}),
	)
}
