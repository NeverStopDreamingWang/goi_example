package role

import (
	"time"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
)

func init() {
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.Migrate("goi_example", MenuModel{})     // 菜单表
	sqlite3DB.Migrate("goi_example", RoleMenuModel{}) // 角色-菜单表
	sqlite3DB.Migrate("goi_example", RoleModel{})     // 角色表
}

// 菜单表
type MenuModel struct {
	Id        *int64  `field_name:"id" field_type:"INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT" json:"id"` // ID
	Parent_id *int64  `field_name:"parent_id" field_type:"INTEGER" json:"parent_id"`                      // 父级
	Name      *string `field_name:"name" field_type:"TEXT NOT NULL" json:"name"`                          // 名称
	Icon      *string `field_name:"icon" field_type:"TEXT" json:"icon"`                                   // 图标
	Path      *string `field_name:"path" field_type:"TEXT NOT NULL" json:"path"`                          // 路由
}

func (MenuModel) ModelSet() *sqlite3.Settings {
	modelSettings := sqlite3.Settings{
		MigrationsHandler: sqlite3.MigrationsHandler{ // 迁移时处理函数
			BeforeHandler: nil,      // 迁移之前处理函数
			AfterHandler:  initMenu, // 迁移之后处理函数
		},

		TABLE_NAME: "tb_menu", // 设置表名

		// 自定义配置
		Settings: goi.Params{},
	}
	return &modelSettings
}

// 初始化菜单
func initMenu() error {
	initMenuList := [][]interface{}{
		{int64(1), nil, "首页", "home", "/home"},
		{int64(2), int64(1), "仪表盘", "", "dashboard"},
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.SetModel(MenuModel{})
	total, err := sqlite3DB.Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	for _, item := range initMenuList {
		var (
			Id   = item[0].(int64)
			Name = item[2].(string)
			Icon = item[3].(string)
			Path = item[4].(string)
		)

		sqlite3DB.SetModel(MenuModel{})
		total, err = sqlite3DB.Where("`id` = ?", Id).Count()
		if err != nil {
			return err
		}
		if total != 0 {
			continue
		}
		menu := MenuModel{
			Id:        &Id,
			Parent_id: nil,
			Name:      &Name,
			Icon:      &Icon,
			Path:      &Path,
		}
		if item[1] != nil {
			Parent_id := item[1].(int64)
			menu.Parent_id = &Parent_id
		}
		sqlite3DB.SetModel(MenuModel{})
		_, err = sqlite3DB.Insert(&menu)
		if err != nil {
			return err
		}
	}
	return nil
}

// 角色表
type RoleModel struct {
	Id          *int64   `field_name:"id" field_type:"INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT" json:"id"` // ID
	Name        *string  `field_name:"name" field_type:"TEXT NOT NULL" json:"name"`                          // 用户名
	Create_Time *string  `field_name:"create_time" field_type:"DATETIME NOT NULL" json:"create_time"`        // 创建时间
	Update_Time *string  `field_name:"update_time" field_type:"DATETIME" json:"update_time"`                 // 更新时间
	Menu_List   []*int64 `json:"-"`
}

func (RoleModel) ModelSet() *sqlite3.Settings {
	modelSettings := sqlite3.Settings{
		MigrationsHandler: sqlite3.MigrationsHandler{ // 迁移时处理函数
			BeforeHandler: nil,      // 迁移之前处理函数
			AfterHandler:  initRole, // 迁移之后处理函数
		},

		TABLE_NAME: "tb_role", // 设置表名

		// 自定义配置
		Settings: goi.Params{},
	}
	return &modelSettings
}

// 初始化角色
func initRole() error {
	initRoleList := [][]interface{}{
		{int64(1), "超级管理员"},
	}
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	var menu_List []*MenuModel
	sqlite3DB.SetModel(MenuModel{})
	err := sqlite3DB.Select(&menu_List)
	if err != nil {
		return err
	}

	for _, item := range initRoleList {
		var (
			Id          = item[0].(int64)
			Name        = item[1].(string)
			Create_Time = goi.GetTime().Format(time.DateTime)
		)

		role := RoleModel{
			Id:          &Id,
			Name:        &Name,
			Create_Time: &Create_Time,
			Update_Time: nil,
		}
		for _, menu := range menu_List {
			role.Menu_List = append(role.Menu_List, menu.Id)
		}
		err := role.Validate()
		if err != nil {
			return err
		}
		err = role.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

// 角色-菜单表
type RoleMenuModel struct {
	Role_Id *int64 `field_name:"role_id" field_type:"INTEGER NOT NULL" json:"role_id"` // 角色ID
	Menu_Id *int64 `field_name:"menu_id" field_type:"INTEGER NOT NULL" json:"menu_id"` // 菜单ID
}

func (RoleMenuModel) ModelSet() *sqlite3.Settings {
	modelSettings := sqlite3.Settings{
		MigrationsHandler: sqlite3.MigrationsHandler{ // 迁移时处理函数
			BeforeHandler: nil, // 迁移之前处理函数
			AfterHandler:  nil, // 迁移之后处理函数
		},

		TABLE_NAME: "tb_role_menu", // 设置表名

		// 自定义配置
		Settings: goi.Params{},
	}
	return &modelSettings
}
