package catalog

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
)

type CourseNode struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Icon     *string      `json:"icon"`
	ParentID int          `json:"parent_id"`
	Sort     int          `json:"sort"`
	Children []CourseNode `json:"children"`
}

type PathNode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CourseTree(ctx context.Context, db *gorm.DB, rdb *redis.Client) ([]CourseNode, error) {
	var tree []CourseNode
	err := cache.Remember(ctx, rdb, "course_category_tree", &tree, func() (any, error) {
		return loadCourseTree(ctx, repo.NewCatalog(db))
	})
	if tree == nil {
		tree = []CourseNode{}
	}
	return tree, err
}

func PathList(ctx context.Context, db *gorm.DB, rdb *redis.Client) ([]PathNode, error) {
	var list []PathNode
	err := cache.Remember(ctx, rdb, "path_category_list", &list, func() (any, error) {
		return loadPathList(ctx, repo.NewCatalog(db))
	})
	if list == nil {
		list = []PathNode{}
	}
	return list, err
}

func Forget(ctx context.Context, rdb *redis.Client) {
	cache.Drop(ctx, rdb,
		"course_category_tree",
		"path_category_list",
		"home_boot_v1",
		"course_list_hot_p1",
		"path_list_p1",
	)
}

func loadCourseTree(ctx context.Context, catalog repo.CatalogRepository) ([]CourseNode, error) {
	parents, err := catalog.ListShownCourseParents(ctx)
	if err != nil {
		return nil, err
	}
	children, err := catalog.ListShownCourseChildren(ctx)
	if err != nil {
		return nil, err
	}
	grouped := map[int][]CourseNode{}
	for _, row := range children {
		grouped[row.ParentID] = append(grouped[row.ParentID], CourseNode{
			ID: row.ID, Name: row.Name, Icon: row.Icon, ParentID: row.ParentID, Sort: row.Sort, Children: []CourseNode{},
		})
	}
	tree := make([]CourseNode, 0, len(parents))
	for _, row := range parents {
		kids := grouped[row.ID]
		if kids == nil {
			kids = []CourseNode{}
		}
		tree = append(tree, CourseNode{
			ID: row.ID, Name: row.Name, Icon: row.Icon, ParentID: row.ParentID, Sort: row.Sort, Children: kids,
		})
	}
	return tree, nil
}

func loadPathList(ctx context.Context, catalog repo.CatalogRepository) ([]PathNode, error) {
	rows, err := catalog.ListPathRoots(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]PathNode, 0, len(rows))
	for _, row := range rows {
		list = append(list, PathNode{ID: row.ID, Name: row.Name})
	}
	return list, nil
}
