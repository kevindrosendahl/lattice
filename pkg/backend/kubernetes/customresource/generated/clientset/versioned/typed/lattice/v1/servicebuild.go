package v1

import (
	v1 "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	scheme "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/generated/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ServiceBuildsGetter has a method to return a ServiceBuildInterface.
// A group's client should implement this interface.
type ServiceBuildsGetter interface {
	ServiceBuilds(namespace string) ServiceBuildInterface
}

// ServiceBuildInterface has methods to work with ServiceBuild resources.
type ServiceBuildInterface interface {
	Create(*v1.ServiceBuild) (*v1.ServiceBuild, error)
	Update(*v1.ServiceBuild) (*v1.ServiceBuild, error)
	UpdateStatus(*v1.ServiceBuild) (*v1.ServiceBuild, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.ServiceBuild, error)
	List(opts meta_v1.ListOptions) (*v1.ServiceBuildList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceBuild, err error)
	ServiceBuildExpansion
}

// serviceBuilds implements ServiceBuildInterface
type serviceBuilds struct {
	client rest.Interface
	ns     string
}

// newServiceBuilds returns a ServiceBuilds
func newServiceBuilds(c *LatticeV1Client, namespace string) *serviceBuilds {
	return &serviceBuilds{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the serviceBuild, and returns the corresponding serviceBuild object, and an error if there is any.
func (c *serviceBuilds) Get(name string, options meta_v1.GetOptions) (result *v1.ServiceBuild, err error) {
	result = &v1.ServiceBuild{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("servicebuilds").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceBuilds that match those selectors.
func (c *serviceBuilds) List(opts meta_v1.ListOptions) (result *v1.ServiceBuildList, err error) {
	result = &v1.ServiceBuildList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("servicebuilds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceBuilds.
func (c *serviceBuilds) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("servicebuilds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a serviceBuild and creates it.  Returns the server's representation of the serviceBuild, and an error, if there is any.
func (c *serviceBuilds) Create(serviceBuild *v1.ServiceBuild) (result *v1.ServiceBuild, err error) {
	result = &v1.ServiceBuild{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("servicebuilds").
		Body(serviceBuild).
		Do().
		Into(result)
	return
}

// Update takes the representation of a serviceBuild and updates it. Returns the server's representation of the serviceBuild, and an error, if there is any.
func (c *serviceBuilds) Update(serviceBuild *v1.ServiceBuild) (result *v1.ServiceBuild, err error) {
	result = &v1.ServiceBuild{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("servicebuilds").
		Name(serviceBuild.Name).
		Body(serviceBuild).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *serviceBuilds) UpdateStatus(serviceBuild *v1.ServiceBuild) (result *v1.ServiceBuild, err error) {
	result = &v1.ServiceBuild{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("servicebuilds").
		Name(serviceBuild.Name).
		SubResource("status").
		Body(serviceBuild).
		Do().
		Into(result)
	return
}

// Delete takes name of the serviceBuild and deletes it. Returns an error if one occurs.
func (c *serviceBuilds) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("servicebuilds").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceBuilds) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("servicebuilds").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched serviceBuild.
func (c *serviceBuilds) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceBuild, err error) {
	result = &v1.ServiceBuild{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("servicebuilds").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
