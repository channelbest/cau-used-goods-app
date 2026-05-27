package report

type Report struct {
	ID           uint64  `json:"id"`
	ReporterID   uint64  `json:"reporterId"`
	TargetType   string  `json:"targetType"`
	TargetID     uint64  `json:"targetId"`
	ReasonType   string  `json:"reasonType"`
	Description  *string `json:"description,omitempty"`
	Status       string  `json:"status"`
	HandleResult *string `json:"handleResult,omitempty"`
	HandlerID    *uint64 `json:"handlerId,omitempty"`
	HandleTime   *string `json:"handleTime,omitempty"`
	CreateTime   string  `json:"createTime"`
	UpdateTime   string  `json:"updateTime"`
}

type ReportDetail struct {
	Report
	ReporterNickname *string `json:"reporterNickname,omitempty"`
	Images           []string `json:"images,omitempty"`
}
