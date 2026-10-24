package adminmember

import (
	"context"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
)

type categoryMini struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type courseNested struct {
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	Thumb         string        `json:"thumb"`
	Type          int           `json:"type,omitempty"`
	Desc          string        `json:"desc,omitempty"`
	IsShow        int           `json:"is_show,omitempty"`
	CategoryID    int           `json:"category_id"`
	Category      *categoryMini `json:"category"`
	ChaptersCount int           `json:"chapters_count,omitempty"`
}

type pathNested struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Thumb        string        `json:"thumb"`
	Desc         string        `json:"desc"`
	IsShow       int           `json:"is_show"`
	CategoryID   int           `json:"category_id"`
	Category     *categoryMini `json:"category"`
	CoursesCount int           `json:"courses_count"`
}

type chapterMini struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type paperCourseMini struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}

type paperNested struct {
	ID        int              `json:"id"`
	Type      string           `json:"type"`
	Score     int              `json:"score"`
	PassScore int              `json:"pass_score"`
	CourseID  *int             `json:"course_id"`
	ChapterID *int             `json:"chapter_id"`
	Course    *paperCourseMini `json:"course"`
	Chapter   *chapterMini     `json:"chapter"`
}

type userCourseView struct {
	ID            int           `json:"id"`
	CourseID      int           `json:"course_id"`
	UserID        int           `json:"user_id"`
	CreatedAt     *string       `json:"created_at"`
	HasEndChapter int           `json:"has_end_chapter"`
	IsEnd         int           `json:"is_end"`
	Course        *courseNested `json:"course"`
}

type userPathView struct {
	ID           int         `json:"id"`
	PathID       int         `json:"path_id"`
	CreatedAt    *string     `json:"created_at"`
	HasEndCourse int         `json:"has_end_course"`
	IsEnd        int         `json:"is_end"`
	Path         *pathNested `json:"path"`
}

type userChapterView struct {
	ID        int           `json:"id"`
	CourseID  *int          `json:"course_id"`
	ChapterID int           `json:"chapter_id"`
	UserID    int           `json:"user_id"`
	CreatedAt *string       `json:"created_at"`
	IsEnd     int           `json:"is_end"`
	UnLock    int           `json:"un_lock"`
	Course    *courseNested `json:"course"`
	Chapter   *chapterMini  `json:"chapter"`
}

type recordView struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Field     string  `json:"field"`
	Sum       int     `json:"sum"`
	Remark    string  `json:"remark"`
	CreatedAt *string `json:"created_at"`
}

type userPaperView struct {
	ID            int          `json:"id"`
	UserID        int          `json:"user_id"`
	PaperID       int          `json:"paper_id"`
	Score         int          `json:"score"`
	SubmitCount   int          `json:"submit_count"`
	CreatedAt     *string      `json:"created_at"`
	IsPass        int          `json:"is_pass"`
	IsPeopleCheck int          `json:"is_people_check"`
	Paper         *paperNested `json:"paper"`
}

type titleMini struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type certNameMini struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type userCertView struct {
	ID        int           `json:"id"`
	UserID    int           `json:"user_id"`
	CertID    int           `json:"cert_id"`
	PathID    *int          `json:"path_id"`
	CourseID  int           `json:"course_id"`
	CertNo    string        `json:"cert_no"`
	Title     string        `json:"title"`
	Image     string        `json:"image"`
	IsValid   int           `json:"is_valid"`
	IssuedAt  *string       `json:"issued_at"`
	CreatedAt *string       `json:"created_at"`
	Path      *titleMini    `json:"path"`
	Course    *titleMini    `json:"course"`
	Cert      *certNameMini `json:"cert"`
}

func (s *Service) GetUserCourseList(c *gin.Context) {
	userID := atoi(c.Query("userId"))
	if userID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	title := c.Query("title")
	current, size := page.Clamp(atoi(c.Query("page")), atoi(c.Query("size")))
	ctx := c.Request.Context()

	rows, total, err := s.members.ListUserCourses(ctx, userID, title, (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	courseIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		courseIDs = append(courseIDs, r.CourseID)
	}
	courses, err := s.loadCoursesFull(ctx, courseIDs)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]userCourseView, 0, len(rows))
	for _, r := range rows {
		var course *courseNested
		if v, ok := courses[r.CourseID]; ok {
			cp := v
			course = &cp
		}
		list = append(list, userCourseView{
			ID: r.ID, CourseID: r.CourseID, UserID: r.UserID,
			CreatedAt:     model.FormatTime(r.CreatedAt),
			HasEndChapter: r.HasEndChapter, IsEnd: r.IsEnd, Course: course,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserPathList(c *gin.Context) {
	userID := atoi(c.Query("userId"))
	if userID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	title := c.Query("title")
	current, size := page.Clamp(atoi(c.Query("page")), atoi(c.Query("size")))
	ctx := c.Request.Context()

	rows, total, err := s.members.ListUserPaths(ctx, userID, title, (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	pathIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		pathIDs = append(pathIDs, r.PathID)
	}
	paths, err := s.loadPathsFull(ctx, pathIDs)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]userPathView, 0, len(rows))
	for _, r := range rows {
		var path *pathNested
		if v, ok := paths[r.PathID]; ok {
			cp := v
			path = &cp
		}
		list = append(list, userPathView{
			ID: r.ID, PathID: r.PathID, CreatedAt: model.FormatTime(r.CreatedAt),
			HasEndCourse: r.HasEndCourse, IsEnd: r.IsEnd, Path: path,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserCourseChapterList(c *gin.Context) {
	userID := atoi(c.Query("userId"))
	if userID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	title := c.Query("title")
	current, size := page.Clamp(atoi(c.Query("page")), atoi(c.Query("size")))
	ctx := c.Request.Context()

	rows, total, err := s.members.ListUserChapters(ctx, userID, title, (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	courseIDs := make([]int, 0, len(rows))
	chapterIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		if r.CourseID != nil && *r.CourseID > 0 {
			courseIDs = append(courseIDs, *r.CourseID)
		}
		chapterIDs = append(chapterIDs, r.ChapterID)
	}
	courses, err := s.loadCoursesBrief(ctx, courseIDs)
	if response.WriteErr(c, err) {
		return
	}
	chapters, err := s.loadChapters(ctx, chapterIDs)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]userChapterView, 0, len(rows))
	for _, r := range rows {
		var course *courseNested
		if r.CourseID != nil {
			if v, ok := courses[*r.CourseID]; ok {
				cp := v
				course = &cp
			}
		}
		var chapter *chapterMini
		if v, ok := chapters[r.ChapterID]; ok {
			cp := v
			chapter = &cp
		}
		list = append(list, userChapterView{
			ID: r.ID, CourseID: r.CourseID, ChapterID: r.ChapterID, UserID: r.UserID,
			CreatedAt: model.FormatTime(r.CreatedAt), IsEnd: r.IsEnd, UnLock: r.UnLock,
			Course: course, Chapter: chapter,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserCreditRecordList(c *gin.Context) {
	s.recordList(c, "user_credit_records")
}

func (s *Service) GetUserCoinsRecordList(c *gin.Context) {
	s.recordList(c, "user_coins_records")
}

func (s *Service) recordList(c *gin.Context, table string) {
	var body struct {
		UserID bind.Int `json:"user_id"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.UserID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	rows, total, err := s.members.ListUserRecords(ctx, table, body.UserID.Int(), (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]recordView, 0, len(rows))
	for _, r := range rows {
		list = append(list, recordView{
			ID: r.ID, UserID: r.UserID, Field: r.Field, Sum: r.Sum, Remark: r.Remark,
			CreatedAt: model.FormatTime(r.CreatedAt),
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserPaperList(c *gin.Context) {
	var body struct {
		UserID        bind.Int `json:"user_id"`
		Type          string   `json:"type"`
		IsPass        any      `json:"is_pass"`
		IsPeopleCheck any      `json:"is_people_check"`
		Page          bind.Int `json:"page"`
		Size          bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.UserID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	var isPass, isPeopleCheck *int
	if bind.Present(body.IsPass) {
		v := bind.AsInt(body.IsPass)
		isPass = &v
	}
	if bind.Present(body.IsPeopleCheck) {
		v := bind.AsInt(body.IsPeopleCheck)
		isPeopleCheck = &v
	}
	rows, total, err := s.members.ListUserPapers(ctx, body.UserID.Int(), body.Type, isPass, isPeopleCheck, (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	paperIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		paperIDs = append(paperIDs, r.PaperID)
	}
	papers, err := s.loadPapers(ctx, paperIDs)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]userPaperView, 0, len(rows))
	for _, r := range rows {
		var paper *paperNested
		if v, ok := papers[r.PaperID]; ok {
			cp := v
			paper = &cp
		}
		list = append(list, userPaperView{
			ID: r.ID, UserID: r.UserID, PaperID: r.PaperID, Score: r.Score,
			SubmitCount: r.SubmitCount, CreatedAt: model.FormatTime(r.CreatedAt),
			IsPass: r.IsPass, IsPeopleCheck: r.IsPeopleCheck, Paper: paper,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) GetUserCertList(c *gin.Context) {
	var body struct {
		UserID bind.Int `json:"user_id"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.UserID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	rows, total, err := s.members.ListUserCerts(ctx, body.UserID.Int(), (current-1)*size, size)
	if response.WriteErr(c, err) {
		return
	}
	pathIDs, courseIDs, certIDs := []int{}, []int{}, []int{}
	for _, r := range rows {
		if r.PathID != nil && *r.PathID > 0 {
			pathIDs = append(pathIDs, *r.PathID)
		}
		if r.CourseID > 0 {
			courseIDs = append(courseIDs, r.CourseID)
		}
		if r.CertID > 0 {
			certIDs = append(certIDs, r.CertID)
		}
	}
	paths, err := s.loadTitles(ctx, "paths", pathIDs, true)
	if response.WriteErr(c, err) {
		return
	}
	courses, err := s.loadTitles(ctx, "courses", courseIDs, true)
	if response.WriteErr(c, err) {
		return
	}
	certs, err := s.loadCertNames(ctx, certIDs)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]userCertView, 0, len(rows))
	for _, r := range rows {
		var path *titleMini
		if r.PathID != nil {
			if v, ok := paths[*r.PathID]; ok {
				cp := v
				path = &cp
			}
		}
		var course *titleMini
		if v, ok := courses[r.CourseID]; ok {
			cp := v
			course = &cp
		}
		var cert *certNameMini
		if v, ok := certs[r.CertID]; ok {
			cp := v
			cert = &cp
		}
		list = append(list, userCertView{
			ID: r.ID, UserID: r.UserID, CertID: r.CertID, PathID: r.PathID, CourseID: r.CourseID,
			CertNo: r.CertNo, Title: r.Title, Image: r.Image, IsValid: r.IsValid,
			IssuedAt: model.FormatTime(r.IssuedAt), CreatedAt: model.FormatTime(r.CreatedAt),
			Path: path, Course: course, Cert: cert,
		})
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) loadCoursesFull(ctx context.Context, ids []int) (map[int]courseNested, error) {
	out := map[int]courseNested{}
	rows, err := s.members.ListCoursesFull(ctx, ids)
	if err != nil {
		return nil, err
	}
	catIDs := make([]int, 0, len(rows))
	courseIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		courseIDs = append(courseIDs, r.ID)
		if r.CategoryID > 0 {
			catIDs = append(catIDs, r.CategoryID)
		}
	}
	cats, err := s.loadCourseCategories(ctx, catIDs)
	if err != nil {
		return nil, err
	}
	chapterCounts, err := s.countChapters(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		var cat *categoryMini
		if v, ok := cats[r.CategoryID]; ok {
			cp := v
			cat = &cp
		}
		out[r.ID] = courseNested{
			ID: r.ID, Title: r.Title, Thumb: r.Thumb, Type: r.Type, Desc: r.Desc,
			IsShow: r.IsShow, CategoryID: r.CategoryID, Category: cat,
			ChaptersCount: chapterCounts[r.ID],
		}
	}
	return out, nil
}

func (s *Service) loadCoursesBrief(ctx context.Context, ids []int) (map[int]courseNested, error) {
	out := map[int]courseNested{}
	rows, err := s.members.ListCoursesBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	catIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		if r.CategoryID > 0 {
			catIDs = append(catIDs, r.CategoryID)
		}
	}
	cats, err := s.loadCourseCategories(ctx, catIDs)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		var cat *categoryMini
		if v, ok := cats[r.CategoryID]; ok {
			cp := v
			cat = &cp
		}
		out[r.ID] = courseNested{
			ID: r.ID, Title: r.Title, Thumb: r.Thumb, CategoryID: r.CategoryID, Category: cat,
		}
	}
	return out, nil
}

func (s *Service) loadPathsFull(ctx context.Context, ids []int) (map[int]pathNested, error) {
	out := map[int]pathNested{}
	rows, err := s.members.ListPathsFull(ctx, ids)
	if err != nil {
		return nil, err
	}
	catIDs := make([]int, 0, len(rows))
	pathIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		pathIDs = append(pathIDs, r.ID)
		if r.CategoryID > 0 {
			catIDs = append(catIDs, r.CategoryID)
		}
	}
	cats, err := s.loadPathCategories(ctx, catIDs)
	if err != nil {
		return nil, err
	}
	courseCounts, err := s.countPathCourses(ctx, pathIDs)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		var cat *categoryMini
		if v, ok := cats[r.CategoryID]; ok {
			cp := v
			cat = &cp
		}
		out[r.ID] = pathNested{
			ID: r.ID, Title: r.Title, Thumb: r.Thumb, Desc: r.Desc, IsShow: r.IsShow,
			CategoryID: r.CategoryID, Category: cat, CoursesCount: courseCounts[r.ID],
		}
	}
	return out, nil
}

func (s *Service) loadChapters(ctx context.Context, ids []int) (map[int]chapterMini, error) {
	out := map[int]chapterMini{}
	rows, err := s.members.ListChaptersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = chapterMini{ID: r.ID, Title: r.Title}
	}
	return out, nil
}

func (s *Service) loadPapers(ctx context.Context, ids []int) (map[int]paperNested, error) {
	out := map[int]paperNested{}
	rows, err := s.members.ListPapersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	courseIDs, chapterIDs := []int{}, []int{}
	for _, r := range rows {
		if r.CourseID != nil && *r.CourseID > 0 {
			courseIDs = append(courseIDs, *r.CourseID)
		}
		if r.ChapterID != nil && *r.ChapterID > 0 {
			chapterIDs = append(chapterIDs, *r.ChapterID)
		}
	}
	courses := map[int]paperCourseMini{}
	cRows, err := s.members.ListCoursesThumb(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	for _, r := range cRows {
		courses[r.ID] = paperCourseMini{ID: r.ID, Title: r.Title, Thumb: r.Thumb}
	}
	chapters, err := s.loadChapters(ctx, chapterIDs)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		var course *paperCourseMini
		if r.CourseID != nil {
			if v, ok := courses[*r.CourseID]; ok {
				cp := v
				course = &cp
			}
		}
		var chapter *chapterMini
		if r.ChapterID != nil {
			if v, ok := chapters[*r.ChapterID]; ok {
				cp := v
				chapter = &cp
			}
		}
		out[r.ID] = paperNested{
			ID: r.ID, Type: r.Type, Score: r.Score, PassScore: r.PassScore,
			CourseID: r.CourseID, ChapterID: r.ChapterID, Course: course, Chapter: chapter,
		}
	}
	return out, nil
}

func (s *Service) loadCourseCategories(ctx context.Context, ids []int) (map[int]categoryMini, error) {
	return s.loadNamed(ctx, "course_categories", ids)
}

func (s *Service) loadPathCategories(ctx context.Context, ids []int) (map[int]categoryMini, error) {
	return s.loadNamed(ctx, "path_categories", ids)
}

func (s *Service) loadNamed(ctx context.Context, table string, ids []int) (map[int]categoryMini, error) {
	out := map[int]categoryMini{}
	rows, err := s.members.ListNamedByIDs(ctx, table, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = categoryMini{ID: r.ID, Name: r.Name}
	}
	return out, nil
}

func (s *Service) loadTitles(ctx context.Context, table string, ids []int, softDelete bool) (map[int]titleMini, error) {
	out := map[int]titleMini{}
	rows, err := s.members.ListTitlesByIDs(ctx, table, ids, softDelete)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = titleMini{ID: r.ID, Title: r.Title}
	}
	return out, nil
}

func (s *Service) loadCertNames(ctx context.Context, ids []int) (map[int]certNameMini, error) {
	out := map[int]certNameMini{}
	rows, err := s.members.ListCertNames(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = certNameMini{ID: r.ID, Name: r.Name}
	}
	return out, nil
}

func (s *Service) countChapters(ctx context.Context, courseIDs []int) (map[int]int, error) {
	out := map[int]int{}
	if len(courseIDs) == 0 {
		return out, nil
	}
	type row struct {
		CourseID int `gorm:"column:course_id"`
		Total    int `gorm:"column:total"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT course_id, COUNT(*) AS total FROM course_chapters
		WHERE course_id IN ? GROUP BY course_id
	`, courseIDs).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.CourseID] = r.Total
	}
	return out, nil
}

func (s *Service) countPathCourses(ctx context.Context, pathIDs []int) (map[int]int, error) {
	out := map[int]int{}
	if len(pathIDs) == 0 {
		return out, nil
	}
	type row struct {
		PathID int `gorm:"column:path_id"`
		Total  int `gorm:"column:total"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT pc.path_id, COUNT(*) AS total
		FROM path_courses pc
		INNER JOIN courses c ON c.id = pc.course_id AND c.deleted_at IS NULL
		WHERE pc.path_id IN ?
		GROUP BY pc.path_id
	`, pathIDs).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.PathID] = r.Total
	}
	return out, nil
}
