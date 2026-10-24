package adminops

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/scope"
)

type Service struct {
	ops repo.AdminOpsRepository
	db  *gorm.DB
	rdb *redis.Client
}

func New(ops repo.AdminOpsRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{ops: ops, db: db, rdb: rdb}
}

func (s *Service) CourseEnrollList(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Title    string   `json:"title"`
		IsEnd    any      `json:"is_end"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	f := repo.OpsListFilter{
		NickName: strings.TrimSpace(body.NickName),
		Title:    strings.TrimSpace(body.Title),
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	if isEnd, ok := asOptionalInt(body.IsEnd); ok {
		f.IsEnd = &isEnd
	}
	rows, total, err := s.ops.ListCourseEnrolls(ctx, f)
	if response.WriteErr(c, err) {
		return
	}
	userIDs, courseIDs := make([]int, 0, len(rows)), make([]int, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		courseIDs = append(courseIDs, r.CourseID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if response.WriteErr(c, err) {
		return
	}
	courses, err := s.ops.LoadTitleThumbs(ctx, "courses", uniquePositive(courseIDs))
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		u := users[r.UserID]
		course := courses[r.CourseID]
		list = append(list, gin.H{
			"id": r.ID, "user_id": r.UserID, "course_id": r.CourseID,
			"has_end_chapter": r.HasEndChapter, "is_end": r.IsEnd,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user": u, "course": course,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) PathEnrollList(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Title    string   `json:"title"`
		IsEnd    any      `json:"is_end"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	f := repo.OpsListFilter{
		NickName: strings.TrimSpace(body.NickName),
		Title:    strings.TrimSpace(body.Title),
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	if isEnd, ok := asOptionalInt(body.IsEnd); ok {
		f.IsEnd = &isEnd
	}
	rows, total, err := s.ops.ListPathEnrolls(ctx, f)
	if response.WriteErr(c, err) {
		return
	}
	userIDs, pathIDs := []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		pathIDs = append(pathIDs, r.PathID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if response.WriteErr(c, err) {
		return
	}
	paths, err := s.ops.LoadTitleThumbs(ctx, "paths", uniquePositive(pathIDs))
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		list = append(list, gin.H{
			"id": r.ID, "user_id": r.UserID, "path_id": r.PathID,
			"has_end_course": r.HasEndCourse, "is_end": r.IsEnd,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user": users[r.UserID], "path": paths[r.PathID],
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) LearnList(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Title    string   `json:"title"`
		IsEnd    any      `json:"is_end"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	f := repo.OpsListFilter{
		NickName: strings.TrimSpace(body.NickName),
		Title:    strings.TrimSpace(body.Title),
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	if isEnd, ok := asOptionalInt(body.IsEnd); ok {
		f.IsEnd = &isEnd
	}
	rows, total, err := s.ops.ListLearns(ctx, f)
	if response.WriteErr(c, err) {
		return
	}
	userIDs, courseIDs, chapterIDs := []int{}, []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		courseIDs = append(courseIDs, r.CourseID)
		chapterIDs = append(chapterIDs, r.ChapterID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if response.WriteErr(c, err) {
		return
	}
	courses, err := s.ops.LoadTitles(ctx, "courses", uniquePositive(courseIDs))
	if response.WriteErr(c, err) {
		return
	}
	chapters, err := s.ops.LoadTitles(ctx, "course_chapters", uniquePositive(chapterIDs))
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		list = append(list, gin.H{
			"id": r.ID, "user_id": r.UserID, "course_id": r.CourseID, "chapter_id": r.ChapterID,
			"is_end": r.IsEnd, "un_lock": r.UnLock,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user": users[r.UserID], "course": courses[r.CourseID], "chapter": chapters[r.ChapterID],
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) CreditList(c *gin.Context) {
	s.amountList(c, "user_credit_records")
}

func (s *Service) CoinsList(c *gin.Context) {
	s.amountList(c, "user_coins_records")
}

func (s *Service) amountList(c *gin.Context, table string) {
	var body struct {
		NickName string   `json:"nick_name"`
		Field    string   `json:"field"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	rows, total, err := s.ops.ListAmountRecords(ctx, table, repo.OpsAmountFilter{
		NickName: strings.TrimSpace(body.NickName),
		Field:    strings.TrimSpace(body.Field),
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	})
	if response.WriteErr(c, err) {
		return
	}
	userIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		list = append(list, gin.H{
			"id": r.ID, "user_id": r.UserID, "field": r.Field, "sum": r.Sum, "remark": r.Remark,
			"created_at": model.FormatTime(r.CreatedAt), "user": users[r.UserID],
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func asOptionalInt(v any) (int, bool) {
	if !bind.Present(v) {
		return 0, false
	}
	return bind.AsInt(v), true
}

func uniquePositive(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
