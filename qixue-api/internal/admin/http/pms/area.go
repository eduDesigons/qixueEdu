package pms

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/platform/cache"
)

var codePattern = regexp.MustCompile(`^[a-z0-9]{1,30}$`)

type Area struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewArea(db *gorm.DB, rdb *redis.Client) *Area {
	return &Area{db: db, rdb: rdb}
}

type record struct {
	ID         int            `gorm:"column:id;primaryKey"`
	Name       string         `gorm:"column:name"`
	ParentID   int            `gorm:"column:parent_id"`
	ParentName string         `gorm:"column:parent_name"`
	Code       string         `gorm:"column:code"`
	Type       int            `gorm:"column:type"`
	Sort       int            `gorm:"column:sort"`
	IsShow     int            `gorm:"column:is_show"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (record) TableName() string { return "areas" }

type node struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ParentID   int    `json:"parent_id"`
	ParentName string `json:"parent_name"`
	Code       string `json:"code"`
	Sort       int    `json:"sort"`
	IsShow     int    `json:"is_show"`
	UserCount  int    `json:"user_count"`
	Children   []node `json:"children,omitempty"`
}

func (h *Area) List(c *gin.Context) {
	ctx := c.Request.Context()
	var rows []record
	if err := h.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	areaUsers := map[int]int{}
	cityUsers := map[int]int{}
	type countRow struct {
		ID    int `gorm:"column:aid"`
		Total int `gorm:"column:total"`
	}
	var areaCounts []countRow
	_ = h.db.WithContext(ctx).Raw("SELECT area_id AS aid, COUNT(*) AS total FROM users WHERE area_id > 0 GROUP BY area_id").Scan(&areaCounts)
	for _, item := range areaCounts {
		areaUsers[item.ID] = item.Total
	}
	var cityCounts []countRow
	_ = h.db.WithContext(ctx).Raw("SELECT city_id AS aid, COUNT(*) AS total FROM users WHERE city_id > 0 GROUP BY city_id").Scan(&cityCounts)
	for _, item := range cityCounts {
		cityUsers[item.ID] = item.Total
	}

	parents := map[int]*node{}
	var order []int
	var children []node
	for _, row := range rows {
		item := node{
			ID: row.ID, Name: row.Name, ParentID: row.ParentID, ParentName: row.ParentName,
			Code: row.Code, Sort: row.Sort, IsShow: showBit(row.IsShow),
			UserCount: areaUsers[row.ID] + cityUsers[row.ID],
		}
		if row.ParentID == 0 {
			item.Children = []node{}
			copied := item
			parents[row.ID] = &copied
			order = append(order, row.ID)
			continue
		}
		children = append(children, item)
	}
	for _, child := range children {
		if parent, ok := parents[child.ParentID]; ok {
			parent.Children = append(parent.Children, child)
		}
	}
	list := make([]node, 0, len(order))
	for _, id := range order {
		list = append(list, *parents[id])
	}
	response.OK(c, "", list)
}

func (h *Area) Save(c *gin.Context) {
	var body struct {
		ID       bind.Int `json:"id"`
		Name     string   `json:"name"`
		ParentID bind.Int `json:"parent_id"`
		Sort     bind.Int `json:"sort"`
		IsShow   any      `json:"is_show"`
		Code     string   `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "请填写名称")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		response.Fail(c, 1, "请填写名称")
		return
	}
	if utf8.RuneCountInString(name) > 32 {
		response.Fail(c, 1, "名称不能超过 32 个字")
		return
	}
	isShow := 1
	if bind.Present(body.IsShow) && bind.AsInt(body.IsShow) == 0 {
		isShow = 0
	}
	ctx := c.Request.Context()
	parentName := ""
	code := strings.ToLower(strings.TrimSpace(body.Code))
	if body.ParentID.Int() > 0 {
		var parent record
		err := h.db.WithContext(ctx).Where("id = ? AND parent_id = 0", body.ParentID.Int()).Take(&parent).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "请选择一级区域")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		parentName = parent.Name
		code = strings.ToLower(strings.TrimSpace(parent.Code))
	}
	if !codePattern.MatchString(code) {
		if body.ParentID.Int() > 0 {
			response.Fail(c, 1, "请先给一级区域填写编码")
		} else {
			response.Fail(c, 1, "编码只能是字母和数字")
		}
		return
	}

	err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if body.ID.Int() > 0 {
			var row record
			if err := tx.Where("id = ?", body.ID.Int()).Take(&row).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return response.BizErr("区域不存在")
				}
				return err
			}
			if row.ParentID == 0 && body.ParentID.Int() > 0 {
				var children int64
				if err := tx.Model(&record{}).Where("parent_id = ?", body.ID.Int()).Count(&children).Error; err != nil {
					return err
				}
				if children > 0 {
					return response.BizErr("下面还有二级区域，不能改成二级")
				}
			}
			if err := tx.Model(&record{}).Where("id = ?", body.ID.Int()).Updates(map[string]any{
				"name": name, "code": code, "parent_id": body.ParentID.Int(), "parent_name": parentName,
				"type": areaType(body.ParentID.Int()), "sort": body.Sort.Int(), "is_show": isShow,
			}).Error; err != nil {
				return err
			}
			if body.ParentID.Int() == 0 {
				if err := tx.Model(&record{}).Where("parent_id = ?", body.ID.Int()).Updates(map[string]any{
					"parent_name": name, "code": code,
				}).Error; err != nil {
					return err
				}
				return tx.Exec("UPDATE users SET city_code = ? WHERE city_id = ?", code, body.ID.Int()).Error
			}
			return tx.Exec("UPDATE users SET city_id = ?, city_code = ? WHERE area_id = ?", body.ParentID.Int(), code, body.ID.Int()).Error
		}
		row := record{
			Name: name, Code: code, ParentID: body.ParentID.Int(), ParentName: parentName,
			Type: areaType(body.ParentID.Int()), Sort: body.Sort.Int(), IsShow: isShow,
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		body.ID = bind.Int(row.ID)
		return nil
	})
	if response.WriteErr(c, err) {
		return
	}
	if h.rdb != nil {
		cache.Drop(c.Request.Context(), h.rdb, "area_site_codes")
	}
	response.OK(c, "保存成功", gin.H{"id": body.ID.Int()})
}

func (h *Area) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	var row record
	err := h.db.WithContext(ctx).Where("id = ?", body.ID.Int()).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "区域不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var children int64
	if err := h.db.WithContext(ctx).Model(&record{}).Where("parent_id = ?", body.ID.Int()).Count(&children).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if children > 0 {
		response.Fail(c, 1, "请先删除下面的二级区域")
		return
	}
	var used int64
	if err := h.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE area_id = ? OR city_id = ?", body.ID.Int(), body.ID.Int()).Scan(&used).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	var admins int64
	if err := h.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM administrator WHERE city_id = ?", body.ID.Int()).Scan(&admins).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if used > 0 || admins > 0 {
		response.Fail(c, 1, "该区域已被使用，不能删除")
		return
	}
	if response.WriteErr(c, h.db.WithContext(ctx).Delete(&record{}, body.ID.Int()).Error) {
		return
	}
	if h.rdb != nil {
		cache.Drop(ctx, h.rdb, "area_site_codes")
	}
	response.OK(c, "删除成功", []any{})
}

func showBit(value int) int {
	if value == 0 {
		return 0
	}
	return 1
}

func areaType(parentID int) int {
	if parentID > 0 {
		return 2
	}
	return 1
}
