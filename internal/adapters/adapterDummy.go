package adapters

/*
 * remove the dummy struct and replace the references with a proper kubernetes API resource
 */

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type dummy struct {
	metaV1.TypeMeta
	metaV1.ObjectMeta
}

func (d dummy) GetObjectKind() schema.ObjectKind {
	return nil
}

func (d dummy) DeepCopyObject() runtime.Object {
	return dummy{}
}

type adapterCoreV1Dummy struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Dummy{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1Dummy) GetType() reflect.Type {
	return reflect.TypeOf(&dummy{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1Dummy) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*dummy)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1Dummy) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1Dummy) Configure(statefulGraph StatefulGraph) error {
	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}
	for resourceName, resourceObject := range objects {
		resource := resourceObject.(*dummy)
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		// do something with each resource
		log.Printf("dummy resource %s, node %s", resource.Name, resourceNode.Name())
	}
	return nil
}
