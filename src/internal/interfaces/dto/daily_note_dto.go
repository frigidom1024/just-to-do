package dto

import (
	"time"

	"todolist/internal/domain/daily_note"
)

// DailyNoteDTO 每日笔记数据传输对象
type DailyNoteDTO struct {
	// ID 每日笔记唯一标识
	ID int64 `json:"id"`

	// UserID 用户ID
	UserID int64 `json:"user_id"`

	// NoteDate 笔记日期
	NoteDate time.Time `json:"note_date"`

	// Content 笔记内容
	Content string `json:"content"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationDTO 分页信息数据传输对象
type PaginationDTO struct {
	// Total 总记录数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页大小
	PageSize int `json:"page_size"`

	// TotalPages 总页数
	TotalPages int `json:"total_pages"`
}

// DailyNotePageDTO 每日笔记分页结果数据传输对象
type DailyNotePageDTO struct {
	// Data 每日笔记列表
	Data []DailyNoteDTO `json:"data"`

	// Pagination 分页信息
	Pagination PaginationDTO `json:"pagination"`
}

// ToDailyNoteDTO 将每日笔记领域实体转换为DTO
func ToDailyNoteDTO(entity daily_note.DailyNoteEntity) DailyNoteDTO {
	return DailyNoteDTO{
		ID:        entity.GetID(),
		UserID:    entity.GetUserID(),
		NoteDate:  entity.GetNoteDate(),
		Content:   entity.GetContent(),
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
	}
}

// ToDailyNotePageDTO 将每日笔记领域实体列表转换为分页DTO
func ToDailyNotePageDTO(entities []daily_note.DailyNoteEntity, total int64, page, pageSize int) DailyNotePageDTO {
	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	// 转换实体列表为DTO列表
	dtos := make([]DailyNoteDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = ToDailyNoteDTO(entity)
	}

	return DailyNotePageDTO{
		Data: dtos,
		Pagination: PaginationDTO{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}
