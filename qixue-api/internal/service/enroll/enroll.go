package enroll

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
)

// ForgetListCache clears cached course-list keys for a user.
func ForgetListCache(ctx context.Context, rdb *redis.Client, uid int) {
	if uid <= 0 {
		return
	}
	s := strconv.Itoa(uid)
	cache.Drop(ctx, rdb,
		"my_course_"+s+"_1_1_10",
		"my_course_"+s+"_2_1_10",
		"my_course_"+s+"_3_1_10",
		"learn_num_"+s,
		"user_my_count_"+s,
	)
}

func courses(db *gorm.DB) repo.CourseRepository {
	return repo.NewCourse(db)
}

// EnrollUsersToCourse inserts missing user_course rows and syncs chapter progress.
func EnrollUsersToCourse(ctx context.Context, db *gorm.DB, rdb *redis.Client, courseID int, userIDs []int) error {
	userIDs = normalizeIDs(userIDs)
	if courseID <= 0 || len(userIDs) == 0 {
		return nil
	}
	cr := courses(db)
	chapterCount, err := cr.CountChapters(ctx, courseID)
	if err != nil {
		return err
	}
	if chapterCount == 0 {
		return nil
	}

	existed, err := cr.ListUserIDsInCourseAmong(ctx, courseID, userIDs)
	if err != nil {
		return err
	}
	existedSet := map[int]struct{}{}
	for _, id := range existed {
		existedSet[id] = struct{}{}
	}
	now := time.Now()
	var rows []map[string]any
	for _, uid := range userIDs {
		if _, ok := existedSet[uid]; ok {
			continue
		}
		rows = append(rows, repo.NewUserCourseRow(uid, courseID, now))
	}
	if err := cr.CreateUserCoursesIgnore(ctx, rows); err != nil {
		return err
	}
	for _, uid := range userIDs {
		if _, err := EnsureUserCourseChapters(ctx, db, rdb, uid, courseID); err != nil {
			return err
		}
		ForgetListCache(ctx, rdb, uid)
	}
	return nil
}

// EnsureUserCourseChapters backfills missing chapter rows and syncs unlock/end state.
func EnsureUserCourseChapters(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int) (int, error) {
	if uid <= 0 || courseID <= 0 {
		return 0, nil
	}
	cr := courses(db)
	chapterIDs, err := cr.ListChapterIDs(ctx, courseID)
	if err != nil {
		return 0, err
	}
	if len(chapterIDs) == 0 {
		return 0, nil
	}
	existed, err := cr.ListUserChapterIDs(ctx, uid, courseID)
	if err != nil {
		return 0, err
	}
	existedSet := map[int]struct{}{}
	for _, id := range existed {
		existedSet[id] = struct{}{}
	}
	now := time.Now()
	var missing []map[string]any
	for _, chapterID := range chapterIDs {
		if _, ok := existedSet[chapterID]; ok {
			continue
		}
		missing = append(missing, repo.NewUserChapterRow(uid, courseID, chapterID, now))
	}
	for start := 0; start < len(missing); start += 400 {
		end := start + 400
		if end > len(missing) {
			end = len(missing)
		}
		if err := cr.CreateUserChaptersIgnore(ctx, missing[start:end]); err != nil {
			return 0, err
		}
	}
	if err := SyncUserCourseStatus(ctx, db, rdb, uid, courseID); err != nil {
		return 0, err
	}
	return len(missing), nil
}

// SyncUserCourseStatus realigns chapter unlock flags and course completion counts.
func SyncUserCourseStatus(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int) error {
	if uid <= 0 || courseID <= 0 {
		return nil
	}
	cr := courses(db)
	chapterIDs, err := cr.ListChapterIDs(ctx, courseID)
	if err != nil {
		return err
	}
	if len(chapterIDs) == 0 {
		return nil
	}
	rows, err := cr.ListUserChapterProgress(ctx, uid, courseID)
	if err != nil {
		return err
	}
	byChapter := map[int]repo.UserChapterProgress{}
	for _, row := range rows {
		byChapter[row.ChapterID] = row
	}

	prefixDone := true
	ended := 0
	var lockIDs, unlockIDs []int
	for _, chapterID := range chapterIDs {
		row, ok := byChapter[chapterID]
		if !ok {
			prefixDone = false
			continue
		}
		isEnd := row.IsEnd == 1
		watched := row.HasEndChapter == 1
		wantUnlock := 0
		if isEnd {
			ended++
			wantUnlock = 1
		} else if !prefixDone {
			wantUnlock = 0
		} else if watched {
			wantUnlock = 1
		} else {
			wantUnlock = row.UnLock
		}
		if row.UnLock != wantUnlock {
			if wantUnlock == 1 {
				unlockIDs = append(unlockIDs, row.ID)
			} else {
				lockIDs = append(lockIDs, row.ID)
			}
		}
		if !isEnd {
			prefixDone = false
		}
	}
	if err := cr.SetChapterUnlock(ctx, unlockIDs, 1); err != nil {
		return err
	}
	if err := cr.SetChapterUnlock(ctx, lockIDs, 0); err != nil {
		return err
	}

	userCourse, err := cr.FindUserCourseStatus(ctx, uid, courseID)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	wasEnd := userCourse.IsEnd == 1
	total := len(chapterIDs)
	isEnd := 0
	if wasEnd && ended == total {
		isEnd = 1
	}
	if userCourse.HasEndChapter != ended || userCourse.IsEnd != isEnd {
		if err := cr.UpdateUserCourseStatus(ctx, userCourse.ID, ended, isEnd); err != nil {
			return err
		}
		ForgetListCache(ctx, rdb, uid)
	}
	if wasEnd && isEnd == 0 {
		return reopenPaths(ctx, db, uid, courseID)
	}
	return nil
}

func reopenPaths(ctx context.Context, db *gorm.DB, uid, courseID int) error {
	cr := courses(db)
	pathIDs, err := cr.ListPathIDsByCourse(ctx, courseID)
	if err != nil {
		return err
	}
	for _, pathID := range pathIDs {
		userPath, err := cr.FindUserPathStatus(ctx, uid, pathID)
		if err == gorm.ErrRecordNotFound {
			continue
		}
		if err != nil {
			return err
		}
		courseIDs, err := cr.ListCourseIDsByPath(ctx, pathID)
		if err != nil {
			return err
		}
		total := len(courseIDs)
		ended := 0
		if total > 0 {
			ended, err = cr.CountEndedUserCourses(ctx, uid, courseIDs)
			if err != nil {
				return err
			}
		}
		isEnd := 0
		if total > 0 && ended >= total {
			isEnd = userPath.IsEnd
		}
		if userPath.HasEndCourse == ended && userPath.IsEnd == isEnd {
			continue
		}
		if err := cr.UpdateUserPathStatus(ctx, userPath.ID, ended, isEnd); err != nil {
			return err
		}
	}
	return nil
}

// SyncEnrolledChapters backfills new chapters for everyone already in the course.
func SyncEnrolledChapters(ctx context.Context, db *gorm.DB, rdb *redis.Client, courseID int) error {
	if courseID <= 0 {
		return nil
	}
	userIDs, err := courses(db).ListUserIDsInCourse(ctx, courseID)
	if err != nil {
		return err
	}
	for _, uid := range userIDs {
		if _, err := EnsureUserCourseChapters(ctx, db, rdb, uid, courseID); err != nil {
			return err
		}
		ForgetListCache(ctx, rdb, uid)
	}
	return nil
}

// EnrollUsersToPath inserts missing user_path rows and enrolls into all path courses.
func EnrollUsersToPath(ctx context.Context, db *gorm.DB, rdb *redis.Client, pathID int, userIDs []int) error {
	userIDs = normalizeIDs(userIDs)
	if pathID <= 0 || len(userIDs) == 0 {
		return nil
	}
	cr := courses(db)
	existed, err := cr.ListUserIDsOnPathAmong(ctx, pathID, userIDs)
	if err != nil {
		return err
	}
	existedSet := map[int]struct{}{}
	for _, id := range existed {
		existedSet[id] = struct{}{}
	}
	now := time.Now()
	var rows []map[string]any
	for _, uid := range userIDs {
		if _, ok := existedSet[uid]; ok {
			continue
		}
		rows = append(rows, repo.NewUserPathRow(uid, pathID, now))
	}
	if err := cr.CreateUserPathsIgnore(ctx, rows); err != nil {
		return err
	}

	courseIDs, err := cr.ListCourseIDsByPath(ctx, pathID)
	if err != nil {
		return err
	}
	for _, courseID := range courseIDs {
		if err := EnrollUsersToCourse(ctx, db, rdb, courseID, userIDs); err != nil {
			return err
		}
	}
	return nil
}

// CourseCreditAlreadyGranted reports whether graduation credit was already granted.
func CourseCreditAlreadyGranted(ctx context.Context, db *gorm.DB, uid int, remark string) (bool, error) {
	return courses(db).CreditGranted(ctx, uid, "课程发放", remark)
}

func normalizeIDs(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
