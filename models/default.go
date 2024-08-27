package models

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

// Snowflake 노드 초기화
var node *snowflake.Node

func init() {
	var err error
	// 노드 번호 1로 Snowflake 노드 생성
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

type Model struct {
	ID int64 `json:"id,string,omitempty" gorm:"primaryKey,autoIncrement=false"`
	UpdateModel
}

// BeforeCreate 훅 정의
func (m *Model) BeforeCreate() error {
	m.ID = node.Generate().Int64()
	return nil
}

type UpdateModel struct {
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
