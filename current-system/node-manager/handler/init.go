package handler

import (
	"github.com/dipghoshraj/media-service/node-manager/model"
)

type DBHandler struct {
	DbManager *model.DbManager
}

func NewDBHandler(DbManager *model.DbManager) *DBHandler {
	return &DBHandler{DbManager: DbManager}
}

// func NewNMHandler(NodeManager *handler.DBHandler) *NMHandler {
// 	return &NMHandler{DbManager: NodeManager}
// }
