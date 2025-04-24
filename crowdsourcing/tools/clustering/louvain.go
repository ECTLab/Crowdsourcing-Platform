package clustering

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/simple"
)

// maxWeight: nodes connect only if their weight is less than maxWeoght
// resolution: less than 1 favours larger communities. if more than 1, the algorithm favours smaller communities
func FindClusters(nodeRowIndex, nodeColIndex []int64, weights [][]float64, maxWeight, resolution float64) ([][]int64, []int64) {
	// g := simple.NewWeightedDirectedGraph(0, math.Inf(1))
	g := simple.NewUndirectedGraph()

	// Creating a weighted directed graph representing the matrix, connecting to nodes only when W is less than Max weight
	for _, nodeI := range nodeRowIndex {
		for _, nodeJ := range nodeColIndex {
			if nodeI == nodeJ {
				continue
			}

			// w := weights[nodeI][nodeJ]
			w := min(weights[nodeI][nodeJ], weights[nodeJ][nodeI])
			if w > maxWeight {
				continue
			}

			// edge := simple.WeightedEdge{
			// 	F: simple.Node(nodeI),
			// 	T: simple.Node(nodeJ),
			// 	W: w,
			// }
			// g.SetWeightedEdge(edge)

			edge := simple.Edge{
				F: simple.Node(nodeI),
				T: simple.Node(nodeJ),
			}
			g.SetEdge(edge)
		}
	}

	// running louvain community detection
	rg := community.Modularize(g, resolution, nil)

	// putting communities and their center node to returnable format
	var comms = [][]int64{}
	var centers = []int64{}
	clusteredNodes := map[int64]struct{}{}
	gcomms := rg.Communities()
	for _, gcomm := range gcomms {
		// extracting communities
		comm := []int64{}
		for _, node := range gcomm {
			nid := node.ID()
			comm = append(comm, nid)
			clusteredNodes[nid] = struct{}{}
		}
		comms = append(comms, comm)

		// extracting a center node for each community
		var maxDegree int
		var centerNode graph.Node
		for _, node := range gcomm {
			// degree := g.From(node.ID()).Len() + g.To(node.ID()).Len() // Directed
			degree := g.From(node.ID()).Len() // Undirected
			if degree > maxDegree {
				maxDegree = degree
				centerNode = node
			}
		}
		if centerNode == nil {
			centerNode = gcomm[0]
		}
		centers = append(centers, centerNode.ID())
	}

	// append nodes that was not clustered
	newComms := [][]int64{}
	newCenters := []int64{}
	for _, nid := range nodeRowIndex {
		if _, wasClustered := clusteredNodes[nid]; !wasClustered {
			newComms = append(newComms, []int64{nid})
			newCenters = append(newCenters, nid)
		}
	}
	comms = append(comms, newComms...)
	centers = append(centers, newCenters...)

	return comms, centers
}
