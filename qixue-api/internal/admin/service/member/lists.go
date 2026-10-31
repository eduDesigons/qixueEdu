package adminmember

import (
	"context"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/page"
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

type UserSubListQuery struct {
	UserID int
	Title  string
	Page   int
	Size   int
}

type RecordListQuery struct {
	UserID int
	Page   int
	Size   int
}

type PaperListQuery struct {
	UserID        int
	Type          string
	IsPass        any
	IsPeopleCheck any
	Page          int
	Size          int
}

func (s *Service) GetUserCourseList(ctx context.Context, q UserSubListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)

	rows, total, err := s.members.ListUserCourses(ctx, q.UserID, q.Title, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
	}
	courseIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		courseIDs = append(courseIDs, r.CourseID)
	}
	courses, err := s.loadCoursesFull(ctx, courseIDs)
	if err != nil {
		return page.Body{}, err
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
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserPathList(ctx context.Context, q UserSubListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)

	rows, total, err := s.members.ListUserPaths(ctx, q.UserID, q.Title, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
	}
	pathIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		pathIDs = append(pathIDs, r.PathID)
	}
	paths, err := s.loadPathsFull(ctx, pathIDs)
	if err != nil {
		return page.Body{}, err
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
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserCourseChapterList(ctx context.Context, q UserSubListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)

	rows, total, err := s.members.ListUserChapters(ctx, q.UserID, q.Title, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
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
	if err != nil {
		return page.Body{}, err
	}
	chapters, err := s.loadChapters(ctx, chapterIDs)
	if err != nil {
		return page.Body{}, err
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
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserCreditRecordList(ctx context.Context, q RecordListQuery) (page.Body, error) {
	return s.recordList(ctx, "user_credit_records", q)
}

func (s *Service) GetUserCoinsRecordList(ctx context.Context, q RecordListQuery) (page.Body, error) {
	return s.recordList(ctx, "user_coins_records", q)
}

func (s *Service) recordList(ctx context.Context, table string, q RecordListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)
	rows, total, err := s.members.ListUserRecords(ctx, table, q.UserID, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
	}
	list := make([]recordView, 0, len(rows))
	for _, r := range rows {
		list = append(list, recordView{
			ID: r.ID, UserID: r.UserID, Field: r.Field, Sum: r.Sum, Remark: r.Remark,
			CreatedAt: model.FormatTime(r.CreatedAt),
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserPaperList(ctx context.Context, q PaperListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)
	var isPass, isPeopleCheck *int
	if bind.Present(q.IsPass) {
		v := bind.AsInt(q.IsPass)
		isPass = &v
	}
	if bind.Present(q.IsPeopleCheck) {
		v := bind.AsInt(q.IsPeopleCheck)
		isPeopleCheck = &v
	}
	rows, total, err := s.members.ListUserPapers(ctx, q.UserID, q.Type, isPass, isPeopleCheck, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
	}
	paperIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		paperIDs = append(paperIDs, r.PaperID)
	}
	papers, err := s.loadPapers(ctx, paperIDs)
	if err != nil {
		return page.Body{}, err
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
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserCertList(ctx context.Context, q RecordListQuery) (page.Body, error) {
	if q.UserID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(q.Page, q.Size)
	rows, total, err := s.members.ListUserCerts(ctx, q.UserID, (current-1)*size, size)
	if err != nil {
		return page.Body{}, err
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
	if err != nil {
		return page.Body{}, err
	}
	courses, err := s.loadTitles(ctx, "courses", courseIDs, true)
	if err != nil {
		return page.Body{}, err
	}
	certs, err := s.loadCertNames(ctx, certIDs)
	if err != nil {
		return page.Body{}, err
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
	return page.New(current, size, total, list), nil
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
