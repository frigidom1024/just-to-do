package daily_note

import domainerr "todolist/internal/pkg/domainerr"

// 领域错误定义
var (
	// ErrDailyNoteNotFound 表示每日笔记不存在
	ErrDailyNoteNotFound = domainerr.BusinessError{
		Code:    "DAILY_NOTE_NOT_FOUND",
		Type:    domainerr.NotFoundError,
		Message: "每日笔记不存在",
	}

	// ErrDailyNoteContentEmpty 表示每日笔记内容为空
	ErrDailyNoteContentEmpty = domainerr.BusinessError{
		Code:    "DAILY_NOTE_CONTENT_EMPTY",
		Type:    domainerr.ValidationError,
		Message: "每日笔记内容不能为空",
	}

	// ErrDailyNoteAlreadyExists 表示当日已存在每日笔记
	ErrDailyNoteAlreadyExists = domainerr.BusinessError{
		Code:    "DAILY_NOTE_ALREADY_EXISTS",
		Type:    domainerr.ConflictError,
		Message: "当日已存在每日笔记",
	}

	// ErrDailyNoteUpdateFailed 表示每日笔记更新失败
	ErrDailyNoteUpdateFailed = domainerr.BusinessError{
		Code:    "DAILY_NOTE_UPDATE_FAILED",
		Type:    domainerr.InternalError,
		Message: "每日笔记更新失败",
	}

	// ErrDailyNoteDeleteFailed 表示每日笔记删除失败
	ErrDailyNoteDeleteFailed = domainerr.BusinessError{
		Code:    "DAILY_NOTE_DELETE_FAILED",
		Type:    domainerr.InternalError,
		Message: "每日笔记删除失败",
	}
)
