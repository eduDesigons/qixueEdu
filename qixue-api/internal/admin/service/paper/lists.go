package adminpaper

import (
	"context"
	"strings"
	"time"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

type userPaperRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	PaperID       int        `gorm:"column:paper_id"`
	Score         int        `gorm:"column:score"`
	IsPass        int        `gorm:"column:is_pass"`
	IsPeopleCheck int        `gorm:"column:is_people_check"`
	SubmitCount   int        `gorm:"column:submit_count"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
}

type userPaperView struct {
	ID            int            `json:"id"`
	UserID        int            `json:"user_id"`
	PaperID       int            `json:"paper_id"`
	Score         int            `json:"score"`
	IsPass        int            `json:"is_pass"`
	IsPeopleCheck int            `json:"is_people_check"`
	SubmitCount   int            `json:"submit_count"`
	CreatedAt     *string        `json:"created_at"`
	User          map[string]any `json:"user"`
	Paper         map[string]any `json:"paper"`
}

type UserPaperListQuery struct {
	CourseID      int
	ChapterID     int
	Title         string
	IsPass        any
	IsPeopleCheck any
	NickName      string
	Sort          string
	Order         string
	Page          int
	Size          int
}

type PaperLogsQuery struct {
	PaperID  int
	NickName string
	Page     int
	Size     int
}

func (s *Service) GetExamList(ctx context.Context, adminID int, q UserPaperListQuery) (page.Body, error) {
	return s.listUserPapers(ctx, adminID, "course", true, q)
}

func (s *Service) GetAssessmentList(ctx context.Context, adminID int, q UserPaperListQuery) (page.Body, error) {
	return s.listUserPapers(ctx, adminID, "chapter", false, q)
}

func (s *Service) listUserPapers(ctx context.Context, adminID int, paperType string, exam bool, q UserPaperListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	sort := q.Sort
	if sort == "" {
		sort = "id"
	}
	allowed := map[string]bool{"id": true, "score": true, "created_at": true, "is_pass": true, "submit_count": true}
	if !allowed[sort] {
		sort = "id"
	}
	order := "DESC"
	if strings.EqualFold(q.Order, "asc") {
		order = "ASC"
	}
	f := adminrepo.UserPaperListFilter{
		PaperType: paperType, Exam: exam,
		CourseID: q.CourseID, ChapterID: q.ChapterID,
		Title: strings.TrimSpace(q.Title), NickName: strings.TrimSpace(q.NickName),
		CityID: cityID, Limited: limited, Sort: sort, Order: order,
		Offset: (current - 1) * size, Limit: size,
	}
	if q.IsPass != nil && q.IsPass != "" {
		n := asInt(q.IsPass)
		f.IsPass = &n
	}
	if q.IsPeopleCheck != nil && q.IsPeopleCheck != "" {
		n := asInt(q.IsPeopleCheck)
		f.IsPeopleCheck = &n
	}
	rows, total, err := s.papers.ListUserPapers(ctx, f)
	if err != nil {
		return page.Body{}, err
	}
	list, err := s.decorateUserPapers(ctx, toUserPaperRows(rows), !exam)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetPaperLogs(ctx context.Context, adminID int, q PaperLogsQuery) (page.Body, error) {
	if q.PaperID <= 0 {
		return page.Body{}, biz("缺少试卷")
	}
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	rows, total, err := s.papers.ListPaperLogs(ctx, adminrepo.UserPaperLogFilter{
		PaperID: q.PaperID, NickName: strings.TrimSpace(q.NickName),
		CityID: cityID, Limited: limited,
		Offset: (current - 1) * size, Limit: size,
	})
	if err != nil {
		return page.Body{}, err
	}
	list, err := s.decorateUserPapers(ctx, toUserPaperRows(rows), true)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func toUserPaperRows(rows []adminrepo.AdminUserPaperRow) []userPaperRow {
	out := make([]userPaperRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, userPaperRow{
			ID: row.ID, UserID: row.UserID, PaperID: row.PaperID, Score: row.Score,
			IsPass: row.IsPass, IsPeopleCheck: row.IsPeopleCheck, SubmitCount: row.SubmitCount,
			CreatedAt: row.CreatedAt,
		})
	}
	return out
}

func (s *Service) decorateUserPapers(ctx context.Context, rows []userPaperRow, withChapter bool) ([]userPaperView, error) {
	if len(rows) == 0 {
		return []userPaperView{}, nil
	}
	userIDs := map[int]struct{}{}
	paperIDs := map[int]struct{}{}
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		paperIDs[row.PaperID] = struct{}{}
	}
	users := map[int]map[string]any{}
	ids := mapKeys(userIDs)
	if len(ids) > 0 {
		userRows, err := s.papers.ListUsersBrief(ctx, ids)
		if err != nil {
			return nil, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{"id": u.ID, "nick_name": u.NickName, "avatar": u.Avatar, "credit": u.Credit}
		}
	}
	pids := mapKeys(paperIDs)
	papers, err := s.papers.ListPapersBrief(ctx, pids)
	if err != nil {
		return nil, err
	}
	courseIDs := map[int]struct{}{}
	chapterIDs := map[int]struct{}{}
	paperMap := map[int]adminrepo.PaperBrief{}
	for _, p := range papers {
		paperMap[p.ID] = p
		if p.CourseID != nil && *p.CourseID > 0 {
			courseIDs[*p.CourseID] = struct{}{}
		}
		if withChapter && p.ChapterID != nil && *p.ChapterID > 0 {
			chapterIDs[*p.ChapterID] = struct{}{}
		}
	}
	courses := map[int]map[string]any{}
	if cids := mapKeys(courseIDs); len(cids) > 0 {
		courseRows, err := s.papers.ListCoursesThumb(ctx, cids)
		if err != nil {
			return nil, err
		}
		for _, row := range courseRows {
			courses[row.ID] = map[string]any{"id": row.ID, "title": row.Title, "thumb": row.Thumb}
		}
	}
	chapters := map[int]map[string]any{}
	if chids := mapKeys(chapterIDs); len(chids) > 0 {
		chapterRows, err := s.papers.ListChaptersTitle(ctx, chids)
		if err != nil {
			return nil, err
		}
		for _, row := range chapterRows {
			chapters[row.ID] = map[string]any{"id": row.ID, "title": row.Title}
		}
	}
	list := make([]userPaperView, 0, len(rows))
	for _, row := range rows {
		p := paperMap[row.PaperID]
		var course any
		if p.CourseID != nil {
			course = courses[*p.CourseID]
		}
		paperView := map[string]any{"id": p.ID, "course_id": p.CourseID, "chapter_id": p.ChapterID, "course": course}
		if withChapter && p.ChapterID != nil {
			paperView["chapter"] = chapters[*p.ChapterID]
		}
		list = append(list, userPaperView{
			ID: row.ID, UserID: row.UserID, PaperID: row.PaperID, Score: row.Score,
			IsPass: row.IsPass, IsPeopleCheck: row.IsPeopleCheck, SubmitCount: row.SubmitCount,
			CreatedAt: model.FormatTime(row.CreatedAt), User: users[row.UserID], Paper: paperView,
		})
	}
	return list, nil
}

func mapKeys(set map[int]struct{}) []int {
	out := make([]int, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}
