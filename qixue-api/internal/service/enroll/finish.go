package enroll

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/repo"
	"qixue-api/internal/service/cert"
	"qixue-api/internal/service/coin"
)

// GrantCourseEnd issues course cert/path sync and graduation credit/coins.
func GrantCourseEnd(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int, title string, credit int) {
	if uid <= 0 || courseID <= 0 {
		return
	}
	cert.NewIssuer(repo.NewCert(db)).OnCourseEnded(ctx, uid, courseID)
	if title == "" {
		title = "课程"
	}
	if credit > 0 {
		remark := title + "课程学习结束发放"
		granted, err := CourseCreditAlreadyGranted(ctx, db, uid, remark)
		if err == nil && !granted {
			_ = coin.AddCredit(ctx, db, uid, "课程发放", credit, remark)
		}
	}
	coin.OnCourseEnd(ctx, db, rdb, uid, courseID, title)
}

// FinishCourseIfReady marks a course ended when every chapter is done and there is no exam.
func FinishCourseIfReady(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int) {
	if uid <= 0 || courseID <= 0 {
		return
	}
	var info struct {
		Title   string `gorm:"column:title"`
		Credit  int    `gorm:"column:credit"`
		PaperID int    `gorm:"column:paper_id"`
	}
	err := db.WithContext(ctx).Table("courses").
		Select("title, credit, IFNULL(paper_id, 0) AS paper_id").
		Where("id = ? AND deleted_at IS NULL", courseID).Take(&info).Error
	if err != nil || info.PaperID > 0 {
		return
	}
	cr := courses(db)
	ids, err := cr.ListChapterIDs(ctx, courseID)
	if err != nil || len(ids) == 0 {
		return
	}
	st, err := cr.FindUserCourseStatus(ctx, uid, courseID)
	if err != nil || st.IsEnd == 1 {
		return
	}
	rows, err := cr.ListUserChapterProgress(ctx, uid, courseID)
	if err != nil {
		return
	}
	ended := 0
	for _, row := range rows {
		if row.IsEnd == 1 {
			ended++
		}
	}
	if ended < len(ids) {
		return
	}
	if err := cr.UpdateUserCourseStatus(ctx, st.ID, ended, 1); err != nil {
		return
	}
	ForgetListCache(ctx, rdb, uid)
	GrantCourseEnd(ctx, db, rdb, uid, courseID, info.Title, info.Credit)
}
