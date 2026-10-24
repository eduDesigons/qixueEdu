package adminhome

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Service) buildCourseRank(c *gin.Context) ([]map[string]any, error) {
	cityID, limited, err := s.cityScope(c)
	if err != nil {
		return nil, err
	}
	rows, err := s.home.CourseLearningRanks(c.Request.Context(), cityID, limited, 10)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(rows))
	for _, r := range rows {
		if r.ID > 0 {
			ids = append(ids, r.ID)
		}
	}
	titles, err := s.home.CourseTitles(c.Request.Context(), ids)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for i, r := range rows {
		title := titles[r.ID]
		if title == "" {
			title = "-"
		}
		list = append(list, map[string]any{
			"rank":           i + 1,
			"course_id":      r.ID,
			"title":          title,
			"enroll_count":   r.Count,
			"learning_count": r.Count,
			"complete_count": 0,
		})
	}
	return list, nil
}

func (s *Service) buildChapterRank(c *gin.Context) ([]map[string]any, error) {
	since := time.Now().In(shanghai()).AddDate(0, 0, -30).Format("2006-01-02 15:04:05")
	cityID, limited, err := s.cityScope(c)
	if err != nil {
		return nil, err
	}
	rows, err := s.home.ChapterCompleteRanks(c.Request.Context(), since, cityID, limited, 10)
	if err != nil {
		return nil, err
	}
	chapterIDs, courseIDs := []int{}, []int{}
	for _, r := range rows {
		if r.ChapterID > 0 {
			chapterIDs = append(chapterIDs, r.ChapterID)
		}
		if r.CourseID > 0 {
			courseIDs = append(courseIDs, r.CourseID)
		}
	}
	chapters, err := s.home.ChapterTitles(c.Request.Context(), chapterIDs)
	if err != nil {
		return nil, err
	}
	courses, err := s.home.CourseTitles(c.Request.Context(), courseIDs)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for i, r := range rows {
		ct := chapters[r.ChapterID]
		if ct == "" {
			ct = "-"
		}
		cot := courses[r.CourseID]
		if cot == "" {
			cot = "-"
		}
		list = append(list, map[string]any{
			"rank":           i + 1,
			"chapter_id":     r.ChapterID,
			"course_id":      r.CourseID,
			"chapter_title":  ct,
			"course_title":   cot,
			"complete_count": r.Count,
		})
	}
	return list, nil
}

func (s *Service) buildCreditRank(c *gin.Context) ([]map[string]any, error) {
	since := time.Now().In(shanghai()).AddDate(0, 0, -14).Format("2006-01-02 15:04:05")
	cityID, limited, err := s.cityScope(c)
	if err != nil {
		return nil, err
	}
	rows, err := s.home.CoinsRanks(c.Request.Context(), since, cityID, limited, 8)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		if r.ID > 0 {
			userIDs = append(userIDs, r.ID)
		}
	}
	users, err := s.home.UserAvatars(c.Request.Context(), userIDs)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for i, r := range rows {
		u := users[r.ID]
		list = append(list, map[string]any{
			"rank":      i + 1,
			"user_id":   r.ID,
			"avatar":    u.Avatar,
			"nick_name": u.NickName,
			"credit":    r.Count,
		})
	}
	return list, nil
}
