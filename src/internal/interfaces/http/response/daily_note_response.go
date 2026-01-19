// Package response 提供 HTTP 响应的 DTO 结构。
//
// 所有的响应结构都用于序列化为 JSON 返回给客户端。
// 这些结构与领域实体分离，避免领域模型泄露到接口层。
package response

import (
	"time"

	"todolist/internal/interfaces/dto"
)

// DailyNoteResponse 每日笔记响应。
//
// 包含每日笔记的详细信息。
type DailyNoteResponse struct {
	// ID 每日笔记唯一标识
	ID int64 `json:"id"`

	// UserID 所属用户ID
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

// DailyNoteListResponse 每日笔记列表响应。
//
// 包含每日笔记列表和分页信息。
type DailyNoteListResponse struct {
	// Data 每日笔记列表
	Data []DailyNoteResponse `json:"data"`

	// Pagination 分页信息
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationResponse 分页信息响应。
//
// 包含分页查询的元数据。
type PaginationResponse struct {
	// Total 总记录数
	Total int64 `json:"total"`

	// Page 当前页码
	Page int `json:"page"`

	// PageSize 每页大小
	PageSize int `json:"page_size"`

	// TotalPages 总页数
	TotalPages int `json:"total_pages"`
}

// ToDailyNoteResponse 将每日笔记DTO转换为响应对象。
//
// 参数：
//
//	dailyNoteDTO - 每日笔记数据传输对象
//
// 返回：
//
//	DailyNoteResponse - HTTP 响应对象
func ToDailyNoteResponse(dailyNoteDTO dto.DailyNoteDTO) DailyNoteResponse {
	return DailyNoteResponse{
		ID:        dailyNoteDTO.ID,
		UserID:    dailyNoteDTO.UserID,
		NoteDate:  dailyNoteDTO.NoteDate,
		Content:   dailyNoteDTO.Content,
		CreatedAt: dailyNoteDTO.CreatedAt,
		UpdatedAt: dailyNoteDTO.UpdatedAt,
	}
}

// ToDailyNoteListResponse 将每日笔记分页DTO转换为响应对象。
//
// 参数：
//
//	dailyNotePageDTO - 每日笔记分页数据传输对象
//
// 返回：
//
//	DailyNoteListResponse - HTTP 响应对象
func ToDailyNoteListResponse(dailyNotePageDTO dto.DailyNotePageDTO) DailyNoteListResponse {
	// 转换数据列表
	data := make([]DailyNoteResponse, len(dailyNotePageDTO.Data))
	for i, dto := range dailyNotePageDTO.Data {
		data[i] = ToDailyNoteResponse(dto)
	}

	// 转换分页信息
	pagination := PaginationResponse{
		Total:      dailyNotePageDTO.Pagination.Total,
		Page:       dailyNotePageDTO.Pagination.Page,
		PageSize:   dailyNotePageDTO.Pagination.PageSize,
		TotalPages: dailyNotePageDTO.Pagination.TotalPages,
	}

	return DailyNoteListResponse{
		Data:       data,
		Pagination: pagination,
	}
}
