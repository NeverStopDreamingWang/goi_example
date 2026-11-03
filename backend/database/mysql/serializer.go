package mysql

import (
	"context"
	"errors"
	"fmt"

	"goi_example/backend/utils"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/auth"
	"github.com/NeverStopDreamingWang/goi/v3/db"
	"github.com/NeverStopDreamingWang/goi/v3/db/mysql"
)

func (self UserModel) Validate(ctx context.Context) error {
	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)

	if self.Id != nil {
		mysqlDB = mysqlDB.Where("`id` != ?", self.Id)
	}

	if self.Username != nil {
		flag, err := mysqlDB.Where("`username` = ?", self.Username).Exists(ctx)
		if err != nil {
			return fmt.Errorf("查询数据库错误: %w", err)
		}
		if flag == true {
			return errors.New("用户名重复")
		}
	}
	if self.Email != nil {
		flag, err := mysqlDB.Where("`email` = ?", self.Email).Exists(ctx)
		if err != nil {
			return fmt.Errorf("查询数据库错误: %w", err)
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

func (self *UserModel) Create(ctx context.Context) error {
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
	result, err := mysqlDB.Insert(ctx, self, true)
	if err != nil {
		return fmt.Errorf("添加用户错误: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("添加用户错误: %w", err)
	}
	self.Id = &id
	return nil
}

func (self *UserModel) Update(ctx context.Context, validated_data *UserModel) error {
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

	_, err := mysqlDB.Where("`id` = ?", self.Id).Update(ctx, validated_data, true)
	if err != nil {
		return fmt.Errorf("修改用户错误: %w", err)
	}
	utils.Update(self, validated_data)
	return nil
}

func (self UserModel) Delete(ctx context.Context) error {
	mysqlDB := db.Connect[*mysql.Engine]("mysql")
	mysqlDB.SetModel(self)
	_, err := mysqlDB.Where("`id` = ?", self.Id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("删除用户错误: %w", err)
	}
	return nil
}
