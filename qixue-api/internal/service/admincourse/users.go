package admincourse

import (
	"context"
	"html"
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

func (s *Service) GetCourseUserList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	courseID := asInt(raw["course_id"])
	if courseID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	ctx := c.Request.Context()
	cityID, err := s.resolveCity(c, asInt(raw["city_id"]))
	if response.WriteErr(c, err) {
		return
	}

	rows, total, err := s.courses.ListCourseUsers(ctx, repo.CourseUserListFilter{
		CourseID: courseID, CityID: cityID,
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
		userRows, err := s.courses.ListUsersFull(ctx, userIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, u := range userRows {
			users[u.ID] = gin.H{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
		certRows, err := s.courses.ListCourseUserCerts(ctx, courseID, userIDs)
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
		var user any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		var userCert any
		if found, ok := certs[row.UserID]; ok {
			userCert = found
		}
		list = append(list, gin.H{
			"id": row.ID, "user_id": row.UserID, "course_id": row.CourseID,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"has_end_chapter": row.HasEndChapter, "is_end": row.IsEnd,
			"user": user, "user_cert": userCert,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) GetCourseUserSearch(c *gin.Context) {
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
	rows, err := s.courses.SearchUsersByNick(c.Request.Context(), nick, cityID, 50)
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

func (s *Service) CourseUserListSave(c *gin.Context) {
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
	courseID := asInt(raw["course_id"])
	if courseID <= 0 {
		response.Fail(c, 1, "请选择课程")
		return
	}
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&courseRow{}).Where("id = ?", courseID).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "课程不存在")
		return
	}
	ok, err := s.courses.ChaptersExist(ctx, courseID)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.Fail(c, 1, "请先添加课节后再加入学员")
		return
	}
	count, err = s.courses.CountEnrolled(ctx, courseID, userList)
	if response.WriteErr(c, err) {
		return
	}
	if count > 0 {
		response.Fail(c, 1, "存在用户已经入过该班级了")
		return
	}
	if response.WriteErr(c, enroll.EnrollUsersToCourse(ctx, s.db, s.rdb, courseID, userList)) {
		return
	}
	response.OK(c, "保存成功", []any{})
}

func (s *Service) CourseUserDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	id := bind.RequestInt(c, raw, "id")
	if id <= 0 {
		response.Fail(c, 1, "用户不存在")
		return
	}
	ctx := c.Request.Context()
	uid, courseID, err := s.courses.FindUserCourse(ctx, id)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "用户不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if response.WriteErr(c, s.courses.DeleteCourseUserCascade(ctx, id, uid, courseID)) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) GetCourseUserApplyList(c *gin.Context) {
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

	f := repo.CourseApplyListFilter{
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
	if response.WriteErr(c, err) {
		return
	}

	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		courseIDs[row.CourseID] = struct{}{}
	}
	users := map[int]gin.H{}
	courses := map[int]gin.H{}
	if len(userIDs) > 0 {
		userRows, err := s.courses.ListUsersFull(ctx, keys(userIDs))
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
	if len(courseIDs) > 0 {
		courseRows, err := s.courses.ListCourseThumbs(ctx, keys(courseIDs))
		if response.WriteErr(c, err) {
			return
		}
		for _, co := range courseRows {
			courses[co.ID] = gin.H{"id": co.ID, "title": co.Title, "thumb": co.Thumb}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var user, course any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		if found, ok := courses[row.CourseID]; ok {
			course = found
		}
		list = append(list, gin.H{
			"id": row.ID, "user_id": row.UserID, "course_id": row.CourseID,
			"reason": row.Reason, "status": row.Status,
			"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
			"user": user, "course": course,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) CourseUserApplyStatusChange(c *gin.Context) {
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
		applied, err := s.applyCourseUserStatus(ctx, id, status)
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
		if status == 1 {
			courseIDs, err := s.courses.PluckPendingApplyCourseIDs(ctx, ids)
			if response.WriteErr(c, err) {
				return
			}
			for _, courseID := range courseIDs {
				exists, err := s.courses.ChaptersExist(ctx, courseID)
				if response.WriteErr(c, err) {
					return
				}
				if !exists {
					response.Fail(c, 1, "请先添加课节后再通过申请")
					return
				}
			}
		}
		response.Fail(c, 1, "申请已处理或不存在")
		return
	}
	msg := "ok"
	if skip > 0 {
		msg = "已处理 " + strconv.Itoa(ok) + " 条，跳过 " + strconv.Itoa(skip) + " 条"
	}
	response.OK(c, msg, []any{})
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

func (s *Service) CourseUserApplyDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	id := bind.RequestInt(c, raw, "id")
	if id <= 0 {
		response.Fail(c, 1, "申请不存在")
		return
	}
	ctx := c.Request.Context()
	count, err := s.courses.CountCourseApply(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if count == 0 {
		response.Fail(c, 1, "申请不存在")
		return
	}
	if response.WriteErr(c, s.courses.DeleteCourseApply(ctx, id)) {
		return
	}
	response.OK(c, "ok", []any{})
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
