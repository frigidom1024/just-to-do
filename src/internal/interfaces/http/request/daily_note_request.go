package request

// DailyNoteRequest 每日笔记请求结构
//
// 用于创建和更新每日笔记
// 包含每日笔记的内容

type DailyNoteRequest struct {
	// Content 笔记内容，不能为空
	Content string `json:"content" validate:"required"`
}

// DailyNoteListRequest 每日笔记列表请求结构
//
// 用于分页查询每日笔记列表
// 包含分页查询参数

type DailyNoteListRequest struct {
	// Page 页码，默认为1
	Page int `json:"page" form:"page"`

	// PageSize 每页大小，默认为10，最大为50
	PageSize int `json:"page_size" form:"page_size"`
}

// EmptyRequest 空请求结构
//
// 用于不需要请求体的请求

type EmptyRequest struct {}

