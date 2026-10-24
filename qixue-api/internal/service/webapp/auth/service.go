package webauth

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/repo"
	svcauth "qixue-api/internal/service/auth"
	"qixue-api/internal/service/captcha"
	"qixue-api/internal/service/site"
)

var (
	ErrUsernameRequired = errors.New("请填写用户名")
	ErrPasswordRequired = errors.New("请填写密码")
	ErrBadCredential    = errors.New("用户名或密码错误")
	ErrBlocked          = errors.New("错误次数过多，请稍后再试")
	ErrNeedCaptcha      = errors.New("请完成滑动验证")
	ErrCaptchaExpired   = errors.New("验证码已失效，请重新验证")
	ErrInactive         = errors.New("账号未认证")
	ErrLocked           = errors.New("账号已被拉黑")
	ErrCaptchaFail      = errors.New("验证失败")
	ErrCaptchaParam     = errors.New("参数错误")
	ErrServer           = errors.New("服务器发生异常")
)

type Service struct {
	users repo.UserRepository
	db    *gorm.DB
	rdb   *redis.Client
	jwt   svcauth.Issuer
	slide *captcha.Slide
	limit *svcauth.Limiter
}

func New(users repo.UserRepository, db *gorm.DB, rdb *redis.Client, jwt svcauth.Issuer, slide *captcha.Slide) *Service {
	return &Service{
		users: users, db: db, rdb: rdb, jwt: jwt, slide: slide,
		limit: svcauth.NewLimiter(rdb),
	}
}

func (s *Service) Captcha(ctx context.Context) (captcha.View, error) {
	return s.slide.Issue(ctx)
}

func (s *Service) CheckCaptcha(ctx context.Context, token string, x, y int) (string, error) {
	if token == "" {
		return "", ErrCaptchaParam
	}
	pass, err := s.slide.Verify(ctx, token, x, y)
	if errors.Is(err, captcha.ErrNotFound) || errors.Is(err, captcha.ErrMismatch) {
		return "", ErrCaptchaFail
	}
	if err != nil {
		return "", ErrServer
	}
	return pass, nil
}

type LoginResult struct {
	Token       string
	NeedCaptcha bool
}

func (s *Service) Login(ctx context.Context, host, username, password, captchaVerification string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrUsernameRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if len(username) > 64 || len(password) > 128 {
		return &LoginResult{NeedCaptcha: true}, ErrBadCredential
	}
	if s.limit.Blocked(ctx, "webapp", username) {
		return nil, ErrBlocked
	}
	if s.limit.Count(ctx, "webapp", username) >= 1 {
		if captchaVerification == "" {
			return nil, ErrNeedCaptcha
		}
		ok, err := s.slide.Consume(ctx, captchaVerification)
		if err != nil {
			return nil, ErrServer
		}
		if !ok {
			return nil, ErrCaptchaExpired
		}
	}

	host = strings.ToLower(host)
	if hname, _, err := net.SplitHostPort(host); err == nil {
		host = hname
	}
	skipCity := strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1")
	codes, err := site.Codes(ctx, s.db, s.rdb)
	if err != nil {
		return nil, ErrServer
	}
	siteCode := site.RequestCode(host, codes)

	var user *repo.LoginUser
	if skipCity || siteCode != "" {
		user, err = s.users.FindLoginByAccount(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServer
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = nil
		}
	}
	hash := ""
	if user != nil {
		hash = user.Password
	}
	passwordOK := svcauth.PasswordOK(hash, password)
	credentialOK := user != nil && passwordOK && (skipCity || strings.TrimSpace(deref(user.CityCode)) == siteCode)
	if !credentialOK {
		s.limit.RecordTTL(ctx, "webapp", username, 5*time.Minute)
		return &LoginResult{NeedCaptcha: true}, ErrBadCredential
	}
	if user.IsActive == -1 {
		return nil, ErrInactive
	}
	if user.IsLock == 1 {
		return nil, ErrLocked
	}
	s.limit.Clear(ctx, "webapp", username)
	token, _, err := s.jwt.Sign(user.ID, time.Now())
	if err != nil {
		return nil, ErrServer
	}
	return &LoginResult{Token: token}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
