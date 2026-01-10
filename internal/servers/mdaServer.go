package servers

import (
	"fme_backend/internal/models"
	"fme_backend/internal/schemas"

	"gorm.io/gorm"
)

type MdaServer struct{}

func (m *MdaServer) CreateMda(db *gorm.DB, data schemas.MdaCreateSchema) error {
	return nil
}

func (m *MdaServer) GetAllMdas(db *gorm.DB, limit, offset int) ([]schemas.GetAllMdaSchema, error) {
	return nil, nil
}

func (m *MdaServer) GetMdaById(db *gorm.DB, id int) (schemas.GetMdaSchema, error) {
	return schemas.GetMdaSchema{}, nil
}

func (m *MdaServer) UpdateMda(db *gorm.DB, data schemas.UpdateMdaSchema) (models.Mda, error) {
	return models.Mda{}, nil
}

func (m *MdaServer) SearchMda(db *gorm.DB, query string) ([]schemas.GetMdaSchema, error) {
	return nil, nil
}

func (m *MdaServer) FilterMdaAscending(db *gorm.DB, query string) ([]schemas.GetMdaSchema, error) {
	return nil, nil
}

func (m *MdaServer) SuspendMda(db *gorm.DB, id string) error {
	return nil
}

func (m *MdaServer) ActivateMda(db *gorm.DB, id string) error {
	return nil
}

func (m *MdaServer) MdaTotal(db *gorm.DB) (any, error) {
	return nil, nil
}

func (m *MdaServer) GetMdaProfile(db *gorm.DB, id int) ([]schemas.GetMdaSchema, error) {
	return nil, nil
}

func (m *MdaServer) DownloadMdaCsv(db *gorm.DB, id string) ([]schemas.GetMdaSchema, error) {
	return nil, nil
}

func (m *MdaServer) EditMdaData(db *gorm.DB, data schemas.MdaCreateSchema, id int) error {
	return nil
}
