package admincategory

import (
	"strings"
	"time"
	"unicode/utf8"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/catalog"
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

type courseNode struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Icon          string       `json:"icon"`
	ParentID      int          `json:"parent_id"`
	IsShow        int          `json:"is_show"`
	Sort          int          `json:"sort"`
	MinNeedCredit int          `json:"min_need_credit"`
	BaseCredit    int          `json:"base_credit"`
	Count         int          `json:"count"`
	Children      []courseNode `json:"children,omitempty"`
}

type pathNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Sort     int    `json:"sort"`
	Count    int    `json:"count"`
}

type optionChild struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ParentID      int    `json:"parent_id"`
	MinNeedCredit int    `json:"min_need_credit"`
}

type optionParent struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	ParentID      int           `json:"parent_id"`
	MinNeedCredit int           `json:"min_need_credit"`
	BaseCredit    int           `json:"base_credit"`
	Children      []optionChild `json:"children"`
}

func (s *Service) List(c *gin.Context) {
	if c.Query("type") == "path" {
		list, err := s.pathList(c)
		if response.WriteErr(c, err) {
			return
		}
		response.OK(c, "", list)
		return
	}
	list, err := s.courseList(c)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (s *Service) Save(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	name := strings.TrimSpace(asString(raw["name"]))
	if name == "" {
		response.Fail(c, 1, "请填写分类名称")
		return
	}
	if utf8.RuneCountInString(name) > 32 {
		response.Fail(c, 1, "分类名称不能超过 32 个字")
		return
	}
	if asString(raw["type"]) == "path" {
		s.savePath(c, asInt(raw["id"]), name, asInt(raw["sort"]))
		return
	}
	s.saveCourse(c, raw, name)
}

func (s *Service) Delete(c *gin.Context) {
	var body struct {
		Type string   `json:"type"`
		ID   bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	if body.Type == "path" {
		var count int64
		if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", body.ID.Int()).Count(&count).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if count == 0 {
			response.Fail(c, 1, "分类不存在")
			return
		}
		var err error
		count, err = s.categories.CountPathsByCategory(ctx, body.ID.Int())
		if err != nil {
			response.WriteErr(c, err)
			return
		}
		if count > 0 {
			response.Fail(c, 1, "该分类下还有系列，不能删除")
			return
		}
		if response.WriteErr(c, s.db.WithContext(ctx).Delete(&pathRow{}, body.ID.Int()).Error) {
			return
		}
		catalog.Forget(ctx, s.rdb)
		response.OK(c, "删除成功", []any{})
		return
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", body.ID.Int()).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "分类不存在")
		return
	}
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ?", body.ID.Int()).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count > 0 {
		response.Fail(c, 1, "请先删除下面的二级分类")
		return
	}
	var err error
	count, err = s.categories.CountCoursesByCategory(ctx, body.ID.Int())
	if err != nil {
		response.WriteErr(c, err)
		return
	}
	if count > 0 {
		response.Fail(c, 1, "该分类下还有课程，不能删除")
		return
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Delete(&courseRow{}, body.ID.Int()).Error) {
		return
	}
	catalog.Forget(ctx, s.rdb)
	response.OK(c, "删除成功", []any{})
}

func (s *Service) Options(c *gin.Context) {
	ctx := c.Request.Context()
	var parents []courseRow
	err := s.db.WithContext(ctx).Where("parent_id = 0").Order("sort ASC, id ASC").Find(&parents).Error
	if response.WriteErr(c, err) {
		return
	}
	var children []courseRow
	if len(parents) > 0 {
		ids := make([]int, 0, len(parents))
		for _, row := range parents {
			ids = append(ids, row.ID)
		}
		err = s.db.WithContext(ctx).Where("parent_id IN ?", ids).Order("sort ASC, id ASC").Find(&children).Error
		if response.WriteErr(c, err) {
			return
		}
	}
	grouped := map[int][]optionChild{}
	for _, row := range children {
		grouped[row.ParentID] = append(grouped[row.ParentID], optionChild{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID, MinNeedCredit: row.MinNeedCredit,
		})
	}
	list := make([]optionParent, 0, len(parents))
	for _, row := range parents {
		kids := grouped[row.ID]
		if kids == nil {
			kids = []optionChild{}
		}
		list = append(list, optionParent{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID,
			MinNeedCredit: row.MinNeedCredit, BaseCredit: row.BaseCredit, Children: kids,
		})
	}
	response.OK(c, "", list)
}

func (s *Service) courseList(c *gin.Context) ([]courseNode, error) {
	counts, err := s.counts(c, "courses")
	if err != nil {
		return nil, err
	}
	var rows []courseRow
	if err := s.db.WithContext(c.Request.Context()).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	children := map[int][]courseNode{}
	for _, row := range rows {
		if row.ParentID == 0 {
			continue
		}
		children[row.ParentID] = append(children[row.ParentID], courseView(row, counts))
	}
	list := make([]courseNode, 0)
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

func (s *Service) pathList(c *gin.Context) ([]pathNode, error) {
	counts, err := s.counts(c, "paths")
	if err != nil {
		return nil, err
	}
	var rows []pathRow
	if err := s.db.WithContext(c.Request.Context()).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]pathNode, 0, len(rows))
	for _, row := range rows {
		list = append(list, pathNode{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID, Sort: row.Sort, Count: counts[row.ID],
		})
	}
	return list, nil
}

func (s *Service) counts(c *gin.Context, table string) (map[int]int, error) {
	return s.categories.CategoryItemCounts(c.Request.Context(), table)
}

func (s *Service) saveCourse(c *gin.Context, raw map[string]any, name string) {
	parentID := asInt(raw["parent_id"])
	if parentID < 0 {
		parentID = 0
	}
	isShow := 1
	if _, ok := raw["is_show"]; ok && asInt(raw["is_show"]) != 1 {
		isShow = 0
	}
	icon := strings.TrimSpace(asString(raw["icon"]))
	if utf8.RuneCountInString(icon) > 255 {
		response.Fail(c, 1, "图标地址过长")
		return
	}
	ctx := c.Request.Context()
	if parentID > 0 {
		var parent courseRow
		err := s.db.WithContext(ctx).Where("id = ?", parentID).Take(&parent).Error
		if err == gorm.ErrRecordNotFound || parent.ParentID != 0 {
			response.Fail(c, 1, "上级分类不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		used, cerr := s.categories.CountCoursesByCategory(ctx, parentID)
		if cerr != nil {
			response.WriteErr(c, cerr)
			return
		}
		if used > 0 {
			response.Fail(c, 1, "这个一级分类下已经有课程，不能再加二级分类")
			return
		}
	}
	id := asInt(raw["id"])
	var dup int64
	query := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ? AND name = ?", parentID, name)
	if id > 0 {
		query = query.Where("id <> ?", id)
	}
	if err := query.Count(&dup).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if dup > 0 {
		response.Fail(c, 1, "同级已有这个分类名")
		return
	}
	now := time.Now()
	if id > 0 {
		var row courseRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "分类不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		if parentID == id {
			response.Fail(c, 1, "不能把自己设为上级")
			return
		}
		var children int64
		if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("parent_id = ?", id).Count(&children).Error; err != nil {
			response.WriteErr(c, err)
			return
		}
		if children > 0 && parentID != 0 {
			response.Fail(c, 1, "已有二级分类，不能再改成二级")
			return
		}
		if parentID == 0 {
			used, err := s.categories.CountCoursesByCategory(ctx, id)
			if err != nil {
				response.WriteErr(c, err)
				return
			}
			if used > 0 {
				response.Fail(c, 1, "已有课程挂在这个分类上，不能改成一级")
				return
			}
		}
		updates := map[string]any{
			"name": name, "parent_id": parentID, "sort": asInt(raw["sort"]), "is_show": isShow, "icon": icon, "updated_at": now,
		}
		if err := applyCredit(raw, updates); err != nil {
			response.WriteErr(c, err)
			return
		}
		if response.WriteErr(c, s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", id).Updates(updates).Error) {
			return
		}
	} else {
		row := courseRow{
			Name: name, Icon: icon, ParentID: parentID, IsShow: isShow, Sort: asInt(raw["sort"]),
			CreatedAt: &now, UpdatedAt: &now,
		}
		credits := map[string]any{}
		if err := applyCredit(raw, credits); err != nil {
			response.WriteErr(c, err)
			return
		}
		if value, ok := credits["min_need_credit"]; ok {
			row.MinNeedCredit = value.(int)
		}
		if value, ok := credits["base_credit"]; ok {
			row.BaseCredit = value.(int)
		}
		if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
			return
		}
		id = row.ID
	}
	catalog.Forget(ctx, s.rdb)
	response.OK(c, "保存成功", gin.H{"id": id})
}

func (s *Service) savePath(c *gin.Context, id int, name string, sort int) {
	ctx := c.Request.Context()
	var dup int64
	query := s.db.WithContext(ctx).Model(&pathRow{}).Where("name = ?", name)
	if id > 0 {
		query = query.Where("id <> ?", id)
	}
	if err := query.Count(&dup).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if dup > 0 {
		response.Fail(c, 1, "已有这个分类名")
		return
	}
	now := time.Now()
	if id > 0 {
		result := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).Updates(map[string]any{
			"name": name, "parent_id": 0, "sort": sort, "updated_at": now,
		})
		if response.WriteErr(c, result.Error) {
			return
		}
		if result.RowsAffected == 0 {
			var count int64
			if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", id).Count(&count).Error; err != nil {
				response.WriteErr(c, err)
				return
			}
			if count == 0 {
				response.Fail(c, 1, "分类不存在")
				return
			}
		}
	} else {
		row := pathRow{ParentID: 0, Name: name, Sort: sort, CreatedAt: &now, UpdatedAt: &now}
		if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
			return
		}
		id = row.ID
	}
	catalog.Forget(ctx, s.rdb)
	response.OK(c, "保存成功", gin.H{"id": id})
}

func courseView(row courseRow, counts map[int]int) courseNode {
	return courseNode{
		ID: row.ID, Name: row.Name, Icon: row.Icon, ParentID: row.ParentID, IsShow: row.IsShow,
		Sort: row.Sort, MinNeedCredit: row.MinNeedCredit, BaseCredit: row.BaseCredit, Count: counts[row.ID],
	}
}

func applyCredit(raw map[string]any, updates map[string]any) error {
	fields := map[string]string{"min_need_credit": "限制学分", "base_credit": "基础学分"}
	for field, label := range fields {
		if _, ok := raw[field]; !ok {
			continue
		}
		credit := asInt(raw[field])
		if credit < 0 {
			return response.BizErr(label + "不能为负数")
		}
		updates[field] = credit
	}
	return nil
}

func asString(value any) string {
	text, _ := value.(string)
	return text
}

func asInt(value any) int { return bind.AsInt(value) }
