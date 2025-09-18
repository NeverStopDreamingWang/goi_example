package role

func get_children_menu(menuList []*menuListModel) []*menuListModel {
	// 创建一个ID映射，以便查找节点
	nodes := make(map[int64]*menuListModel)
	for _, item := range menuList {
		nodes[*item.Id] = item
	}
	// 用来存放根节点
	var tree []*menuListModel

	// 遍历数据，构建树形结构
	for _, item := range menuList {
		// 判断父节点ID，如果为 nil 说明是根节点
		if item.Parent_id == nil {
			tree = append(tree, item)
		} else {
			// 查找父节点
			parent, exists := nodes[*item.Parent_id]
			if exists {
				// 将当前节点添加到父节点的Children字段
				parent.Children = append(parent.Children, item)
			}
		}
	}
	return tree
}
