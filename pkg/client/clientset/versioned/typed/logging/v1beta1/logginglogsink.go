// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/logging/v1beta1"
	scheme "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// LoggingLogSinksGetter has a method to return a LoggingLogSinkInterface.
// A group's client should implement this interface.
type LoggingLogSinksGetter interface {
	LoggingLogSinks(namespace string) LoggingLogSinkInterface
}

// LoggingLogSinkInterface has methods to work with LoggingLogSink resources.
type LoggingLogSinkInterface interface {
	Create(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.CreateOptions) (*v1beta1.LoggingLogSink, error)
	Update(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.UpdateOptions) (*v1beta1.LoggingLogSink, error)
	UpdateStatus(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.UpdateOptions) (*v1beta1.LoggingLogSink, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.LoggingLogSink, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.LoggingLogSinkList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.LoggingLogSink, err error)
	LoggingLogSinkExpansion
}

// loggingLogSinks implements LoggingLogSinkInterface
type loggingLogSinks struct {
	client rest.Interface
	ns     string
}

// newLoggingLogSinks returns a LoggingLogSinks
func newLoggingLogSinks(c *LoggingV1beta1Client, namespace string) *loggingLogSinks {
	return &loggingLogSinks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the loggingLogSink, and returns the corresponding loggingLogSink object, and an error if there is any.
func (c *loggingLogSinks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.LoggingLogSink, err error) {
	result = &v1beta1.LoggingLogSink{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("logginglogsinks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of LoggingLogSinks that match those selectors.
func (c *loggingLogSinks) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.LoggingLogSinkList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.LoggingLogSinkList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("logginglogsinks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested loggingLogSinks.
func (c *loggingLogSinks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("logginglogsinks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a loggingLogSink and creates it.  Returns the server's representation of the loggingLogSink, and an error, if there is any.
func (c *loggingLogSinks) Create(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.CreateOptions) (result *v1beta1.LoggingLogSink, err error) {
	result = &v1beta1.LoggingLogSink{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("logginglogsinks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loggingLogSink).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a loggingLogSink and updates it. Returns the server's representation of the loggingLogSink, and an error, if there is any.
func (c *loggingLogSinks) Update(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.UpdateOptions) (result *v1beta1.LoggingLogSink, err error) {
	result = &v1beta1.LoggingLogSink{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("logginglogsinks").
		Name(loggingLogSink.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loggingLogSink).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *loggingLogSinks) UpdateStatus(ctx context.Context, loggingLogSink *v1beta1.LoggingLogSink, opts v1.UpdateOptions) (result *v1beta1.LoggingLogSink, err error) {
	result = &v1beta1.LoggingLogSink{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("logginglogsinks").
		Name(loggingLogSink.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loggingLogSink).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the loggingLogSink and deletes it. Returns an error if one occurs.
func (c *loggingLogSinks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("logginglogsinks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *loggingLogSinks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("logginglogsinks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched loggingLogSink.
func (c *loggingLogSinks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.LoggingLogSink, err error) {
	result = &v1beta1.LoggingLogSink{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("logginglogsinks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}