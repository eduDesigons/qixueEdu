package repo

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AdminMemberRepository covers admin member (学员) Table-based lookups.
type AdminMemberRepository interface {
	ListAreasByIDs(ctx context.Context, ids []int, withParent bool) ([]MemberAreaRow, error)
	FillCityFromArea(ctx context.Context, areaID int) (cityID *int, cityCode *string, err error)
	TopCategoryBaseCredit(ctx context.Context, categoryID int) (credit int, ok bool, err error)

	ListUserCourses(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserCourseRow, int64, error)
	ListUserPaths(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserPathRow, int64, error)
	ListUserChapters(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserChapterRow, int64, error)
	ListUserRecords(ctx context.Context, table string, userID, offset, limit int) ([]MemberRecordRow, int64, error)
	ListUserPapers(ctx context.Context, userID int, paperType string, isPass, isPeopleCheck *int, offset, limit int) ([]MemberUserPaperRow, int64, error)
	ListUserCerts(ctx context.Context, userID, offset, limit int) ([]MemberUserCertRow, int64, error)

	ListCoursesFull(ctx context.Context, ids []int) ([]MemberCourseFullRow, error)
	ListCoursesBrief(ctx context.Context, ids []int) ([]MemberCourseBriefRow, error)
	ListCoursesThumb(ctx context.Context, ids []int) ([]MemberCourseThumbRow, error)
	ListPathsFull(ctx context.Context, ids []int) ([]MemberPathFullRow, error)
	ListChaptersByIDs(ctx context.Context, ids []int) ([]MemberIDTitle, error)
	ListPapersByIDs(ctx context.Context, ids []int) ([]MemberPaperRow, error)
	ListNamedByIDs(ctx context.Context, table string, ids []int) ([]MemberIDName, error)
	ListTitlesByIDs(ctx context.Context, table string, ids []int, softDelete bool) ([]MemberIDTitle, error)
	ListCertNames(ctx context.Context, ids []int) ([]MemberIDName, error)
}

type MemberAreaRow struct {
	ID         int    `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	ParentName string `gorm:"column:parent_name"`
}

type MemberUserCourseRow struct {
	ID            int        `gorm:"column:id"`
	CourseID      int        `gorm:"column:course_id"`
	UserID        int        `gorm:"column:user_id"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	HasEndChapter int        `gorm:"column:has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end"`
}

type MemberUserPathRow struct {
	ID           int        `gorm:"column:id"`
	PathID       int        `gorm:"column:path_id"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	HasEndCourse int        `gorm:"column:has_end_course"`
	IsEnd        int        `gorm:"column:is_end"`
}

type MemberUserChapterRow struct {
	ID        int        `gorm:"column:id"`
	CourseID  *int       `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	UserID    int        `gorm:"column:user_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	IsEnd     int        `gorm:"column:is_end"`
	UnLock    int        `gorm:"column:un_lock"`
}

type MemberRecordRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	Field     string     `gorm:"column:field"`
	Sum       int        `gorm:"column:sum"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}

type MemberUserPaperRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	PaperID       int        `gorm:"column:paper_id"`
	Score         int        `gorm:"column:score"`
	SubmitCount   int        `gorm:"column:submit_count"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	IsPass        int        `gorm:"column:is_pass"`
	IsPeopleCheck int        `gorm:"column:is_people_check"`
}

type MemberUserCertRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CertID    int        `gorm:"column:cert_id"`
	PathID    *int       `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	CertNo    string     `gorm:"column:cert_no"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   int        `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}

type MemberCourseFullRow struct {
	ID         int    `gorm:"column:id"`
	Title      string `gorm:"column:title"`
	Thumb      string `gorm:"column:thumb"`
	Type       int    `gorm:"column:type"`
	Desc       string `gorm:"column:desc"`
	IsShow     int    `gorm:"column:is_show"`
	CategoryID int    `gorm:"column:category_id"`
}

type MemberCourseBriefRow struct {
	ID         int    `gorm:"column:id"`
	Title      string `gorm:"column:title"`
	Thumb      string `gorm:"column:thumb"`
	CategoryID int    `gorm:"column:category_id"`
}

type MemberCourseThumbRow struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Thumb string `gorm:"column:thumb"`
}

type MemberPathFullRow struct {
	ID         int    `gorm:"column:id"`
	Title      string `gorm:"column:title"`
	Thumb      string `gorm:"column:thumb"`
	Desc       string `gorm:"column:desc"`
	IsShow     int    `gorm:"column:is_show"`
	CategoryID int    `gorm:"column:category_id"`
}

type MemberPaperRow struct {
	ID        int    `gorm:"column:id"`
	Type      string `gorm:"column:type"`
	Score     int    `gorm:"column:score"`
	PassScore int    `gorm:"column:pass_score"`
	CourseID  *int   `gorm:"column:course_id"`
	ChapterID *int   `gorm:"column:chapter_id"`
}

type MemberIDTitle struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
}

type MemberIDName struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type adminMemberRepo struct{ db *gorm.DB }

func NewAdminMember(db *gorm.DB) AdminMemberRepository {
	return &adminMemberRepo{db: db}
}

func (r *adminMemberRepo) ListAreasByIDs(ctx context.Context, ids []int, withParent bool) ([]MemberAreaRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	cols := "id, name"
	if withParent {
		cols = "id, name, parent_name"
	}
	var rows []MemberAreaRow
	err := r.db.WithContext(ctx).Table("areas").Select(cols).Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) FillCityFromArea(ctx context.Context, areaID int) (*int, *string, error) {
	if areaID <= 0 {
		return nil, nil, nil
	}
	type areaRow struct {
		ID       int    `gorm:"column:id"`
		ParentID int    `gorm:"column:parent_id"`
		Code     string `gorm:"column:code"`
	}
	var area areaRow
	err := r.db.WithContext(ctx).Table("areas").Select("id, parent_id, code").Where("id = ?", areaID).Take(&area).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if area.ParentID > 0 {
		cityID := area.ParentID
		var parentCode string
		_ = r.db.WithContext(ctx).Table("areas").Select("code").Where("id = ?", area.ParentID).Scan(&parentCode).Error
		code := strings.ToLower(strings.TrimSpace(parentCode))
		if code == "" {
			code = strings.ToLower(strings.TrimSpace(area.Code))
		}
		return &cityID, &code, nil
	}
	cityID := area.ID
	code := strings.ToLower(strings.TrimSpace(area.Code))
	return &cityID, &code, nil
}

func (r *adminMemberRepo) TopCategoryBaseCredit(ctx context.Context, categoryID int) (int, bool, error) {
	if categoryID <= 0 {
		return 0, false, nil
	}
	var row struct {
		BaseCredit int `gorm:"column:base_credit"`
	}
	err := r.db.WithContext(ctx).Table("course_categories").
		Select("base_credit").
		Where("id = ? AND parent_id = 0", categoryID).
		Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	credit := row.BaseCredit
	if credit < 0 {
		credit = 0
	}
	return credit, true, nil
}

func (r *adminMemberRepo) ListUserCourses(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserCourseRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_course").Where("user_id = ?", userID)
	if title != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM courses c WHERE c.id = user_course.course_id AND c.deleted_at IS NULL AND c.title LIKE ?)",
			"%"+title+"%",
		)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberUserCourseRow
	err := query.Select("id, course_id, user_id, created_at, has_end_chapter, is_end").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListUserPaths(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserPathRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_path").Where("user_id = ? AND deleted_at IS NULL", userID)
	if title != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM paths p WHERE p.id = user_path.path_id AND p.deleted_at IS NULL AND p.title LIKE ?)",
			"%"+title+"%",
		)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberUserPathRow
	err := query.Select("id, path_id, created_at, has_end_course, is_end").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListUserChapters(ctx context.Context, userID int, title string, offset, limit int) ([]MemberUserChapterRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_course_chapter").Where("user_id = ?", userID)
	if title != "" {
		like := "%" + title + "%"
		query = query.Where(`(
			EXISTS (SELECT 1 FROM courses c WHERE c.id = user_course_chapter.course_id AND c.deleted_at IS NULL AND c.title LIKE ?)
			OR EXISTS (SELECT 1 FROM course_chapters ch WHERE ch.id = user_course_chapter.chapter_id AND ch.title LIKE ?)
		)`, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberUserChapterRow
	err := query.Select("id, course_id, chapter_id, user_id, created_at, is_end, un_lock").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListUserRecords(ctx context.Context, table string, userID, offset, limit int) ([]MemberRecordRow, int64, error) {
	query := r.db.WithContext(ctx).Table(table).Where("user_id = ? AND deleted_at IS NULL", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberRecordRow
	err := query.Select("id, user_id, field, `sum`, remark, created_at").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListUserPapers(ctx context.Context, userID int, paperType string, isPass, isPeopleCheck *int, offset, limit int) ([]MemberUserPaperRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_paper").Where("user_id = ? AND deleted_at IS NULL", userID)
	if paperType != "" {
		query = query.Where(
			"EXISTS (SELECT 1 FROM papers p WHERE p.id = user_paper.paper_id AND p.deleted_at IS NULL AND p.type = ?)",
			paperType,
		)
	}
	if isPass != nil {
		query = query.Where("is_pass = ?", *isPass)
	}
	if isPeopleCheck != nil {
		query = query.Where("is_people_check = ?", *isPeopleCheck)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberUserPaperRow
	err := query.Select("id, user_id, paper_id, score, submit_count, created_at, is_pass, is_people_check").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListUserCerts(ctx context.Context, userID, offset, limit int) ([]MemberUserCertRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_cert").Where("user_id = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []MemberUserCertRow
	err := query.Select("id, user_id, cert_id, path_id, course_id, cert_no, title, image, is_valid, issued_at, created_at").
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminMemberRepo) ListCoursesFull(ctx context.Context, ids []int) ([]MemberCourseFullRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberCourseFullRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, type, `desc`, is_show, category_id").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListCoursesBrief(ctx context.Context, ids []int) ([]MemberCourseBriefRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberCourseBriefRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, category_id").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListCoursesThumb(ctx context.Context, ids []int) ([]MemberCourseThumbRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberCourseThumbRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListPathsFull(ctx context.Context, ids []int) ([]MemberPathFullRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberPathFullRow
	err := r.db.WithContext(ctx).Table("paths").
		Select("id, title, thumb, `desc`, is_show, category_id").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListChaptersByIDs(ctx context.Context, ids []int) ([]MemberIDTitle, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberIDTitle
	err := r.db.WithContext(ctx).Table("course_chapters").Select("id, title").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListPapersByIDs(ctx context.Context, ids []int) ([]MemberPaperRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberPaperRow
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, type, score, pass_score, course_id, chapter_id").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListNamedByIDs(ctx context.Context, table string, ids []int) ([]MemberIDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberIDName
	err := r.db.WithContext(ctx).Table(table).Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListTitlesByIDs(ctx context.Context, table string, ids []int, softDelete bool) ([]MemberIDTitle, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Table(table).Select("id, title").Where("id IN ?", ids)
	if softDelete {
		q = q.Where("deleted_at IS NULL")
	}
	var rows []MemberIDTitle
	err := q.Find(&rows).Error
	return rows, err
}

func (r *adminMemberRepo) ListCertNames(ctx context.Context, ids []int) ([]MemberIDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []MemberIDName
	err := r.db.WithContext(ctx).Table("certs").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}
