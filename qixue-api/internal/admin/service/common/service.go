package admincommon

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

type Service struct {
	areas adminrepo.AdminAreaCommonRepository
	db    *gorm.DB
	rdb   *redis.Client
}

func New(areas adminrepo.AdminAreaCommonRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{areas: areas, db: db, rdb: rdb}
}

type AreaTree struct {
	ID       int                  `json:"id"`
	Name     string               `json:"name"`
	ParentID int                  `json:"parent_id"`
	Children []adminrepo.AreaItem `json:"children"`
}

func (s *Service) List(ctx context.Context, adminID int) ([]adminrepo.AreaItem, error) {
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return nil, err
	}
	list, err := s.areas.ListVisible(ctx, cityID, limited)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []adminrepo.AreaItem{}
	}
	return list, nil
}

func (s *Service) Tree(ctx context.Context, adminID int) ([]AreaTree, error) {
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return nil, err
	}
	parents, err := s.areas.ListParents(ctx, cityID, limited)
	if err != nil {
		return nil, err
	}
	children, err := s.areas.ListChildren(ctx, cityID, limited)
	if err != nil {
		return nil, err
	}
	grouped := map[int][]adminrepo.AreaItem{}
	for _, child := range children {
		grouped[child.ParentID] = append(grouped[child.ParentID], child)
	}
	list := make([]AreaTree, 0, len(parents))
	for _, parent := range parents {
		kids := grouped[parent.ID]
		if kids == nil {
			kids = []adminrepo.AreaItem{}
		}
		list = append(list, AreaTree{
			ID: parent.ID, Name: parent.Name, ParentID: parent.ParentID, Children: kids,
		})
	}
	return list, nil
}
