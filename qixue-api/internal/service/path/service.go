package path

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/catalog"
	"qixue-api/internal/service/enroll"
)

var (
	ErrNotFound      = errors.New("学习路径不存在")
	ErrAlreadyJoined = errors.New("已加入过了")
	ErrCreditLow     = errors.New("学分不足,可主动申请本系列课程")
	ErrReasonEmpty   = errors.New("请填写申请理由")
	ErrApplyExists   = errors.New("已提交过申请")
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

func (s *Service) Categories(ctx context.Context) ([]catalog.PathNode, error) {
	return catalog.PathList(ctx, s.db, s.rdb)
}

func (s *Service) List(ctx context.Context, categoryID, pageNum, limit int) (page.MoreBody, error) {
	if categoryID == 0 && pageNum == 1 {
		var cached page.MoreBody
		if cache.GetJSON(ctx, s.rdb, "path_list_p1", &cached) && cached.List != nil {
			return cached, nil
		}
	}

	total, rows, err := s.courses.ListShownPaths(ctx, categoryID, pageNum, limit)
	if err != nil {
		return page.MoreBody{}, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": r.Desc,
			"is_show": r.IsShow, "category_id": r.CategoryID, "is_rec": r.IsRec,
			"is_open": r.IsOpen, "min_need_credit": r.MinNeedCredit, "cert_id": r.CertID,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"created_id": r.CreatedID, "category": nil,
		}
		if r.CateName != nil {
			item["category"] = map[string]any{"id": r.CategoryID, "name": *r.CateName}
		}
		list = append(list, item)
	}
	payload := page.NewMore(pageNum, limit, total, list)
	if categoryID == 0 && pageNum == 1 {
		cache.SetJSON(ctx, s.rdb, "path_list_p1", payload, 20*time.Second)
	}
	return payload, nil
}

func (s *Service) Detail(ctx context.Context, id, uid int) (map[string]any, error) {
	path, err := s.courses.FindPathDetail(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	courses, _ := s.courses.ListPathCourses(ctx, id)
	courseList := make([]map[string]any, 0, len(courses))
	for _, r := range courses {
		item := map[string]any{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": r.Desc,
			"credit": r.Credit, "type": r.Type, "category_id": r.CategoryID,
			"chapter_count": r.ChapterCount, "category": nil,
		}
		if r.CateName != nil {
			item["category"] = map[string]any{"id": r.CategoryID, "name": *r.CateName}
		}
		courseList = append(courseList, item)
	}

	usersCount, _ := s.courses.CountPathUsers(ctx, id)

	userPath, err := s.courses.FindUserPath(ctx, uid, id)
	hasJoin := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	showApply := false
	applyPending := false
	if !hasJoin {
		credit, _, _ := s.users.CreditCoins(ctx, uid)
		if credit < path.MinNeedCredit {
			showApply = true
		}
		n, _ := s.courses.CountPendingPathApply(ctx, uid, id)
		applyPending = n > 0
	}

	out := map[string]any{
		"id": path.ID, "title": path.Title, "thumb": path.Thumb, "desc": path.Desc,
		"content": path.Content, "category_id": path.CategoryID,
		"min_need_credit": path.MinNeedCredit, "cert_id": path.CertID,
		"created_at": model.FormatTime(path.CreatedAt), "updated_at": model.FormatTime(path.UpdatedAt),
		"created_id": path.CreatedID, "users_count": usersCount,
		"courses": courseList, "category": nil, "cert": nil,
		"has_join": hasJoin, "is_end": 0, "show_apply": showApply,
		"apply_pending": applyPending, "my_cert": nil,
	}
	if hasJoin {
		out["is_end"] = userPath.IsEnd
	}
	if path.CateName != nil {
		out["category"] = map[string]any{"id": path.CategoryID, "name": *path.CateName}
	}
	if path.CertID != nil && path.CertName != nil {
		out["cert"] = map[string]any{"id": *path.CertID, "name": *path.CertName}
	}
	if hasJoin && path.CertID != nil && *path.CertID > 0 {
		if myCert, err := s.courses.FindValidPathCert(ctx, uid, id); err == nil {
			out["my_cert"] = myCert
		}
	}
	return out, nil
}

func (s *Service) Join(ctx context.Context, id, uid int) error {
	_, minCredit, err := s.courses.FindPathMinCredit(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	n, err := s.courses.CountUserOnPath(ctx, uid, id)
	if err != nil {
		return err
	}
	if n > 0 {
		_ = enroll.EnrollUsersToPath(ctx, s.db, s.rdb, id, []int{uid})
		return ErrAlreadyJoined
	}
	credit, _, _ := s.users.CreditCoins(ctx, uid)
	if credit < minCredit {
		return ErrCreditLow
	}
	return enroll.EnrollUsersToPath(ctx, s.db, s.rdb, id, []int{uid})
}

func (s *Service) Apply(ctx context.Context, id, uid int, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrReasonEmpty
	}
	ok, err := s.courses.PathExists(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	n, err := s.courses.CountUserOnPath(ctx, uid, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrAlreadyJoined
	}
	n, err = s.courses.CountPendingPathApply(ctx, uid, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrApplyExists
	}
	return s.courses.CreatePathApply(ctx, uid, id, reason)
}
