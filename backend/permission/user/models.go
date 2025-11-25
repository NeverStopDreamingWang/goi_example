package user

import (
	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
)

func init() {
	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.Migrate(UserModel{})
}

type UserStatusType uint16

const (
	DISABLE UserStatusType = iota // 禁用
	NORMAL                        // 正常
)

// 用户表
type UserModel struct {
	Id            *int64          `field_name:"id" field_type:"INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT" json:"id"` // ID
	Username      *string         `field_name:"username" field_type:"TEXT NOT NULL" json:"username"`                  // 用户名
	Password      *string         `field_name:"password" field_type:"TEXT NOT NULL" json:"-"`                         // 密码
	Email         *string         `field_name:"email" field_type:"TEXT NOT NULL" json:"email"`                        // 邮箱
	Status        *UserStatusType `field_name:"status" field_type:"INTEGER NOT NULL DEFAULT 1" json:"status"`         // 状态
	RoleId        *int64          `field_name:"role_id" field_type:"INTEGER NOT NULL" json:"role_id"`                 // 角色ID
	LastLoginTime *string         `field_name:"last_login_time" field_type:"DATETIME" json:"last_login_time"`         // 最后登录时间
	CreateTime    *string         `field_name:"create_time" field_type:"DATETIME NOT NULL" json:"create_time"`        // 创建时间
	UpdateTime    *string         `field_name:"update_time" field_type:"DATETIME" json:"update_time"`                 // 更新时间
}

// 设置表配置
func (userModel UserModel) ModelSet() *sqlite3.Settings {
	modelSettings := &sqlite3.Settings{
		MigrationsHandler: sqlite3.MigrationsHandler{ // 迁移时处理函数
			BeforeHandler: nil,      // 迁移之前处理函数
			AfterHandler:  initUser, // 迁移之后处理函数
		},

		TABLE_NAME: "tb_user", // 设置表名
		// 自定义配置
		Settings: goi.Params{},
	}

	return modelSettings
}

// 初始化数据
func initUser() error {
	initUserList := [][]any{
		{"admin", "admin", "admin@qq.com", UserStatusType(1), int64(1), int64(1)},
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	sqlite3DB.SetModel(UserModel{})
	total, err := sqlite3DB.Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	for _, item := range initUserList {
		var (
			Username = item[0].(string)
			Password = item[1].(string)
			Email    = item[2].(string)
			Status   = item[3].(UserStatusType)
			RoleId   = item[4].(int64)
		)
		user := &UserModel{
			Username: &Username,
			Password: &Password,
			Email:    &Email,
			Status:   &Status,
			RoleId:   &RoleId,
		}
		// 参数验证
		err = user.Validate()
		if err != nil {
			return err
		}
		// 添加
		err = user.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

type UserInfo struct {
	Id         *int64  `json:"id" bson:"id"`
	Username   *string `json:"username" bson:"username"`
	Email      *string `json:"email" bson:"email"`
	CreateTime *string `json:"create_time" bson:"create_Time"`
	UpdateTime *string `json:"update_time" bson:"update_Time"`
}

func (user UserModel) ToUserInfo() *UserInfo {
	return &UserInfo{
		Id:         user.Id,
		Username:   user.Username,
		Email:      user.Email,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}
