package http

import (
	"github.com/gin-gonic/gin"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"net/http"
)

type NodeHandler struct {
	nodePool pool.WorkerNode
}

func ProvideNodeHandler(nodePool pool.WorkerNode) NodeHandler {
	return NodeHandler{
		nodePool: nodePool,
	}
}

func (h *NodeHandler) GetNodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"poolSize": h.nodePool.PoolSize(),
		"nodeList": h.nodePool.GetAllWorkerNode(),
	})
}
