package main

import (
	"flag"
	"fmt"
	"strings"
)

type NodeType int64

const (
	SOURCE NodeType = iota
	SINK
	PEOPLE
)

type Node struct {
	ID        int64
	NType     NodeType
	EdgesTo   map[int64]*Edge // edges going out from this node
	EdgesFrom map[int64]*Edge // edges to this node, should not really be useful but nice to have I guess
	Name      string
}

type Edge struct {
	Capacity     int64
	UsedCapacity int64
	From         int64
	To           int64
	Note         string
}

type Graph struct {
	SourceNode *Node
	SinkNode   *Node
	People     map[int64]*Node // key = nodeId and value is the node
}

func InitGraph() *Graph {
	var names string
	var num int64
	flag.StringVar(&names, "p", "", "a comma separated list of people")
	flag.Int64Var(&num, "n", 0, "number of people")

	flag.Parse()

	if names == "" && num == 0 {
		panic("provide some input")
	}

	parseNum := true
	if names != "" && num != 0 {
		parseNum = false
	} else if names != "" && num == 0 {
		parseNum = false
	}

	var peopleMap map[int64]*Node

	if parseNum {
		peopleMap = initFromNum(num)
	} else {
		peopleMap = initFromNames(names)
	}

	sourceNode := &Node{
		ID:        -1,
		NType:     SOURCE,
		EdgesFrom: nil,
		EdgesTo:   nil,
		Name:      "Source",
	}

	sinkNode := &Node{
		ID:        -2,
		NType:     SINK,
		EdgesFrom: nil,
		EdgesTo:   nil,
		Name:      "Sink",
	}

	graph := &Graph{
		SourceNode: sourceNode,
		SinkNode:   sinkNode,
		People:     peopleMap,
	}

	initializeEdges(graph)

	return graph
}

func initFromNames(csNames string) map[int64]*Node {
	names := strings.Split(csNames, ",")

	if len(names) == 0 || len(names) == 1 {
		panic("names is either empty or you only have one name, either ways the input is invalid")
	}

	peopleNodeMap := make(map[int64]*Node)

	// create all people nodes
	for idx, name := range names {
		peopleNodeMap[int64(idx)] = &Node{
			ID:        int64(idx) + 1,
			NType:     PEOPLE,
			EdgesFrom: nil,
			EdgesTo:   nil,
			Name:      name,
		}
	}

	return peopleNodeMap
}

func initFromNum(num int64) map[int64]*Node {
	if num < 2 {
		panic("less than 2 people, go figure out the solution yourself")
	}

	peopleNodeMap := make(map[int64]*Node)

	// create all people nodes
	for id := int64(1); id < num+1; id++ {
		peopleNodeMap[id] = &Node{
			ID:        id,
			NType:     PEOPLE,
			EdgesFrom: nil,
			EdgesTo:   nil,
			Name:      fmt.Sprintf("Person %d", id),
		}
	}

	return peopleNodeMap
}

/*
*
initializeEdges.
Each person gets an edge from source with a capacity of 12, representing one apple per month.
There will be an edge from each person to the other person with capacity 4 (1 for every 3 months)
Each person will also have an edge from themselves to the sink with a capacity of 12 since they can give 12 apples per year

Calling this function will reset every existing edge
*/
func initializeEdges(g *Graph) {
	sourceNode := g.SourceNode
	sinkNode := g.SinkNode

	sourceNode.EdgesFrom = make(map[int64]*Edge)
	sourceNode.EdgesTo = make(map[int64]*Edge)

	sinkNode.EdgesFrom = make(map[int64]*Edge)
	sinkNode.EdgesTo = make(map[int64]*Edge)

	peopleMap := g.People

	// source and sink edges
	for _, people := range peopleMap {
		people.EdgesFrom = make(map[int64]*Edge)
		people.EdgesTo = make(map[int64]*Edge)

		sourceEdge := &Edge{
			Capacity:     12,
			UsedCapacity: 0,
			From:         sourceNode.ID,
			To:           people.ID,
			Note:         fmt.Sprintf("Edge from source to person %s", people.Name),
		}

		sinkEdge := &Edge{
			Capacity:     12,
			UsedCapacity: 0,
			From:         people.ID,
			To:           sinkNode.ID,
			Note:         fmt.Sprintf("Edge from person %s to sink node", people.Name),
		}

		sourceNode.EdgesTo[people.ID] = sourceEdge
		sinkNode.EdgesFrom[people.ID] = sinkEdge

		people.EdgesFrom[sourceNode.ID] = sourceEdge
		people.EdgesTo[sinkNode.ID] = sinkEdge
	}

	// edges between people
	for _, person1 := range peopleMap {
		for _, person2 := range peopleMap {
			if person1.ID == person2.ID {
				continue
			}

			// edge present from this person to the other person
			if _, present := person1.EdgesTo[person2.ID]; present {
				continue
			}

			// sanity check
			if _, present := person2.EdgesFrom[person1.ID]; present {
				panic("invalid edges this if should not be possible")
			}

			edge := &Edge{
				Capacity:     4,
				UsedCapacity: 0,
				From:         person1.ID,
				To:           person2.ID,
				Note:         fmt.Sprintf("Edge from person %s to person %s", person1.Name, person2.Name),
			}

			person1.EdgesTo[person2.ID] = edge
			person2.EdgesFrom[person1.ID] = edge
		}
	}
}
