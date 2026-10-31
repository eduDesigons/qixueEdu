package course

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/bg"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/catalog"
	"qixue-api/internal/service/enroll"
	webapprepo "qixue-api/internal/webapp/repo"
)

var (
	ErrNotFound        = errors.New("课程不存在")
	ErrChapterNotFound = errors.New("章节不存在")
)

type Service struct {
	courses  repo.CourseRepository
	users    repo.UserRepository
	db       *gorm.DB
	rdb      *redis.Client
	learnEnd func(ctx context.Context, uid, chapterID int) (*webapprepo.UserChapterFullRow, error)
}

func New(courses repo.CourseRepository, users repo.UserRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{courses: courses, users: users, db: db, rdb: rdb}
}

func (s *Service) SetLearnEnd(fn func(ctx context.Context, uid, chapterID int) (*webapprepo.UserChapterFullRow, error)) {
	s.learnEnd = fn
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

	var (
		userCourse   *repo.UserCourseStatus
		likeCount    int64
		attaches     []repo.CourseAttachRow
		chapters     []repo.ChapterListRow
		progresses   []repo.ChapterProgressRow
		credit       int
		pending      int64
		myCert       *repo.UserCertBrief
		commentTotal int64
	)
	jobs := []func() error{
		func() error {
			var callErr error
			userCourse, callErr = s.courses.FindUserCourseStatus(ctx, uid, id)
			return async.SkipNotFound(callErr)
		},
		func() error {
			likeCount, _ = s.courses.CountCourseLike(ctx, uid, id)
			return nil
		},
		func() error {
			attaches, _ = s.courses.ListCourseAttaches(ctx, id)
			return nil
		},
		func() error {
			var callErr error
			chapters, callErr = s.courses.ListChaptersWithMedia(ctx, id)
			return callErr
		},
		func() error {
			progresses, _ = s.courses.ListUserChapterProgress(ctx, uid, id)
			return nil
		},
		func() error {
			credit, _, _ = s.users.CreditCoins(ctx, uid)
			return nil
		},
		func() error {
			pending, _ = s.courses.CountPendingCourseApply(ctx, uid, id)
			return nil
		},
		func() error {
			return s.db.WithContext(ctx).Model(&model.CourseComment{}).
				Where("course_id = ? AND pid = 0", id).Count(&commentTotal).Error
		},
	}
	if course.CertID != nil && *course.CertID > 0 {
		jobs = append(jobs, func() error {
			if cert, err := s.courses.FindValidCourseCert(ctx, uid, id); err == nil {
				myCert = cert
			}
			return nil
		})
	}
	if err := async.Run(jobs...); err != nil {
		return nil, err
	}
	isJoin := userCourse != nil
	canJoin := true
	isEnd := false
	applyPending := false
	if !isJoin {
		canJoin = credit >= course.MinNeedCredit
		applyPending = pending > 0
	} else {
		isEnd = userCourse.IsEnd == 1
		if err := enroll.UnlockReadyChapters(ctx, s.db, uid, id); err == nil {
			progresses, _ = s.courses.ListUserChapterProgress(ctx, uid, id)
		}
	}

	showPaper := false
	if isJoin && !isEnd && userCourse.HasEndChapter == course.ChapterCount {
		showPaper = true
	}
	if !refresh {
		bg.Go(ctx, func(ctx context.Context) {
			_ = s.courses.IncrCourseView(ctx, id)
		})
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
	if isJoin && myCert != nil {
		info["my_cert"] = myCert
	}

	chapterList := chapterViews(chapters, progresses)
	var pendingPaperID any
	if isJoin {
		pendingPaperID = attachPendingPapers(ctx, s.courses, uid, chapterList, course.PaperID)
	}

	return map[string]any{
		"info": info, "isJoin": isJoin, "isEnd": isEnd, "canJoin": canJoin,
		"applyPending": applyPending, "isCollect": likeCount > 0, "showPaper": showPaper,
		"pendingPaperId": pendingPaperID, "attachList": attaches, "chapters": chapterList,
		"commentTotal": commentTotal,
	}, nil
}

func (s *Service) ChapterList(ctx context.Context, courseID, uid int) ([]map[string]any, error) {
	var (
		rows       []repo.ChapterListRow
		progresses []repo.ChapterProgressRow
	)
	if err := async.Run(
		func() error {
			var callErr error
			rows, callErr = s.courses.ListChaptersWithMedia(ctx, courseID)
			return callErr
		},
		func() error {
			progresses, _ = s.courses.ListUserChapterProgress(ctx, uid, courseID)
			return nil
		},
	); err != nil {
		return nil, err
	}
	if uid > 0 {
		if syncErr := enroll.UnlockReadyChapters(ctx, s.db, uid, courseID); syncErr == nil {
			progresses, _ = s.courses.ListUserChapterProgress(ctx, uid, courseID)
		}
	}
	list := chapterViews(rows, progresses)
	if uid > 0 {
		attachPendingPapers(ctx, s.courses, uid, list, nil)
	}
	return list, nil
}

func chapterViews(rows []repo.ChapterListRow, progresses []repo.ChapterProgressRow) []map[string]any {
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
			"is_people_check": 0, "user_paper_id": nil, "media": nil,
		}
		if r.MediaID != nil {
			item["media"] = map[string]any{"id": *r.MediaID, "duration": r.MDur, "type": r.MType}
		}
		out = append(out, item)
	}
	return out
}

func (s *Service) ChapterDetail(ctx context.Context, id, uid int, refresh, prefetch bool) (map[string]any, error) {
	chapter, err := s.courses.FindChapterDetail(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChapterNotFound
	}
	if err != nil {
		return nil, err
	}

	var (
		userChapter  *repo.CatalogUserChapterRow
		likeCount    int64
		prevRow      *repo.ChapterNeighborRow
		nextRow      *repo.ChapterNeighborRow
		latest       *repo.LatestUserPaper
		commentTotal int64
	)
	paperID := 0
	if chapter.PaperID != nil {
		paperID = *chapter.PaperID
	}
	jobs := []func() error{
		func() error {
			var callErr error
			userChapter, callErr = s.courses.FindUserChapter(ctx, uid, id)
			return async.SkipNotFound(callErr)
		},
		func() error {
			likeCount, _ = s.courses.CountChapterLike(ctx, uid, id)
			return nil
		},
		func() error {
			var callErr error
			prevRow, callErr = s.courses.FindPrevChapter(ctx, chapter.CourseID, chapter.ID, chapter.Sort)
			return async.SkipNotFound(callErr)
		},
		func() error {
			var callErr error
			nextRow, callErr = s.courses.FindNextChapter(ctx, chapter.CourseID, chapter.ID, chapter.Sort)
			return async.SkipNotFound(callErr)
		},
		func() error {
			return s.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
				Where("chapter_id = ? AND pid = 0", id).Count(&commentTotal).Error
		},
	}
	if paperID > 0 {
		jobs = append(jobs, func() error {
			latest, _ = s.courses.FindLatestUserPaper(ctx, uid, paperID)
			return nil
		})
	}
	if err := async.Run(jobs...); err != nil {
		return nil, err
	}
	hasUser := userChapter != nil
	if !refresh {
		bg.Go(ctx, func(ctx context.Context) {
			_ = s.courses.IncrChapterView(ctx, id)
		})
	}

	prev, next := s.neighborViews(ctx, uid, prevRow, nextRow)

	if hasUser && userChapter.UnLock != 1 {
		isFirst := prev == nil
		prevEnded := prev != nil && asInt(prev["is_end"]) == 1
		if isFirst || prevEnded {
			_ = s.courses.UnlockUserChapter(ctx, uid, id)
			userChapter.UnLock = 1
		}
	}

	if hasUser && !prefetch && userChapter.UnLock == 1 && s.learnEnd != nil {
		if row, err := s.learnEnd(ctx, uid, id); err == nil && row != nil {
			userChapter.HasEndChapter = row.HasEndChapter
			userChapter.IsEnd = row.IsEnd
		}
	}

	var userPassPaperID any
	isPeopleCheck := 0
	if paperID > 0 && hasUser && latest != nil && latest.ID > 0 {
		if latest.IsPass == 1 || latest.IsPeopleCheck == 1 {
			userPassPaperID = latest.ID
		}
		if latest.IsPeopleCheck == 1 && latest.IsPass != 1 {
			isPeopleCheck = 1
		}
	}

	info := map[string]any{
		"id": chapter.ID, "course_id": chapter.CourseID, "title": chapter.Title,
		"type": chapter.Type, "sort": chapter.Sort, "content": chapter.Content,
		"media_url": chapter.MediaURL, "media_id": chapter.MediaID,
		"paper_id": chapter.PaperID,
		"course": nil, "media": nil,
	}
	if chapter.CTitle != nil {
		thumb := ""
		if chapter.CThumb != nil {
			thumb = *chapter.CThumb
		}
		info["course"] = map[string]any{"id": chapter.CourseID, "title": *chapter.CTitle, "thumb": thumb}
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
		"userPassPaperId": userPassPaperID, "isPeopleCheck": isPeopleCheck,
		"prevChapter": prev, "nextChapter": next,
		"commentTotal": commentTotal,
	}, nil
}

func (s *Service) neighborViews(ctx context.Context, uid int, prevRow, nextRow *repo.ChapterNeighborRow) (prev, next map[string]any) {
	ids := make([]int, 0, 2)
	if prevRow != nil {
		ids = append(ids, prevRow.ID)
	}
	if nextRow != nil {
		ids = append(ids, nextRow.ID)
	}
	prog := map[int]struct{ UnLock, IsEnd int }{}
	if len(ids) > 0 {
		rows, _ := s.courses.ListNeighborProgress(ctx, uid, ids)
		for _, r := range rows {
			prog[r.ChapterID] = struct{ UnLock, IsEnd int }{r.UnLock, r.IsEnd}
		}
	}
	format := func(r *repo.ChapterNeighborRow) map[string]any {
		if r == nil {
			return nil
		}
		p := prog[r.ID]
		return map[string]any{
			"id": r.ID, "title": r.Title, "type": r.Type, "sort": r.Sort,
			"un_lock": p.UnLock, "is_end": p.IsEnd,
		}
	}
	return format(prevRow), format(nextRow)
}

func asInt(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case *int:
		if t != nil {
			return *t
		}
	}
	return 0
}

func attachPendingPapers(ctx context.Context, courses repo.CourseRepository, uid int, chapters []map[string]any, coursePaperID *int) any {
	ids := make([]int, 0, len(chapters)+1)
	seen := map[int]struct{}{}
	add := func(id int) {
		if id <= 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	for _, ch := range chapters {
		add(asInt(ch["paper_id"]))
	}
	if coursePaperID != nil {
		add(*coursePaperID)
	}
	rows, err := courses.ListPendingUserPapers(ctx, uid, ids)
	if err != nil || len(rows) == 0 {
		return nil
	}
	latestByPaper := map[int]int{}
	for _, r := range rows {
		if _, ok := latestByPaper[r.PaperID]; !ok {
			latestByPaper[r.PaperID] = r.ID
		}
	}
	for _, ch := range chapters {
		pid := asInt(ch["paper_id"])
		if uidPaper, ok := latestByPaper[pid]; ok {
			ch["is_people_check"] = 1
			ch["user_paper_id"] = uidPaper
		}
	}
	if coursePaperID != nil {
		if id, ok := latestByPaper[*coursePaperID]; ok {
			return id
		}
	}
	return nil
}
