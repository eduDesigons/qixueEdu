package cert

import (
	"context"
	"crypto/rand"
	"fmt"
	"html"
	"math/big"
	"strings"
	"time"

	"qixue-api/internal/repo"
	"qixue-api/internal/service/bunny"
)

// Issuer issues certificates and renders images asynchronously.
type Issuer struct {
	certs   repo.CertRepository
	storage *bunny.Storage
}

func NewIssuer(certs repo.CertRepository) *Issuer {
	return &Issuer{certs: certs, storage: bunny.NewStorage(bunny.LoadConfig())}
}

// UserCert is the issued certificate row (alias of repo type for callers).
type UserCert = repo.UserCert

type templateMeta struct {
	ID     int
	Name   string
	IsShow int
}

// IssueForCourse issues a course cert record and schedules async image render.
func (s *Issuer) IssueForCourse(ctx context.Context, userID, courseID int, force bool) (*UserCert, error) {
	if userID <= 0 || courseID <= 0 {
		return nil, nil
	}
	course, err := s.certs.GetCourseBrief(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, nil
	}
	if course.CertID == nil || *course.CertID <= 0 {
		return nil, nil
	}
	tpl, err := s.loadTemplate(ctx, *course.CertID)
	if err != nil || tpl == nil {
		return nil, err
	}

	existing, err := s.findUserCert(ctx, userID, 0, courseID)
	if err != nil {
		return nil, err
	}
	if existing != nil && !force {
		if strings.TrimSpace(existing.Image) == "" {
			s.scheduleRender(existing.ID)
		}
		return existing, nil
	}
	if !force && tpl.IsShow != 1 {
		return nil, nil
	}
	if !force {
		ended, err := s.certs.GetUserCourseEnded(ctx, userID, courseID)
		if err != nil {
			return nil, err
		}
		if ended != 1 {
			return nil, nil
		}
	}

	ok, err := s.certs.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	certNo := GenerateNo(userID, courseID)
	if existing != nil && strings.TrimSpace(existing.CertNo) != "" {
		certNo = existing.CertNo
	}
	title := tpl.Name
	if title == "" {
		title = course.Title + "结业证书"
	}
	now := time.Now()
	payload := map[string]any{
		"user_id": userID, "cert_id": *course.CertID, "path_id": 0, "course_id": courseID,
		"cert_no": certNo, "title": title, "is_valid": 1, "issued_at": now, "updated_at": now,
	}

	if existing != nil {
		if err := s.certs.UpdateUserCert(ctx, existing.ID, payload); err != nil {
			return nil, err
		}
		record, err := s.findUserCert(ctx, userID, 0, courseID)
		if err != nil {
			return nil, err
		}
		if record != nil {
			s.scheduleRender(record.ID)
		}
		return record, nil
	}

	payload["image"] = ""
	payload["created_at"] = now
	row := UserCert{
		UserID: userID, CertID: *course.CertID, PathID: 0, CourseID: courseID,
		CertNo: certNo, Title: title, Image: "", IsValid: 1, IssuedAt: &now,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.certs.CreateUserCert(ctx, &row); err != nil {
		if again, findErr := s.findUserCert(ctx, userID, 0, courseID); findErr == nil && again != nil {
			s.scheduleRender(again.ID)
			return again, nil
		}
		return nil, err
	}
	_ = s.notify(ctx, userID, course.Title, row.ID, "课程")
	s.scheduleRender(row.ID)
	return &row, nil
}

// IssueForPath issues a path cert record and schedules async image render.
func (s *Issuer) IssueForPath(ctx context.Context, userID, pathID int, force bool) (*UserCert, error) {
	if userID <= 0 || pathID <= 0 {
		return nil, nil
	}
	path, err := s.certs.GetPathBrief(ctx, pathID)
	if err != nil {
		return nil, err
	}
	if path == nil {
		return nil, nil
	}
	if path.CertID == nil || *path.CertID <= 0 {
		return nil, nil
	}
	tpl, err := s.loadTemplate(ctx, *path.CertID)
	if err != nil || tpl == nil {
		return nil, err
	}

	existing, err := s.findUserCert(ctx, userID, pathID, 0)
	if err != nil {
		return nil, err
	}
	if existing != nil && !force {
		if strings.TrimSpace(existing.Image) == "" {
			s.scheduleRender(existing.ID)
		}
		return existing, nil
	}
	if !force && tpl.IsShow != 1 {
		return nil, nil
	}
	if !force {
		ended, err := s.certs.GetUserPathEnded(ctx, userID, pathID)
		if err != nil {
			return nil, err
		}
		if ended != 1 {
			return nil, nil
		}
	}

	ok, err := s.certs.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	certNo := GenerateNo(userID, pathID)
	if existing != nil && strings.TrimSpace(existing.CertNo) != "" {
		certNo = existing.CertNo
	}
	title := tpl.Name
	if title == "" {
		title = path.Title + "结业证书"
	}
	now := time.Now()
	payload := map[string]any{
		"user_id": userID, "cert_id": *path.CertID, "path_id": pathID, "course_id": 0,
		"cert_no": certNo, "title": title, "is_valid": 1, "issued_at": now, "updated_at": now,
	}

	if existing != nil {
		if err := s.certs.UpdateUserCert(ctx, existing.ID, payload); err != nil {
			return nil, err
		}
		record, err := s.findUserCert(ctx, userID, pathID, 0)
		if err != nil {
			return nil, err
		}
		if record != nil {
			s.scheduleRender(record.ID)
		}
		return record, nil
	}

	payload["image"] = ""
	payload["created_at"] = now
	row := UserCert{
		UserID: userID, CertID: *path.CertID, PathID: pathID, CourseID: 0,
		CertNo: certNo, Title: title, Image: "", IsValid: 1, IssuedAt: &now,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.certs.CreateUserCert(ctx, &row); err != nil {
		if again, findErr := s.findUserCert(ctx, userID, pathID, 0); findErr == nil && again != nil {
			s.scheduleRender(again.ID)
			return again, nil
		}
		return nil, err
	}
	_ = s.notify(ctx, userID, path.Title, row.ID, "系列")
	s.scheduleRender(row.ID)
	return &row, nil
}

// BackfillForCourse issues certs for ended course users without a valid cert.
func (s *Issuer) BackfillForCourse(ctx context.Context, courseID int, limit int) (int, error) {
	if courseID <= 0 {
		return 0, nil
	}
	if limit <= 0 {
		limit = 500
	}
	certID, err := s.certs.GetCourseCertID(ctx, courseID)
	if err != nil {
		return 0, err
	}
	if certID == nil || *certID <= 0 {
		return 0, nil
	}
	userIDs, err := s.certs.ListEndedCourseUserIDs(ctx, courseID, limit)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, uid := range userIDs {
		existing, err := s.findUserCert(ctx, uid, 0, courseID)
		if err != nil {
			continue
		}
		if existing != nil && existing.IsValid == 1 {
			continue
		}
		force := existing != nil
		record, err := s.IssueForCourse(ctx, uid, courseID, force)
		if err != nil || record == nil {
			continue
		}
		n++
	}
	return n, nil
}

// BackfillForPath issues certs for ended path users without a valid cert.
func (s *Issuer) BackfillForPath(ctx context.Context, pathID int, limit int) (int, error) {
	if pathID <= 0 {
		return 0, nil
	}
	if limit <= 0 {
		limit = 500
	}
	certID, err := s.certs.GetPathCertID(ctx, pathID)
	if err != nil {
		return 0, err
	}
	if certID == nil || *certID <= 0 {
		return 0, nil
	}
	userIDs, err := s.certs.ListEndedPathUserIDs(ctx, pathID, limit)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, uid := range userIDs {
		existing, err := s.findUserCert(ctx, uid, pathID, 0)
		if err != nil {
			continue
		}
		if existing != nil && existing.IsValid == 1 {
			continue
		}
		force := existing != nil
		record, err := s.IssueForPath(ctx, uid, pathID, force)
		if err != nil || record == nil {
			continue
		}
		n++
	}
	return n, nil
}

func (s *Issuer) loadTemplate(ctx context.Context, certID int) (*templateMeta, error) {
	tpl, err := s.certs.LoadTemplateMeta(ctx, certID)
	if err != nil || tpl == nil {
		return nil, err
	}
	return &templateMeta{ID: tpl.ID, Name: tpl.Name, IsShow: tpl.IsShow}, nil
}

func (s *Issuer) findUserCert(ctx context.Context, userID, pathID, courseID int) (*UserCert, error) {
	return s.certs.FindUserCert(ctx, userID, pathID, courseID)
}

func (s *Issuer) notify(ctx context.Context, userID int, subjectTitle string, certID int, kind string) error {
	if userID <= 0 {
		return nil
	}
	content := "<p>同学，你已完成" + html.EscapeString(kind) +
		"「" + html.EscapeString(subjectTitle) +
		"」，结业证书已发放，请到「我的证书」查看。</p>"
	return s.certs.CreateCertNotice(ctx, userID, certID, "恭喜获得结业证书", content, time.Now())
}

// OnCourseEnded issues the course cert and syncs related path progress
func (s *Issuer) OnCourseEnded(ctx context.Context, userID, courseID int) {
	if userID <= 0 || courseID <= 0 {
		return
	}
	_, _ = s.IssueForCourse(ctx, userID, courseID, false)
	pathIDs, err := s.certs.ListPathIDsByCourse(ctx, courseID)
	if err != nil {
		return
	}
	for _, pathID := range pathIDs {
		_ = s.SyncPathProgress(ctx, userID, pathID)
	}
}

// SyncPathProgress updates path completion and issues path cert when fully ended.
func (s *Issuer) SyncPathProgress(ctx context.Context, userID, pathID int) error {
	if userID <= 0 || pathID <= 0 {
		return nil
	}
	userPath, err := s.certs.GetUserPathProgress(ctx, userID, pathID)
	if err != nil {
		return err
	}
	if userPath == nil {
		return nil
	}
	courseIDs, err := s.certs.ListActivePathCourseIDs(ctx, pathID)
	if err != nil {
		return err
	}
	total := len(courseIDs)
	ended := 0
	if total > 0 {
		count, err := s.certs.CountEndedUserCourses(ctx, userID, courseIDs)
		if err != nil {
			return err
		}
		ended = int(count)
	}
	isEnd := userPath.IsEnd
	if total > 0 && ended >= total {
		isEnd = 1
	}
	if userPath.HasEndCourse != ended || userPath.IsEnd != isEnd {
		if err := s.certs.UpdateUserPathProgress(ctx, userPath.ID, ended, isEnd); err != nil {
			return err
		}
	}
	if isEnd == 1 {
		_, _ = s.IssueForPath(ctx, userID, pathID, false)
	}
	return nil
}

// GenerateNo builds a certificate serial number.
func GenerateNo(userID, refID int) string {
	n, err := rand.Int(rand.Reader, big.NewInt(90))
	suffix := 10
	if err == nil {
		suffix = int(n.Int64()) + 10
	}
	return fmt.Sprintf("QX%s%04d%04d%02d",
		time.Now().Format("20060102"),
		userID%10000,
		refID%10000,
		suffix,
	)
}
