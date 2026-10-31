package repo

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CourseListFilter struct {
	CategoryIDs []int
	ForceNone   bool
	Type        int
	Title       string
	SortType    string
	Page        int
	Limit       int
}

type CourseListRow struct {
	ID            int     `gorm:"column:id"`
	Title         string  `gorm:"column:title"`
	Thumb         string  `gorm:"column:thumb"`
	Desc          string  `gorm:"column:desc"`
	Type          int     `gorm:"column:type"`
	Credit        int     `gorm:"column:credit"`
	MinNeedCredit int     `gorm:"column:min_need_credit"`
	CategoryID    int     `gorm:"column:category_id"`
	CateName      *string `gorm:"column:cate_name"`
}

type CourseDetailRow struct {
	ID            int     `gorm:"column:id"`
	Title         string  `gorm:"column:title"`
	Thumb         string  `gorm:"column:thumb"`
	Desc          string  `gorm:"column:desc"`
	Type          int     `gorm:"column:type"`
	Credit        int     `gorm:"column:credit"`
	Content       *string `gorm:"column:content"`
	PaperID       *int    `gorm:"column:paper_id"`
	CertID        *int    `gorm:"column:cert_id"`
	CategoryID    int     `gorm:"column:category_id"`
	MinNeedCredit int     `gorm:"column:min_need_credit"`
	ViewCount     int     `gorm:"column:view_count"`
	ChapterCount  int     `gorm:"column:chapter_count"`
	CateName      *string `gorm:"column:cate_name"`
	CertName      *string `gorm:"column:cert_name"`
}

type CourseAttachRow struct {
	ID            int64      `gorm:"column:id" json:"id"`
	CourseID      int        `gorm:"column:course_id" json:"course_id"`
	Name          string     `gorm:"column:name" json:"name"`
	Path          string     `gorm:"column:path" json:"path"`
	Size          int        `gorm:"column:size" json:"size"`
	DownloadTimes int        `gorm:"column:download_times" json:"download_times"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type UserCertBrief struct {
	ID       int        `gorm:"column:id" json:"id"`
	UserID   int        `gorm:"column:user_id" json:"user_id"`
	CertID   int        `gorm:"column:cert_id" json:"cert_id"`
	PathID   int        `gorm:"column:path_id" json:"path_id"`
	CourseID int        `gorm:"column:course_id" json:"course_id"`
	CertNo   string     `gorm:"column:cert_no" json:"cert_no"`
	Title    string     `gorm:"column:title" json:"title"`
	Image    string     `gorm:"column:image" json:"image"`
	IsValid  int        `gorm:"column:is_valid" json:"is_valid"`
	IssuedAt *time.Time `gorm:"column:issued_at" json:"issued_at"`
}

type ChapterListRow struct {
	ID      int     `gorm:"column:id"`
	Title   string  `gorm:"column:title"`
	Type    int     `gorm:"column:type"`
	Sort    int     `gorm:"column:sort"`
	MediaID *int    `gorm:"column:media_id"`
	PaperID *int    `gorm:"column:paper_id"`
	MDur    *int    `gorm:"column:m_duration"`
	MType   *string `gorm:"column:m_type"`
}

type ChapterDetailRow struct {
	ID       int     `gorm:"column:id"`
	CourseID int     `gorm:"column:course_id"`
	Title    string  `gorm:"column:title"`
	Type     int     `gorm:"column:type"`
	Sort     int     `gorm:"column:sort"`
	Content  *string `gorm:"column:content"`
	MediaURL *string `gorm:"column:media_url"`
	MediaID  *int    `gorm:"column:media_id"`
	PaperID  *int    `gorm:"column:paper_id"`
	CTitle   *string `gorm:"column:c_title"`
	CThumb   *string `gorm:"column:c_thumb"`
	MDur     *int    `gorm:"column:m_duration"`
	MType    *string `gorm:"column:m_type"`
	MSize    *int64  `gorm:"column:m_size"`
	MThumb   *string `gorm:"column:m_thumb"`
}

type CatalogUserChapterRow struct {
	ID            int        `gorm:"column:id" json:"id"`
	UserID        int        `gorm:"column:user_id" json:"user_id"`
	CourseID      int        `gorm:"column:course_id" json:"course_id"`
	ChapterID     int        `gorm:"column:chapter_id" json:"chapter_id"`
	UnLock        int        `gorm:"column:un_lock" json:"un_lock"`
	HasEndChapter int        `gorm:"column:has_end_chapter" json:"has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end" json:"is_end"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ChapterNeighborRow struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Type  int    `gorm:"column:type"`
	Sort  int    `gorm:"column:sort"`
}

type NeighborProgressRow struct {
	ChapterID int `gorm:"column:chapter_id"`
	UnLock    int `gorm:"column:un_lock"`
	IsEnd     int `gorm:"column:is_end"`
}

type PathListRow struct {
	ID            int        `gorm:"column:id"`
	Title         string     `gorm:"column:title"`
	Thumb         string     `gorm:"column:thumb"`
	Desc          string     `gorm:"column:desc"`
	IsShow        int        `gorm:"column:is_show"`
	CategoryID    int        `gorm:"column:category_id"`
	IsRec         int        `gorm:"column:is_rec"`
	MinNeedCredit int        `gorm:"column:min_need_credit"`
	CertID        *int       `gorm:"column:cert_id"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	CreatedID     int        `gorm:"column:created_id"`
	CateName      *string    `gorm:"column:cate_name"`
}

type PathDetailRow struct {
	ID            int        `gorm:"column:id"`
	Title         string     `gorm:"column:title"`
	Thumb         string     `gorm:"column:thumb"`
	Desc          string     `gorm:"column:desc"`
	Content       *string    `gorm:"column:content"`
	CategoryID    int        `gorm:"column:category_id"`
	MinNeedCredit int        `gorm:"column:min_need_credit"`
	CertID        *int       `gorm:"column:cert_id"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	CreatedID     int        `gorm:"column:created_id"`
	CateName      *string    `gorm:"column:cate_name"`
	CertName      *string    `gorm:"column:cert_name"`
}

type PathCourseBrief struct {
	ID           int     `gorm:"column:id"`
	Title        string  `gorm:"column:title"`
	Thumb        string  `gorm:"column:thumb"`
	Desc         string  `gorm:"column:desc"`
	Credit       int     `gorm:"column:credit"`
	Type         int     `gorm:"column:type"`
	CategoryID   int     `gorm:"column:category_id"`
	ChapterCount int     `gorm:"column:chapter_count"`
	CateName     *string `gorm:"column:cate_name"`
}

type UserPathBrief struct {
	ID    int `gorm:"column:id"`
	IsEnd int `gorm:"column:is_end"`
}

func (r *courseRepo) ListShownCourses(ctx context.Context, f CourseListFilter) (int64, []CourseListRow, error) {
	query := r.db.WithContext(ctx).Table("courses AS c").
		Select("c.id, c.title, c.thumb, c.`desc`, c.type, c.credit, c.min_need_credit, c.category_id, cat.name AS cate_name").
		Joins("LEFT JOIN course_categories cat ON cat.id = c.category_id").
		Where("c.is_show = 1 AND c.deleted_at IS NULL")
	if f.ForceNone {
		query = query.Where("0 = 1")
	} else if len(f.CategoryIDs) > 0 {
		query = query.Where("c.category_id IN ?", f.CategoryIDs)
	}
	if f.Type > 0 {
		query = query.Where("c.type = ?", f.Type)
	}
	if title := strings.TrimSpace(f.Title); title != "" {
		query = query.Where("c.title LIKE ?", "%"+title+"%")
	}
	switch f.SortType {
	case "hot":
		query = query.Order("c.view_count DESC")
	case "new":
		query = query.Order("c.id DESC")
	default:
		query = query.Order("c.id ASC")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []CourseListRow
	err := query.Offset((f.Page - 1) * f.Limit).Limit(f.Limit).Find(&rows).Error
	return total, rows, err
}

func (r *courseRepo) FindCourseDetail(ctx context.Context, id int) (*CourseDetailRow, error) {
	var row CourseDetailRow
	err := r.db.WithContext(ctx).Table("courses AS c").
		Select("c.id, c.title, c.thumb, c.`desc`, c.type, c.credit, c.content, c.paper_id, c.cert_id, c.category_id, c.min_need_credit, c.view_count, c.chapter_count, cat.name AS cate_name, cert.name AS cert_name").
		Joins("LEFT JOIN course_categories cat ON cat.id = c.category_id").
		Joins("LEFT JOIN certs cert ON cert.id = c.cert_id").
		Where("c.id = ? AND c.deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) CountPendingCourseApply(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course_apply").
		Where("user_id = ? AND course_id = ? AND status != 2", uid, courseID).
		Count(&n).Error
	return n, err
}

func (r *courseRepo) CountCourseLike(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_like_courses").
		Where("user_id = ? AND course_id = ?", uid, courseID).Count(&n).Error
	return n, err
}

func (r *courseRepo) ListCourseAttaches(ctx context.Context, courseID int) ([]CourseAttachRow, error) {
	var rows []CourseAttachRow
	err := r.db.WithContext(ctx).Table("course_attach").
		Where("course_id = ? AND deleted_at IS NULL", courseID).
		Order("created_at ASC").Find(&rows).Error
	if rows == nil {
		rows = []CourseAttachRow{}
	}
	return rows, err
}

func (r *courseRepo) FindValidCourseCert(ctx context.Context, uid, courseID int) (*UserCertBrief, error) {
	var row UserCertBrief
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("user_id = ? AND course_id = ? AND path_id = 0 AND is_valid = 1", uid, courseID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) IncrCourseView(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("courses").Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *courseRepo) ListChaptersWithMedia(ctx context.Context, courseID int) ([]ChapterListRow, error) {
	var rows []ChapterListRow
	err := r.db.WithContext(ctx).Table("course_chapters AS cc").
		Select("cc.id, cc.title, cc.type, cc.sort, cc.media_id, cc.paper_id, mf.duration AS m_duration, mf.type AS m_type").
		Joins("LEFT JOIN media_file mf ON mf.id = cc.media_id").
		Where("cc.course_id = ? AND cc.deleted_at IS NULL", courseID).
		Order("cc.sort ASC, cc.id ASC").Find(&rows).Error
	return rows, err
}

func (r *courseRepo) FindChapterDetail(ctx context.Context, id int) (*ChapterDetailRow, error) {
	var row ChapterDetailRow
	err := r.db.WithContext(ctx).Table("course_chapters AS cc").
		Select("cc.id, cc.course_id, cc.title, cc.type, cc.sort, cc.content, cc.media_url, cc.media_id, cc.paper_id, c.title AS c_title, c.thumb AS c_thumb, mf.duration AS m_duration, mf.type AS m_type, mf.size AS m_size, mf.thumb AS m_thumb").
		Joins("LEFT JOIN courses c ON c.id = cc.course_id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN media_file mf ON mf.id = cc.media_id").
		Where("cc.id = ? AND cc.deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) FindUserChapter(ctx context.Context, uid, chapterID int) (*CatalogUserChapterRow, error) {
	var row CatalogUserChapterRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) CountChapterLike(ctx context.Context, uid, chapterID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_like_chapters").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).Count(&n).Error
	return n, err
}

func (r *courseRepo) IncrChapterView(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("course_chapters").Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *courseRepo) FindLatestUserPaper(ctx context.Context, uid, paperID int) (*LatestUserPaper, error) {
	if uid <= 0 || paperID <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var row LatestUserPaper
	err := r.db.WithContext(ctx).Table("user_paper").
		Select("id, is_pass, is_people_check").
		Where("user_id = ? AND paper_id = ? AND deleted_at IS NULL", uid, paperID).
		Order("id DESC").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) ListPendingUserPapers(ctx context.Context, uid int, paperIDs []int) ([]PendingUserPaper, error) {
	if uid <= 0 || len(paperIDs) == 0 {
		return nil, nil
	}
	var rows []PendingUserPaper
	err := r.db.WithContext(ctx).Table("user_paper AS up").
		Select("up.id, up.paper_id").
		Joins(`INNER JOIN (
			SELECT paper_id, MAX(id) AS id
			FROM user_paper
			WHERE user_id = ? AND deleted_at IS NULL AND paper_id IN ?
			GROUP BY paper_id
		) latest ON latest.id = up.id`, uid, paperIDs).
		Where("up.is_people_check = 1 AND up.is_pass = 0").
		Find(&rows).Error
	return rows, err
}

func (r *courseRepo) FindPrevChapter(ctx context.Context, courseID, chapterID, sort int) (*ChapterNeighborRow, error) {
	var row ChapterNeighborRow
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title, type, sort").
		Where("course_id = ? AND deleted_at IS NULL AND (sort < ? OR (sort = ? AND id < ?))", courseID, sort, sort, chapterID).
		Order("sort DESC, id DESC").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) FindNextChapter(ctx context.Context, courseID, chapterID, sort int) (*ChapterNeighborRow, error) {
	var row ChapterNeighborRow
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title, type, sort").
		Where("course_id = ? AND deleted_at IS NULL AND (sort > ? OR (sort = ? AND id > ?))", courseID, sort, sort, chapterID).
		Order("sort ASC, id ASC").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) ListNeighborProgress(ctx context.Context, uid int, ids []int) ([]NeighborProgressRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []NeighborProgressRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id IN ?", uid, ids).Find(&rows).Error
	return rows, err
}

func (r *courseRepo) UnlockUserChapter(ctx context.Context, uid, chapterID int) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).
		Updates(map[string]any{"un_lock": 1}).Error
}

func (r *courseRepo) ListShownPaths(ctx context.Context, categoryID, page, limit int) (int64, []PathListRow, error) {
	query := r.db.WithContext(ctx).Table("paths AS p").
		Select("p.id, p.title, p.thumb, p.`desc`, p.is_show, p.category_id, p.is_rec, p.min_need_credit, p.cert_id, p.created_at, p.updated_at, p.created_id, cat.name AS cate_name").
		Joins("LEFT JOIN path_categories cat ON cat.id = p.category_id").
		Where("p.is_show = 1 AND p.deleted_at IS NULL")
	if categoryID > 0 {
		query = query.Where("p.category_id = ?", categoryID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []PathListRow
	err := query.Order("p.updated_at DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return total, rows, err
}

func (r *courseRepo) FindPathDetail(ctx context.Context, id int) (*PathDetailRow, error) {
	var row PathDetailRow
	err := r.db.WithContext(ctx).Table("paths AS p").
		Select("p.id, p.title, p.thumb, p.`desc`, p.content, p.category_id, p.min_need_credit, p.cert_id, p.created_at, p.updated_at, p.created_id, cat.name AS cate_name, cert.name AS cert_name").
		Joins("LEFT JOIN path_categories cat ON cat.id = p.category_id").
		Joins("LEFT JOIN certs cert ON cert.id = p.cert_id").
		Where("p.id = ? AND p.deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) ListPathCourses(ctx context.Context, pathID int) ([]PathCourseBrief, error) {
	var rows []PathCourseBrief
	err := r.db.WithContext(ctx).Table("path_courses AS pc").
		Select("c.id, c.title, c.thumb, c.`desc`, c.credit, c.type, c.category_id, c.chapter_count, cat.name AS cate_name").
		Joins("JOIN courses c ON c.id = pc.course_id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN course_categories cat ON cat.id = c.category_id").
		Where("pc.path_id = ? AND pc.deleted_at IS NULL", pathID).
		Order("pc.sort ASC").Find(&rows).Error
	return rows, err
}

func (r *courseRepo) CountPathUsers(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").
		Where("path_id = ? AND deleted_at IS NULL", pathID).Count(&n).Error
	return n, err
}

func (r *courseRepo) FindUserPath(ctx context.Context, uid, pathID int) (*UserPathBrief, error) {
	var row UserPathBrief
	err := r.db.WithContext(ctx).Table("user_path").
		Where("user_id = ? AND path_id = ? AND deleted_at IS NULL", uid, pathID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) CountPendingPathApply(ctx context.Context, uid, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path_apply").
		Where("user_id = ? AND path_id = ? AND status != 2", uid, pathID).
		Count(&n).Error
	return n, err
}

func (r *courseRepo) FindValidPathCert(ctx context.Context, uid, pathID int) (*UserCertBrief, error) {
	var row UserCertBrief
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("user_id = ? AND path_id = ? AND is_valid = 1", uid, pathID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) FindPathMinCredit(ctx context.Context, id int) (int, int, error) {
	var row struct {
		ID            int `gorm:"column:id"`
		MinNeedCredit int `gorm:"column:min_need_credit"`
	}
	err := r.db.WithContext(ctx).Table("paths").
		Select("id, min_need_credit").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return 0, 0, err
	}
	return row.ID, row.MinNeedCredit, nil
}

func (r *courseRepo) CountUserOnPath(ctx context.Context, uid, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").
		Where("user_id = ? AND path_id = ? AND deleted_at IS NULL", uid, pathID).Count(&n).Error
	return n, err
}

func (r *courseRepo) PathExists(ctx context.Context, id int) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("paths").
		Where("id = ? AND deleted_at IS NULL", id).Count(&n).Error
	return n > 0, err
}

func (r *courseRepo) CreatePathApply(ctx context.Context, uid, pathID int, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Table("user_path_apply").Create(map[string]any{
		"user_id": uid, "path_id": pathID, "status": 0, "reason": reason,
		"created_at": now, "updated_at": now,
	}).Error
}
