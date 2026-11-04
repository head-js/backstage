//go:build ignore

package main

import (
	restfulecho "github.com/labstack/echo/v5"
	restfulhttp "net/http"
)

// Yaegi loads this file before main.go; receive the client from the entry point.
type whiteboardRESTful struct {
	httpClient *vikunjaHTTPClient
}

func (r *whiteboardRESTful) registerRoutes(g *restfulecho.Group) {
	g.GET("/whiteboard/workspaces", r.handleWorkspaces)
	g.GET("/whiteboard/visit-nodes", r.handleVisitNodes)
	g.GET("/whiteboard/visit-edges", r.handleVisitEdges)
	g.GET("/whiteboard/visit-nodes.read/:visitNodeId", r.handleReadVisitNode)
	g.POST("/whiteboard/visit-nodes.create", r.handleCreateVisitNode)
	g.DELETE("/whiteboard/visit-nodes.delete/:id", r.handleDeleteVisitNode)
	g.POST("/whiteboard/visit-edge.create", r.handleCreateVisitEdge)
	g.POST("/whiteboard/visit-edges.delete", r.handleDeleteVisitEdge)
}

func (r *whiteboardRESTful) handleDeleteVisitNode(c *restfulecho.Context) error {
	request := c.Request()
	if err := deleteVisitNode(
		r.httpClient,
		request.Context(),
		c.Param("id"),
		request.Header.Get("Authorization"),
	); err != nil {
		return err
	}
	return c.NoContent(restfulhttp.StatusNoContent)
}

func (r *whiteboardRESTful) handleDeleteVisitEdge(c *restfulecho.Context) error {
	var input VisitEdge
	if err := c.Bind(&input); err != nil {
		return err
	}
	request := c.Request()
	if err := deleteVisitEdge(
		r.httpClient,
		request.Context(),
		input,
		request.Header.Get("Authorization"),
	); err != nil {
		return err
	}
	return c.NoContent(restfulhttp.StatusNoContent)
}

func (r *whiteboardRESTful) handleCreateVisitEdge(c *restfulecho.Context) error {
	var input VisitEdge
	if err := c.Bind(&input); err != nil {
		return err
	}
	request := c.Request()
	visitEdge, err := createVisitEdge(
		r.httpClient,
		request.Context(),
		input,
		request.Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}
	return c.JSON(restfulhttp.StatusCreated, visitEdge)
}

func (r *whiteboardRESTful) handleVisitEdges(c *restfulecho.Context) error {
	request := c.Request()
	visitEdges, err := listVisitEdges(
		r.httpClient,
		request.Context(),
		request.Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}
	return c.JSON(restfulhttp.StatusOK, visitEdges)
}

func (r *whiteboardRESTful) handleCreateVisitNode(c *restfulecho.Context) error {
	var input VisitNodeCreateRequest
	if err := c.Bind(&input); err != nil {
		return err
	}
	request := c.Request()
	visitNode, err := createVisitNode(
		r.httpClient,
		request.Context(),
		input,
		request.Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}
	return c.JSON(restfulhttp.StatusCreated, visitNode)
}

func (r *whiteboardRESTful) handleWorkspaces(c *restfulecho.Context) error {
	request := c.Request()
	authorization := request.Header.Get("Authorization")

	workspaces, err := listWorkspaces(r.httpClient, request.Context(), authorization)
	if err != nil {
		return err
	}
	return c.JSON(restfulhttp.StatusOK, workspaces)
}

func (r *whiteboardRESTful) handleReadVisitNode(c *restfulecho.Context) error {
	request := c.Request()
	visitNode, err := getVisitNode(
		r.httpClient,
		request.Context(),
		c.Param("visitNodeId"),
		request.Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	return c.JSON(restfulhttp.StatusOK, visitNode)
}

func (r *whiteboardRESTful) handleVisitNodes(c *restfulecho.Context) error {
	request := c.Request()
	visitNodes, err := listVisitNodes(
		r.httpClient,
		request.Context(),
		request.Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	return c.JSON(restfulhttp.StatusOK, visitNodes)
}
