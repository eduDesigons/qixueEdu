package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/captcha"
	commentsvc "qixue-api/internal/service/webapp/comment"
	coursesvc "qixue-api/internal/service/webapp/course"
	feedbacksvc "qixue-api/internal/service/webapp/feedback"
	homesvc "qixue-api/internal/service/webapp/home"
	learnsvc "qixue-api/internal/service/webapp/learn"
	papersvc "qixue-api/internal/service/webapp/paper"
	pathsvc "qixue-api/internal/service/webapp/path"
	ranksvc "qixue-api/internal/service/webapp/rank"
	usersvc "qixue-api/internal/service/webapp/user"
	webauthsvc "qixue-api/internal/service/webapp/auth"
	webauth "qixue-api/internal/webapp/auth"
	webcomment "qixue-api/internal/webapp/comment"
	webcommon "qixue-api/internal/webapp/common"
	webcourse "qixue-api/internal/webapp/course"
	webfeedback "qixue-api/internal/webapp/feedback"
	webhome "qixue-api/internal/webapp/home"
	weblearn "qixue-api/internal/webapp/learn"
	webpaper "qixue-api/internal/webapp/paper"
	webpath "qixue-api/internal/webapp/path"
	webrank "qixue-api/internal/webapp/rank"
	websys "qixue-api/internal/webapp/sys"
	webuser "qixue-api/internal/webapp/user"
)

func registerWebapp(r *gin.Engine, db *gorm.DB, rdb *redis.Client, slide *captcha.Slide, userJWT auth.Issuer) {
	userRepo := repo.NewUser(db)
	studentAuthSvc := webauthsvc.New(userRepo, db, rdb, userJWT, slide)
	studentAuth := webauth.New(studentAuthSvc)

	webOpen := r.Group("/webapp/auth")
	webOpen.GET("/captcha", studentAuth.Captcha)
	webOpen.POST("/captchaCheck", studentAuth.CheckCaptcha)
	webOpen.POST("/login", studentAuth.Login)

	student := r.Group("/webapp", middleware.WebappAuth(db, rdb, userJWT))

	sysExtra := websys.New(db, rdb)
	student.GET("/sys/getQuestion", sysExtra.GetQuestion)
	student.GET("/sys/about", sysExtra.About)
	student.GET("/sys/getTools", sysExtra.GetTools)

	learnerRepo := repo.NewLearner(db)
	courseRepo := repo.NewCourse(db)
	homeRepo := repo.NewHome(db)
	learnRepo := repo.NewLearn(db)
	paperRepo := repo.NewPaper(db)
	commentRepo := repo.NewComment(db)
	feedbackRepo := repo.NewFeedback(db)

	homeSvc := homesvc.New(homeRepo, db, rdb)
	homes := webhome.New(homeSvc)
	student.GET("/home/boot", homes.Boot)
	student.GET("/home/getHomeInfo", homes.GetHomeInfo)
	student.GET("/home/getBannerDetail", homes.GetBannerDetail)
	student.GET("/home/getNoticeDetail", homes.GetNoticeDetail)
	student.GET("/home/getHotClassListByCateId", homes.GetHotClassListByCateId)
	student.GET("/home/getHotPathList", homes.GetHotPathList)
	student.GET("/home/getHotCourseList", homes.GetHotCourseList)

	courseSvc := coursesvc.New(courseRepo, userRepo, db, rdb)
	learnCourse := webcourse.New(courseSvc)
	student.GET("/course/getCourseCategories", learnCourse.Categories)
	student.POST("/course/getCourseList", learnCourse.GetCourseList)
	student.GET("/course/courseDetail", learnCourse.CourseDetail)
	student.GET("/course/getCourseChapterList", learnCourse.GetCourseChapterList)
	student.GET("/course/chapterDetail", learnCourse.ChapterDetail)

	pathSvc := pathsvc.New(courseRepo, userRepo, db, rdb)
	series := webpath.New(pathSvc)
	student.GET("/path/getPathCategories", series.Categories)
	student.GET("/path/getPathList", series.GetPathList)
	student.GET("/path/getPathDetail", series.GetPathDetail)
	student.POST("/path/joinPath", series.JoinPath)
	student.POST("/path/applyPath", series.ApplyPath)

	learnSvc := learnsvc.New(learnRepo, userRepo, db, rdb)
	study := weblearn.New(learnSvc)
	student.GET("/learn/myCourse", study.MyCourse)
	student.GET("/learn/myTask", study.MyTask)
	student.GET("/learn/myHistory", study.MyHistory)
	student.GET("/learn/learnNum", study.LearnNum)
	student.GET("/learn/soonExpire", study.SoonExpire)
	student.POST("/learn/soonExpire", study.SoonExpire)
	student.POST("/learn/collect", study.Collect)
	student.POST("/learn/collectChapter", study.CollectChapter)
	student.POST("/learn/joinCourse", study.JoinCourse)
	student.POST("/learn/applyCourse", study.ApplyCourse)
	student.POST("/learn/unlockChapter", study.UnlockChapter)
	student.POST("/learn/learnEnd", study.LearnEnd)
	student.POST("/learn/endCourse", study.EndCourse)

	meSvc := usersvc.New(userRepo, learnerRepo, rdb)
	me := webuser.New(meSvc)
	student.GET("/user/info", me.Info)
	student.POST("/user/edit", me.Edit)
	student.POST("/user/setPwd", me.SetPwd)
	student.GET("/user/myCount", me.MyCount)
	student.GET("/user/myMessages", me.MyMessages)
	student.GET("/user/unreadMessageCount", me.UnreadMessageCount)
	student.POST("/user/readMessage", me.ReadMessage)
	student.GET("/user/myLikeCourse", me.MyLikeCourse)
	student.GET("/user/myLikeChapter", me.MyLikeChapter)
	student.GET("/user/myCourseComment", me.MyCourseComment)
	student.GET("/user/myChapterComment", me.MyChapterComment)
	student.GET("/user/myExam", me.MyExam)
	student.GET("/user/myTest", me.MyTest)
	student.GET("/user/myCoinsLog", me.MyCoinsLog)
	student.GET("/user/myCreditLog", me.MyCreditLog)
	student.GET("/user/myPath", me.MyPath)
	student.GET("/user/myCert", me.MyCert)

	paperSvc := papersvc.New(paperRepo, db, rdb)
	exams := webpaper.New(paperSvc)
	student.GET("/paper/getPaperInfo", exams.GetPaperInfo)
	student.GET("/paper/getPaperLogInfo", exams.GetPaperLogInfo)
	student.POST("/paper/submitPaper", exams.SubmitPaper)

	commentSvc := commentsvc.New(commentRepo, db, rdb)
	talk := webcomment.New(commentSvc)
	student.POST("/comment/getChapterCommentList", talk.GetChapterCommentList)
	student.POST("/comment/chapterComment", talk.ChapterComment)
	student.POST("/comment/getCourseCommentList", talk.GetCourseCommentList)
	student.POST("/comment/courseComment", talk.CourseComment)
	student.GET("/comment/getChapterCommentList", talk.GetChapterCommentList)
	student.GET("/comment/getCourseCommentList", talk.GetCourseCommentList)

	feedbackSvc := feedbacksvc.New(feedbackRepo)
	ideas := webfeedback.New(feedbackSvc)
	student.POST("/feedback/submit", ideas.Submit)
	student.GET("/feedback/list", ideas.List)

	rankSvc := ranksvc.New(repo.NewRank(db), rdb)
	boards := webrank.New(rankSvc)
	student.GET("/rank/boards", boards.Boards)
	student.GET("/rank/history", boards.History)

	shared := webcommon.New()
	student.POST("/common/upload", shared.Upload)
}
