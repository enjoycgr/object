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
	objectDataSetTreeFieldNames          = builder.RawFieldNames(&ObjectDataSetTree{})
	objectDataSetTreeRows                = strings.Join(objectDataSetTreeFieldNames, ",")
	objectDataSetTreeRowsExpectAutoSet   = strings.Join(stringx.Remove(objectDataSetTreeFieldNames, "`create_time`", "`update_time`"), ",")
	objectDataSetTreeRowsWithPlaceHolder = strings.Join(stringx.Remove(objectDataSetTreeFieldNames, "`jing_uuid`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ObjectDataSetTreeModel interface {
		Insert(data *ObjectDataSetTree) (sql.Result, error)
		FindOne(jingUuid string) (*ObjectDataSetTree, error)
		Update(data *ObjectDataSetTree) error
		Delete(jingUuid string) error
	}

	defaultObjectDataSetTreeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ObjectDataSetTree struct {
		JingUuid         string       `db:"jing_uuid"`
		Mid              int64        `db:"mid"`               // 租户id，等于0时为公共数据集
		ParentId         string       `db:"parent_id"`         // 上级id
		Code             string       `db:"code"`              // 数据编码
		Value            string       `db:"value"`             // 数据名称
		Order            int64        `db:"order"`             // 当前层级排序
		Level            int64        `db:"level"`             // 层级
		Remark           string       `db:"remark"`            // 备注
		DataSetCode      string       `db:"data_set_code"`     // jing_object_v2_data_set的code
		CreatedTimestamp time.Time    `db:"created_timestamp"` // 创建时间
		UpdatedTimestamp time.Time    `db:"updated_timestamp"` // 更新时间
		DeletedAt        sql.NullTime `db:"deleted_at"`        // 软删除时间
	}
)

func NewObjectDataSetTreeModel(conn sqlx.SqlConn) ObjectDataSetTreeModel {
	return &defaultObjectDataSetTreeModel{
		conn:  conn,
		table: "`object_data_set_tree`",
	}
}

func (m *defaultObjectDataSetTreeModel) Insert(data *ObjectDataSetTree) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, objectDataSetTreeRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.JingUuid, data.Mid, data.ParentId, data.Code, data.Value, data.Order, data.Level, data.Remark, data.DataSetCode, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt)
	return ret, err
}

func (m *defaultObjectDataSetTreeModel) FindOne(jingUuid string) (*ObjectDataSetTree, error) {
	query := fmt.Sprintf("select %s from %s where `jing_uuid` = ? limit 1", objectDataSetTreeRows, m.table)
	var resp ObjectDataSetTree
	err := m.conn.QueryRow(&resp, query, jingUuid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultObjectDataSetTreeModel) Update(data *ObjectDataSetTree) error {
	query := fmt.Sprintf("update %s set %s where `jing_uuid` = ?", m.table, objectDataSetTreeRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Mid, data.ParentId, data.Code, data.Value, data.Order, data.Level, data.Remark, data.DataSetCode, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt, data.JingUuid)
	return err
}

func (m *defaultObjectDataSetTreeModel) Delete(jingUuid string) error {
	query := fmt.Sprintf("delete from %s where `jing_uuid` = ?", m.table)
	_, err := m.conn.Exec(query, jingUuid)
	return err
}
