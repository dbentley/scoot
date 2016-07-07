package cluster_membership

import (
	"fmt"
	"strings"
	"testing"
)

/*
 * Verify Creating an Empty Dynamic Cluster
 */
func TestCreateEmptyDynamicCluster(t *testing.T) {

	var emptyNodes []Node
	dc, _ := DynamicClusterFactory(emptyNodes)
	members := dc.Members()

	if len(members) != 0 {
		t.Error("Empty Dynamic Cluster should have 0 nodes")
	}
}

/*
 * Verify Creating a Dynamic Cluster
 */
func TestCreateDynamicCluster(t *testing.T) {
	testNodes := GenerateTestNodes(10)
	dc, _ := DynamicClusterFactory(testNodes)
	members := dc.Members()

	if len(members) != len(testNodes) {
		t.Error("number of nodes supplied at creation differs from the number of nodes in dynamic cluster")
	}

	testNodeMap := make(map[string]Node)
	for _, node := range testNodes {
		testNodeMap[node.Id()] = node
	}

	//Ensure all nodes passed into factory are present in the cluster
	for _, node := range members {

		_, exists := testNodeMap[node.Id()]
		if !exists {
			t.Error(fmt.Sprintf("node %s is not in created cluster", node.Id()))
		}
	}
}

/*
 * Verify Add Node to Cluster
 */
func TestAddNodesToDynamicCluster(t *testing.T) {

	var emptyNodes []Node
	dc, clusterState := DynamicClusterFactory(emptyNodes)
	members := clusterState.InitialMembers
	updateCh := clusterState.Updates

	if len(members) != 0 {
		t.Error("Empty Dynamic Cluster should have 0 nodes")
	}

	tNode := TestNode{
		id: "testNode1",
	}

	dc.AddNode(&tNode)
	members = dc.Members()
	if len(members) != 1 {
		t.Error("Dynamic Cluster should have 1 node")
	}

	if members[0].Id() != tNode.Id() {
		t.Error(fmt.Sprintf("Dynamic Cluster should have 1 node with id %s", tNode.Id()))
	}

	select {
	case nodeUpdate := <-updateCh:
		if nodeUpdate.UpdateType != NodeAdded {
			t.Error("Expected UpdateType to be Add")
		}
		if strings.Compare(nodeUpdate.Node.Id(), tNode.id) != 0 {
			t.Error("Unexpected NodeId Added to AddCh", nodeUpdate.Node.Id(), tNode.id)
		}
	default:
		t.Error("Expected Added Channel to contain Just Added Node")
	}

	tNode2 := TestNode{
		id: "testNode2",
	}
	dc.AddNode(&tNode2)
	members = dc.Members()

	if len(members) != 2 {
		t.Error("Dynamic Cluster should have 2 node")
	}

	select {
	case nodeUpdate := <-updateCh:
		if nodeUpdate.UpdateType != NodeAdded {
			t.Error("Expected UpdateType to be Add")
		}
		if strings.Compare(nodeUpdate.Node.Id(), tNode2.id) != 0 {
			t.Error("Unexpected NodeId Added to AddCh", nodeUpdate.Node.Id(), tNode2.id)
		}
	default:
		t.Error("Expected Added Channel to contain Just Added Node")
	}

	if members[0].Id() != tNode.Id() {
		t.Error(fmt.Sprintf("Dynamic Cluster should have node with id %s", tNode.Id()))
	}

	if members[1].Id() != tNode2.Id() {
		t.Error(fmt.Sprintf("Dynamic Cluster should have node with id %s", tNode2.Id()))
	}

}

/*
 * Verify Add Node to Cluster, Node Already in Cluster,
 * Add should be Idempotent
 */
func TestAddNodeToClusterThatAlreadyExists(t *testing.T) {
	var emptyNodes []Node
	dc, clusterState := DynamicClusterFactory(emptyNodes)
	updateCh := clusterState.Updates

	members := dc.Members()
	tNode := TestNode{
		id: "testNode1",
	}

	dc.AddNode(&tNode)
	members = dc.Members()
	if len(members) != 1 {
		t.Error("Dynamic Cluster should have 1 node")
	}

	dc.AddNode(&tNode)
	members = dc.Members()
	if len(members) != 1 {
		t.Error(fmt.Sprintf("Dynamic Cluster should have 1 node not %d nodes", len(members)))
	}

	select {
	case <-updateCh:
	default:
		t.Error("Expected added node to be added to channel")
	}

	select {
	case <-updateCh:
		t.Error("Expected already existing node not to be added to the channel")
	default:
	}
}

/*
 * Verify that Delete is Idempotent, a non existant node can
 * be deleted from an empty cluster
 */
func TestDeleteNodeFromEmptyCluster(t *testing.T) {
	var emptyNodes []Node
	dc, clusterState := DynamicClusterFactory(emptyNodes)
	updateCh := clusterState.Updates

	dc.RemoveNode("node_X")
	if len(dc.Members()) != 0 {
		t.Error("Dynamic Cluster should have 0 nodes after Delete of Non Existant Node")
	}

	select {
	case <-updateCh:
		t.Error("Expected no updates on channel when non existant node removed")
	default:
	}
}

/*
 * Verify that nodes can successfully be removed from the cluster
 */
func TestDeleteNodeFromCluster(t *testing.T) {
	tNode := TestNode{
		id: "testNode1",
	}

	nodes := make([]Node, 1)
	nodes[0] = &tNode
	dc, clusterState := DynamicClusterFactory(nodes)
	updateCh := clusterState.Updates

	if len(dc.Members()) != 1 {
		t.Error("Dynamic Cluster should have 1 node")
	}

	dc.RemoveNode(tNode.Id())

	if len(dc.Members()) != 0 {
		t.Error("Dynamic CLuster Should have 0 nodes")
	}

	select {
	case nodeUpdate := <-updateCh:
		if nodeUpdate.UpdateType != NodeRemoved {
			t.Error("Expected update type to be Removed")
		}
	default:
		t.Error("Expected removed node to be added to channel")
	}
}
