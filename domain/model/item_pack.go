package model

import "github.com/muhafs/go-restaurant-management/domain/entity"

type ItemPack struct {
	TableID string        `json:"table_id"`
	Items   []entity.Item `json:"items"`
}
