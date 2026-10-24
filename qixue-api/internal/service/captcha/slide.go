package captcha

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/slide"
)

const (
	padding = 5
	ttl     = 5 * time.Minute
)

var (
	ErrNotFound = errors.New("captcha not found")
	ErrMismatch = errors.New("captcha mismatch")
)

type View struct {
	Token       string `json:"token"`
	Image       string `json:"image"`
	Thumb       string `json:"thumb"`
	ThumbX      int    `json:"thumbX"`
	ThumbY      int    `json:"thumbY"`
	ThumbWidth  int    `json:"thumbWidth"`
	ThumbHeight int    `json:"thumbHeight"`
}

type Slide struct {
	engine slide.Captcha
	rdb    *redis.Client
}

type answer struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func New(rdb *redis.Client) (*Slide, error) {
	backgrounds, err := imagesv2.GetImages()
	if err != nil {
		return nil, err
	}
	graphs, err := tiles.GetTiles()
	if err != nil {
		return nil, err
	}
	images := make([]*slide.GraphImage, 0, len(graphs))
	for _, graph := range graphs {
		images = append(images, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			ShadowImage:  graph.ShadowImage,
			MaskImage:    graph.MaskImage,
		})
	}

	builder := slide.NewBuilder(slide.WithImageSize(option.Size{Width: 300, Height: 220}))
	builder.SetResources(
		slide.WithGraphImages(images),
		slide.WithBackgrounds(backgrounds),
	)
	return &Slide{engine: builder.Make(), rdb: rdb}, nil
}

func (s *Slide) Issue(ctx context.Context) (View, error) {
	data, err := s.engine.Generate()
	if err != nil {
		return View{}, err
	}
	block := data.GetData()
	if block == nil {
		return View{}, errors.New("empty captcha")
	}
	image64, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return View{}, err
	}
	thumb64, err := data.GetTileImage().ToBase64()
	if err != nil {
		return View{}, err
	}
	token, err := randomToken()
	if err != nil {
		return View{}, err
	}
	raw, err := json.Marshal(answer{X: block.X, Y: block.Y})
	if err != nil {
		return View{}, err
	}
	if err := s.rdb.Set(ctx, slideKey(token), raw, ttl).Err(); err != nil {
		return View{}, err
	}
	return View{
		Token:       token,
		Image:       image64,
		Thumb:       thumb64,
		ThumbX:      block.DX,
		ThumbY:      block.DY,
		ThumbWidth:  block.Width,
		ThumbHeight: block.Height,
	}, nil
}

func (s *Slide) Verify(ctx context.Context, token string, x, y int) (string, error) {
	key := slideKey(token)
	raw, err := s.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		return "", err
	}
	var ans answer
	if err := json.Unmarshal(raw, &ans); err != nil {
		return "", err
	}
	if !slide.Validate(x, y, ans.X, ans.Y, padding) {
		return "", ErrMismatch
	}
	pass, err := randomToken()
	if err != nil {
		return "", err
	}
	if err := s.rdb.Set(ctx, passKey(pass), "1", ttl).Err(); err != nil {
		return "", err
	}
	return pass, nil
}

func (s *Slide) Consume(ctx context.Context, pass string) (bool, error) {
	n, err := s.rdb.Del(ctx, passKey(pass)).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func slideKey(token string) string { return "qx:captcha:" + token }

func passKey(token string) string { return "qx:captcha_pass:" + token }

func randomToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
