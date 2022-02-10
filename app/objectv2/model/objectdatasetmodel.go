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
	objectDataSetFieldNames          = builder.RawFieldNames(&ObjectDataSet{})
	objectDataSetRows                = strings.Join(objectDataSetFieldNames, ",")
	objectDataSetRowsExpectAutoSet   = strings.Join(stringx.Remove(objectDataSetFieldNames, "`create_time`", "`update_time`"), ",")
	objectDataSetRowsWithPlaceHolder = strings.Join(stringx.Remove(objectDataSetFieldNames, "`jing_uuid`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ObjectDataSetModel interface {
		Insert(data *ObjectDataSet) (sql.Result, error)
		FindOne(jingUuid string) (*ObjectDataSet, error)
		Update(data *ObjectDataSet) error
		Delete(jingUuid string) error
	}

	defaultObjectDataSetModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ObjectDataSet struct {
		JingUuid         string       `db:"jing_uuid"`         // 主键ID
		Mid              int64        `db:"mid"`               // 租户id，等于0时为公共数据集
		Name             string       `db:"name"`              // 数据集名称，比如行业，地区
		EnglishName      string       `db:"english_name"`      // 数据集英文名称
		Code             string       `db:"code"`              // 数据集code，比如行业=>industry，地区=>region
		Tp               int64        `db:"tp"`                // 数据集类型，1-tree
		Remark           string       `db:"remark"`            // 备注
		CreatedTimestamp time.Time    `db:"created_timestamp"` // 创建时间
		UpdatedTimestamp time.Time    `db:"updated_timestamp"` // 更新时间
		DeletedAt        sql.NullTime `db:"deleted_at"`        // 软删除时间
	}
)

func NewObjectDataSetModel(conn sqlx.SqlConn) ObjectDataSetModel {
	return &defaultObjectDataSetModel{
		conn:  conn,
		table: "`object_data_set`",
	}
}

func (m *defaultObjectDataSetModel) Insert(data *ObjectDataSet) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, objectDataSetRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.JingUuid, data.Mid, data.Name, data.EnglishName, data.Code, data.Tp, data.Remark, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt)
	return ret, err
}

func (m *defaultObjectDataSetModel) FindOne(jingUuid string) (*ObjectDataSet, error) {
	query := fmt.Sprintf("select %s from %s where `jing_uuid` = ? limit 1", objectDataSetRows, m.table)
	var resp ObjectDataSet
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

func (m *defaultObjectDataSetModel) Update(data *ObjectDataSet) error {
	query := fmt.Sprintf("update %s set %s where `jing_uuid` = ?", m.table, objectDataSetRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Mid, data.Name, data.EnglishName, data.Code, data.Tp, data.Remark, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt, data.JingUuid)
	return err
}

func (m *defaultObjectDataSetModel) Delete(jingUuid string) error {
	query := fmt.Sprintf("delete from %s where `jing_uuid` = ?", m.table)
	_, err := m.conn.Exec(query, jingUuid)
	return err
}
