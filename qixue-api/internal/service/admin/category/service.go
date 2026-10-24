package admincategory

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/repo"
	"qixue-api/internal/service/catalog"
)

var (
	ErrBadParams             = errors.New("参数错误")
	ErrEmptyName             = errors.New("请填写分类名称")
	ErrNameTooLong           = errors.New("分类名称不能超过 32 个字")
	ErrNotFound              = errors.New("分类不存在")
	ErrHasPaths              = errors.New("该分类下还有系列，不能删除")
	ErrHasChildren           = errors.New("请先删除下面的二级分类")
	ErrHasCourses            = errors.New("该分类下还有课程，不能删除")
	ErrIconTooLong           = errors.New("图标地址过长")
	ErrParentMissing         = errors.New("上级分类不存在")
	ErrParentHasCourses      = errors.New("这个一级分类下已经有课程，不能再加二级分类")
	ErrDupName               = errors.New("同级已有这个分类名")
	ErrSelfParent            = errors.New("不能把自己设为上级")
	ErrHasKidsCantDemote     = errors.New("已有二级分类，不能再改成二级")
	ErrHasCoursesCantPromote = errors.New("已有课程挂在这个分类上，不能改成一级")
	ErrPathDupName           = errors.New("已有这个分类名")
	ErrMinCreditNegative     = errors.New("限制学分不能为负数")
	ErrBaseCreditNegative    = errors.New("基础学分不能为负数")
)

type Service struct {
	categories repo.AdminCategoryRepository
	db         *gorm.DB
	rdb        *redis.Client
}

func New(categories repo.AdminCategoryRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{categories: categories, db: db, rdb: rdb}
}

type courseRow struct {
	ID            int    `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name"`
	Icon          string `gorm:"column:icon"`
	ParentID      int    `gorm:"column:parent_id"`
	IsShow        int    `gorm:"column:is_show"`
	Sort          int    `gorm:"column:sort"`
	MinNeedCredit int    `gorm:"column:min_need_credit"`
	BaseCredit    int    `gorm:"column:base_credit"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (courseRow) TableName() string { return "course_categories" }

type pathRow struct {
	ID        int    `gorm:"column:id;primaryKey"`
	ParentID  int    `gorm:"column:parent_id"`
	Name      string `gorm:"column:name"`
	Sort      int    `gorm:"column:sort"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (pathRow) TableName() string { return "path_categories" }

type CourseNode struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Icon          string       `json:"icon"`
	ParentID      int          `json:"parent_id"`
	IsShow        int          `json:"is_show"`
	Sort          int          `json:"sort"`
	MinNeedCredit int          `json:"min_need_credit"`
	BaseCredit    int          `json:"base_credit"`
	Count         int          `json:"count"`
	Children      []CourseNode `json:"children,omitempty"`
}

type PathNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Sort     int    `json:"sort"`
	Count    int    `json:"count"`
}

type OptionChild struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ParentID      int    `json:"parent_id"`
	MinNeedCredit int    `json:"min_need_credit"`
}

type OptionParent struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	ParentID      int           `json:"parent_id"`
	MinNeedCredit int           `json:"min_need_credit"`
	BaseCredit    int           `json:"base_credit"`
	Children      []OptionChild `json:"children"`
}

type SaveInput struct {
	Type          string
	ID            int
	Name          string
	Sort          int
	ParentID      int
	IsShow        int
	Icon          string
	MinNeedCredit *int
	BaseCredit    *int
}

type SaveResult struct {
	ID int `json:"id"`
}

func (s *Service) List(ctx context.Context, typ string) (any, error) {
	if typ == "path" {
		return s.pathList(ctx)
	}
	return s.courseList(ctx)
}

func (s *Service) Save(ctx context.Context, in SaveInput) (SaveResult, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return SaveResult{}, ErrEmptyName
	}
	if utf8.RuneCountInString(name) > 32 {
		return SaveResult{}, ErrNameTooLong
	}
	if in.Type == "path" {
		id, err := s.savePath(ctx, in.ID, name, in.Sort)
		return SaveResult{ID: id}, err
	}
	id, err := s.saveCourse(ctx, in, name)
	return SaveResult{ID: id}, err
}

func (s *Service) Delete(ctx context.Context, typ string, id int) error {
	if id <= 0 {
		return ErrBadParams
	}
	if typ == "path" {
		var count int64
		if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return ErrNotFound
		}
		count, err := s.categories.CountPathsByCategory(ctx, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrHasPaths
		}
		if err := s.db.WithContext(ctx).Delete(&pathRow{}, id).Error; err != nil {
			return err
		}
		catalog.Forget(ctx, s.rdb)
		return nil
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrHasChildren
	}
	count, err := s.categories.CountCoursesByCategory(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrHasCourses
	}
	if err := s.db.WithContext(ctx).Delete(&courseRow{}, id).Error; err != nil {
		return err
	}
	catalog.Forget(ctx, s.rdb)
	return nil
}

func (s *Service) Options(ctx context.Context) ([]OptionParent, error) {
	var parents []courseRow
	if err := s.db.WithContext(ctx).Where("parent_id = 0").Order("sort ASC, id ASC").Find(&parents).Error; err != nil {
		return nil, err
	}
	var children []courseRow
	if len(parents) > 0 {
		ids := make([]int, 0, len(parents))
		for _, row := range parents {
			ids = append(ids, row.ID)
		}
		if err := s.db.WithContext(ctx).Where("parent_id IN ?", ids).Order("sort ASC, id ASC").Find(&children).Error; err != nil {
			return nil, err
		}
	}
	grouped := map[int][]OptionChild{}
	for _, row := range children {
		grouped[row.ParentID] = append(grouped[row.ParentID], OptionChild{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID, MinNeedCredit: row.MinNeedCredit,
		})
	}
	list := make([]OptionParent, 0, len(parents))
	for _, row := range parents {
		kids := grouped[row.ID]
		if kids == nil {
			kids = []OptionChild{}
		}
		list = append(list, OptionParent{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID,
			MinNeedCredit: row.MinNeedCredit, BaseCredit: row.BaseCredit, Children: kids,
		})
	}
	return list, nil
}

func (s *Service) courseList(ctx context.Context) ([]CourseNode, error) {
	counts, err := s.categories.CategoryItemCounts(ctx, "courses")
	if err != nil {
		return nil, err
	}
	var rows []courseRow
	if err := s.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	children := map[int][]CourseNode{}
	for _, row := range rows {
		if row.ParentID == 0 {
			continue
		}
		children[row.ParentID] = append(children[row.ParentID], courseView(row, counts))
	}
	list := make([]CourseNode, 0)
	seen := map[int]struct{}{}
	for _, row := range rows {
		if row.ParentID != 0 {
			continue
		}
		node := courseView(row, counts)
		if kids := children[row.ID]; len(kids) > 0 {
			node.Children = kids
		}
		list = append(list, node)
		seen[row.ID] = struct{}{}
	}
	for parentID, kids := range children {
		if _, ok := seen[parentID]; ok {
			continue
		}
		list = append(list, kids...)
	}
	return list, nil
}

func (s *Service) pathList(ctx context.Context) ([]PathNode, error) {
	counts, err := s.categories.CategoryItemCounts(ctx, "paths")
	if err != nil {
		return nil, err
	}
	var rows []pathRow
	if err := s.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]PathNode, 0, len(rows))
	for _, row := range rows {
		list = append(list, PathNode{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID, Sort: row.Sort, Count: counts[row.ID],
		})
	}
	return list, nil
}

func (s *Service) saveCourse(ctx context.Context, in SaveInput, name string) (int, error) {
	parentID := in.ParentID
	if parentID < 0 {
		parentID = 0
	}
	isShow := in.IsShow
	if isShow != 1 {
		isShow = 0
	}
	icon := strings.TrimSpace(in.Icon)
	if utf8.RuneCountInString(icon) > 255 {
		return 0, ErrIconTooLong
	}
	if parentID > 0 {
		var parent courseRow
		err := s.db.WithContext(ctx).Where("id = ?", parentID).Take(&parent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) || parent.ParentID != 0 {
			return 0, ErrParentMissing
		}
		if err != nil {
			return 0, err
		}
		used, cerr := s.categories.CountCoursesByCategory(ctx, parentID)
		if cerr != nil {
			return 0, cerr
		}
		if used > 0 {
			return 0, ErrParentHasCourses
		}
	}
	id := in.ID
	var dup int64
	query := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ? AND name = ?", parentID, name)
	if id > 0 {
		query = query.Where("id <> ?", id)
	}
	if err := query.Count(&dup).Error; err != nil {
		return 0, err
	}
	if dup > 0 {
		return 0, ErrDupName
	}
	now := time.Now()
	if id > 0 {
		var row courseRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		if err != nil {
			return 0, err
		}
		if parentID == id {
			return 0, ErrSelfParent
		}
		var children int64
		if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ?", id).Count(&children).Error; err != nil {
			return 0, err
		}
		if children > 0 && parentID != 0 {
			return 0, ErrHasKidsCantDemote
		}
		if parentID == 0 {
			used, err := s.categories.CountCoursesByCategory(ctx, id)
			if err != nil {
				return 0, err
			}
			if used > 0 {
				return 0, ErrHasCoursesCantPromote
			}
		}
		updates := map[string]any{
			"name": name, "parent_id": parentID, "sort": in.Sort, "is_show": isShow, "icon": icon, "updated_at": now,
		}
		if err := applyCredit(in, updates); err != nil {
			return 0, err
		}
		if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return 0, err
		}
	} else {
		row := courseRow{
			Name: name, Icon: icon, ParentID: parentID, IsShow: isShow, Sort: in.Sort,
			CreatedAt: &now, UpdatedAt: &now,
		}
		credits := map[string]any{}
		if err := applyCredit(in, credits); err != nil {
			return 0, err
		}
		if value, ok := credits["min_need_credit"]; ok {
			row.MinNeedCredit = value.(int)
		}
		if value, ok := credits["base_credit"]; ok {
			row.BaseCredit = value.(int)
		}
		if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
			return 0, err
		}
		id = row.ID
	}
	catalog.Forget(ctx, s.rdb)
	return id, nil
}

func (s *Service) savePath(ctx context.Context, id int, name string, sort int) (int, error) {
	var dup int64
	query := s.db.WithContext(ctx).Model(&pathRow{}).Where("name = ?", name)
	if id > 0 {
		query = query.Where("id <> ?", id)
	}
	if err := query.Count(&dup).Error; err != nil {
		return 0, err
	}
	if dup > 0 {
		return 0, ErrPathDupName
	}
	now := time.Now()
	if id > 0 {
		result := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).Updates(map[string]any{
			"name": name, "parent_id": 0, "sort": sort, "updated_at": now,
		})
		if result.Error != nil {
			return 0, result.Error
		}
		if result.RowsAffected == 0 {
			var count int64
			if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).Count(&count).Error; err != nil {
				return 0, err
			}
			if count == 0 {
				return 0, ErrNotFound
			}
		}
	} else {
		row := pathRow{ParentID: 0, Name: name, Sort: sort, CreatedAt: &now, UpdatedAt: &now}
		if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
			return 0, err
		}
		id = row.ID
	}
	catalog.Forget(ctx, s.rdb)
	return id, nil
}

func courseView(row courseRow, counts map[int]int) CourseNode {
	return CourseNode{
		ID: row.ID, Name: row.Name, Icon: row.Icon, ParentID: row.ParentID, IsShow: row.IsShow,
		Sort: row.Sort, MinNeedCredit: row.MinNeedCredit, BaseCredit: row.BaseCredit, Count: counts[row.ID],
	}
}

func applyCredit(in SaveInput, updates map[string]any) error {
	if in.MinNeedCredit != nil {
		if *in.MinNeedCredit < 0 {
			return ErrMinCreditNegative
		}
		updates["min_need_credit"] = *in.MinNeedCredit
	}
	if in.BaseCredit != nil {
		if *in.BaseCredit < 0 {
			return ErrBaseCreditNegative
		}
		updates["base_credit"] = *in.BaseCredit
	}
	return nil
}
