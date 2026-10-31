package menutree

import "qixue-api/internal/model"

type Node struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Type        string  `json:"type"`
	ParentID    *int    `json:"parentId"`
	Path        *string `json:"path"`
	Redirect    *string `json:"redirect"`
	Icon        *string `json:"icon"`
	Component   *string `json:"component"`
	Layout      *string `json:"layout"`
	KeepAlive   *int    `json:"keepAlive"`
	Method      *string `json:"method"`
	Description *string `json:"description"`
	Show        int     `json:"show"`
	Enable      int     `json:"enable"`
	Order       *int    `json:"order"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	Children    []Node  `json:"children,omitempty"`
}

// Tree 把 parentId 为空或 0 的菜单当作根。父级不在本次结果里的节点不返回。
func Tree(rows []model.AdminMenu) []Node {
	nodes := make(map[int]Node, len(rows))
	childrenOf := make(map[int][]int, len(rows))
	var roots []int
	for _, row := range rows {
		nodes[row.ID] = ToNode(row)
		parentID := 0
		if row.ParentID != nil {
			parentID = *row.ParentID
		}
		if parentID == 0 {
			roots = append(roots, row.ID)
			continue
		}
		childrenOf[parentID] = append(childrenOf[parentID], row.ID)
	}

	var build func(id int) Node
	build = func(id int) Node {
		node := nodes[id]
		for _, childID := range childrenOf[id] {
			if _, ok := nodes[childID]; !ok {
				continue
			}
			node.Children = append(node.Children, build(childID))
		}
		return node
	}

	out := make([]Node, 0, len(roots))
	for _, id := range roots {
		out = append(out, build(id))
	}
	return out
}

func ToNode(row model.AdminMenu) Node {
	return Node{
		ID:          row.ID,
		Name:        row.Name,
		Code:        row.Code,
		Type:        row.Type,
		ParentID:    row.ParentID,
		Path:        row.Path,
		Redirect:    row.Redirect,
		Icon:        row.Icon,
		Component:   row.Component,
		Layout:      row.Layout,
		KeepAlive:   row.KeepAlive,
		Method:      row.Method,
		Description: row.Description,
		Show:        row.Show,
		Enable:      row.Enable,
		Order:       row.Order,
		CreatedAt:   model.FormatTime(row.CreatedAt),
		UpdatedAt:   model.FormatTime(row.UpdatedAt),
	}
}
