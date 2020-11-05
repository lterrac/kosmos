// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/lterrac/system-autoscaler/pkg/apis/systemautoscaler/v1beta1"
	scheme "github.com/lterrac/system-autoscaler/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodScalesGetter has a method to return a PodScaleInterface.
// A group's client should implement this interface.
type PodScalesGetter interface {
	PodScales(namespace string) PodScaleInterface
}

// PodScaleInterface has methods to work with PodScale resources.
type PodScaleInterface interface {
	Create(ctx context.Context, podScale *v1beta1.PodScale, opts v1.CreateOptions) (*v1beta1.PodScale, error)
	Update(ctx context.Context, podScale *v1beta1.PodScale, opts v1.UpdateOptions) (*v1beta1.PodScale, error)
	UpdateStatus(ctx context.Context, podScale *v1beta1.PodScale, opts v1.UpdateOptions) (*v1beta1.PodScale, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.PodScale, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.PodScaleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.PodScale, err error)
	PodScaleExpansion
}

// podScales implements PodScaleInterface
type podScales struct {
	client rest.Interface
	ns     string
}

// newPodScales returns a PodScales
func newPodScales(c *SystemautoscalerV1beta1Client, namespace string) *podScales {
	return &podScales{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podScale, and returns the corresponding podScale object, and an error if there is any.
func (c *podScales) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.PodScale, err error) {
	result = &v1beta1.PodScale{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podscales").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodScales that match those selectors.
func (c *podScales) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.PodScaleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.PodScaleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podscales").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podScales.
func (c *podScales) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podscales").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a podScale and creates it.  Returns the server's representation of the podScale, and an error, if there is any.
func (c *podScales) Create(ctx context.Context, podScale *v1beta1.PodScale, opts v1.CreateOptions) (result *v1beta1.PodScale, err error) {
	result = &v1beta1.PodScale{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podscales").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podScale).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a podScale and updates it. Returns the server's representation of the podScale, and an error, if there is any.
func (c *podScales) Update(ctx context.Context, podScale *v1beta1.PodScale, opts v1.UpdateOptions) (result *v1beta1.PodScale, err error) {
	result = &v1beta1.PodScale{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podscales").
		Name(podScale.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podScale).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *podScales) UpdateStatus(ctx context.Context, podScale *v1beta1.PodScale, opts v1.UpdateOptions) (result *v1beta1.PodScale, err error) {
	result = &v1beta1.PodScale{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podscales").
		Name(podScale.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podScale).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the podScale and deletes it. Returns an error if one occurs.
func (c *podScales) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podscales").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podScales) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podscales").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched podScale.
func (c *podScales) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.PodScale, err error) {
	result = &v1beta1.PodScale{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podscales").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
