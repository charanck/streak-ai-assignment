package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Request struct {
	Edges [][]int `json:"edges"`
	Start int     `json:"start"`
	End   int     `json:"end"`
}

type Response struct {
	Paths [][]int `json:"paths"`
}

func main() {
	e := echo.New()
	e.POST("/find-paths", findPathsHandler)
	e.Logger.Fatal(e.Start(":3000"))
}

func findPathsHandler(c echo.Context) error {
	request := &Request{}
	err := json.NewDecoder(c.Request().Body).Decode(request)
	if err != nil {
		fmt.Println("error while decoding the request: ", err)
		return c.String(http.StatusBadRequest, "invalid request body")
	}
	graph := buildGraph(request.Edges)
	if _, ok := graph[request.Start]; !ok {
		fmt.Println("Start node doesn't exist in the graph")
		return c.String(http.StatusBadRequest, "Start node doesn't exist in the graph")
	}
	if _, ok := graph[request.End]; !ok {
		fmt.Println("End node doesn't exist in the graph")
		return c.String(http.StatusBadRequest, "End node doesn't exist in the graph")
	}
	response := Response{}
	paths := [][]int{}
	search(graph, request.Start, request.End, []int{request.Start}, &paths)
	response.Paths = paths
	return c.JSON(http.StatusAccepted, response)
}

func buildGraph(edges [][]int) map[int][]int {
	graph := map[int][]int{}
	for _, edge := range edges {
		currNode, nextNode := edge[0], edge[1]
		if adjList, ok := graph[currNode]; ok {
			graph[currNode] = append(adjList, nextNode)
		} else {
			graph[currNode] = []int{nextNode}
		}
		if _, ok := graph[nextNode]; !ok {
			graph[nextNode] = []int{}
		}
	}
	return graph
}

func search(graph map[int][]int, currNode, targetNode int, currPath []int, paths *[][]int) {
	if currNode == targetNode {
		*paths = append(*paths, currPath)
		return
	}
	if connectedNodes := graph[currNode]; len(connectedNodes) == 0 {
		return
	}
	for _, connectedNode := range graph[currNode] {
		search(graph, connectedNode, targetNode, append(currPath, connectedNode), paths)
	}
}
