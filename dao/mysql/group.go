package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Group struct {
	ID       int    `db:"id"`
	Group    string `db:"group"`
	Leverage int    `db:"leverage"`
}

func GetGroups() ([]Group, error) {
	var groups []Group
	query := "SELECT * FROM `group`"
	err := db.Select(&groups, query)
	if err != nil {
		zap.L().Error("query groups failed", zap.Error(err))
		return nil, err
	}
	return groups, nil
}

func GetGroupByGroup(groupName string) (Group, error) {
	var group Group
	query := "SELECT * FROM `group` WHERE `group` = ?"
	err := db.Get(&group, query, groupName)
	if err != nil {
		zap.L().Error("query group by group failed", zap.Error(err))
		return Group{}, err
	}
	return group, nil
}

// UpdateGroup 更新 Group 记录
func UpdateGroup(g *Group) error {
	query := "UPDATE `group` SET "
	var params []interface{}

	if g.Group != "" {
		query += "`group`=?, "
		params = append(params, g.Group)
	}
	if g.Leverage != 0 {
		query += "leverage=?, "
		params = append(params, g.Leverage)
	}

	// 移除最后的逗号和空格
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id=%d", g.ID)

	query, params, err := sqlx.In(query, params...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)
	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
