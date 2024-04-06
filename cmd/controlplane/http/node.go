package http

import (
	"github.com/gin-gonic/gin"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"net/http"
)

type NodeHandler struct {
	nodePool    pool.WorkerNode
	taskWatcher service.CPTaskWatcher
}

func ProvideNodeHandler(nodePool pool.WorkerNode, taskWatcher service.CPTaskWatcher) NodeHandler {
	return NodeHandler{
		nodePool:    nodePool,
		taskWatcher: taskWatcher,
	}
}

func (h *NodeHandler) GetNodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"poolSize": h.nodePool.PoolSize(),
		"nodeList": h.nodePool.GetAllWorkerNode(),
	})
}

func (h *NodeHandler) GetTaskWatcher(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"taskUnderWatch": h.taskWatcher.GetTaskUnderWatch(),
	})
}
