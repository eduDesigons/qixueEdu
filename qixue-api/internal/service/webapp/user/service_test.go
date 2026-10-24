package user

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

type stubUsers struct {
	profile map[string]any
	err     error
	hash    string
	exists  bool
}

func (s stubUsers) FindByID(ctx context.Context, id int) (*model.User, error) {
	return &model.User{ID: uint(id)}, s.err
}
func (s stubUsers) FindProfile(ctx context.Context, id int) (map[string]any, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.profile, nil
}
func (s stubUsers) Exists(ctx context.Context, id int) (bool, error) { return s.exists, nil }
func (s stubUsers) Update(ctx context.Context, id int, fields map[string]any) error {
	return nil
}
func (s stubUsers) FindPassword(ctx context.Context, id int) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.hash, nil
}
func (s stubUsers) CreditCoins(ctx context.Context, id int) (int, int, error) { return 1, 2, nil }
func (s stubUsers) FindLoginByAccount(ctx context.Context, account string) (*repo.LoginUser, error) {
	return nil, gorm.ErrRecordNotFound
}

type stubLearner struct{}

func (stubLearner) CountUserCourses(ctx context.Context, uid int) (int64, error) { return 0, nil }
func (stubLearner) CountUnreadMessages(ctx context.Context, uid int) (int64, error) {
	return 0, nil
}
func (stubLearner) ListMessages(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	return 0, nil, nil
}
func (stubLearner) MarkMessagesRead(ctx context.Context, uid, id int) (int64, error) { return 0, nil }
func (stubLearner) CountLikeCourses(ctx context.Context, uid int) (int64, error)     { return 0, nil }
func (stubLearner) ListLikeCourses(ctx context.Context, uid, page, limit int) ([]repo.LikeCourseRow, error) {
	return nil, nil
}
func (stubLearner) CountLikeChapters(ctx context.Context, uid int) (int64, error) { return 0, nil }
func (stubLearner) ListLikeChapters(ctx context.Context, uid, page, limit int) ([]repo.LikeChapterRow, error) {
	return nil, nil
}
func (stubLearner) CountComments(ctx context.Context, table string, uid int) (int64, error) {
	return 0, nil
}
func (stubLearner) ListComments(ctx context.Context, table string, uid, page, limit int) ([]repo.CommentRow, error) {
	return nil, nil
}
func (stubLearner) ListCommentChildren(ctx context.Context, table string, pids []int) ([]map[string]any, error) {
	return nil, nil
}
func (stubLearner) ListUserBriefs(ctx context.Context, ids []int) ([]map[string]any, error) {
	return nil, nil
}
func (stubLearner) ListCourseBriefs(ctx context.Context, ids []int) ([]repo.CourseBrief, error) {
	return nil, nil
}
func (stubLearner) ListChapterBriefs(ctx context.Context, ids []int) ([]repo.ChapterBrief, error) {
	return nil, nil
}
func (stubLearner) CountUserPapers(ctx context.Context, uid int, paperType string) (int64, error) {
	return 0, nil
}
func (stubLearner) ListUserPapers(ctx context.Context, uid int, paperType string, page, limit int) ([]repo.UserPaperRow, error) {
	return nil, nil
}
func (stubLearner) CountRecords(ctx context.Context, table string, uid int) (int64, error) {
	return 0, nil
}
func (stubLearner) ListRecords(ctx context.Context, table string, uid, page, limit int) ([]map[string]any, error) {
	return nil, nil
}
func (stubLearner) CountUserPaths(ctx context.Context, uid int) (int64, error) { return 0, nil }
func (stubLearner) ListUserPaths(ctx context.Context, uid, page, limit int) ([]repo.UserPathRow, error) {
	return nil, nil
}
func (stubLearner) ListValidPathCerts(ctx context.Context, uid int, pathIDs []int) ([]map[string]any, error) {
	return nil, nil
}
func (stubLearner) CountValidCerts(ctx context.Context, uid int) (int64, error) { return 0, nil }
func (stubLearner) ListValidCerts(ctx context.Context, uid, page, limit int) ([]repo.UserCertRow, error) {
	return nil, nil
}

func TestInfo_NotFound(t *testing.T) {
	svc := New(stubUsers{err: gorm.ErrRecordNotFound}, stubLearner{}, nil)
	_, err := svc.Info(context.Background(), 1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v want ErrNotFound", err)
	}
}

func TestInfo_OK(t *testing.T) {
	svc := New(stubUsers{profile: map[string]any{"id": 1, "nick_name": "a"}}, stubLearner{}, nil)
	row, err := svc.Info(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if row["nick_name"] != "a" {
		t.Fatalf("unexpected %#v", row)
	}
}

func TestSetPwd_Required(t *testing.T) {
	svc := New(stubUsers{}, stubLearner{}, nil)
	err := svc.SetPwd(context.Background(), 1, "", "abcdef")
	if !errors.Is(err, ErrPwdRequired) {
		t.Fatalf("got %v", err)
	}
}
