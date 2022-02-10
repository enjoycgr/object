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
	addressFieldNames          = builder.RawFieldNames(&Address{})
	addressRows                = strings.Join(addressFieldNames, ",")
	addressRowsExpectAutoSet   = strings.Join(stringx.Remove(addressFieldNames, "`create_time`", "`update_time`"), ",")
	addressRowsWithPlaceHolder = strings.Join(stringx.Remove(addressFieldNames, "`jing_uuid`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	AddressModel interface {
		Insert(data *Address) (sql.Result, error)
		FindOne(jingUuid string) (*Address, error)
		Update(data *Address) error
		Delete(jingUuid string) error
	}

	defaultAddressModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Address struct {
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

func NewAddressModel(conn sqlx.SqlConn) AddressModel {
	return &defaultAddressModel{
		conn:  conn,
		table: "`address`",
	}
}

func (m *defaultAddressModel) Insert(data *Address) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, addressRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.JingUuid, data.Mid, data.ParentId, data.Code, data.Value, data.Order, data.Level, data.Remark, data.DataSetCode, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt)
	return ret, err
}

func (m *defaultAddressModel) FindOne(jingUuid string) (*Address, error) {
	query := fmt.Sprintf("select %s from %s where `jing_uuid` = ? limit 1", addressRows, m.table)
	var resp Address
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

func (m *defaultAddressModel) Update(data *Address) error {
	query := fmt.Sprintf("update %s set %s where `jing_uuid` = ?", m.table, addressRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Mid, data.ParentId, data.Code, data.Value, data.Order, data.Level, data.Remark, data.DataSetCode, data.CreatedTimestamp, data.UpdatedTimestamp, data.DeletedAt, data.JingUuid)
	return err
}

func (m *defaultAddressModel) Delete(jingUuid string) error {
	query := fmt.Sprintf("delete from %s where `jing_uuid` = ?", m.table)
	_, err := m.conn.Exec(query, jingUuid)
	return err
}
