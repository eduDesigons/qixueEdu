package adminops

import (
	"context"
	"errors"
	"strings"

	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/enroll"
)

var (
	ErrNotSuper = errors.New("仅超管可操作")
	ErrBadKind  = errors.New("请指定巡检类型")
)

type InspectIssue struct {
	enroll.Issue
	User  adminrepo.OpsUserBrief `json:"user"`
	Title string                 `json:"title"`
}

type InspectResult struct {
	Kind      string             `json:"kind"`
	Label     string             `json:"label"`
	Total     int                `json:"total"`
	Truncated bool               `json:"truncated"`
	Summary   []enroll.KindCount `json:"summary"`
	Issues    []InspectIssue     `json:"issues"`
}

func (s *Service) requireSuper(ctx context.Context, adminID int) error {
	var n int64
	err := s.db.WithContext(ctx).
		Table("administrator_role_relation AS rel").
		Joins("JOIN administrator_roles role ON role.id = rel.roleId").
		Where("rel.administratorId = ? AND role.code = ?", adminID, "SUPER_ADMIN").
		Count(&n).Error
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotSuper
	}
	return nil
}

func (s *Service) InspectKinds(ctx context.Context, adminID int) ([]enroll.KindCount, error) {
	if err := s.requireSuper(ctx, adminID); err != nil {
		return nil, err
	}
	return enroll.KindSteps(), nil
}

func (s *Service) Inspect(ctx context.Context, adminID int, kind string) (InspectResult, error) {
	if err := s.requireSuper(ctx, adminID); err != nil {
		return InspectResult{}, err
	}
	if strings.TrimSpace(kind) == "" {
		return InspectResult{}, ErrBadKind
	}
	raw, err := enroll.ScanKind(ctx, s.db, kind)
	if err != nil {
		if errors.Is(err, enroll.ErrUnknownKind) {
			return InspectResult{}, ErrBadKind
		}
		return InspectResult{}, err
	}
	issues, err := s.decorateIssues(ctx, raw.Issues)
	if err != nil {
		return InspectResult{}, err
	}
	return InspectResult{
		Kind:      raw.Kind,
		Label:     raw.Label,
		Total:     raw.Total,
		Truncated: raw.Truncated,
		Summary:   raw.Summary,
		Issues:    issues,
	}, nil
}

func (s *Service) Repair(ctx context.Context, adminID int, kind string) (enroll.RepairResult, error) {
	if err := s.requireSuper(ctx, adminID); err != nil {
		return enroll.RepairResult{}, err
	}
	if strings.TrimSpace(kind) == "" {
		return enroll.RepairResult{}, ErrBadKind
	}
	out, err := enroll.RepairKind(ctx, s.db, s.rdb, kind)
	if err != nil && errors.Is(err, enroll.ErrUnknownKind) {
		return enroll.RepairResult{}, ErrBadKind
	}
	return out, err
}

func (s *Service) decorateIssues(ctx context.Context, rows []enroll.Issue) ([]InspectIssue, error) {
	userIDs := make([]int, 0, len(rows))
	courseIDs := make([]int, 0, len(rows))
	pathIDs := make([]int, 0, len(rows))
	chapterIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
		if row.CourseID > 0 {
			courseIDs = append(courseIDs, row.CourseID)
		}
		if row.PathID > 0 {
			pathIDs = append(pathIDs, row.PathID)
		}
		if row.ChapterID > 0 {
			chapterIDs = append(chapterIDs, row.ChapterID)
		}
	}
	users, err := s.ops.LoadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return nil, err
	}
	courses, err := s.ops.LoadTitles(ctx, "courses", uniquePositive(courseIDs))
	if err != nil {
		return nil, err
	}
	paths, err := s.ops.LoadTitles(ctx, "paths", uniquePositive(pathIDs))
	if err != nil {
		return nil, err
	}
	chapters, err := s.ops.LoadTitles(ctx, "course_chapters", uniquePositive(chapterIDs))
	if err != nil {
		return nil, err
	}
	out := make([]InspectIssue, 0, len(rows))
	for _, row := range rows {
		title := ""
		if row.ChapterID > 0 {
			title = chapters[row.ChapterID].Title
		}
		if title == "" && row.CourseID > 0 {
			title = courses[row.CourseID].Title
		}
		if title == "" && row.PathID > 0 {
			title = paths[row.PathID].Title
		}
		out = append(out, InspectIssue{Issue: row, User: users[row.UserID], Title: title})
	}
	return out, nil
}
