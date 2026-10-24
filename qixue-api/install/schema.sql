SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS `administrator`;
CREATE TABLE `administrator` (
  `id` int NOT NULL AUTO_INCREMENT,
  `city_id` int DEFAULT NULL,
  `area_id` int DEFAULT NULL COMMENT '关联区域id',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `enable` tinyint NOT NULL DEFAULT '1',
  `createTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updateTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `gender` int DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `nickName` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `is_super` int NOT NULL DEFAULT '0' COMMENT '是否是超管',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_78a916df40e02a9deb1c4b75ed` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='后台用户';

DROP TABLE IF EXISTS `administrator_menus`;
CREATE TABLE `administrator_menus` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `parentId` int DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `redirect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `layout` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `keepAlive` tinyint DEFAULT NULL,
  `method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `show` tinyint NOT NULL DEFAULT '1' COMMENT '是否展示在页面菜单',
  `enable` tinyint NOT NULL DEFAULT '1',
  `order` int DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_30e166e8c6359970755c5727a2` (`code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='菜单';

DROP TABLE IF EXISTS `administrator_permissions`;
CREATE TABLE `administrator_permissions` (
  `roleId` int NOT NULL,
  `permissionId` int NOT NULL,
  PRIMARY KEY (`roleId`,`permissionId`) USING BTREE,
  KEY `IDX_b36cb2e04bc353ca4ede00d87b` (`roleId`) USING BTREE,
  KEY `IDX_bfbc9e263d4cea6d7a8c9eb3ad` (`permissionId`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色权限表';

DROP TABLE IF EXISTS `administrator_role_relation`;
CREATE TABLE `administrator_role_relation` (
  `administratorId` int NOT NULL,
  `roleId` int NOT NULL,
  PRIMARY KEY (`administratorId`,`roleId`) USING BTREE,
  KEY `IDX_5f9286e6c25594c6b88c108db7` (`administratorId`) USING BTREE,
  KEY `IDX_4be2f7adf862634f5f803d246b` (`roleId`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户角色关联表';

DROP TABLE IF EXISTS `administrator_roles`;
CREATE TABLE `administrator_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `enable` tinyint NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `IDX_ee999bb389d7ac0fd967172c41` (`code`) USING BTREE,
  UNIQUE KEY `IDX_ae4578dcaed5adff96595e6166` (`name`) USING BTREE,
  UNIQUE KEY `administrator_roles_pk` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色表';

DROP TABLE IF EXISTS `app_config`;
CREATE TABLE `app_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(188) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '键',
  `value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '值',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `app_config_key_unique` (`key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='配置';

DROP TABLE IF EXISTS `areas`;
CREATE TABLE `areas` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '区域名字',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父id',
  `parent_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '父名称',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '类型：1城市 2区域',
  `is_show` tinyint NOT NULL DEFAULT '1' COMMENT '是否显示,1显示,0不显示',
  `sort` int NOT NULL DEFAULT '0' COMMENT '升序',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `code` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'code',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='区域';

DROP TABLE IF EXISTS `banners`;
CREATE TABLE `banners` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `img_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `link_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT ' 1=仅展示, 2=链接, 3=内容',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1启用 2停用',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '内容',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_banners_status_sort` (`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='轮播图设置';

DROP TABLE IF EXISTS `certs`;
CREATE TABLE `certs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '证书名',
  `template_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '证书模板',
  `preview_image` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '参数',
  `is_show` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='证书';

DROP TABLE IF EXISTS `course_attach`;
CREATE TABLE `course_attach` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '附件名',
  `path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '路径',
  `size` int NOT NULL DEFAULT '0' COMMENT '单位：byte',
  `download_times` int NOT NULL DEFAULT '0' COMMENT '下载次数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `course_attach_course_id_index` (`course_id`) USING BTREE,
  KEY `idx_ca_course` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课程课件';

DROP TABLE IF EXISTS `course_categories`;
CREATE TABLE `course_categories` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '展示图标',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父id',
  `is_show` tinyint NOT NULL DEFAULT '1' COMMENT '是否显示,1显示,0不显示',
  `sort` int NOT NULL DEFAULT '0' COMMENT '升序',
  `min_need_credit` int NOT NULL DEFAULT '0' COMMENT '新建课程默认限制学分',
  `base_credit` int NOT NULL DEFAULT '0' COMMENT '学员等级基础学分',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ccat_parent_sort` (`parent_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课程分类';

DROP TABLE IF EXISTS `course_chapter_comments`;
CREATE TABLE `course_chapter_comments` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pid` int NOT NULL DEFAULT '0' COMMENT '回复楼上id',
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `chapter_id` int NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ccmt_chapter_pid` (`chapter_id`,`pid`,`id`),
  KEY `idx_ccmt_user_id` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课节评论';

DROP TABLE IF EXISTS `course_chapters`;
CREATE TABLE `course_chapters` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `media_id` int DEFAULT NULL,
  `media_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '媒体资源地址',
  `type` int NOT NULL DEFAULT '1' COMMENT '课程类型',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '详情',
  `limit_day` int DEFAULT NULL COMMENT '限制天数（超期扣积分）',
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `paper_id` int DEFAULT NULL COMMENT '关联测试id',
  `sort` int NOT NULL DEFAULT '1',
  `view_count` int NOT NULL DEFAULT '0' COMMENT '浏览量',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_cc_course_sort` (`course_id`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课程小节表';

DROP TABLE IF EXISTS `course_comments`;
CREATE TABLE `course_comments` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pid` int NOT NULL DEFAULT '0' COMMENT '回复楼上评论',
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始内容',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_cmt_course_pid` (`course_id`,`pid`,`id`),
  KEY `idx_cmt_user_id` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课程评论';

DROP TABLE IF EXISTS `courses`;
CREATE TABLE `courses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名',
  `thumb` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '封面',
  `desc` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '简短介绍',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '课程类型 1视频 2音频 3图文',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '介绍内容',
  `is_show` tinyint NOT NULL COMMENT '1显示,-1隐藏',
  `category_id` int NOT NULL DEFAULT '0' COMMENT '分类id',
  `is_rec` tinyint NOT NULL DEFAULT '0' COMMENT '推荐,0否,1是',
  `is_open` tinyint NOT NULL DEFAULT '0' COMMENT '是否对外',
  `min_need_credit` int NOT NULL DEFAULT '0' COMMENT '要求学分',
  `credit` int NOT NULL DEFAULT '0' COMMENT '奖励学分',
  `paper_id` int DEFAULT NULL COMMENT '考试id',
  `cert_id` int unsigned DEFAULT NULL,
  `created_id` int NOT NULL,
  `view_count` int NOT NULL DEFAULT '0' COMMENT '浏览量',
  `chapter_count` int DEFAULT '0' COMMENT '章节数量',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_courses_show_rec_view` (`is_show`,`is_rec`,`view_count`),
  KEY `idx_courses_show_cate` (`is_show`,`category_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='课程库';

DROP TABLE IF EXISTS `media_cate`;
CREATE TABLE `media_cate` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类名字',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件类型',
  `is_init` tinyint NOT NULL DEFAULT '0' COMMENT '是否内置(内置的不能删除)',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `media_cate_type_index` (`type`) USING BTREE
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='媒体资源分类';

DROP TABLE IF EXISTS `media_file`;
CREATE TABLE `media_file` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频名',
  `thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `duration` int NOT NULL DEFAULT '0' COMMENT '时长',
  `size` bigint NOT NULL DEFAULT '0' COMMENT '单位：byte',
  `storage_file_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '存储特征值',
  `transcode_status` tinyint NOT NULL DEFAULT '0' COMMENT '0未转码,3转码中,5转码成功,7转码失败',
  `ref_count` int NOT NULL DEFAULT '0' COMMENT '引用次数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `is_open` tinyint NOT NULL DEFAULT '0' COMMENT '开放[0:否,1:是].开放的话则可以公开播放',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型[video,audio,image,document]',
  `cate_id` int NOT NULL DEFAULT '0' COMMENT '分类',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件路径',
  `vid` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '视频上传返回uid 必存',
  `download_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频文件下载路径',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `media_file_type_index` (`type`) USING BTREE
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='媒体资源';

DROP TABLE IF EXISTS `notices`;
CREATE TABLE `notices` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告标题',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `publish_time` timestamp NOT NULL COMMENT '发布时间',
  `admin_id` int NOT NULL DEFAULT '0',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `view_times` int NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_notices_status_time` (`status`,`publish_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='公告';

DROP TABLE IF EXISTS `paper_questions`;
CREATE TABLE `paper_questions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `paper_id` int DEFAULT NULL COMMENT '考试id',
  `sort` int NOT NULL DEFAULT '1' COMMENT '序号',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型 single multiple completion text uploadImage uploadDocument',
  `title` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '内容',
  `answer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '答案',
  `score` tinyint NOT NULL DEFAULT '0' COMMENT '分数',
  `options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `qt` (`type`) USING BTREE,
  KEY `paper_questions_paper_id_index` (`paper_id`) USING BTREE,
  KEY `idx_pq_paper_sort` (`paper_id`,`sort`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='考试题目';

DROP TABLE IF EXISTS `papers`;
CREATE TABLE `papers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `course_id` int DEFAULT NULL,
  `chapter_id` int DEFAULT NULL,
  `score` int NOT NULL DEFAULT '0' COMMENT '总分',
  `pass_score` int NOT NULL DEFAULT '0' COMMENT '及格分',
  `created_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `is_skip_mark` tinyint NOT NULL DEFAULT '0' COMMENT '人工阅卷',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'course' COMMENT '试卷类型 course chapter',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_papers_type` (`type`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='试卷 课节/课程';

DROP TABLE IF EXISTS `path_categories`;
CREATE TABLE `path_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父类',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类名',
  `sort` int NOT NULL DEFAULT '0' COMMENT '升序',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_pcat_parent_sort` (`parent_id`,`sort`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系列课程分类';

DROP TABLE IF EXISTS `path_courses`;
CREATE TABLE `path_courses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `path_id` int NOT NULL,
  `course_id` int NOT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_pc_path_sort` (`path_id`,`sort`),
  KEY `idx_pc_course` (`course_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=FIXED COMMENT='系列-课程关联表';

DROP TABLE IF EXISTS `paths`;
CREATE TABLE `paths` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `thumb` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `desc` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '介绍 内容',
  `is_open` tinyint NOT NULL DEFAULT '0' COMMENT '是否对外',
  `is_show` tinyint NOT NULL DEFAULT '0' COMMENT '发布,1是,0否',
  `is_rec` int NOT NULL DEFAULT '0' COMMENT '是否推荐',
  `created_at` timestamp NULL DEFAULT NULL,
  `created_id` int DEFAULT NULL COMMENT '创建人',
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `category_id` int NOT NULL DEFAULT '0' COMMENT '分类id',
  `cert_id` int unsigned DEFAULT NULL COMMENT '证书id',
  `min_need_credit` int NOT NULL DEFAULT '0' COMMENT '要求学分',
  `view_count` int NOT NULL DEFAULT '0' COMMENT '浏览量',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_paths_show_updated` (`is_show`,`updated_at`),
  KEY `idx_paths_show_cate` (`is_show`,`category_id`,`updated_at`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='系列课';

DROP TABLE IF EXISTS `rank_snapshots`;
CREATE TABLE `rank_snapshots` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `site_code` varchar(16) NOT NULL DEFAULT '',
  `type` varchar(16) NOT NULL,
  `period` varchar(8) NOT NULL,
  `period_key` varchar(16) NOT NULL,
  `period_label` varchar(64) NOT NULL DEFAULT '',
  `snapshot_date` date NOT NULL,
  `payload` mediumtext NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rank_snap` (`site_code`,`type`,`period`,`period_key`,`snapshot_date`),
  KEY `idx_rank_type` (`site_code`,`type`,`period`,`snapshot_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_cert`;
CREATE TABLE `user_cert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cert_id` bigint NOT NULL DEFAULT '0' COMMENT '证书id',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `cert_no` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '证书编号',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '证书图片地址',
  `path_id` int DEFAULT NULL COMMENT '系列id',
  `course_id` int unsigned NOT NULL DEFAULT '0',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `image` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_valid` tinyint unsigned NOT NULL DEFAULT '1',
  `issued_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `cert_user_cert_no_unique` (`cert_no`) USING BTREE,
  UNIQUE KEY `uk_user_path_course` (`user_id`,`path_id`,`course_id`),
  KEY `c_u_id` (`cert_id`,`user_id`) USING BTREE,
  KEY `idx_course` (`course_id`),
  KEY `idx_user_valid_id` (`user_id`,`is_valid`,`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户证书记录';

DROP TABLE IF EXISTS `user_coins_records`;
CREATE TABLE `user_coins_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `field` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '字段',
  `sum` int NOT NULL DEFAULT '0' COMMENT '变动额度',
  `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '变动说明',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_coins_records_user_id_field_index` (`user_id`,`field`) USING BTREE,
  KEY `idx_ucoins_user_id` (`user_id`,`id`),
  KEY `idx_ucoins_created` (`created_at`,`user_id`,`sum`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户积分记录';

DROP TABLE IF EXISTS `user_course`;
CREATE TABLE `user_course` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `course_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `has_end_chapter` tinyint NOT NULL DEFAULT '0' COMMENT '已经学完课次',
  `is_end` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否结业（考试通过）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_uc_user_course` (`user_id`,`course_id`),
  KEY `user_course_course_id_index` (`course_id`) USING BTREE,
  KEY `user_course_user_id_index` (`user_id`) USING BTREE,
  KEY `user_course_user_id_course_id_index` (`user_id`,`course_id`),
  KEY `idx_uc_user_updated` (`user_id`,`updated_at`),
  KEY `idx_uc_user_created` (`user_id`,`created_at`),
  KEY `idx_uc_user_end` (`user_id`,`is_end`,`updated_at`),
  KEY `idx_uc_user_course` (`user_id`,`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户课程订阅表';

DROP TABLE IF EXISTS `user_course_apply`;
CREATE TABLE `user_course_apply` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `course_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `reason` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '申请原因',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '申请状态 0 待处理 1 通过 2 拒绝',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_course_apply_id_index` (`course_id`) USING BTREE,
  KEY `user_course_apply_user_id_index` (`user_id`) USING BTREE,
  KEY `idx_uca_user_course` (`user_id`,`course_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户课程申请加入表';

DROP TABLE IF EXISTS `user_course_chapter`;
CREATE TABLE `user_course_chapter` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `chapter_id` int NOT NULL COMMENT '小节',
  `course_id` int DEFAULT NULL COMMENT '课程',
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `has_end_chapter` tinyint(1) DEFAULT '0' COMMENT '是否学完 0未学完 1已学完',
  `is_end` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否结束 测评也结束了',
  `un_lock` int NOT NULL DEFAULT '0' COMMENT '是否解锁(0 未解锁 1已解锁)',
  `learn_time` int NOT NULL DEFAULT '0' COMMENT '分钟',
  `expire_time` timestamp NULL DEFAULT NULL COMMENT '过期时间（解锁时写入）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_ucc_user_chapter` (`user_id`,`chapter_id`),
  KEY `user_course_chapter_chapter_id_index` (`chapter_id`),
  KEY `user_course_chapter_course_id_index` (`course_id`),
  KEY `user_course_chapter_user_id_index` (`user_id`),
  KEY `idx_ucc_user_chapter` (`user_id`,`chapter_id`),
  KEY `idx_ucc_user_course` (`user_id`,`course_id`),
  KEY `idx_ucc_user_end_updated` (`user_id`,`is_end`,`updated_at`),
  KEY `idx_ucc_user_has_end` (`user_id`,`has_end_chapter`,`is_end`),
  KEY `idx_ucc_user_expire` (`user_id`,`is_end`,`expire_time`),
  KEY `idx_ucc_end_updated_ch` (`is_end`,`updated_at`,`chapter_id`,`course_id`),
  KEY `idx_ucc_finished_at` (`has_end_chapter`,`updated_at`,`user_id`,`course_id`,`chapter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户解锁小节';

DROP TABLE IF EXISTS `user_credit_records`;
CREATE TABLE `user_credit_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `field` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '字段',
  `sum` int NOT NULL DEFAULT '0' COMMENT '变动额度',
  `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '变动说明',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_credit_records_user_id_field_index` (`user_id`,`field`) USING BTREE,
  KEY `idx_ucredit_user_id` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户学分记录';

DROP TABLE IF EXISTS `user_feedbacks`;
CREATE TABLE `user_feedbacks` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL,
  `content` varchar(500) NOT NULL DEFAULT '',
  `images` text,
  `status` tinyint unsigned NOT NULL DEFAULT '0',
  `reply` varchar(1000) DEFAULT NULL,
  `reply_admin_id` int unsigned NOT NULL DEFAULT '0',
  `replied_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`,`id`),
  KEY `idx_status_id` (`status`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_like_chapters`;
CREATE TABLE `user_like_chapters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `chapter_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_like_chapter_user_id_course_id_index` (`user_id`,`chapter_id`) USING BTREE,
  KEY `idx_ulch_user_chapter` (`user_id`,`chapter_id`),
  KEY `idx_ulch_user_id` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户-收藏章节';

DROP TABLE IF EXISTS `user_like_courses`;
CREATE TABLE `user_like_courses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `course_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_like_courses_user_id_course_id_index` (`user_id`,`course_id`) USING BTREE,
  KEY `user_like_courses_course_id_index` (`course_id`),
  KEY `user_like_courses_user_id_index` (`user_id`),
  KEY `idx_ulc_user_course` (`user_id`,`course_id`),
  KEY `idx_ulc_user_id` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户-收藏课程';

DROP TABLE IF EXISTS `user_messages`;
CREATE TABLE `user_messages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL,
  `admin_id` int unsigned NOT NULL DEFAULT '0',
  `title` varchar(200) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `type` varchar(32) NOT NULL DEFAULT 'notice',
  `related_type` varchar(32) DEFAULT NULL,
  `related_id` int unsigned DEFAULT NULL,
  `is_read` tinyint unsigned NOT NULL DEFAULT '0',
  `read_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_read` (`user_id`,`is_read`),
  KEY `idx_user_created` (`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_paper`;
CREATE TABLE `user_paper` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户id',
  `paper_id` int NOT NULL DEFAULT '0' COMMENT '练习id',
  `score` int NOT NULL DEFAULT '0' COMMENT '得分',
  `submit_count` int NOT NULL DEFAULT '0' COMMENT '提交数量，表示用户已对答案',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `is_pass` int NOT NULL DEFAULT '0' COMMENT '是否通过 1已通过 0未通过 2待批改',
  `is_people_check` int NOT NULL DEFAULT '0' COMMENT '是否需要人工阅卷',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_paper_user_id_paper_id_index` (`user_id`,`paper_id`) USING BTREE,
  KEY `idx_up_user_paper` (`user_id`,`paper_id`,`id`),
  KEY `idx_up_user_id` (`user_id`,`id`),
  KEY `idx_up_user_paper_pass` (`user_id`,`paper_id`,`is_pass`,`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=FIXED COMMENT='用户考试记录';

DROP TABLE IF EXISTS `user_paper_question_log`;
CREATE TABLE `user_paper_question_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT '0',
  `user_paper_id` int NOT NULL,
  `paper_id` bigint NOT NULL DEFAULT '0',
  `question_id` bigint NOT NULL DEFAULT '0',
  `answer_content` varchar(2550) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '答案',
  `is_correct` tinyint NOT NULL DEFAULT '0' COMMENT '是否正确',
  `score` tinyint NOT NULL DEFAULT '0' COMMENT '获得分数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `need_people_check` int NOT NULL DEFAULT '0' COMMENT '是否需要人工审核',
  `grade_remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '批改备注',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_paper_question_log_user_paper_id_index` (`user_paper_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户答题做题详情';

DROP TABLE IF EXISTS `user_path`;
CREATE TABLE `user_path` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `path_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `has_end_course` tinyint NOT NULL DEFAULT '0' COMMENT '已完成的课程',
  `is_end` tinyint NOT NULL DEFAULT '0' COMMENT '是否完成系列课',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_upath_user_path` (`user_id`,`path_id`),
  KEY `idx_upath_user_path` (`user_id`,`path_id`),
  KEY `idx_upath_user_id` (`user_id`,`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=FIXED COMMENT='用户系列课加入记录';

DROP TABLE IF EXISTS `user_path_apply`;
CREATE TABLE `user_path_apply` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `path_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `reason` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '申请原因',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '申请状态 0 待处理 1 通过 2 拒绝',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_path_apply_id_index` (`path_id`) USING BTREE,
  KEY `user_path_apply_user_id_index` (`user_id`) USING BTREE,
  KEY `idx_upa_user_path` (`user_id`,`path_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户系列申请加入表';

DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE `user_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `real_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
  `age` tinyint NOT NULL DEFAULT '0' COMMENT '年龄',
  `gender` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '性别',
  `birthday` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '生日',
  `profession` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '职业',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '住址',
  `graduated_school` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '毕业院校',
  `diploma` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '毕业证书',
  `id_number` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '身份证号',
  `id_frontend_thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '身份证正面照',
  `id_backend_thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '身份证反面照',
  `id_hand_thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手持身份证照',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `is_verify` tinyint NOT NULL DEFAULT '0' COMMENT '实名认证[0:否,1:是]',
  `verify_image_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '认证图片URL',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_profiles_user_id_unique` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='用户个人信息';

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `avatar` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `nick_name` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mobile` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `passwordbf` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_active` tinyint NOT NULL DEFAULT '-1' COMMENT '是否认证1:active,-1:unactive',
  `is_lock` tinyint NOT NULL DEFAULT '-1' COMMENT '是否拉黑1:lock,-1:unlock',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `invite_user_id` int NOT NULL DEFAULT '0' COMMENT '邀请人id',
  `register_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '注册ip',
  `area_id` int NOT NULL COMMENT '注册地区',
  `city_id` int DEFAULT NULL,
  `city_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_login_id` int NOT NULL DEFAULT '0' COMMENT '最后一次登录记录id',
  `level` int NOT NULL DEFAULT '1' COMMENT '等级',
  `credit` int NOT NULL DEFAULT '0' COMMENT '学分',
  `coins` int NOT NULL DEFAULT '100' COMMENT '积分',
  `account` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '账号',
  `gender` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '性别',
  `age` tinyint NOT NULL DEFAULT '0' COMMENT '年龄',
  `degree` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '学历：小学，初中，高中，专科，本科，硕士，博士',
  `profession` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '专业',
  `job` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '当前工作',
  `marital` varchar(11) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT NULL COMMENT '婚姻状况 1已婚 2未婚',
  `target` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `users_city_code_index` (`city_code`),
  KEY `idx_users_account` (`account`),
  KEY `idx_users_city_code` (`city_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='前台用户';

DROP TABLE IF EXISTS `verse`;
CREATE TABLE `verse` (
  `id` int NOT NULL AUTO_INCREMENT,
  `content` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
SET FOREIGN_KEY_CHECKS=1;

