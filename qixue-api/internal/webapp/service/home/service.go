package home

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/service/catalog"
	webapprepo "qixue-api/internal/webapp/repo"
)

var (
	ErrBannerNotFound = errors.New("banner不存在")
	ErrNoticeNotFound = errors.New("公告不存在")
)

type Service struct {
	home webapprepo.HomeRepository
	db   *gorm.DB
	rdb  *redis.Client
}

func New(home webapprepo.HomeRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{home: home, db: db, rdb: rdb}
}

func (s *Service) GetHomeInfo(ctx context.Context) (map[string]any, error) {
	var (
		notice  *model.Notice
		swipers []map[string]any
	)
	if err := async.Run(
		func() error {
			var callErr error
			notice, callErr = s.home.LatestNoticeMeta(ctx)
			return callErr
		},
		func() error {
			var callErr error
			swipers, callErr = s.swiperListMeta(ctx)
			return callErr
		},
	); err != nil {
		return nil, err
	}
	var noticeOut any
	if notice != nil {
		noticeOut = noticeMetaView(notice)
	}
	return map[string]any{
		"notice":     noticeOut,
		"swiperList": swipers,
	}, nil
}

func (s *Service) Boot(ctx context.Context) (map[string]any, error) {
	var payload map[string]any
	err := cache.RememberTTL(ctx, s.rdb, "home_boot_v2", &payload, 5*time.Minute, func() (any, error) {
		return s.buildBoot(ctx)
	})
	return payload, err
}

func (s *Service) buildBoot(ctx context.Context) (map[string]any, error) {
	// 分类树走缓存，先取出首个二级分类，分类榜才能和其它指标同一波查询。
	categories, err := catalog.CourseTree(ctx, s.db, s.rdb)
	if err != nil {
		return nil, err
	}
	rankCateID := 0
	if len(categories) > 0 && len(categories[0].Children) > 0 {
		rankCateID = categories[0].Children[0].ID
	}

	var (
		notice   *model.Notice
		swipers  []map[string]any
		paths    []map[string]any
		hotRows  []map[string]any
		rankList []map[string]any
	)
	if err := async.Run(
		func() error {
			var callErr error
			notice, callErr = s.home.LatestNoticeMeta(ctx)
			return callErr
		},
		func() error {
			var callErr error
			swipers, callErr = s.swiperListMeta(ctx)
			return callErr
		},
		func() error {
			var callErr error
			paths, callErr = s.HotPaths(ctx, 10)
			return callErr
		},
		func() error {
			var callErr error
			hotRows, callErr = s.HotCoursesPage(ctx, 1, 11)
			return callErr
		},
		func() error {
			if rankCateID <= 0 {
				return nil
			}
			var callErr error
			rankList, callErr = s.HotByCate(ctx, rankCateID, 10)
			return callErr
		},
	); err != nil {
		return nil, err
	}
	hasMore := len(hotRows) > 10
	if hasMore {
		hotRows = hotRows[:10]
	}
	if rankList == nil {
		rankList = []map[string]any{}
	}
	var noticeOut any
	if notice != nil {
		noticeOut = noticeMetaView(notice)
	}
	return map[string]any{
		"notice":     noticeOut,
		"swiperList": swipers,
		"categories": categories,
		"pathList":   paths,
		"hotCourses": map[string]any{
			"list":    hotRows,
			"hasMore": hasMore,
		},
		"rankList": rankList,
	}, nil
}

func (s *Service) BannerDetail(ctx context.Context, id int) (map[string]any, error) {
	if id <= 0 {
		return nil, ErrBannerNotFound
	}
	row, err := s.home.FindBanner(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBannerNotFound
	}
	if err != nil {
		return nil, err
	}
	out := bannerMetaView(row)
	out["content"] = row.Content
	return out, nil
}

func (s *Service) NoticeDetail(ctx context.Context, id int) (map[string]any, error) {
	if id <= 0 {
		return nil, ErrNoticeNotFound
	}
	row, err := s.home.FindNotice(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoticeNotFound
	}
	if err != nil {
		return nil, err
	}
	out := noticeMetaView(row)
	out["content"] = row.Content
	return out, nil
}

func (s *Service) HotByCate(ctx context.Context, cateID, limit int) ([]map[string]any, error) {
	rows, err := s.home.ListHotCoursesByCate(ctx, cateID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]any{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": r.Desc,
			"type": r.Type, "credit": r.Credit, "min_need_credit": r.MinNeedCredit,
		})
	}
	return out, nil
}

func (s *Service) HotPaths(ctx context.Context, limit int) ([]map[string]any, error) {
	rows, err := s.home.ListHotPaths(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		out = append(out, map[string]any{
			"id": row.ID, "title": row.Title, "thumb": row.Thumb, "desc": row.Desc,
		})
	}
	return out, nil
}

func (s *Service) HotCoursesPage(ctx context.Context, pageNum, limit int) ([]map[string]any, error) {
	rows, err := s.home.ListHotCoursesPage(ctx, pageNum, limit)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": r.Desc,
			"category_id": r.CategoryID, "category": nil,
		}
		if r.CateName != nil {
			item["category"] = map[string]any{"id": r.CategoryID, "name": *r.CateName}
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) HotCoursesList(ctx context.Context, pageNum, limit int) (list []map[string]any, hasMore bool, total int64, err error) {
	list, err = s.HotCoursesPage(ctx, pageNum, limit+1)
	if err != nil {
		return nil, false, 0, err
	}
	hasMore = len(list) > limit
	if hasMore {
		list = list[:limit]
	}
	total = int64((pageNum-1)*limit + len(list))
	if hasMore {
		total++
	}
	return list, hasMore, total, nil
}

func (s *Service) swiperListMeta(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.home.ListActiveBannersMeta(ctx)
	out := make([]map[string]any, 0, len(rows))
	for i := range rows {
		out = append(out, bannerMetaView(&rows[i]))
	}
	return out, err
}

func noticeMetaView(row *model.Notice) map[string]any {
	return map[string]any{
		"id": row.ID, "title": row.Title, "status": row.Status,
		"publish_time": row.PublishTime.Format("2006-01-02 15:04:05"),
		"admin_id":     row.AdminID, "content": row.Content, "view_times": row.ViewTimes,
		"created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt),
	}
}

func bannerMetaView(row *model.Banner) map[string]any {
	return map[string]any{
		"id": row.ID, "name": row.Name, "sort": row.Sort,
		"img_url": row.ImgURL, "link_url": row.LinkURL,
		"type": row.Type, "status": row.Status,
		"created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt),
	}
}
