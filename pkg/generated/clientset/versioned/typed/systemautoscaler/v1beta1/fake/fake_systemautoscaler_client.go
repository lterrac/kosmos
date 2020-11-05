// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/lterrac/system-autoscaler/pkg/generated/clientset/versioned/typed/systemautoscaler/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSystemautoscalerV1beta1 struct {
	*testing.Fake
}

func (c *FakeSystemautoscalerV1beta1) PodScales(namespace string) v1beta1.PodScaleInterface {
	return &FakePodScales{c, namespace}
}

func (c *FakeSystemautoscalerV1beta1) ServiceLevelAgreements(namespace string) v1beta1.ServiceLevelAgreementInterface {
	return &FakeServiceLevelAgreements{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSystemautoscalerV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
