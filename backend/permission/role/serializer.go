package role

import (
	"database/sql"
	"errors"

	"goi_example/backend/utils"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
)

func (self RoleModel) Validate() error {
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.SetModel(self)

	if self.Id != nil {
		sqlite3DB = sqlite3DB.Where("`id` != ?", self.Id)
	}

	if self.Name != nil {
		flag, err := sqlite3DB.Where("`name` = ?", self.Name).Exists()
		if err != nil {
			return errors.New("查询数据库错误")
		}
		if flag == true {
			return errors.New("角色名重复")
		}
	}
	return nil
}

func (self *RoleModel) Create() error {
	if self.CreateTime == nil {
		CreateTime := goi.GetTime()
		self.CreateTime = &CreateTime
	}

	err := sqlite3.Validate(self, true)
	if err != nil {
		return err
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	err = sqlite3DB.WithTransaction(func(engine *sqlite3.Engine, args ...any) error {
		engine.SetModel(self)
		result, err := engine.Insert(self)
		if err != nil {
			return errors.New("添加角色错误")
		}
		id, err := result.LastInsertId()
		if err != nil {
			return errors.New("添加角色错误")
		}
		self.Id = &id

		for _, menu_id := range self.MenuList {
			roleMenu := RoleMenuModel{
				RoleId: self.Id,
				MenuId: menu_id,
			}
			err = roleMenu.Create(engine)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (self *RoleModel) Update(validated_data *RoleModel) error {
	UpdateTime := goi.GetTime()
	validated_data.UpdateTime = &UpdateTime

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	err := sqlite3DB.WithTransaction(func(engine *sqlite3.Engine, args ...any) error {
		engine.SetModel(self)
		_, err := engine.Where("`id` = ?", self.Id).Update(validated_data)
		if err != nil {
			return errors.New("修改角色错误")
		}

		if validated_data.MenuList == nil {
			return nil
		}

		engine.SetModel(RoleMenuModel{})
		_, err = engine.Where("`role_id` = ?", self.Id).Delete()
		if err != nil {
			return errors.New("修改角色错误")
		}
		for _, menu_id := range validated_data.MenuList {
			roleMenu := RoleMenuModel{
				RoleId: self.Id,
				MenuId: menu_id,
			}
			err = roleMenu.Create(engine)
			if err != nil {
				return err
			}
		}
		return nil
	})
	utils.Update(self, validated_data)
	return err
}

func (self RoleModel) Delete() error {
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	err := sqlite3DB.WithTransaction(func(engine *sqlite3.Engine, args ...any) error {
		engine.SetModel(self)
		_, err := engine.Where("`id` = ?", self.Id).Delete()
		if err != nil {
			return err
		}

		engine.SetModel(RoleMenuModel{})
		_, err = engine.Where("`role_id` = ?", self.Id).Delete()
		if err != nil {
			return errors.New("删除角色错误")
		}
		return nil
	})
	return err
}

func (self *RoleMenuModel) Create(engine *sqlite3.Engine) error {
	err := sqlite3.Validate(self, true)
	if err != nil {
		return err
	}
	// 关联角色权限
	engine.SetModel(self)
	flag, err := engine.Where("`role_id`=? and `menu_id`=?", self.RoleId, self.MenuId).Exists()
	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		return errors.New("添加角色菜单错误")
	}
	if flag == true {
		return nil
	}
	engine.SetModel(self)
	_, err = engine.Insert(self)
	if err != nil {
		return errors.New("添加角色菜单错误")
	}
	return nil
}
