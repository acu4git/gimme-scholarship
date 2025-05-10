package database

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/acu4git/gimme-scholarship/internal/domain/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
)

const (
	tableUsers              = "users"
	tableScholarships       = "scholarships"
	tableScholarshipTargets = "scholarship_targets"
	tableEducationLevels    = "education_levels"
)

var (
	host     string
	port     string
	username string
	password string
	dbname   string
)

type Database struct {
	sess *dbr.Session
}

func NewDatabase() (*Database, error) {
	if host = os.Getenv("DB_HOST"); host == "" {
		host = "localhost"
	}
	if port = os.Getenv("DB_PORT"); port == "" {
		port = "3308"
	}
	if username = os.Getenv("DB_USER"); username == "" {
		username = "root"
	}
	if password = os.Getenv("DB_PASSWORD"); password == "" {
		password = "root"
	}
	if dbname = os.Getenv("DB_NAME"); dbname == "" {
		dbname = "gimme_scholarship"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	conn, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		log.Println("failed at dbr.Open()")
		return nil, err
	}

	sess := conn.NewSession(nil)

	return &Database{sess: sess}, nil
}

func (db *Database) CreateUser(user model.User) error {
	return nil
}

func (db *Database) GetScholarships() ([]model.Scholarship, error) {
	res := make([]model.Scholarship, 0)

	type scholarshipWithTarget struct {
		scholarship
		Target string `db:"target"`
	}
	flat := make([]scholarshipWithTarget, 0)
	stmt := db.sess.Select(fmt.Sprintf("%s.*, %s.name AS target", tableScholarships, tableEducationLevels)).
		From(tableScholarships).
		Join(tableScholarshipTargets, fmt.Sprintf("%s.id = %s.scholarship_id", tableScholarships, tableScholarshipTargets)).
		Join(tableEducationLevels, fmt.Sprintf("%s.id = %s.education_level_id", tableEducationLevels, tableScholarshipTargets))

	if _, err := stmt.Load(&flat); err != nil {
		return nil, err
	}

	scholarshipMap := make(map[int64]*model.Scholarship)
	for _, f := range flat {
		s, ok := scholarshipMap[f.ID]
		if !ok {
			s = &model.Scholarship{
				ID:             f.ID,
				Name:           f.Name,
				Address:        f.Address,
				TargetDetail:   f.TargetDetail,
				AmountDetail:   f.AmountDetail,
				TypeDetail:     f.TypeDetail,
				CapacityDetail: f.CapacityDetail,
				DeadlineDetail: f.DeadlineDetail,
				ContactPoint:   f.ContactPoint,
				Remark:         f.Remark,
				PostingDate:    f.PostingDate.Format("2006-01-02"), // 文字列に変換
			}
			scholarshipMap[f.ID] = s
		}
		s.Targets = append(s.Targets, f.Target)
	}

	for _, s := range scholarshipMap {
		res = append(res, *s)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res, nil
}
