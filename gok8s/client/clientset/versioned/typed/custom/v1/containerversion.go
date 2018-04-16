/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "github.com/nearmap/cvmanager/gok8s/apis/custom/v1"
	scheme "github.com/nearmap/cvmanager/gok8s/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ContainerVersionsGetter has a method to return a ContainerVersionInterface.
// A group's client should implement this interface.
type ContainerVersionsGetter interface {
	ContainerVersions(namespace string) ContainerVersionInterface
}

// ContainerVersionInterface has methods to work with ContainerVersion resources.
type ContainerVersionInterface interface {
	Create(*v1.ContainerVersion) (*v1.ContainerVersion, error)
	Update(*v1.ContainerVersion) (*v1.ContainerVersion, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.ContainerVersion, error)
	List(opts meta_v1.ListOptions) (*v1.ContainerVersionList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ContainerVersion, err error)
	ContainerVersionExpansion
}

// containerVersions implements ContainerVersionInterface
type containerVersions struct {
	client rest.Interface
	ns     string
}

// newContainerVersions returns a ContainerVersions
func newContainerVersions(c *CustomV1Client, namespace string) *containerVersions {
	return &containerVersions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the containerVersion, and returns the corresponding containerVersion object, and an error if there is any.
func (c *containerVersions) Get(name string, options meta_v1.GetOptions) (result *v1.ContainerVersion, err error) {
	result = &v1.ContainerVersion{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("containerversions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ContainerVersions that match those selectors.
func (c *containerVersions) List(opts meta_v1.ListOptions) (result *v1.ContainerVersionList, err error) {
	result = &v1.ContainerVersionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("containerversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested containerVersions.
func (c *containerVersions) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("containerversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a containerVersion and creates it.  Returns the server's representation of the containerVersion, and an error, if there is any.
func (c *containerVersions) Create(containerVersion *v1.ContainerVersion) (result *v1.ContainerVersion, err error) {
	result = &v1.ContainerVersion{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("containerversions").
		Body(containerVersion).
		Do().
		Into(result)
	return
}

// Update takes the representation of a containerVersion and updates it. Returns the server's representation of the containerVersion, and an error, if there is any.
func (c *containerVersions) Update(containerVersion *v1.ContainerVersion) (result *v1.ContainerVersion, err error) {
	result = &v1.ContainerVersion{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("containerversions").
		Name(containerVersion.Name).
		Body(containerVersion).
		Do().
		Into(result)
	return
}

// Delete takes name of the containerVersion and deletes it. Returns an error if one occurs.
func (c *containerVersions) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("containerversions").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *containerVersions) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("containerversions").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched containerVersion.
func (c *containerVersions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ContainerVersion, err error) {
	result = &v1.ContainerVersion{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("containerversions").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}