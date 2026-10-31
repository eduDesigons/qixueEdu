package adminops

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

type Service struct {
	ops adminrepo.AdminOpsRepository
	db  *gorm.DB
	rdb *redis.Client
}

func New(ops adminrepo.AdminOpsRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{ops: ops, db: db, rdb: rdb}
}

type EnrollListQuery struct {
	NickName string
	Title    string
	IsEnd    *int
	Page     int
	Size     int
}

type AmountListQuery struct {
	NickName string
	Field    string
	Page     int
	Size     int
}

type CourseEnrollView struct {
	ID            int                     `json:"id"`
	UserID        int                     `json:"user_id"`
	CourseID      int                     `json:"course_id"`
	HasEndChapter int                     `json:"has_end_chapter"`
	IsEnd         int                     `json:"is_end"`
	CreatedAt     *string                 `json:"created_at"`
	UpdatedAt     *string                 `json:"updated_at"`
	User          adminrepo.OpsUserBrief  `json:"user"`
	Course        adminrepo.OpsTitleThumb `json:"course"`
}

type PathEnrollView struct {
	ID           int                     `json:"id"`
	UserID       int                     `json:"user_id"`
	PathID       int                     `json:"path_id"`
	HasEndCourse int                     `json:"has_end_course"`
	IsEnd        int                     `json:"is_end"`
	CreatedAt    *string                 `json:"created_at"`
	UpdatedAt    *string                 `json:"updated_at"`
	User         adminrepo.OpsUserBrief  `json:"user"`
	Path         adminrepo.OpsTitleThumb `json:"path"`
}

type LearnView struct {
	ID        int                    `json:"id"`
	UserID    int                    `json:"user_id"`
	CourseID  int                    `json:"course_id"`
	ChapterID int                    `json:"chapter_id"`
	IsEnd     int                    `json:"is_end"`
	UnLock    int                    `json:"un_lock"`
	CreatedAt *string                `json:"created_at"`
	UpdatedAt *string                `json:"updated_at"`
	User      adminrepo.OpsUserBrief `json:"user"`
	Course    adminrepo.OpsTitleOnly `json:"course"`
	Chapter   adminrepo.OpsTitleOnly `json:"chapter"`
}

type AmountView struct {
	ID        int                    `json:"id"`
	UserID    int                    `json:"user_id"`
	Field     string                 `json:"field"`
	Sum       int                    `json:"sum"`
	Remark    string                 `json:"remark"`
	CreatedAt *string                `json:"created_at"`
	User      adminrepo.OpsUserBrief `json:"user"`
}

func (s *Service) CourseEnrollList(ctx context.Context, adminID int, q EnrollListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	f := adminrepo.OpsListFilter{
		NickName: strings.TrimSpace(q.NickName),
		Title:    strings.TrimSpace(q.Title),
		IsEnd:    q.IsEnd,
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	rows, total, err := s.ops.ListCourseEnrolls(ctx, f)
	if err != nil {
		return page.Body{}, err
	}
	userIDs, courseIDs := make([]int, 0, len(rows)), make([]int, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		courseIDs = append(courseIDs, r.CourseID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return page.Body{}, err
	}
	courses, err := s.ops.LoadTitleThumbs(ctx, "courses", uniquePositive(courseIDs))
	if err != nil {
		return page.Body{}, err
	}
	list := make([]CourseEnrollView, 0, len(rows))
	for _, r := range rows {
		list = append(list, CourseEnrollView{
			ID: r.ID, UserID: r.UserID, CourseID: r.CourseID,
			HasEndChapter: r.HasEndChapter, IsEnd: r.IsEnd,
			CreatedAt: model.FormatTime(r.CreatedAt), UpdatedAt: model.FormatTime(r.UpdatedAt),
			User: users[r.UserID], Course: courses[r.CourseID],
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) PathEnrollList(ctx context.Context, adminID int, q EnrollListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	f := adminrepo.OpsListFilter{
		NickName: strings.TrimSpace(q.NickName),
		Title:    strings.TrimSpace(q.Title),
		IsEnd:    q.IsEnd,
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	rows, total, err := s.ops.ListPathEnrolls(ctx, f)
	if err != nil {
		return page.Body{}, err
	}
	userIDs, pathIDs := []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		pathIDs = append(pathIDs, r.PathID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return page.Body{}, err
	}
	paths, err := s.ops.LoadTitleThumbs(ctx, "paths", uniquePositive(pathIDs))
	if err != nil {
		return page.Body{}, err
	}
	list := make([]PathEnrollView, 0, len(rows))
	for _, r := range rows {
		list = append(list, PathEnrollView{
			ID: r.ID, UserID: r.UserID, PathID: r.PathID,
			HasEndCourse: r.HasEndCourse, IsEnd: r.IsEnd,
			CreatedAt: model.FormatTime(r.CreatedAt), UpdatedAt: model.FormatTime(r.UpdatedAt),
			User: users[r.UserID], Path: paths[r.PathID],
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) LearnList(ctx context.Context, adminID int, q EnrollListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	f := adminrepo.OpsListFilter{
		NickName: strings.TrimSpace(q.NickName),
		Title:    strings.TrimSpace(q.Title),
		IsEnd:    q.IsEnd,
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	}
	rows, total, err := s.ops.ListLearns(ctx, f)
	if err != nil {
		return page.Body{}, err
	}
	userIDs, courseIDs, chapterIDs := []int{}, []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		courseIDs = append(courseIDs, r.CourseID)
		chapterIDs = append(chapterIDs, r.ChapterID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return page.Body{}, err
	}
	courses, err := s.ops.LoadTitles(ctx, "courses", uniquePositive(courseIDs))
	if err != nil {
		return page.Body{}, err
	}
	chapters, err := s.ops.LoadTitles(ctx, "course_chapters", uniquePositive(chapterIDs))
	if err != nil {
		return page.Body{}, err
	}
	list := make([]LearnView, 0, len(rows))
	for _, r := range rows {
		list = append(list, LearnView{
			ID: r.ID, UserID: r.UserID, CourseID: r.CourseID, ChapterID: r.ChapterID,
			IsEnd: r.IsEnd, UnLock: r.UnLock,
			CreatedAt: model.FormatTime(r.CreatedAt), UpdatedAt: model.FormatTime(r.UpdatedAt),
			User: users[r.UserID], Course: courses[r.CourseID], Chapter: chapters[r.ChapterID],
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) CreditList(ctx context.Context, adminID int, q AmountListQuery) (page.Body, error) {
	return s.amountList(ctx, adminID, "user_credit_records", q)
}

func (s *Service) CoinsList(ctx context.Context, adminID int, q AmountListQuery) (page.Body, error) {
	return s.amountList(ctx, adminID, "user_coins_records", q)
}

func (s *Service) amountList(ctx context.Context, adminID int, table string, q AmountListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	rows, total, err := s.ops.ListAmountRecords(ctx, table, adminrepo.OpsAmountFilter{
		NickName: strings.TrimSpace(q.NickName),
		Field:    strings.TrimSpace(q.Field),
		CityID:   cityID,
		Limited:  limited,
		Offset:   (current - 1) * size,
		Limit:    size,
	})
	if err != nil {
		return page.Body{}, err
	}
	userIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return page.Body{}, err
	}
	list := make([]AmountView, 0, len(rows))
	for _, r := range rows {
		list = append(list, AmountView{
			ID: r.ID, UserID: r.UserID, Field: r.Field, Sum: r.Sum, Remark: r.Remark,
			CreatedAt: model.FormatTime(r.CreatedAt), User: users[r.UserID],
		})
	}
	return page.New(current, size, total, list), nil
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
