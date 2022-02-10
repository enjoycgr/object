package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	objectFieldNames          = builder.RawFieldNames(&Object{})
	objectRows                = strings.Join(objectFieldNames, ",")
	objectRowsExpectAutoSet   = strings.Join(stringx.Remove(objectFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	objectRowsWithPlaceHolder = strings.Join(stringx.Remove(objectFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ObjectModel interface {
		Insert(data *Object) (sql.Result, error)
		FindOne(id int64) (*Object, error)
		Update(data *Object) error
		Delete(id int64) error
	}

	defaultObjectModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Object struct {
		Id               int64        `db:"id"`
		JingUuid         string       `db:"jing_uuid"`
		Mid              int64        `db:"mid"`               // 租户id
		Name             string       `db:"name"`              // 对象名称
		Label            string       `db:"label"`             // 对象展示名称
		EnglishLabel     string       `db:"english_label"`     // 对象英文展示名称
		Description      string       `db:"description"`       // 对象描述
		IsPreset         int64        `db:"is_preset"`         // 是否系统预设，1-是，0-否
		Status           int64        `db:"status"`            // 对象状态，1-开启，0-关闭
		CreatedBy        string       `db:"created_by"`        // 创建者id
		UpdatedBy        string       `db:"updated_by"`        // 最近修改人
		CreatedTimestamp time.Time    `db:"created_timestamp"` // 创建时间
		UpdatedTimestamp time.Time    `db:"updated_timestamp"` // 更新时间
		DeletedAt        sql.NullTime `db:"deleted_at"`        // 软删除时间
	}
)

func NewObjectModel(conn sqlx.SqlConn) ObjectModel {
	return &defaultObjectModel{
		conn:  conn,
		table: "`object`",
	}
}

func (m *defaultObjectModel) Insert(data *Object) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, objectRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.JingUuid, data.Mid, data.Name, data.Label, data.EnglishLabel, data.Description, data.IsPreset, data.Status, data.CreatedBy, data.UpdatedBy, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt)
	return ret, err
}

func (m *defaultObjectModel) FindOne(id int64) (*Object, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", objectRows, m.table)
	var resp Object
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultObjectModel) Update(data *Object) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, objectRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.JingUuid, data.Mid, data.Name, data.Label, data.EnglishLabel, data.Description, data.IsPreset, data.Status, data.CreatedBy, data.UpdatedBy, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt, data.Id)
	return err
}

func (m *defaultObjectModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
