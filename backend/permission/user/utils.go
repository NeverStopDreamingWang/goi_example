package user

import (
	"context"
	"database/sql"

	"github.com/NeverStopDreamingWang/goi/v3/db"
	"github.com/NeverStopDreamingWang/goi/v3/db/sqlite3"
)

import (
	"errors"
)

func GetUser(ctx context.Context, pk any) (*UserModel, error) {
	if pk == nil {
		return nil, nil
	}
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.SetModel(UserModel{})
	instance := &UserModel{}
	err := sqlite3DB.Where("`id` = ?", pk).First(ctx, instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return instance, nil
}
func GetUserInfo(ctx context.Context, pk any) (*UserInfo, error) {
	if pk == nil {
		return nil, nil
	}
	var instance *UserModel
	instance, err := GetUser(ctx, pk)
	if err != nil {
		return nil, err
	}
	return instance.ToUserInfo(), err
}

func GetUserMap(ctx context.Context, pk any) (map[string]any, error) {
	if pk == nil {
		return nil, nil
	}
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.SetModel(UserModel{})

	var instance = map[string]any{}
	err := sqlite3DB.Where("`id` = ?", pk).First(ctx, instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return instance, nil
}

func get_children_menu(menuList []*userMenuList) []*userMenuList {
	// 创建一个ID映射，以便查找节点
	nodes := make(map[int64]*userMenuList)
	for _, item := range menuList {
		nodes[*item.Id] = item
	}
	// 用来存放根节点
	var tree []*userMenuList

	// 遍历数据，构建树形结构
	for _, item := range menuList {
		// 判断父节点ID，如果为 nil 说明是根节点
		if item.ParentId == nil {
			tree = append(tree, item)
		} else {
			// 查找父节点
			parent, exists := nodes[*item.ParentId]
			if exists {
				// 将当前节点添加到父节点的Children字段
				parent.Children = append(parent.Children, item)
			}
		}
	}
	return tree
}
