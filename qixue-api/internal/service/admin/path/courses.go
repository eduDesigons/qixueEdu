package adminpath

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/http/page"
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

func (s *Service) GetPathCourseSearch(ctx context.Context, raw map[string]any) (page.Body, error) {
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

	rows, total, err := s.paths.SearchCourses(ctx, repo.PathCourseSearchFilter{
		Title: title, Category: cid, Type: typ, PathID: pathID,
		Offset: (current - 1) * size, Limit: size,
	})
	if err != nil {
		return page.Body{}, err
	}

	cateIDs := map[int]struct{}{}
	for _, row := range rows {
		if row.CategoryID > 0 {
			cateIDs[row.CategoryID] = struct{}{}
		}
	}
	categories := map[int]map[string]any{}
	if len(cateIDs) > 0 {
		cats, err := s.paths.ListCourseCategories(ctx, keys(cateIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, cat := range cats {
			categories[cat.ID] = map[string]any{"id": cat.ID, "name": cat.Name}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var cat any
		if found, ok := categories[row.CategoryID]; ok {
			cat = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb,
			"min_need_credit": row.MinNeedCredit, "credit": row.Credit,
			"type": row.Type, "category_id": row.CategoryID, "category": cat,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetPathCourseList(ctx context.Context, pathID int) ([]map[string]any, error) {
	if pathID <= 0 {
		return nil, biz("系列不存在")
	}
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		return nil, err
	}
	if pathCount == 0 {
		return nil, biz("系列不存在")
	}

	var rows []pathCourseRow
	if err := s.db.WithContext(ctx).
		Select("id, path_id, course_id, sort, created_at").
		Where("path_id = ?", pathID).Order("sort ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	courseIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		courseIDs = append(courseIDs, row.CourseID)
	}
	courses := map[int]map[string]any{}
	if len(courseIDs) > 0 {
		courseRows, err := s.paths.ListCoursesDetail(ctx, courseIDs)
		if err != nil {
			return nil, err
		}
		cateIDs := map[int]struct{}{}
		for _, co := range courseRows {
			if co.CategoryID > 0 {
				cateIDs[co.CategoryID] = struct{}{}
			}
		}
		categories := map[int]map[string]any{}
		if len(cateIDs) > 0 {
			cats, err := s.paths.ListCourseCategories(ctx, keys(cateIDs))
			if err != nil {
				return nil, err
			}
			for _, cat := range cats {
				categories[cat.ID] = map[string]any{"id": cat.ID, "name": cat.Name}
			}
		}
		chapterCounts := map[int]int{}
		userCounts := map[int]int{}
		chapters, err := s.paths.ChapterCountsByCourseIDs(ctx, courseIDs)
		if err != nil {
			return nil, err
		}
		for _, row := range chapters {
			chapterCounts[row.Key] = row.Total
		}
		users, err := s.paths.UserCountsByCourseIDs(ctx, courseIDs)
		if err != nil {
			return nil, err
		}
		for _, row := range users {
			userCounts[row.Key] = row.Total
		}
		for _, co := range courseRows {
			var cat any
			if found, ok := categories[co.CategoryID]; ok {
				cat = found
			}
			courses[co.ID] = map[string]any{
				"id": co.ID, "title": co.Title, "thumb": co.Thumb, "type": co.Type,
				"desc": co.Desc, "is_show": co.IsShow, "category_id": co.CategoryID,
				"is_rec": co.IsRec, "is_open": co.IsOpen, "min_need_credit": co.MinNeedCredit,
				"credit": co.Credit, "category": cat,
				"chapters_count": chapterCounts[co.ID], "users_count": userCounts[co.ID],
			}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var course any
		if found, ok := courses[row.CourseID]; ok {
			course = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "path_id": row.PathID, "course_id": row.CourseID,
			"sort": row.Sort, "created_at": model.FormatTime(row.CreatedAt), "course": course,
		})
	}
	return list, nil
}

func (s *Service) PathCourseListSave(ctx context.Context, raw map[string]any) error {
	if raw == nil {
		return biz("参数错误")
	}
	pathID := asInt(raw["path_id"])
	list := toIntSlice(raw["course_list"])
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		return err
	}
	if pathCount == 0 {
		return biz("系列不存在")
	}
	if len(list) == 0 {
		return biz("请选择课程")
	}
	existInPath, err := s.paths.CountPathCoursesIn(ctx, pathID, list)
	if err != nil {
		return err
	}
	if existInPath > 0 {
		return biz("存在课程已经入过该系列了")
	}
	existIDs, err := s.paths.PluckExistingCourseIDs(ctx, list)
	if err != nil {
		return err
	}
	if len(existIDs) != len(list) {
		return biz("课程不存在")
	}
	sort, err := s.paths.CountPathCoursesActive(ctx, pathID)
	if err != nil {
		return err
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
	return s.paths.CreatePathCourses(ctx, rows)
}

func (s *Service) PathCourseDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("记录不存在")
	}
	var row pathCourseRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return biz("记录不存在")
	}
	if err != nil {
		return err
	}
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", row.PathID).Count(&pathCount).Error; err != nil {
		return err
	}
	if pathCount == 0 {
		return biz("系列不存在")
	}
	userCount, err := s.paths.CountAnyPathUsers(ctx, row.PathID)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return biz("学员已报名，无法删除")
	}
	return s.db.WithContext(ctx).Delete(&row).Error
}
