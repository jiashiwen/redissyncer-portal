package model

type RemoveNodeModel struct {
	NodeType          string
	NodeID            string
	TasksOnNodePolice string //位于该节点任务的政策，销毁或迁移 'destroy'、'move'
}
