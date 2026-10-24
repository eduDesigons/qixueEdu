package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// CoinRepository covers user coin/credit grants and idempotency checks.
type CoinRepository interface {
	CountFieldInRange(ctx context.Context, uid int, field, start, end string) (int64, error)
	CountFieldRemarkLike(ctx context.Context, uid int, field, like string) (int64, error)
	AddCoins(ctx context.Context, uid int, field string, coins int, remark string) error
	AddCredit(ctx context.Context, uid int, field string, credit int, remark string) error
}

type coinRepo struct{ db *gorm.DB }

func NewCoin(db *gorm.DB) CoinRepository {
	return &coinRepo{db: db}
}

func (r *coinRepo) CountFieldInRange(ctx context.Context, uid int, field, start, end string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_coins_records").
		Where("user_id = ? AND field = ? AND deleted_at IS NULL AND created_at >= ? AND created_at <= ?",
			uid, field, start, end).
		Count(&n).Error
	return n, err
}

func (r *coinRepo) CountFieldRemarkLike(ctx context.Context, uid int, field, like string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_coins_records").
		Where("user_id = ? AND field = ? AND remark LIKE ? AND deleted_at IS NULL", uid, field, like).
		Count(&n).Error
	return n, err
}

func (r *coinRepo) AddCoins(ctx context.Context, uid int, field string, coins int, remark string) error {
	if uid <= 0 || coins <= 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET coins = coins + ? WHERE id = ?", coins, uid).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Exec(
			"INSERT INTO user_coins_records (user_id, field, `sum`, remark, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			uid, field, coins, remark, now, now,
		).Error
	})
}

func (r *coinRepo) AddCredit(ctx context.Context, uid int, field string, credit int, remark string) error {
	if uid <= 0 || credit == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET credit = credit + ? WHERE id = ?", credit, uid).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Exec(
			"INSERT INTO user_credit_records (user_id, field, `sum`, remark, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			uid, field, credit, remark, now, now,
		).Error
	})
}
