package model

type Node struct {
	Address   string `json:"address" bson:"address"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Token     string `json:"token" bson:"token"`
	Trust     int    `json:"trust" bson:"trust"`
	Type      int    `json:"type" bson:"type"`
	Version   string `json:"version" bson:"version"`
	Name      string `json:"name" bson:"name"`
}

type NodeGetter interface {
	GetTypeName() string
}

func (node *Node) GetTypeName() string {
	switch node.Type {
	case 0:
		return "未定义节点"
	case 1:
		return "主节点"
	case 2:
		return "从节点"
	case 3:
		return "互通节点"
	case 4:
		return "黑名单"
	default:
		return "未知节点模式"
	}
}
