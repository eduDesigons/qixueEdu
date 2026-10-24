package course

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/catalog"
)

var (
	ErrNotFound        = errors.New("课程不存在")
	ErrChapterNotFound = errors.New("章节不存在")
)

type Service struct {
	courses repo.CourseRepository
	users   repo.UserRepository
	db      *gorm.DB
	rdb     *redis.Client
}

func New(courses repo.CourseRepository, users repo.UserRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{courses: courses, users: users, db: db, rdb: rdb}
}

func (s *Service) Categories(ctx context.Context) ([]catalog.CourseNode, error) {
	return catalog.CourseTree(ctx, s.db, s.rdb)
}

type ListFilter struct {
	CateID1  int
	CateID2  int
	Type     int
	SortType string
	Title    string
	Page     int
	Limit    int
}

func (s *Service) List(ctx context.Context, f ListFilter) (page.MoreBody, error) {
	if f.SortType == "" {
		f.SortType = "hot"
	}
	isDefaultHot := f.Page == 1 &&
		f.CateID1 == 0 && f.CateID2 == 0 && f.Type == 0 &&
		strings.TrimSpace(f.Title) == "" && f.SortType == "hot"

	if isDefaultHot {
		var cached page.MoreBody
		if cache.GetJSON(ctx, s.rdb, "course_list_hot_p1", &cached) && cached.List != nil {
			return cached, nil
		}
	}

	tree, err := catalog.CourseTree(ctx, s.db, s.rdb)
	if err != nil {
		return page.MoreBody{}, err
	}
	leafIDs := []int{}
	childrenOf := map[int][]int{}
	for _, parent := range tree {
		kids := make([]int, 0, len(parent.Children))
		for _, child := range parent.Children {
			leafIDs = append(leafIDs, child.ID)
			kids = append(kids, child.ID)
		}
		childrenOf[parent.ID] = kids
	}

	rf := repo.CourseListFilter{
		Type: f.Type, Title: f.Title, SortType: f.SortType, Page: f.Page, Limit: f.Limit,
	}
	switch {
	case f.CateID2 > 0:
		ok := false
		for _, id := range leafIDs {
			if id == f.CateID2 {
				ok = true
				break
			}
		}
		if !ok {
			rf.ForceNone = true
		} else {
			rf.CategoryIDs = []int{f.CateID2}
		}
	case f.CateID1 > 0:
		kids := childrenOf[f.CateID1]
		if len(kids) == 0 {
			rf.ForceNone = true
		} else {
			rf.CategoryIDs = kids
		}
	default:
		if len(leafIDs) == 0 {
			rf.ForceNone = true
		} else {
			rf.CategoryIDs = leafIDs
		}
	}

	total, rows, err := s.courses.ListShownCourses(ctx, rf)
	if err != nil {
		return page.MoreBody{}, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": r.Desc,
			"type": r.Type, "credit": r.Credit, "min_need_credit": r.MinNeedCredit,
			"category_id": r.CategoryID, "category": nil,
		}
		if r.CateName != nil {
			item["category"] = map[string]any{"id": r.CategoryID, "name": *r.CateName}
		}
		list = append(list, item)
	}
	payload := page.NewMore(f.Page, f.Limit, total, list)
	if isDefaultHot {
		cache.SetJSON(ctx, s.rdb, "course_list_hot_p1", payload, 20*time.Second)
	}
	return payload, nil
}

func (s *Service) Detail(ctx context.Context, id, uid int, refresh bool) (map[string]any, error) {
	course, err := s.courses.FindCourseDetail(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	userCourse, err := s.courses.FindUserCourseStatus(ctx, uid, id)
	isJoin := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	canJoin := true
	isEnd := false
	applyPending := false
	if !isJoin {
		credit, _, _ := s.users.CreditCoins(ctx, uid)
		canJoin = credit >= course.MinNeedCredit
		n, _ := s.courses.CountPendingCourseApply(ctx, uid, id)
		applyPending = n > 0
	} else {
		isEnd = userCourse.IsEnd == 1
	}

	likeCount, _ := s.courses.CountCourseLike(ctx, uid, id)
	attaches, _ := s.courses.ListCourseAttaches(ctx, id)

	showPaper := false
	if isJoin && !isEnd && userCourse.HasEndChapter == course.ChapterCount {
		showPaper = true
	}
	if !refresh {
		_ = s.courses.IncrCourseView(ctx, id)
	}

	typeNames := map[int]string{1: "视频", 2: "音频", 3: "图文"}
	info := map[string]any{
		"id": course.ID, "title": course.Title, "thumb": course.Thumb, "desc": course.Desc,
		"type": course.Type, "credit": course.Credit, "content": course.Content,
		"paper_id": course.PaperID, "cert_id": course.CertID, "category_id": course.CategoryID,
		"min_need_credit": course.MinNeedCredit, "view_count": course.ViewCount,
		"chapter_count": course.ChapterCount, "type_name": typeNames[course.Type],
		"category": nil, "cert": nil, "my_cert": nil,
	}
	if course.CateName != nil {
		info["category"] = map[string]any{"id": course.CategoryID, "name": *course.CateName}
	}
	if course.CertID != nil && course.CertName != nil {
		info["cert"] = map[string]any{"id": *course.CertID, "name": *course.CertName}
	}
	if isJoin && course.CertID != nil && *course.CertID > 0 {
		if myCert, err := s.courses.FindValidCourseCert(ctx, uid, id); err == nil {
			info["my_cert"] = myCert
		}
	}

	chapters, err := s.ChapterList(ctx, id, uid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"info": info, "isJoin": isJoin, "isEnd": isEnd, "canJoin": canJoin,
		"applyPending": applyPending, "isCollect": likeCount > 0, "showPaper": showPaper,
		"attachList": attaches, "chapters": chapters,
	}, nil
}

func (s *Service) ChapterList(ctx context.Context, courseID, uid int) ([]map[string]any, error) {
	rows, err := s.courses.ListChaptersWithMedia(ctx, courseID)
	if err != nil {
		return nil, err
	}
	progresses, _ := s.courses.ListUserChapterProgress(ctx, uid, courseID)
	byChapter := map[int]repo.UserChapterProgress{}
	for _, p := range progresses {
		byChapter[p.ChapterID] = p
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		p := byChapter[r.ID]
		item := map[string]any{
			"id": r.ID, "title": r.Title, "type": r.Type, "sort": r.Sort,
			"media_id": r.MediaID, "paper_id": r.PaperID,
			"is_end": p.IsEnd, "un_lock": p.UnLock, "has_end_chapter": p.HasEndChapter,
			"media": nil,
		}
		if r.MediaID != nil {
			item["media"] = map[string]any{"id": *r.MediaID, "duration": r.MDur, "type": r.MType}
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) ChapterDetail(ctx context.Context, id, uid int, refresh bool) (map[string]any, error) {
	chapter, err := s.courses.FindChapterDetail(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChapterNotFound
	}
	if err != nil {
		return nil, err
	}

	userChapter, err := s.courses.FindUserChapter(ctx, uid, id)
	hasUser := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	likeCount, _ := s.courses.CountChapterLike(ctx, uid, id)
	if !refresh {
		_ = s.courses.IncrChapterView(ctx, id)
	}

	var userPassPaperID any
	if chapter.PaperID != nil && *chapter.PaperID > 0 && hasUser && userChapter.IsEnd == 1 {
		if pid, _ := s.courses.FindPassedPaperID(ctx, uid, *chapter.PaperID); pid > 0 {
			userPassPaperID = pid
		}
	}

	prev, next, err := s.resolveNeighbors(ctx, chapter.CourseID, chapter.ID, chapter.Sort, uid)
	if err != nil {
		return nil, err
	}

	if hasUser && userChapter.UnLock != 1 {
		isFirst := prev == nil
		prevEnded := prev != nil && asInt(prev["is_end"]) == 1
		if isFirst || prevEnded {
			days := 30
			if chapter.LimitDay != nil && *chapter.LimitDay > 0 {
				days = *chapter.LimitDay
			}
			expire := time.Now().AddDate(0, 0, days)
			_ = s.courses.UnlockUserChapter(ctx, uid, id, expire)
			userChapter.ExpireTime = &expire
			userChapter.UnLock = 1
		}
	}

	info := map[string]any{
		"id": chapter.ID, "course_id": chapter.CourseID, "title": chapter.Title,
		"type": chapter.Type, "sort": chapter.Sort, "content": chapter.Content,
		"media_url": chapter.MediaURL, "media_id": chapter.MediaID,
		"limit_day": chapter.LimitDay, "paper_id": chapter.PaperID,
		"course": nil, "media": nil,
	}
	if chapter.CTitle != nil {
		info["course"] = map[string]any{"id": chapter.CourseID, "title": *chapter.CTitle, "thumb": deref(chapter.CThumb)}
	}
	if chapter.MediaID != nil {
		info["media"] = map[string]any{
			"id": *chapter.MediaID, "duration": chapter.MDur, "type": chapter.MType,
			"size": chapter.MSize, "thumb": chapter.MThumb,
		}
	}
	var userOut any
	if hasUser {
		userOut = userChapter
	}
	return map[string]any{
		"info": info, "userChapter": userOut, "isCollect": likeCount > 0,
		"userPassPaperId": userPassPaperID, "prevChapter": prev, "nextChapter": next,
	}, nil
}

func (s *Service) resolveNeighbors(ctx context.Context, courseID, chapterID, sort, uid int) (prev, next map[string]any, err error) {
	prevRow, err := s.courses.FindPrevChapter(ctx, courseID, chapterID, sort)
	hasPrev := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}
	nextRow, err := s.courses.FindNextChapter(ctx, courseID, chapterID, sort)
	hasNext := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}
	ids := []int{}
	if hasPrev {
		ids = append(ids, prevRow.ID)
	}
	if hasNext {
		ids = append(ids, nextRow.ID)
	}
	prog := map[int]struct{ UnLock, IsEnd int }{}
	if len(ids) > 0 {
		rows, _ := s.courses.ListNeighborProgress(ctx, uid, ids)
		for _, r := range rows {
			prog[r.ChapterID] = struct{ UnLock, IsEnd int }{r.UnLock, r.IsEnd}
		}
	}
	format := func(r *repo.ChapterNeighborRow, ok bool) map[string]any {
		if !ok || r == nil {
			return nil
		}
		p := prog[r.ID]
		return map[string]any{
			"id": r.ID, "title": r.Title, "type": r.Type, "sort": r.Sort,
			"un_lock": p.UnLock, "is_end": p.IsEnd,
		}
	}
	return format(prevRow, hasPrev), format(nextRow, hasNext), nil
}

func asInt(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	default:
		return 0
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
