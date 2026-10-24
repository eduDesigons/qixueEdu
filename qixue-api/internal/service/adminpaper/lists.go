package adminpaper

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
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

func (s *Service) GetExamList(c *gin.Context) {
	s.listUserPapers(c, "course", true)
}

func (s *Service) GetAssessmentList(c *gin.Context) {
	s.listUserPapers(c, "chapter", false)
}

func (s *Service) listUserPapers(c *gin.Context, paperType string, exam bool) {
	var body struct {
		CourseID      bind.Int `json:"course_id"`
		ChapterID     bind.Int `json:"chapter_id"`
		Title         string   `json:"title"`
		IsPass        any      `json:"is_pass"`
		IsPeopleCheck any      `json:"is_people_check"`
		NickName      string   `json:"nick_name"`
		Sort          string   `json:"sort"`
		Order         string   `json:"order"`
		Page          bind.Int `json:"page"`
		Size          bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	sort := body.Sort
	if sort == "" {
		sort = "id"
	}
	allowed := map[string]bool{"id": true, "score": true, "created_at": true, "is_pass": true, "submit_count": true}
	if !allowed[sort] {
		sort = "id"
	}
	order := "DESC"
	if strings.EqualFold(body.Order, "asc") {
		order = "ASC"
	}
	f := repo.UserPaperListFilter{
		PaperType: paperType, Exam: exam,
		CourseID: body.CourseID.Int(), ChapterID: body.ChapterID.Int(),
		Title: strings.TrimSpace(body.Title), NickName: strings.TrimSpace(body.NickName),
		CityID: cityID, Limited: limited, Sort: sort, Order: order,
		Offset: (current - 1) * size, Limit: size,
	}
	if body.IsPass != nil && body.IsPass != "" {
		n := bind.AsInt(body.IsPass)
		f.IsPass = &n
	}
	if body.IsPeopleCheck != nil && body.IsPeopleCheck != "" {
		n := bind.AsInt(body.IsPeopleCheck)
		f.IsPeopleCheck = &n
	}
	rows, total, err := s.papers.ListUserPapers(ctx, f)
	if response.WriteErr(c, err) {
		return
	}
	list, err := s.decorateUserPapers(c, toUserPaperRows(rows), !exam)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) GetPaperLogs(c *gin.Context) {
	var body struct {
		PaperID  bind.Int `json:"paper_id"`
		NickName string   `json:"nick_name"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	paperID := body.PaperID.Int()
	if paperID <= 0 {
		response.Fail(c, 1, "缺少试卷")
		return
	}
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	rows, total, err := s.papers.ListPaperLogs(ctx, repo.UserPaperLogFilter{
		PaperID: paperID, NickName: strings.TrimSpace(body.NickName),
		CityID: cityID, Limited: limited,
		Offset: (current - 1) * size, Limit: size,
	})
	if response.WriteErr(c, err) {
		return
	}
	list, err := s.decorateUserPapers(c, toUserPaperRows(rows), true)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func toUserPaperRows(rows []repo.AdminUserPaperRow) []userPaperRow {
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

func (s *Service) decorateUserPapers(c *gin.Context, rows []userPaperRow, withChapter bool) ([]userPaperView, error) {
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
		userRows, err := s.papers.ListUsersBrief(c.Request.Context(), ids)
		if err != nil {
			return nil, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{"id": u.ID, "nick_name": u.NickName, "avatar": u.Avatar, "credit": u.Credit}
		}
	}
	pids := mapKeys(paperIDs)
	papers, err := s.papers.ListPapersBrief(c.Request.Context(), pids)
	if err != nil {
		return nil, err
	}
	courseIDs := map[int]struct{}{}
	chapterIDs := map[int]struct{}{}
	paperMap := map[int]repo.PaperBrief{}
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
		courseRows, err := s.papers.ListCoursesThumb(c.Request.Context(), cids)
		if err != nil {
			return nil, err
		}
		for _, row := range courseRows {
			courses[row.ID] = map[string]any{"id": row.ID, "title": row.Title, "thumb": row.Thumb}
		}
	}
	chapters := map[int]map[string]any{}
	if chids := mapKeys(chapterIDs); len(chids) > 0 {
		chapterRows, err := s.papers.ListChaptersTitle(c.Request.Context(), chids)
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
