package admincourse

import (
	"context"
	"html"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/enroll"
)

func (s *Service) GetCourseUserList(ctx context.Context, adminID int, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	courseID := asInt(raw["course_id"])
	if courseID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return page.Body{}, err
	}

	rows, total, err := s.courses.ListCourseUsers(ctx, adminrepo.CourseUserListFilter{
		CourseID: courseID, CityID: cityID,
		NickName: strings.TrimSpace(asString(raw["nick_name"])),
		Mobile:   strings.TrimSpace(asString(raw["mobile"])),
		Offset:   (current - 1) * size, Limit: size,
	})
	if err != nil {
		return page.Body{}, err
	}

	userIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}
	users := map[int]map[string]any{}
	certs := map[int]map[string]any{}
	if len(userIDs) > 0 {
		userRows, err := s.courses.ListUsersFull(ctx, userIDs)
		if err != nil {
			return page.Body{}, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
		certRows, err := s.courses.ListCourseUserCerts(ctx, courseID, userIDs)
		if err != nil {
			return page.Body{}, err
		}
		for _, cert := range certRows {
			certs[cert.UserID] = map[string]any{
				"id": cert.ID, "user_id": cert.UserID, "cert_id": cert.CertID, "cert_no": cert.CertNo,
				"url": cert.URL, "path_id": cert.PathID, "course_id": cert.CourseID,
				"title": cert.Title, "image": cert.Image, "is_valid": cert.IsValid,
				"issued_at":  model.FormatTime(cert.IssuedAt),
				"created_at": model.FormatTime(cert.CreatedAt), "updated_at": model.FormatTime(cert.UpdatedAt),
			}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var user any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		var userCert any
		if found, ok := certs[row.UserID]; ok {
			userCert = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "user_id": row.UserID, "course_id": row.CourseID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"has_end_chapter": row.HasEndChapter, "is_end": row.IsEnd,
			"user": user, "user_cert": userCert,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetCourseUserSearch(ctx context.Context, adminID int, raw map[string]any) ([]map[string]any, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	nick := strings.TrimSpace(asString(raw["nick_name"]))
	if nick == "" {
		return []map[string]any{}, nil
	}
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return nil, err
	}
	rows, err := s.courses.SearchUsersByNick(ctx, nick, cityID, 50)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		list = append(list, map[string]any{
			"id": row.ID, "nick_name": row.NickName, "credit": row.Credit, "avatar": row.Avatar,
		})
	}
	return list, nil
}

func (s *Service) CourseUserListSave(ctx context.Context, raw map[string]any) error {
	if raw == nil {
		return biz("参数错误")
	}
	userList := toIntSlice(raw["user_list"])
	if len(userList) == 0 {
		return biz("请选择学员")
	}
	courseID := asInt(raw["course_id"])
	if courseID <= 0 {
		return biz("请选择课程")
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", courseID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return biz("课程不存在")
	}
	ok, err := s.courses.ChaptersExist(ctx, courseID)
	if err != nil {
		return err
	}
	if !ok {
		return biz("请先添加课节后再加入学员")
	}
	count, err = s.courses.CountEnrolled(ctx, courseID, userList)
	if err != nil {
		return err
	}
	if count > 0 {
		return biz("存在用户已经入过该班级了")
	}
	return enroll.EnrollUsersToCourse(ctx, s.db, s.rdb, courseID, userList)
}

func (s *Service) CourseUserDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("用户不存在")
	}
	uid, courseID, err := s.courses.FindUserCourse(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return biz("用户不存在")
	}
	if err != nil {
		return err
	}
	return s.courses.DeleteCourseUserCascade(ctx, id, uid, courseID)
}

func (s *Service) GetCourseUserApplyList(ctx context.Context, adminID int, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return page.Body{}, err
	}

	f := adminrepo.CourseApplyListFilter{
		CourseTitle: strings.TrimSpace(asString(raw["course_title"])),
		NickName:    strings.TrimSpace(asString(raw["nick_name"])),
		StartTime:   strings.TrimSpace(asString(raw["start_time"])),
		EndTime:     strings.TrimSpace(asString(raw["end_time"])),
		CityID:      cityID,
		Offset:      (current - 1) * size,
		Limit:       size,
	}
	if v, ok := raw["course_id"]; ok && asString(v) != "" {
		n := asInt(v)
		f.CourseID = &n
	}
	if v, ok := raw["user_id"]; ok && asString(v) != "" {
		n := asInt(v)
		f.UserID = &n
	}
	if v, ok := raw["status"]; ok && asString(v) != "" {
		n := asInt(v)
		f.Status = &n
	}

	rows, total, err := s.courses.ListCourseApplies(ctx, f)
	if err != nil {
		return page.Body{}, err
	}

	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		courseIDs[row.CourseID] = struct{}{}
	}
	users := map[int]map[string]any{}
	courses := map[int]map[string]any{}
	if len(userIDs) > 0 {
		userRows, err := s.courses.ListUsersFull(ctx, keys(userIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
	}
	if len(courseIDs) > 0 {
		courseRows, err := s.courses.ListCourseThumbs(ctx, keys(courseIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, co := range courseRows {
			courses[co.ID] = map[string]any{"id": co.ID, "title": co.Title, "thumb": co.Thumb}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var user, course any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := courses[row.CourseID]; ok {
			course = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "user_id": row.UserID, "course_id": row.CourseID,
			"reason": row.Reason, "status": row.Status,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"user": user, "course": course,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) CourseUserApplyStatusChange(ctx context.Context, raw map[string]any) (string, error) {
	if raw == nil {
		return "", biz("参数错误")
	}
	status := asInt(raw["status"])
	ids := toIntSlice(raw["ids"])
	if len(ids) == 0 {
		if id := asInt(raw["id"]); id > 0 {
			ids = []int{id}
		}
	}
	ids = uniquePositive(ids)
	if len(ids) == 0 {
		return "", biz("请选择申请")
	}
	if status != 1 && status != 2 {
		return "", biz("状态无效")
	}
	ok, skip := 0, 0
	for _, id := range ids {
		applied, err := s.applyCourseUserStatus(ctx, id, status)
		if err != nil {
			return "", err
		}
		if applied {
			ok++
		} else {
			skip++
		}
	}
	if ok == 0 {
		if status == 1 {
			courseIDs, err := s.courses.PluckPendingApplyCourseIDs(ctx, ids)
			if err != nil {
				return "", err
			}
			for _, courseID := range courseIDs {
				exists, err := s.courses.ChaptersExist(ctx, courseID)
				if err != nil {
					return "", err
				}
				if !exists {
					return "", biz("请先添加课节后再通过申请")
				}
			}
		}
		return "", biz("申请已处理或不存在")
	}
	msg := "ok"
	if skip > 0 {
		msg = "已处理 " + strconv.Itoa(ok) + " 条，跳过 " + strconv.Itoa(skip) + " 条"
	}
	return msg, nil
}

func (s *Service) applyCourseUserStatus(ctx context.Context, id, status int) (bool, error) {
	apply, err := s.courses.FindCourseApply(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if apply.Status != 0 {
		return false, nil
	}
	title, _ := s.courses.CourseTitle(ctx, apply.CourseID)
	if title == "" {
		title = "课程"
	}
	safeTitle := html.EscapeString(title)

	if status == 1 {
		exists, err := s.courses.ChaptersExist(ctx, apply.CourseID)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
		if err := enroll.EnrollUsersToCourse(ctx, s.db, s.rdb, apply.CourseID, []int{apply.UserID}); err != nil {
			return false, err
		}
		if _, err := enroll.EnsureUserCourseChapters(ctx, s.db, s.rdb, apply.UserID, apply.CourseID); err != nil {
			return false, err
		}
		enroll.ForgetListCache(ctx, s.rdb, apply.UserID)
		content := "<p>你申请加入的课程「" + safeTitle + "」已审核通过，可以开始学习了。</p>"
		if err := s.courses.SendUserMessage(ctx, apply.UserID, "入课申请已通过", content, "course", apply.CourseID); err != nil {
			return false, err
		}
	} else {
		content := "<p>很抱歉，你申请加入的课程「" + safeTitle + "」未通过审核。</p>"
		if err := s.courses.SendUserMessage(ctx, apply.UserID, "入课申请未通过", content, "course", apply.CourseID); err != nil {
			return false, err
		}
	}
	if err := s.courses.UpdateCourseApplyStatus(ctx, id, status); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) CourseUserApplyDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("申请不存在")
	}
	count, err := s.courses.CountCourseApply(ctx, id)
	if err != nil {
		return err
	}
	if count == 0 {
		return biz("申请不存在")
	}
	return s.courses.DeleteCourseApply(ctx, id)
}

func toIntSlice(v any) []int {
	switch list := v.(type) {
	case []any:
		out := make([]int, 0, len(list))
		for _, item := range list {
			if n := asInt(item); n > 0 {
				out = append(out, n)
			}
		}
		return out
	case []int:
		return list
	case []float64:
		out := make([]int, 0, len(list))
		for _, item := range list {
			if n := int(item); n > 0 {
				out = append(out, n)
			}
		}
		return out
	default:
		return nil
	}
}
