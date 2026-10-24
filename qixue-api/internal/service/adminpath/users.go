package adminpath

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/enroll"
)

func (s *Service) GetUserPathList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	pathID := asInt(raw["path_id"])
	if pathID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	ctx := c.Request.Context()
	cityID, err := s.resolveCity(c, asInt(raw["city_id"]))
	if response.WriteErr(c, err) {
		return
	}

	rows, total, err := s.paths.ListPathUsers(ctx, repo.PathUserListFilter{
		PathID: pathID, CityID: cityID,
		NickName: strings.TrimSpace(asString(raw["nick_name"])),
		Mobile:   strings.TrimSpace(asString(raw["mobile"])),
		Offset:   (current - 1) * size, Limit: size,
	})
	if response.WriteErr(c, err) {
		return
	}

	userIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}
	users := map[int]gin.H{}
	certs := map[int]gin.H{}
	if len(userIDs) > 0 {
		userRows, err := s.paths.ListUsersFull(ctx, userIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, u := range userRows {
			users[u.ID] = gin.H{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
		certRows, err := s.paths.ListPathUserCerts(ctx, pathID, userIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, cert := range certRows {
			certs[cert.UserID] = gin.H{
				"id": cert.ID, "user_id": cert.UserID, "cert_id": cert.CertID, "cert_no": cert.CertNo,
				"url": cert.URL, "path_id": cert.PathID, "course_id": cert.CourseID,
				"title": cert.Title, "image": cert.Image, "is_valid": cert.IsValid,
				"issued_at":  model.FormatTime(cert.IssuedAt),
				"created_at": model.FormatTime(cert.CreatedAt), "updated_at": model.FormatTime(cert.UpdatedAt),
			}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var user, userCert any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := certs[row.UserID]; ok {
			userCert = found
		}
		list = append(list, gin.H{
			"id": row.ID, "user_id": row.UserID, "path_id": row.PathID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"is_end": row.IsEnd, "user": user, "user_cert": userCert,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) GetUserPathSearch(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	nick := strings.TrimSpace(asString(raw["nick_name"]))
	if nick == "" {
		response.OK(c, "", []any{})
		return
	}
	cityID, err := s.resolveCity(c, asInt(raw["city_id"]))
	if response.WriteErr(c, err) {
		return
	}
	rows, err := s.paths.SearchUsersByNick(c.Request.Context(), nick, cityID, 50)
	if response.WriteErr(c, err) {
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{
			"id": row.ID, "nick_name": row.NickName, "credit": row.Credit, "avatar": row.Avatar,
		})
	}
	response.OK(c, "", list)
}

func (s *Service) UserPathListSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	userList := toIntSlice(raw["user_list"])
	if len(userList) == 0 {
		response.Fail(c, 1, "请选择学员")
		return
	}
	pathID := asInt(raw["path_id"])
	if pathID <= 0 {
		response.Fail(c, 1, "请选择系列")
		return
	}
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}
	count, err := s.paths.CountEnrolled(ctx, pathID, userList)
	if response.WriteErr(c, err) {
		return
	}
	if count > 0 {
		response.Fail(c, 1, "存在用户已经入过该系列了")
		return
	}
	if response.WriteErr(c, enroll.EnrollUsersToPath(ctx, s.db, s.rdb, pathID, userList)) {
		return
	}
	response.OK(c, "保存成功", []any{})
}

func (s *Service) UserPathDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	id := bind.RequestInt(c, raw, "id")
	if id <= 0 {
		response.Fail(c, 1, "用户不存在")
		return
	}
	ctx := c.Request.Context()
	uid, pathID, err := s.paths.FindUserPath(ctx, id)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "用户不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var pathCount int64
	if err := s.db.WithContext(ctx).Model(&pathRow{}).Where("id = ?", pathID).Count(&pathCount).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if pathCount == 0 {
		response.Fail(c, 1, "系列不存在")
		return
	}

	courseIDs, err := s.paths.PluckPathCourseIDs(ctx, pathID)
	if response.WriteErr(c, err) {
		return
	}
	if response.WriteErr(c, s.paths.DeleteUserPathCascade(ctx, id, uid, pathID, courseIDs)) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) GetUserPathApplyList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	ctx := c.Request.Context()
	cityID, err := s.resolveCity(c, asInt(raw["city_id"]))
	if response.WriteErr(c, err) {
		return
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
	if response.WriteErr(c, err) {
		return
	}

	userIDs := map[int]struct{}{}
	pathIDs := map[int]struct{}{}
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		pathIDs[row.PathID] = struct{}{}
	}
	users := map[int]gin.H{}
	paths := map[int]gin.H{}
	if len(userIDs) > 0 {
		userRows, err := s.paths.ListUsersFull(ctx, keys(userIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, u := range userRows {
			users[u.ID] = gin.H{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
	}
	if len(pathIDs) > 0 {
		pathRows, err := s.paths.ListPathThumbs(ctx, keys(pathIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, p := range pathRows {
			paths[p.ID] = gin.H{"id": p.ID, "title": p.Title, "thumb": p.Thumb}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var user, path any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := paths[row.PathID]; ok {
			path = found
		}
		list = append(list, gin.H{
			"id": row.ID, "user_id": row.UserID, "path_id": row.PathID,
			"reason": row.Reason, "status": row.Status,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"user": user, "path": path,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) UserPathApplyStatusChange(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
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
		response.Fail(c, 1, "请选择申请")
		return
	}
	if status != 1 && status != 2 {
		response.Fail(c, 1, "状态无效")
		return
	}
	ctx := c.Request.Context()
	ok, skip := 0, 0
	for _, id := range ids {
		applied, err := s.applyUserPathStatus(ctx, id, status)
		if response.WriteErr(c, err) {
			return
		}
		if applied {
			ok++
		} else {
			skip++
		}
	}
	if ok == 0 {
		response.Fail(c, 1, "申请已处理或不存在")
		return
	}
	msg := "ok"
	if skip > 0 {
		msg = "已处理 " + strconv.Itoa(ok) + " 条，跳过 " + strconv.Itoa(skip) + " 条"
	}
	response.OK(c, msg, []any{})
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

func (s *Service) UserPathApplyDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	id := bind.RequestInt(c, raw, "id")
	if id <= 0 {
		response.Fail(c, 1, "申请不存在")
		return
	}
	ctx := c.Request.Context()
	count, err := s.paths.CountPathApply(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if count == 0 {
		response.Fail(c, 1, "申请不存在")
		return
	}
	if response.WriteErr(c, s.paths.DeletePathApply(ctx, id)) {
		return
	}
	response.OK(c, "ok", []any{})
}
