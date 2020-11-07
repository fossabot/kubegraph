package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddStyledNode creates a new styled node with the given resource
func (kgraph KubeGraph) AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (*dot.Node, error) {

	node, err := kgraph.createStyledNode(nodeName, resourceName, icon)
	if err != nil {
		return nil, err
	}

	kgraph.addNode(resourceType, resourceName, node)
	if err := kgraph.addObject(resourceType, resourceName, resourceObject); err != nil {
		return nil, err
	}

	return node, nil
}

// LinkNode links the node to the target node type/name, if it exists
func (kgraph KubeGraph) LinkNode(node *dot.Node, targetNodeType reflect.Type, targetNodeName string) (edge *dot.Edge, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			edge = nil
			err = fmt.Errorf("%++v", recoverErr)
		}
	}()

	targetNode, ok := kgraph.nodes[targetNodeType][targetNodeName]
	// TODO get or create unknown node and link here
	if !ok {
		// log.Printf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
		return nil, fmt.Errorf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
	}

	edgeName := fmt.Sprintf("%s-%s", node.ID(), targetNode.ID())
	edge = kgraph.graph.Edge(node, targetNode, edgeName)
	edge.Label("")
	return edge, nil
}

// GetObjects gets all objects in store
func (kgraph KubeGraph) GetObjects(objectType reflect.Type) (map[string]runtime.Object, error) {
	typeObjects, typeExists := kgraph.objects[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no objects for type %s found", objectType.String())
	}

	return typeObjects, nil
}

// GetNode gets a node by type/name
func (kgraph KubeGraph) GetNode(nodeType reflect.Type, nodeName string) (*dot.Node, error) {
	typeNodes, typeExists := kgraph.nodes[nodeType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", nodeType.String())
	}

	node, nodeExists := typeNodes[nodeName]
	if !nodeExists {
		return nil, fmt.Errorf("node %s/%s not found", nodeType.String(), nodeName)
	}

	return node, nil
}
