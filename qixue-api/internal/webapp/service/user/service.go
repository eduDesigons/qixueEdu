package user

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	webapprepo "qixue-api/internal/webapp/repo"
	"qixue-api/internal/service/auth"
)

var (
	ErrNotFound    = errors.New("用户不存在")
	ErrBadParam    = errors.New("参数错误")
	ErrPwdRequired = errors.New("请填写旧密码和新密码")
	ErrPwdMismatch = errors.New("旧密码错误")
)

type Service struct {
	users    repo.UserRepository
	learner  webapprepo.LearnerRepository
	comments webapprepo.CommentRepository
	rdb      *redis.Client
}

func New(users repo.UserRepository, learner webapprepo.LearnerRepository, comments webapprepo.CommentRepository, rdb *redis.Client) *Service {
	return &Service{users: users, learner: learner, comments: comments, rdb: rdb}
}

func (s *Service) Info(ctx context.Context, uid int) (map[string]any, error) {
	row, err := s.users.FindProfile(ctx, uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return row, err
}

func (s *Service) Edit(ctx context.Context, uid int, body map[string]any) error {
	if body == nil {
		return ErrBadParam
	}
	allowed := []string{"avatar", "nick_name", "mobile", "gender", "age", "degree", "profession", "job", "marital", "target"}
	updates := map[string]any{}
	for _, key := range allowed {
		if v, ok := body[key]; ok {
			updates[key] = v
		}
	}
	ok, err := s.users.Exists(ctx, uid)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	if err := s.users.Update(ctx, uid, updates); err != nil {
		return err
	}
	cache.Drop(ctx, s.rdb, "user_my_count_"+strconv.Itoa(uid))
	return nil
}

func (s *Service) SetPwd(ctx context.Context, uid int, oldPwd, newPwd string) error {
	if strings.TrimSpace(oldPwd) == "" || strings.TrimSpace(newPwd) == "" {
		return ErrPwdRequired
	}
	if err := auth.CheckPassword(newPwd); err != nil {
		return err
	}
	hash, err := s.users.FindPassword(ctx, uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if !auth.PasswordOK(hash, oldPwd) {
		return ErrPwdMismatch
	}
	newHash, err := auth.HashPassword(newPwd)
	if err != nil {
		return err
	}
	return s.users.Update(ctx, uid, map[string]any{"password": newHash})
}

func (s *Service) MyCount(ctx context.Context, uid int) (map[string]any, error) {
	cacheKey := "user_my_count_" + strconv.Itoa(uid)
	var cached map[string]any
	if cache.GetJSON(ctx, s.rdb, cacheKey, &cached) {
		if _, ok := cached["courseCount"]; ok {
			return cached, nil
		}
	}
	var credit, coins int
	var courseCount, unread int64
	if err := async.Run(
		func() error {
			var callErr error
			credit, coins, callErr = s.users.CreditCoins(ctx, uid)
			return callErr
		},
		func() error {
			var callErr error
			courseCount, callErr = s.learner.CountUserCourses(ctx, uid)
			return callErr
		},
		func() error {
			var callErr error
			unread, callErr = s.learner.CountUnreadMessages(ctx, uid)
			return callErr
		},
	); err != nil {
		return nil, err
	}
	data := map[string]any{
		"courseCount":        courseCount,
		"credit":             credit,
		"coins":              coins,
		"unreadMessageCount": unread,
	}
	cache.SetJSON(ctx, s.rdb, cacheKey, data, 20*time.Second)
	return data, nil
}

func (s *Service) MyMessages(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	total, rows, err := s.learner.ListMessages(ctx, uid, page, limit)
	if err != nil {
		return 0, nil, err
	}
	formatTimes(rows)
	return total, rows, nil
}

func (s *Service) UnreadMessageCount(ctx context.Context, uid int) (int64, error) {
	return s.learner.CountUnreadMessages(ctx, uid)
}

func (s *Service) ReadMessage(ctx context.Context, uid, id int) (int64, error) {
	n, err := s.learner.MarkMessagesRead(ctx, uid, id)
	if err != nil {
		return 0, err
	}
	cache.Drop(ctx, s.rdb, "user_my_count_"+strconv.Itoa(uid))
	return n, nil
}

func (s *Service) MyLikeCourse(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []webapprepo.LikeCourseRow
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountLikeCourses(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListLikeCourses(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		courseID := r.CourseID
		if r.CID != nil {
			courseID = *r.CID
		}
		list = append(list, map[string]any{
			"id": r.ID, "user_id": r.UserID, "course_id": r.CourseID,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"course": map[string]any{
				"id": courseID, "title": deref(r.Title), "thumb": deref(r.Thumb), "desc": deref(r.Desc),
			},
		})
	}
	return total, list, nil
}

func (s *Service) MyLikeChapter(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []webapprepo.LikeChapterRow
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountLikeChapters(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListLikeChapters(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		chapterID, courseID := r.ChapterID, r.CourseID
		if r.ChID != nil {
			chapterID = *r.ChID
		}
		if r.CID != nil {
			courseID = *r.CID
		}
		list = append(list, map[string]any{
			"id": r.ID, "user_id": r.UserID, "course_id": r.CourseID, "chapter_id": r.ChapterID,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"chapter": map[string]any{"id": chapterID, "title": deref(r.ChTitle)},
			"course":  map[string]any{"id": courseID, "title": deref(r.CTitle), "thumb": deref(r.CThumb), "desc": deref(r.CDesc)},
		})
	}
	return total, list, nil
}

func (s *Service) MyCourseComments(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []model.CourseComment
	if err := async.Run(
		func() error {
			var err error
			total, err = s.comments.CountUserCourseComments(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.comments.ListUserCourseComments(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list, err := s.decorateMyCourseComments(ctx, rows)
	return total, list, err
}

func (s *Service) MyChapterComments(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []model.CourseChapterComment
	if err := async.Run(
		func() error {
			var err error
			total, err = s.comments.CountUserChapterComments(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.comments.ListUserChapterComments(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list, err := s.decorateMyChapterComments(ctx, rows)
	return total, list, err
}

func (s *Service) MyPapers(ctx context.Context, uid int, paperType string, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []webapprepo.UserPaperRow
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountUserPapers(ctx, uid, paperType)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListUserPapers(ctx, uid, paperType, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		title := ""
		if r.ChapterTitle != nil && *r.ChapterTitle != "" {
			title = *r.ChapterTitle
		} else if r.CourseTitle != nil && *r.CourseTitle != "" {
			title = *r.CourseTitle
		}
		if title == "" {
			if r.PType == "course" {
				title = "结课考试"
			} else {
				title = "章节测评"
			}
		}
		paper := map[string]any{
			"id": r.PaperID, "type": r.PType, "course_id": r.CourseID, "chapter_id": r.ChapterID, "title": title,
		}
		if r.CID != nil || r.CourseTitle != nil {
			paper["course"] = map[string]any{
				"id": derefInt(r.CID, r.CourseID), "title": deref(r.CourseTitle), "thumb": deref(r.CourseThumb),
			}
		}
		if r.ChapterID != nil && *r.ChapterID > 0 {
			paper["chapter"] = map[string]any{"id": derefInt(r.ChID, *r.ChapterID), "title": deref(r.ChapterTitle)}
		}
		list = append(list, map[string]any{
			"id": r.ID, "user_id": r.UserID, "paper_id": r.PaperID, "score": r.Score,
			"is_pass": r.IsPass, "is_people_check": r.IsPeopleCheck,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"paper": paper,
		})
	}
	return total, list, nil
}

func (s *Service) MyRecords(ctx context.Context, table string, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []map[string]any
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountRecords(ctx, table, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListRecords(ctx, table, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	formatTimes(rows)
	return total, rows, nil
}

func (s *Service) MyPath(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []webapprepo.UserPathRow
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountUserPaths(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListUserPaths(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	pathIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		pathIDs = append(pathIDs, r.PathID)
	}
	certs, err := s.loadPathCerts(ctx, uid, pathIDs)
	if err != nil {
		return 0, nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		pathID := r.PathID
		if r.PID != nil {
			pathID = *r.PID
		}
		list = append(list, map[string]any{
			"id": r.ID, "user_id": r.UserID, "path_id": r.PathID,
			"has_end_course": r.HasEndCourse, "is_end": r.IsEnd,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"path": map[string]any{
				"id": pathID, "title": deref(r.PTitle), "thumb": deref(r.PThumb),
				"desc": deref(r.PDesc), "cert_id": r.PCertID,
			},
			"user_cert": certs[r.PathID],
		})
	}
	return total, list, nil
}

func (s *Service) MyCert(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	var total int64
	var rows []webapprepo.UserCertRow
	if err := async.Run(
		func() error {
			var err error
			total, err = s.learner.CountValidCerts(ctx, uid)
			return err
		},
		func() error {
			var err error
			rows, err = s.learner.ListValidCerts(ctx, uid, page, limit)
			return err
		},
	); err != nil {
		return 0, nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "user_id": r.UserID, "cert_id": r.CertID, "path_id": r.PathID, "course_id": r.CourseID,
			"cert_no": r.CertNo, "title": r.Title, "image": r.Image, "is_valid": r.IsValid,
			"issued_at": model.FormatTime(r.IssuedAt), "created_at": model.FormatTime(r.CreatedAt),
			"cert_image": r.Image, "cert_name": r.Title, "name": r.Title,
		}
		if r.PathID > 0 && r.PTitle != nil {
			item["source_type"] = "path"
			item["source_title"] = *r.PTitle
			item["path"] = map[string]any{"id": derefInt(r.PID, r.PathID), "title": *r.PTitle}
			item["course"] = map[string]any{"id": derefInt(r.PID, r.PathID), "title": *r.PTitle}
		} else if r.CourseID > 0 && r.CTitle != nil {
			item["source_type"] = "course"
			item["source_title"] = *r.CTitle
			item["course"] = map[string]any{"id": derefInt(r.CID, r.CourseID), "title": *r.CTitle}
		} else {
			item["source_type"] = ""
			item["source_title"] = ""
		}
		list = append(list, item)
	}
	return total, list, nil
}
