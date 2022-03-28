package model

import (
	"PassageTwo/global"
	"PassageTwo/pkg/setting"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

// 从mysql-->go 的映射
// 公共Model字段
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databasesetting *setting.DataBaseSettings) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databasesetting.DBtype, fmt.Sprintf(
		s,
		databasesetting.UserName,
		databasesetting.Password,
		databasesetting.Host,
		databasesetting.DBName,
		databasesetting.Charset,
		databasesetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	// Debug for Log
	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	db.SingularTable(true)
	// 重新注册回调函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(databasesetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databasesetting.MaxOpenConns)

	return db, nil
}

// 1.新增行为的回调函数
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
	}
}

// 2.更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIsExist(scope.CombinedConditionSql()),
				addExtraSpaceIsExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIsExist(scope.CombinedConditionSql()),
				addExtraSpaceIsExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIsExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
