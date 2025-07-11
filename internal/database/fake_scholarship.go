package database

import (
	"time"

	"github.com/acu4git/gimme-scholarship/internal/domain/model"
)

func (fdb *FakeDatabase) TestInitScholarships() error {
	data := []model.Scholarship{
		{1, "テスト奨学金", "URL: https://example.com", []string{"学部生", "大学院生"}, "", "【月額】10万円", "【給付】無利子", "1名程度", time.Now().Format("2006-01-02"), "当日締切", "", "", "2025-01-01"},
		{2, "テスト財団", "URL: https://example.com", []string{"学部生", "大学院生"}, "", "【月額】10万円", "【給付】無利子", "1名程度", time.Now().AddDate(0, 0, 7).Format("2006-01-02"), "来週締切", "", "", "2025-01-01"},
		{3, "テスト支援金", "URL: https://example.com", []string{"学部生", "大学院生"}, "", "【月額】10万円", "【給付】無利子", "1名程度", time.Now().AddDate(0, 0, 7).Format("2006-01-02"), "来週締切", "", "", "2025-01-01"},
	}

	tx, err := fdb.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, d := range data {
		deadline, err := time.Parse("2006-01-02", d.Deadline)
		if err != nil {
			return err
		}
		postingDate, err := time.Parse("2006-01-02", d.PostingDate)
		if err != nil {
			return err
		}

		s := scholarship{
			ID:             d.ID,
			Name:           d.Name,
			Address:        d.Address,
			TargetDetail:   d.TargetDetail,
			AmountDetail:   d.AmountDetail,
			TypeDetail:     d.TypeDetail,
			CapacityDetail: d.CapacityDetail,
			Deadline:       deadline,
			DeadlineDetail: d.DeadlineDetail,
			ContactPoint:   d.ContactPoint,
			Remark:         d.Remark,
			PostingDate:    postingDate,
		}
		if _, err := tx.InsertInto(tableScholarships).Record(s).Exec(); err != nil {
			return err
		}
		for _, target := range d.Targets {
			el := educationLevel{}
			if err := tx.Select("*").From(tableEducationLevels).Where("name = ?", target).LoadOne(&el); err != nil {
				return err
			}
			if _, err := tx.InsertInto(tableScholarshipTargets).Pair("scholarship_id", d.ID).Pair("education_level_id", el.ID).Exec(); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
