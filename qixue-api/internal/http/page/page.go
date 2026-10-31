package page

type Body struct {
	Current int   `json:"current"`
	List    any   `json:"list"`
	Size    int   `json:"size"`
	Total   int64 `json:"total"`
}

type MoreBody struct {
	Current int   `json:"current"`
	List    any   `json:"list"`
	Size    int   `json:"size"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"hasMore"`
}

func Clamp(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 200 {
		size = 200
	}
	return page, size
}

func New(current, size int, total int64, list any) Body {
	if list == nil {
		list = []any{}
	}
	return Body{Current: current, List: list, Size: size, Total: total}
}

func NewMore(current, size int, total int64, list any) MoreBody {
	if list == nil {
		list = []any{}
	}
	return MoreBody{
		Current: current,
		List:    list,
		Size:    size,
		Total:   total,
		HasMore: int64(current)*int64(size) < total,
	}
}
