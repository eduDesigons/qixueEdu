package adminpath

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

type pathCourseRow struct {
	ID        int64      `gorm:"column:id;primaryKey"`
	PathID    int        `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	Sort      int        `gorm:"column:sort"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (pathCourseRow) TableName() string { return "path_courses" }

func (s *Service) GetPathCourseSearch(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	title := strings.TrimSpace(asString(raw["title"]))
	cid := asInt(raw["cid"])
	typ := asInt(raw["type"])
	pathID := asInt(raw["path_id"])
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	if size > 50 {
		size = 50
	}

	ctx := c.Request.Context()
	rows, total, err := s.paths.SearchCourses(ctx, repo.PathCourseSearchFilter{
		Title: title, Category: cid, Type: typ, PathID: pathID,
		Offset: (current - 1) * size, Limit: size,
	})
	if response.WriteErr(c, err) {
		return
	}

	cateIDs := map[int]struct{}{}
	for _, row := range rows {
		if row.CategoryID > 0 {
			cateIDs[row.CategoryID] = struct{}{}
		}
	}
	categories := map[int]gin.H{}
	if len(cateIDs) > 0 {
		cats, err := s.paths.ListCourseCategories(ctx, keys(cateIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, cat := range cats {
			categories[cat.ID] = gin.H{"id": cat.ID, "name": cat.Name}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var cat any
		if found, ok := categories[row.CategoryID]; ok {
			cat = found
		}
		list = append(list, gin.H{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb,
			"min_need_credit": row.MinNeedCredit, "credit": row.Credit,
			"type": row.Type, "category_id": row.CategoryID, "category": cat,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) GetPathCourseList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	pathID := bind.RequestInt(c, raw, "pathId", "path_id")
	if pathID <= 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	ctx := c.Request.Context()
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if pathCount == 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}

	var rows []pathCourseRow
	if response.WriteErr(c, s.db.WithContext(ctx).
		Select("id, path_id, course_id, sort, created_at").
		Where("path_id = ?", pathID).Order("sort ASC").Find(&rows).Error) {
		return
	}
	courseIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		courseIDs = append(courseIDs, row.CourseID)
	}
	courses := map[int]gin.H{}
	if len(courseIDs) > 0 {
		courseRows, err := s.paths.ListCoursesDetail(ctx, courseIDs)
		if response.WriteErr(c, err) {
			return
		}
		cateIDs := map[int]struct{}{}
		for _, co := range courseRows {
			if co.CategoryID > 0 {
				cateIDs[co.CategoryID] = struct{}{}
			}
		}
		categories := map[int]gin.H{}
		if len(cateIDs) > 0 {
			cats, err := s.paths.ListCourseCategories(ctx, keys(cateIDs))
			if response.WriteErr(c, err) {
				return
			}
			for _, cat := range cats {
				categories[cat.ID] = gin.H{"id": cat.ID, "name": cat.Name}
			}
		}
		chapterCounts := map[int]int{}
		userCounts := map[int]int{}
		chapters, err := s.paths.ChapterCountsByCourseIDs(ctx, courseIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range chapters {
			chapterCounts[row.Key] = row.Total
		}
		users, err := s.paths.UserCountsByCourseIDs(ctx, courseIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range users {
			userCounts[row.Key] = row.Total
		}
		for _, co := range courseRows {
			var cat any
			if found, ok := categories[co.CategoryID]; ok {
				cat = found
			}
			courses[co.ID] = gin.H{
				"id": co.ID, "title": co.Title, "thumb": co.Thumb, "type": co.Type,
				"desc": co.Desc, "is_show": co.IsShow, "category_id": co.CategoryID,
				"is_rec": co.IsRec, "is_open": co.IsOpen, "min_need_credit": co.MinNeedCredit,
				"credit": co.Credit, "category": cat,
				"chapters_count": chapterCounts[co.ID], "users_count": userCounts[co.ID],
			}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var course any
		if found, ok := courses[row.CourseID]; ok {
			course = found
		}
		list = append(list, gin.H{
			"id": row.ID, "path_id": row.PathID, "course_id": row.CourseID,
			"sort": row.Sort, "created_at": model.FormatTime(row.CreatedAt), "course": course,
		})
	}
	response.OK(c, "", list)
}

func (s *Service) PathCourseListSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	pathID := asInt(raw["path_id"])
	list := toIntSlice(raw["course_list"])
	ctx := c.Request.Context()
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if pathCount == 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	if len(list) == 0 {
		response.Fail(c, 1, "请选择课程")
		return
	}
	existInPath, err := s.paths.CountPathCoursesIn(ctx, pathID, list)
	if response.WriteErr(c, err) {
		return
	}
	if existInPath > 0 {
		response.Fail(c, 1, "存在课程已经入过该系列了")
		return
	}
	existIDs, err := s.paths.PluckExistingCourseIDs(ctx, list)
	if response.WriteErr(c, err) {
		return
	}
	if len(existIDs) != len(list) {
		response.Fail(c, 1, "课程不存在")
		return
	}
	sort, err := s.paths.CountPathCoursesActive(ctx, pathID)
	if response.WriteErr(c, err) {
		return
	}
	now := time.Now()
	rows := make([]map[string]any, 0, len(list))
	for _, courseID := range list {
		sort++
		rows = append(rows, map[string]any{
			"path_id": pathID, "course_id": courseID, "sort": sort,
			"created_at": now, "updated_at": now,
		})
	}
	if response.WriteErr(c, s.paths.CreatePathCourses(ctx, rows)) {
		return
	}
	response.OK(c, "课节创建成功", []any{})
}

func (s *Service) PathCourseDelete(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		response.Fail(c, 1, "记录不存在")
		return
	}
	ctx := c.Request.Context()
	var row pathCourseRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", row.PathID).Count(&pathCount).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if pathCount == 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	userCount, err := s.paths.CountAnyPathUsers(ctx, row.PathID)
	if response.WriteErr(c, err) {
		return
	}
	if userCount > 0 {
		response.Fail(c, 1, "学员已报名，无法删除")
		return
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Delete(&row).Error) {
		return
	}
	response.OK(c, "删除成功", []any{})
}
