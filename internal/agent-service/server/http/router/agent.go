package handler

import (
	http_server "github.com/UnicomAI/wanwu/internal/agent-service/pkg/http-server"
	"github.com/UnicomAI/wanwu/internal/agent-service/server/http/handler"
	"net/http"
)

func init() {
	group := http_server.Group("/agent")
	group.Register("/chat", http.MethodPost, handler.AgentChat, "智能体流式问答")
}
