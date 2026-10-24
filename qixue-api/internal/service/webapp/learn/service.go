package learn

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/cert"
	"qixue-api/internal/service/coin"
	"qixue-api/internal/service/enroll"
)

var (
	ErrAlreadyJoined      = errors.New("已经加入过了")
	ErrCourseNotFound     = errors.New("课程不存在")
	ErrNoChapters         = errors.New("课程暂无课节，无法加入")
	ErrCreditShort        = errors.New("学分不足")
	ErrJoinFailed         = errors.New("加入失败，请重试")
	ErrReasonRequired     = errors.New("请填写申请理由")
	ErrApplyNoChapters    = errors.New("课程暂无课节，无法申请")
	ErrAlreadyApplied     = errors.New("已提交过申请")
	ErrChapterNotFound    = errors.New("章节不存在")
	ErrNeedJoinCourse     = errors.New("请先加入课程")
	ErrAlreadyEnded       = errors.New("已经学习结束了")
	ErrAlreadyUnlocked    = errors.New("已经解锁过了")
	ErrNeedPrevChapter    = errors.New("请先完成上一章节")
	ErrNeedUnlock         = errors.New("请先解锁")
	ErrCourseAlreadyEnded = errors.New("课程已经结束")
	ErrCourseNeedExam     = errors.New("课程需要考试")
)

type Service struct {
	learn  repo.LearnRepository
	users  repo.UserRepository
	db     *gorm.DB
	rdb    *redis.Client
	issuer *cert.Issuer
}

func New(learn repo.LearnRepository, users repo.UserRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{learn: learn, users: users, db: db, rdb: rdb, issuer: cert.NewIssuer(repo.NewCert(db))}
}

func (s *Service) MyCourse(ctx context.Context, uid, typ, pageNum, limit int) (map[string]any, error) {
	if typ == 0 {
		typ = 1
	}
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}
	cacheKey := "my_course_" + strconv.Itoa(uid) + "_" + strconv.Itoa(typ) + "_" + strconv.Itoa(pageNum) + "_" + strconv.Itoa(limit)
	if pageNum == 1 {
		var cached map[string]any
		if cache.GetJSON(ctx, s.rdb, cacheKey, &cached) {
			if _, ok := cached["list"]; ok {
				return cached, nil
			}
		}
	}

	rows, err := s.learn.ListMyCourses(ctx, uid, typ, pageNum, limit)
	if err != nil {
		return nil, err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	list := make([]map[string]any, 0, len(rows))
	courseIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		courseID := r.CourseID
		if r.CID != nil && *r.CID > 0 {
			courseID = *r.CID
		}
		courseIDs = append(courseIDs, courseID)
		title, thumb, chapterCount := "", "", 0
		if r.Title != nil {
			title = *r.Title
		}
		if r.Thumb != nil {
			thumb = *r.Thumb
		}
		if r.ChapterCount != nil {
			chapterCount = *r.ChapterCount
		}
		list = append(list, map[string]any{
			"id": r.ID, "has_end_chapter": r.HasEndChapter, "is_end": r.IsEnd,
			"course": map[string]any{
				"id": courseID, "title": title, "thumb": thumb, "chapter_count": chapterCount,
			},
			"continueChapter": nil,
		})
	}
	continueMap, err := s.continueChapters(ctx, uid, courseIDs)
	if err != nil {
		return nil, err
	}
	for i := range list {
		course := list[i]["course"].(map[string]any)
		cid := course["id"].(int)
		if v, ok := continueMap[cid]; ok {
			list[i]["continueChapter"] = v
		}
	}
	total := (pageNum-1)*limit + len(list)
	if hasMore {
		total++
	}
	data := map[string]any{
		"current": pageNum, "list": list, "size": limit, "total": total, "hasMore": hasMore,
	}
	if pageNum == 1 {
		cache.SetJSON(ctx, s.rdb, cacheKey, data, 20*time.Second)
	}
	return data, nil
}

func (s *Service) MyTask(ctx context.Context, uid int) ([]map[string]any, error) {
	rows, err := s.learn.ListMyTasks(ctx, uid, 20)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		title, desc := "", ""
		if r.Title != nil {
			title = *r.Title
		}
		if r.ChapterTitle != nil {
			desc = *r.ChapterTitle
		}
		list = append(list, map[string]any{
			"courseId": r.CourseID, "chapterId": r.ChapterID, "paperId": r.PaperID,
			"title": title, "desc": desc,
		})
	}
	return list, nil
}

func (s *Service) MyHistory(ctx context.Context, uid, pageNum, limit int) (list []map[string]any, total int64, hasMore bool, err error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 30 {
		limit = 30
	}
	total, err = s.learn.CountHistory(ctx, uid)
	if err != nil {
		return nil, 0, false, err
	}
	rows, err := s.learn.ListHistory(ctx, uid, pageNum, limit)
	if err != nil {
		return nil, 0, false, err
	}
	hasMore = len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	list = make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		at := r.UpdatedAt.Format("2006-01-02 15:04:05")
		title, courseTitle, thumb := "课节", "课程", ""
		if r.ChapterTitle != nil && *r.ChapterTitle != "" {
			title = *r.ChapterTitle
		}
		if r.CourseTitle != nil && *r.CourseTitle != "" {
			courseTitle = *r.CourseTitle
		}
		if r.CourseThumb != nil {
			thumb = *r.CourseThumb
		}
		timePart := ""
		if len(at) >= 16 {
			timePart = at[11:16]
		}
		list = append(list, map[string]any{
			"id": r.ID, "chapterId": r.ChapterID, "courseId": r.CourseID,
			"title": title, "courseTitle": courseTitle, "thumb": thumb,
			"happenedAt": at, "date": at[:10], "time": timePart,
		})
	}
	return list, total, hasMore, nil
}

func (s *Service) SoonExpire(ctx context.Context, uid, pageNum, limit int) (list []map[string]any, total int64, err error) {
	deadline := time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02")
	total, err = s.learn.CountSoonExpire(ctx, uid, deadline)
	if err != nil {
		return nil, 0, err
	}
	rows, err := s.learn.ListSoonExpire(ctx, uid, pageNum, limit, deadline)
	if err != nil {
		return nil, 0, err
	}
	list = make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "user_id": r.UserID, "course_id": r.CourseID, "chapter_id": r.ChapterID,
			"is_end": r.IsEnd, "expire_time": r.ExpireTime,
		}
		if r.CourseTitle != nil || r.CourseThumb != nil {
			course := map[string]any{"id": r.CourseID}
			if r.CourseTitle != nil {
				course["title"] = *r.CourseTitle
			}
			if r.CourseThumb != nil {
				course["thumb"] = *r.CourseThumb
			}
			item["course"] = course
		}
		if r.ChapterTitle != nil || r.ChapterSort != nil {
			ch := map[string]any{"id": r.ChapterID}
			if r.ChapterTitle != nil {
				ch["title"] = *r.ChapterTitle
			}
			if r.ChapterSort != nil {
				ch["sort"] = *r.ChapterSort
			}
			item["chapter"] = ch
		}
		list = append(list, item)
	}
	return list, total, nil
}

func (s *Service) LearnNum(ctx context.Context, uid int) (map[string]any, error) {
	cacheKey := "learn_num_" + strconv.Itoa(uid)
	var cached map[string]any
	if cache.GetJSON(ctx, s.rdb, cacheKey, &cached) {
		if _, ok := cached["courseNum"]; ok {
			return cached, nil
		}
	}
	courseNum, _ := s.learn.CountEndedCourses(ctx, uid)
	rows, _ := s.learn.ListChapterProgress(ctx, uid)

	chapterNum, todayNum := 0, 0
	today := time.Now().Format("2006-01-02")
	doneIDs := []int{}
	for _, r := range rows {
		if r.IsEnd == 1 {
			chapterNum++
			if r.UpdatedAt != nil && r.UpdatedAt.Format("2006-01-02") >= today {
				todayNum++
			}
		}
		if r.HasEndChapter == 1 && r.ChapterID > 0 {
			doneIDs = append(doneIDs, r.ChapterID)
		}
	}
	timeSum := 0
	if len(doneIDs) > 0 {
		seen := map[int]struct{}{}
		uniq := make([]int, 0, len(doneIDs))
		for _, id := range doneIDs {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			uniq = append(uniq, id)
		}
		sum, _ := s.learn.SumMediaDuration(ctx, uniq)
		timeSum = int(sum)
	}
	data := map[string]any{
		"courseNum": courseNum, "chapterNum": chapterNum, "todayNum": todayNum, "time": timeSum,
	}
	cache.SetJSON(ctx, s.rdb, cacheKey, data, 30*time.Second)
	return data, nil
}

func (s *Service) Collect(ctx context.Context, uid, courseID int) error {
	id, err := s.learn.FindLikeCourseID(ctx, uid, courseID)
	if err == nil {
		return s.learn.DeleteLikeCourse(ctx, id)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.learn.CreateLikeCourse(ctx, uid, courseID)
	}
	return err
}

func (s *Service) CollectChapter(ctx context.Context, uid, chapterID int) error {
	courseID, err := s.learn.FindChapterCourseID(ctx, chapterID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrChapterNotFound
	}
	if err != nil {
		return err
	}
	id, err := s.learn.FindLikeChapterID(ctx, uid, chapterID)
	if err == nil {
		return s.learn.DeleteLikeChapter(ctx, id)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.learn.CreateLikeChapter(ctx, uid, chapterID, courseID)
	}
	return err
}

func (s *Service) JoinCourse(ctx context.Context, uid, courseID int) error {
	joined, _ := s.learn.CountUserCourse(ctx, uid, courseID)
	if joined > 0 {
		_, _ = enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, uid, courseID)
		return ErrAlreadyJoined
	}
	minNeed, err := s.learn.FindCourseMinCredit(ctx, courseID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCourseNotFound
	}
	if err != nil {
		return err
	}
	chapterCount, _ := s.learn.CountCourseChapters(ctx, courseID)
	if chapterCount == 0 {
		return ErrNoChapters
	}
	if minNeed > 0 {
		credit, _, _ := s.users.CreditCoins(ctx, uid)
		if credit < minNeed {
			return ErrCreditShort
		}
	}
	err = s.learn.JoinCourse(ctx, uid, courseID)
	if errors.Is(err, repo.ErrAlreadyJoined) {
		_, _ = enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, uid, courseID)
		return ErrAlreadyJoined
	}
	if err != nil {
		_, _ = enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, uid, courseID)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrAlreadyJoined
		}
		return ErrJoinFailed
	}
	_, _ = enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, uid, courseID)
	enroll.ForgetListCache(ctx, s.rdb, uid)
	return nil
}

func (s *Service) ApplyCourse(ctx context.Context, uid, id int, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrReasonRequired
	}
	ok, err := s.learn.CourseExists(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCourseNotFound
	}
	n, _ := s.learn.CountCourseChapters(ctx, id)
	if n == 0 {
		return ErrApplyNoChapters
	}
	n, _ = s.learn.CountUserCourse(ctx, uid, id)
	if n > 0 {
		return ErrAlreadyJoined
	}
	n, _ = s.learn.CountPendingCourseApply(ctx, uid, id)
	if n > 0 {
		return ErrAlreadyApplied
	}
	return s.learn.CreateCourseApply(ctx, uid, id, reason)
}

func (s *Service) UnlockChapter(ctx context.Context, uid, chapterID int) error {
	chapter, err := s.learn.FindChapterUnlock(ctx, chapterID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrChapterNotFound
	}
	if err != nil {
		return err
	}
	userChapter, err := s.learn.FindUserChapter(ctx, uid, chapterID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 已加入课程但缺进度行时（后加课节、历史数据）先补齐再解锁
		joined, _ := s.learn.CountUserCourse(ctx, uid, chapter.CourseID)
		if joined == 0 {
			return ErrNeedJoinCourse
		}
		if _, syncErr := enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, uid, chapter.CourseID); syncErr != nil {
			return syncErr
		}
		userChapter, err = s.learn.FindUserChapter(ctx, uid, chapterID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrChapterNotFound
		}
	}
	if err != nil {
		return err
	}
	if userChapter.IsEnd == 1 {
		return ErrAlreadyEnded
	}
	if userChapter.UnLock == 1 {
		return ErrAlreadyUnlocked
	}
	preID, err := s.learn.FindPrevChapterID(ctx, chapter.CourseID, chapter.Sort, chapter.ID)
	if err == nil {
		isEnd, err2 := s.learn.FindUserChapterIsEnd(ctx, uid, preID)
		if err2 != nil || isEnd != 1 {
			return ErrNeedPrevChapter
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	days := 30
	if chapter.LimitDay != nil && *chapter.LimitDay > 0 {
		days = *chapter.LimitDay
	}
	expire := time.Now().AddDate(0, 0, days)
	return s.learn.UnlockUserChapter(ctx, userChapter.ID, expire)
}

func (s *Service) LearnEnd(ctx context.Context, uid, chapterID int) (any, error) {
	chapter, err := s.learn.FindChapterLearn(ctx, chapterID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChapterNotFound
	}
	if err != nil {
		return nil, err
	}
	userChapter, err := s.learn.FindUserChapterFull(ctx, uid, chapterID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChapterNotFound
	}
	if err != nil {
		return nil, err
	}
	if userChapter.UnLock == 0 {
		return nil, ErrNeedUnlock
	}
	firstLearnEnd := userChapter.IsEnd == 0 && userChapter.HasEndChapter == 0
	if userChapter.IsEnd == 1 || userChapter.HasEndChapter == 1 {
		return userChapter, nil
	}
	userChapter.HasEndChapter = 1
	if chapter.PaperID != nil && *chapter.PaperID > 0 {
		pass, _ := s.learn.LatestPaperPass(ctx, uid, *chapter.PaperID)
		if pass == 1 {
			userChapter.IsEnd = 1
		}
	} else {
		userChapter.IsEnd = 1
	}
	if err := s.learn.UpdateUserChapterEnd(ctx, userChapter.ID, userChapter.HasEndChapter, userChapter.IsEnd); err != nil {
		return nil, err
	}
	ended, _ := s.learn.CountEndedChapters(ctx, uid, chapter.CourseID)
	_ = s.learn.UpdateUserCourseHasEnd(ctx, uid, chapter.CourseID, int(ended))
	enroll.ForgetListCache(ctx, s.rdb, uid)
	if firstLearnEnd {
		coin.OnLearnEnd(ctx, s.db, s.rdb, uid, chapter.ID, chapter.Type, chapter.Title)
	}
	return userChapter, nil
}

func (s *Service) EndCourse(ctx context.Context, uid, courseID int) error {
	userCourse, err := s.learn.FindUserCourse(ctx, uid, courseID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCourseNotFound
	}
	if err != nil {
		return err
	}
	if userCourse.IsEnd == 1 {
		return ErrCourseAlreadyEnded
	}
	course, err := s.learn.FindCourseEnd(ctx, courseID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCourseNotFound
	}
	if err != nil {
		return err
	}
	if course.PaperID != nil && *course.PaperID > 0 {
		return ErrCourseNeedExam
	}
	remark := course.Title + "课程学习结束发放"
	if course.Credit > 0 {
		granted, err := enroll.CourseCreditAlreadyGranted(ctx, s.db, uid, remark)
		if err != nil {
			return err
		}
		if !granted {
			_ = coin.AddCredit(ctx, s.db, uid, "课程发放", course.Credit, remark)
		}
	}
	if err := s.learn.MarkUserCourseEnded(ctx, userCourse.ID); err != nil {
		return err
	}
	enroll.ForgetListCache(ctx, s.rdb, uid)
	s.issuer.OnCourseEnded(ctx, uid, courseID)
	coin.OnCourseEnd(ctx, s.db, s.rdb, uid, courseID, course.Title)
	return nil
}

func (s *Service) continueChapters(ctx context.Context, uid int, courseIDs []int) (map[int]map[string]any, error) {
	seen := map[int]struct{}{}
	ids := make([]int, 0, len(courseIDs))
	for _, id := range courseIDs {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return map[int]map[string]any{}, nil
	}
	rows, err := s.learn.ListContinueChapters(ctx, uid, ids)
	if err != nil {
		return nil, err
	}
	grouped := map[int][]repo.ContinueChapterRow{}
	for _, r := range rows {
		grouped[r.CourseID] = append(grouped[r.CourseID], r)
	}
	out := map[int]map[string]any{}
	for courseID, chapters := range grouped {
		var picked *repo.ContinueChapterRow
		for i := range chapters {
			if chapters[i].UnLock == 1 && chapters[i].IsEnd != 1 {
				picked = &chapters[i]
				break
			}
		}
		if picked == nil {
			for i := range chapters {
				if chapters[i].UnLock != 0 {
					continue
				}
				if i == 0 || chapters[i-1].IsEnd == 1 {
					picked = &chapters[i]
					break
				}
			}
		}
		if picked == nil && len(chapters) > 0 {
			picked = &chapters[len(chapters)-1]
		}
		if picked == nil {
			continue
		}
		out[courseID] = map[string]any{
			"id": picked.ChapterID, "title": picked.Title, "course_id": courseID,
			"un_lock": picked.UnLock, "is_end": picked.IsEnd,
		}
	}
	return out, nil
}
