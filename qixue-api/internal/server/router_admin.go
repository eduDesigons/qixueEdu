package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/admin/http/common"
	admincoursehttp "qixue-api/internal/admin/http/course"
	admindecoration "qixue-api/internal/admin/http/decoration"
	adminhome "qixue-api/internal/admin/http/home"
	"qixue-api/internal/admin/http/login"
	adminmedia "qixue-api/internal/admin/http/media"
	adminmemberhttp "qixue-api/internal/admin/http/member"
	adminops "qixue-api/internal/admin/http/ops"
	adminpms "qixue-api/internal/admin/http/pms"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/admin/service/banner"
	"qixue-api/internal/admin/service/category"
	"qixue-api/internal/admin/service/cert"
	svcadminchapter "qixue-api/internal/admin/service/chapter"
	admincommentsvc "qixue-api/internal/admin/service/comment"
	"qixue-api/internal/admin/service/common"
	svcadmincourse "qixue-api/internal/admin/service/course"
	"qixue-api/internal/admin/service/feedback"
	adminhomesvc "qixue-api/internal/admin/service/home"
	"qixue-api/internal/admin/service/login"
	"qixue-api/internal/admin/service/member"
	"qixue-api/internal/admin/service/message"
	"qixue-api/internal/admin/service/notice"
	adminopssvc "qixue-api/internal/admin/service/ops"
	svcadminpaper "qixue-api/internal/admin/service/paper"
	svcadminpath "qixue-api/internal/admin/service/path"
	adminrolesvc "qixue-api/internal/admin/service/role"
	adminstaffsvc "qixue-api/internal/admin/service/staff"
	"qixue-api/internal/admin/service/verse"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/captcha"
)

func registerAdmin(r *gin.Engine, db *gorm.DB, rdb *redis.Client, slide *captcha.Slide, adminJWT auth.Issuer) {
	account := login.New(adminlogin.New(db, rdb, adminJWT, slide), rdb)
	settings := adminpms.NewConfig(db, rdb)

	open := r.Group("/admin/login")
	open.GET("/captcha", account.Captcha)
	open.POST("/captchaCheck", account.CheckCaptcha)
	open.POST("/login", account.Login)

	admin := r.Group("/admin", middleware.AdminAuth(rdb, adminJWT))
	admin.GET("/login/getUserInfo", account.UserInfo)
	admin.GET("/login/getUserPermissions", account.Permissions)
	admin.POST("/login/updateProfile", account.UpdateProfile)
	admin.POST("/login/changePassword", account.ChangePassword)
	admin.POST("/login/logout", account.Logout)

	words := adminpms.NewVerse(adminverse.New(adminrepo.NewAdminVerse(db)))
	admin.POST("/pms/verse/verseList", words.List)
	admin.POST("/pms/verse/save", words.Save)
	admin.POST("/pms/verse/delete", words.Delete)

	regions := adminpms.NewArea(db, rdb)
	admin.POST("/pms/area/list", regions.List)
	admin.POST("/pms/area/save", regions.Save)
	admin.POST("/pms/area/delete", regions.Delete)

	slides := admindecoration.NewBanner(adminbanner.New(adminrepo.NewAdminBanner(db)))
	admin.GET("/decoration/banner/getBannerList", slides.List)
	admin.POST("/decoration/banner/saveBannerInfo", slides.Save)
	admin.GET("/decoration/banner/delBanner", slides.Delete)
	admin.POST("/decoration/banner/updateBannerSort", slides.Sort)

	notes := admindecoration.NewNotice(adminnotice.New(adminrepo.NewAdminNotice(db)))
	admin.POST("/decoration/notice/getNoticeList", notes.List)
	admin.GET("/decoration/notice/getNoticeDetail", notes.Detail)
	admin.POST("/decoration/notice/saveNotice", notes.Save)
	admin.POST("/decoration/notice/updateNoticeStatus", notes.Status)
	admin.GET("/decoration/notice/deleteNotice", notes.Delete)

	menus := adminpms.NewMenu(db)
	admin.GET("/pms/menu/getMenuTree", menus.Tree)
	admin.GET("/pms/menu/getButtons", menus.Buttons)
	admin.POST("/pms/menu/addMenu", menus.Add)
	admin.POST("/pms/menu/saveMenu", menus.Save)
	admin.POST("/pms/menu/delMenu", menus.Delete)

	roles := adminpms.NewRole(adminrolesvc.New(adminrepo.NewAdminRole(db), db))
	admin.POST("/pms/role/addRole", roles.Add)
	admin.POST("/pms/role/editRole", roles.Edit)
	admin.POST("/pms/role/delRole", roles.Delete)
	admin.POST("/pms/role/roleList", roles.List)
	admin.GET("/pms/role/getMenuTree", roles.MenuTree)
	admin.GET("/pms/role/getAllRole", roles.All)

	admins := adminpms.NewStaff(adminstaffsvc.New(adminrepo.NewAdminStaff(db), db, rdb))
	admin.POST("/pms/admin/adminList", admins.List)
	admin.POST("/pms/admin/addRoleToAdmin", admins.AddRole)
	admin.POST("/pms/admin/delRoleFromAdmin", admins.RemoveRole)
	admin.POST("/pms/admin/addAdmin", admins.Add)
	admin.POST("/pms/admin/editAdmin", admins.Edit)
	admin.POST("/pms/admin/delAdmin", admins.Delete)
	admin.POST("/pms/admin/resetPassword", admins.ResetPassword)

	admin.GET("/pms/config/getInfo", settings.Get)
	admin.POST("/pms/config/saveInfo", settings.Save)

	cates := adminmedia.NewCate(db)
	admin.POST("/media/mediaCate/list", cates.List)
	admin.POST("/media/mediaCate/destroy", cates.Destroy)

	files := adminmedia.NewMedia(db)
	admin.POST("/media/media/list", files.List)
	admin.POST("/media/media/upload", files.Upload)
	admin.POST("/media/media/destroy", files.Destroy)
	admin.POST("/media/media/saveUploaded", files.SaveUploaded)
	admin.POST("/media/media/initVideoDirectUpload", files.InitVideoDirectUpload)
	admin.POST("/media/media/completeVideoDirectUpload", files.CompleteVideoDirectUpload)
	admin.POST("/media/media/abortVideoDirectUpload", files.AbortVideoDirectUpload)

	ops := adminops.New(adminopssvc.New(adminrepo.NewAdminOps(db), db, rdb))
	admin.POST("/ops/ops/courseEnrollList", ops.CourseEnrollList)
	admin.POST("/ops/ops/pathEnrollList", ops.PathEnrollList)
	admin.POST("/ops/ops/learnList", ops.LearnList)
	admin.POST("/ops/ops/creditList", ops.CreditList)
	admin.POST("/ops/ops/coinsList", ops.CoinsList)
	admin.POST("/ops/ops/chart", ops.Chart)
	admin.POST("/ops/ops/inspectKinds", ops.InspectKinds)
	admin.POST("/ops/ops/inspect", ops.Inspect)
	admin.POST("/ops/ops/repair", ops.Repair)

	dash := adminhome.New(adminhomesvc.New(adminrepo.NewAdminHome(db), db, rdb))
	admin.GET("/home/dashboard/getVerse", dash.GetVerse)
	admin.GET("/home/dashboard/getPendingStats", dash.GetPendingStats)
	admin.GET("/home/dashboard/getTodayStats", dash.GetTodayStats)
	admin.GET("/home/dashboard/getStats", dash.GetStats)
	admin.GET("/home/dashboard/getChartData", dash.GetChartData)
	admin.GET("/home/dashboard/getCourseLearningRank", dash.GetCourseLearningRank)
	admin.GET("/home/dashboard/getChapterLearningRank", dash.GetChapterLearningRank)
	admin.GET("/home/dashboard/getChapterCompleteUserRank", dash.GetChapterCompleteUserRank)
	admin.GET("/home/dashboard/getCdnTraffic", dash.GetCdnTraffic)

	places := common.New(admincommon.New(adminrepo.NewAdminAreaCommon(db), db, rdb))
	admin.GET("/common/area/getAreaList", places.List)
	admin.GET("/common/area/getAreaTree", places.Tree)

	kinds := admincoursehttp.NewCategory(admincategory.New(adminrepo.NewAdminCategory(db), db, rdb))
	admin.GET("/course/category/list", kinds.List)
	admin.POST("/course/category/save", kinds.Save)
	admin.POST("/course/category/delete", kinds.Delete)
	admin.GET("/course/course/getCourseCategories", kinds.Options)

	courses := admincoursehttp.NewCourse(svcadmincourse.New(adminrepo.NewAdminCourse(db), db, rdb))
	admin.POST("/course/course/getCourseList", courses.GetCourseList)
	admin.GET("/course/course/getCourseDetail", courses.GetCourseDetail)
	admin.POST("/course/course/courseSave", courses.CourseSave)
	admin.GET("/course/course/courseDelete", courses.CourseDelete)
	admin.GET("/course/course/getCourseAttachList", courses.GetCourseAttachList)
	admin.POST("/course/course/courseAttachSave", courses.CourseAttachSave)
	admin.GET("/course/course/courseAttachDelete", courses.CourseAttachDelete)
	admin.POST("/course/course/getCourseUserList", courses.GetCourseUserList)
	admin.GET("/course/course/courseUserDelete", courses.CourseUserDelete)
	admin.POST("/course/course/getCourseUserSearch", courses.GetCourseUserSearch)
	admin.POST("/course/course/courseUserListSave", courses.CourseUserListSave)
	admin.POST("/course/course/getCourseUserApplyList", courses.GetCourseUserApplyList)
	admin.POST("/course/course/courseUserApplyStatusChange", courses.CourseUserApplyStatusChange)
	admin.GET("/course/course/courseUserApplyDelete", courses.CourseUserApplyDelete)

	chapters := admincoursehttp.NewChapter(svcadminchapter.New(adminrepo.NewAdminChapter(db), db, rdb))
	admin.GET("/course/chapter/getChapterList", chapters.GetChapterList)
	admin.GET("/course/chapter/getChapterDetail", chapters.GetChapterDetail)
	admin.GET("/course/chapter/chapterDelete", chapters.ChapterDelete)
	admin.POST("/course/chapter/chapterSave", chapters.ChapterSave)
	admin.POST("/course/chapter/chapterSort", chapters.ChapterSort)
	admin.POST("/course/chapter/getChapterUserList", chapters.GetChapterUserList)
	admin.POST("/course/chapter/unlockChapterUser", chapters.UnlockChapterUser)

	paths := admincoursehttp.NewPath(svcadminpath.New(adminrepo.NewAdminPath(db), db, rdb))
	admin.GET("/course/path/getPathCategories", paths.GetPathCategories)
	admin.POST("/course/path/getPathList", paths.GetPathList)
	admin.GET("/course/path/getPathDetail", paths.GetPathDetail)
	admin.GET("/course/path/pathDelete", paths.PathDelete)
	admin.POST("/course/path/pathSave", paths.PathSave)
	admin.POST("/course/path/getPathCourseList", paths.GetPathCourseList)
	admin.POST("/course/path/pathCourseListSave", paths.PathCourseListSave)
	admin.GET("/course/path/pathCourseDelete", paths.PathCourseDelete)
	admin.POST("/course/path/getPathCourseSearch", paths.GetPathCourseSearch)
	admin.POST("/course/path/getUserPathList", paths.GetUserPathList)
	admin.GET("/course/path/userPathDelete", paths.UserPathDelete)
	admin.POST("/course/path/getUserPathSearch", paths.GetUserPathSearch)
	admin.POST("/course/path/userPathListSave", paths.UserPathListSave)
	admin.POST("/course/path/getUserPathApplyList", paths.GetUserPathApplyList)
	admin.POST("/course/path/userPathApplyStatusChange", paths.UserPathApplyStatusChange)
	admin.GET("/course/path/userPathApplyDelete", paths.UserPathApplyDelete)

	papers := admincoursehttp.NewPaper(svcadminpaper.New(adminrepo.NewAdminPaper(db), db, rdb))
	admin.GET("/course/paper/getPaperDetail", papers.GetPaperDetail)
	admin.POST("/course/paper/paperSave", papers.PaperSave)
	admin.GET("/course/paper/paperDelete", papers.PaperDelete)
	admin.POST("/course/paper/getExamList", papers.GetExamList)
	admin.POST("/course/paper/getAssessmentList", papers.GetAssessmentList)
	admin.POST("/course/paper/getPaperLogs", papers.GetPaperLogs)
	admin.GET("/course/paper/getUserPaperDetail", papers.GetUserPaperDetail)
	admin.POST("/course/paper/paperGrade", papers.PaperGrade)

	comments := admincoursehttp.NewComment(admincommentsvc.New(adminrepo.NewAdminComment(db), db, rdb))
	admin.POST("/course/comment/getCourseCommentList", comments.GetCourseCommentList)
	admin.POST("/course/comment/addCourseCommentReply", comments.AddCourseCommentReply)
	admin.GET("/course/comment/getCourseCommentReplies", comments.GetCourseCommentReplies)
	admin.GET("/course/comment/courseCommentDelete", comments.CourseCommentDelete)
	admin.POST("/course/comment/courseCommentBatchDelete", comments.CourseCommentBatchDelete)
	admin.POST("/course/comment/getCourseChapterCommentList", comments.GetCourseChapterCommentList)
	admin.POST("/course/comment/addChapterCommentReply", comments.AddChapterCommentReply)
	admin.GET("/course/comment/getChapterCommentReplies", comments.GetChapterCommentReplies)
	admin.GET("/course/comment/courseChapterCommentDelete", comments.CourseChapterCommentDelete)
	admin.POST("/course/comment/courseChapterCommentBatchDelete", comments.CourseChapterCommentBatchDelete)

	members := adminmemberhttp.NewUser(adminmember.New(adminrepo.NewAdminMember(db), db, rdb))
	admin.POST("/member/user/getUserList", members.GetUserList)
	admin.GET("/member/user/getUserDetail", members.GetUserDetail)
	admin.POST("/member/user/createUser", members.CreateUser)
	admin.POST("/member/user/updateUser", members.UpdateUser)
	admin.POST("/member/user/setLock", members.SetLock)
	admin.GET("/member/user/deleteUser", members.DeleteUser)
	admin.GET("/member/user/getUserCourseList", members.GetUserCourseList)
	admin.GET("/member/user/getUserPathList", members.GetUserPathList)
	admin.GET("/member/user/getUserCourseChapterList", members.GetUserCourseChapterList)
	admin.POST("/member/user/getUserCreditRecordList", members.GetUserCreditRecordList)
	admin.POST("/member/user/getUserCoinsRecordList", members.GetUserCoinsRecordList)
	admin.POST("/member/user/getUserPaperList", members.GetUserPaperList)
	admin.POST("/member/user/getUserCertList", members.GetUserCertList)

	feedbacks := adminmemberhttp.NewFeedback(adminfeedback.New(adminrepo.NewAdminFeedback(db), db, rdb))
	admin.POST("/member/feedback/list", feedbacks.List)
	admin.POST("/member/feedback/reply", feedbacks.Reply)

	messages := adminmemberhttp.NewMessage(adminmessage.New(adminrepo.NewAdminMessage(db), db, rdb))
	admin.POST("/member/message/send", messages.Send)
	admin.POST("/member/message/list", messages.List)
	admin.POST("/member/message/delete", messages.Delete)

	certs := adminmemberhttp.NewCert(admincert.New(adminrepo.NewAdminCert(db), db, rdb))
	admin.POST("/member/cert/list", certs.List)
	admin.GET("/member/cert/detail", certs.Detail)
	admin.GET("/member/cert/options", certs.Options)
	admin.POST("/member/cert/save", certs.Save)
	admin.POST("/member/cert/setShow", certs.SetShow)
	admin.POST("/member/cert/delete", certs.Delete)
	admin.POST("/member/cert/issue", certs.Issue)
	admin.POST("/member/cert/backfill", certs.Backfill)
	admin.POST("/member/cert/issuedList", certs.IssuedList)
	admin.POST("/member/cert/revoke", certs.Revoke)
}
