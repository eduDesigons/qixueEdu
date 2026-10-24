package adminpath

import (
	"context"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/enroll"
)

func (s *Service) GetUserPathList(ctx context.Context, adminID int, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	pathID := asInt(raw["path_id"])
	if pathID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return page.Body{}, err
	}

	rows, total, err := s.paths.ListPathUsers(ctx, repo.PathUserListFilter{
		PathID: pathID, CityID: cityID,
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
		userRows, err := s.paths.ListUsersFull(ctx, userIDs)
		if err != nil {
			return page.Body{}, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
		certRows, err := s.paths.ListPathUserCerts(ctx, pathID, userIDs)
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
		var user, userCert any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := certs[row.UserID]; ok {
			userCert = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "user_id": row.UserID, "path_id": row.PathID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"is_end": row.IsEnd, "user": user, "user_cert": userCert,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) GetUserPathSearch(ctx context.Context, adminID int, raw map[string]any) ([]map[string]any, error) {
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
	rows, err := s.paths.SearchUsersByNick(ctx, nick, cityID, 50)
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

func (s *Service) UserPathListSave(ctx context.Context, raw map[string]any) error {
	if raw == nil {
		return biz("参数错误")
	}
	userList := toIntSlice(raw["user_list"])
	if len(userList) == 0 {
		return biz("请选择学员")
	}
	pathID := asInt(raw["path_id"])
	if pathID <= 0 {
		return biz("请选择系列")
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return biz("系列不存在")
	}
	count, err := s.paths.CountEnrolled(ctx, pathID, userList)
	if err != nil {
		return err
	}
	if count > 0 {
		return biz("存在用户已经入过该系列了")
	}
	return enroll.EnrollUsersToPath(ctx, s.db, s.rdb, pathID, userList)
}

func (s *Service) UserPathDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("用户不存在")
	}
	uid, pathID, err := s.paths.FindUserPath(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return biz("用户不存在")
	}
	if err != nil {
		return err
	}
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		return err
	}
	if pathCount == 0 {
		return biz("系列不存在")
	}

	courseIDs, err := s.paths.PluckPathCourseIDs(ctx, pathID)
	if err != nil {
		return err
	}
	return s.paths.DeleteUserPathCascade(ctx, id, uid, pathID, courseIDs)
}

func (s *Service) GetUserPathApplyList(ctx context.Context, adminID int, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return page.Body{}, err
	}

	f := repo.PathApplyListFilter{
		PathTitle: strings.TrimSpace(asString(raw["path_title"])),
		NickName:  strings.TrimSpace(asString(raw["nick_name"])),
		StartTime: strings.TrimSpace(asString(raw["start_time"])),
		EndTime:   strings.TrimSpace(asString(raw["end_time"])),
		CityID:    cityID,
		Offset:    (current - 1) * size,
		Limit:     size,
	}
	if v, ok := raw["path_id"]; ok && asString(v) != "" {
		n := asInt(v)
		f.PathID = &n
	}
	if v, ok := raw["user_id"]; ok && asString(v) != "" {
		n := asInt(v)
		f.UserID = &n
	}
	if v, ok := raw["status"]; ok && asString(v) != "" {
		n := asInt(v)
		f.Status = &n
	}

	rows, total, err := s.paths.ListPathApplies(ctx, f)
	if err != nil {
		return page.Body{}, err
	}

	userIDs := map[int]struct{}{}
	pathIDs := map[int]struct{}{}
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		pathIDs[row.PathID] = struct{}{}
	}
	users := map[int]map[string]any{}
	paths := map[int]map[string]any{}
	if len(userIDs) > 0 {
		userRows, err := s.paths.ListUsersFull(ctx, keys(userIDs))
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
	if len(pathIDs) > 0 {
		pathRows, err := s.paths.ListPathThumbs(ctx, keys(pathIDs))
		if err != nil {
			return page.Body{}, err
		}
		for _, p := range pathRows {
			paths[p.ID] = map[string]any{"id": p.ID, "title": p.Title, "thumb": p.Thumb}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var user, path any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := paths[row.PathID]; ok {
			path = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "user_id": row.UserID, "path_id": row.PathID,
			"reason": row.Reason, "status": row.Status,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"user": user, "path": path,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) UserPathApplyStatusChange(ctx context.Context, raw map[string]any) (string, error) {
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
		applied, err := s.applyUserPathStatus(ctx, id, status)
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
		return "", biz("申请已处理或不存在")
	}
	msg := "ok"
	if skip > 0 {
		msg = "已处理 " + strconv.Itoa(ok) + " 条，跳过 " + strconv.Itoa(skip) + " 条"
	}
	return msg, nil
}

func (s *Service) applyUserPathStatus(ctx context.Context, id, status int) (bool, error) {
	apply, err := s.paths.FindPathApply(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if apply.Status != 0 {
		return false, nil
	}
	title, _ := s.paths.PathTitle(ctx, apply.PathID)
	if title == "" {
		title = "系列课程"
	}
	safeTitle := escapeTitle(title)

	if status == 1 {
		if err := enroll.EnrollUsersToPath(ctx, s.db, s.rdb, apply.PathID, []int{apply.UserID}); err != nil {
			return false, err
		}
		enroll.ForgetListCache(ctx, s.rdb, apply.UserID)
		content := "<p>你申请加入的系列「" + safeTitle + "」已审核通过，可以开始学习了。</p>"
		if err := s.paths.SendUserMessage(ctx, apply.UserID, "入系申请已通过", content, "path", apply.PathID); err != nil {
			return false, err
		}
	} else {
		content := "<p>很抱歉，你申请加入的系列「" + safeTitle + "」未通过审核。</p>"
		if err := s.paths.SendUserMessage(ctx, apply.UserID, "入系申请未通过", content, "path", apply.PathID); err != nil {
			return false, err
		}
	}
	if err := s.paths.UpdatePathApplyStatus(ctx, id, status); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) UserPathApplyDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return biz("申请不存在")
	}
	count, err := s.paths.CountPathApply(ctx, id)
	if err != nil {
		return err
	}
	if count == 0 {
		return biz("申请不存在")
	}
	return s.paths.DeletePathApply(ctx, id)
}
