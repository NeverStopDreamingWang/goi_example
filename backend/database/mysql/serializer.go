package mysql

import (
	"errors"

	"goi_example/backend/utils"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/auth"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/mysql"
)

func (self UserModel) Validate() error {
	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)

	if self.Id != nil {
		mysqlDB = mysqlDB.Where("`id` != ?", self.Id)
	}

	if self.Username != nil {
		flag, err := mysqlDB.Where("`username` = ?", self.Username).Exists()
		if err != nil {
			return errors.New("查询数据库错误")
		}
		if flag == true {
			return errors.New("用户名重复")
		}
	}
	if self.Email != nil {
		flag, err := mysqlDB.Where("`email` = ?", self.Email).Exists()
		if err != nil {
			return errors.New("查询数据库错误")
		}
		if flag == true {
			return errors.New("邮箱已使用")
		}
	}
	if self.Status != nil {
		if *self.Status != DISABLE && *self.Status != NORMAL {
			return errors.New("用户状态错误")
		}
	}
	return nil
}

func (self *UserModel) Create() error {
	if self.CreateTime == nil {
		CreateTime := goi.GetTime()
		self.CreateTime = &CreateTime
	}

	// 密码加密
	encryptPassword, err := auth.MakePassword(*self.Password)
	if err != nil {
		return errors.New("密码格式错误")
	}
	self.Password = &encryptPassword

	err = mysql.Validate(self, true)
	if err != nil {
		return err
	}

	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)
	result, err := mysqlDB.Insert(self)
	if err != nil {
		return errors.New("添加用户错误")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.New("添加用户错误")
	}
	self.Id = &id
	return nil
}

func (self *UserModel) Update(validated_data *UserModel) error {
	if validated_data.Password != nil {
		// 密码加密
		encryptPassword, err := auth.MakePassword(*validated_data.Password)
		if err != nil {
			return errors.New("密码格式错误")
		}
		validated_data.Password = &encryptPassword
	}

	UpdateTime := goi.GetTime()
	validated_data.UpdateTime = &UpdateTime

	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)

	_, err := mysqlDB.Where("`id` = ?", self.Id).Update(validated_data)
	if err != nil {
		return errors.New("修改用户错误")
	}
	utils.Update(self, validated_data)
	return nil
}

func (self UserModel) Delete() error {
	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)
	_, err := mysqlDB.Where("`id` = ?", self.Id).Delete()
	if err != nil {
		return errors.New("删除用户错误")
	}
	return nil
}
