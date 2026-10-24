package applog

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

// Setup wires std log、slog、gin 到按天切割的日志文件；同时仍打一份到 stdout（方便 journalctl）。
func Setup(dir string, keepDays int) (*DailyWriter, error) {
	w, err := NewDailyWriter(dir, "qixue-api", keepDays)
	if err != nil {
		return nil, err
	}
	out := io.MultiWriter(os.Stdout, w)

	log.SetOutput(out)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("")

	slog.SetDefault(slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	gin.DefaultWriter = out
	gin.DefaultErrorWriter = out
	return w, nil
}
