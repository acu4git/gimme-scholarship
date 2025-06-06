package handler

import "github.com/acu4git/gimme-scholarship/internal/domain/model"

type scholarshipOutput struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	Targets        []string `json:"targets"`
	TargetDetail   string   `json:"target_detail"`
	AmountDetail   string   `json:"amount_detail"`
	TypeDetail     string   `json:"type_detail"`
	CapacityDetail string   `json:"capacity_detail"`
	Deadline       string   `json:"deadline"`
	DeadlineDetail string   `json:"deadline_detail"`
	ContactPoint   string   `json:"contact_point"`
	Remark         string   `json:"remark"`
	PostingDate    string   `json:"posting_date"`
	IsFavorite     bool     `json:"is_favorite"`
}

type GetScholarshipsOutput []scholarshipOutput

func toGetScholarshipsOutput(scholarships []model.Scholarship, favoriteMap map[int64]bool) GetScholarshipsOutput {
	res := make(GetScholarshipsOutput, 0)
	for _, s := range scholarships {
		res = append(res, scholarshipOutput{
			ID:             s.ID,
			Name:           s.Name,
			Address:        s.Address,
			Targets:        s.Targets,
			TargetDetail:   s.TargetDetail,
			AmountDetail:   s.AmountDetail,
			TypeDetail:     s.TypeDetail,
			CapacityDetail: s.CapacityDetail,
			Deadline:       s.Deadline,
			DeadlineDetail: s.DeadlineDetail,
			ContactPoint:   s.ContactPoint,
			Remark:         s.Remark,
			PostingDate:    s.PostingDate,
			IsFavorite:     favoriteMap[s.ID],
		})
	}
	return res
}
