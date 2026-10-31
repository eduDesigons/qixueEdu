package async

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

func SkipNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// Run 并发执行互不依赖的查询，任一出错则返回第一个错误。
func Run(fns ...func() error) error {
	n := len(fns)
	if n == 0 {
		return nil
	}
	if n == 1 {
		if fns[0] == nil {
			return nil
		}
		return fns[0]()
	}
	errs := make([]error, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := range fns {
		go func(i int) {
			defer wg.Done()
			if fns[i] != nil {
				errs[i] = fns[i]()
			}
		}(i)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
