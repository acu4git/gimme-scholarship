package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/acu4git/gimme-scholarship/internal/domain/model"
	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
)

const (
	tlsKey = "custom"

	tableUsers              = "users"
	tableScholarships       = "scholarships"
	tableScholarshipTargets = "scholarship_targets"
	tableEducationLevels    = "education_levels"
	tableUserFavorites      = "user_favorites"
)

var (
	host     string
	port     string
	username string
	password string
	dbname   string

	registeredTLSKey = registerTLSConfig("RDS_CERT_FILE_PATH")
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
	if len(registeredTLSKey) > 0 {
		dsn = fmt.Sprintf("%s&tls=%s", dsn, registeredTLSKey)
	}
	conn, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		log.Println("failed at dbr.Open()")
		return nil, err
	}

	sess := conn.NewSession(nil)

	return &Database{sess: sess}, nil
}

func registerTLSConfig(envKey string) string {
	pemPath := os.Getenv(envKey)
	if len(pemPath) == 0 {
		return ""
	}

	certPool := x509.NewCertPool()
	pem, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatal(err)
	}

	if ok := certPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("error: failed to append certs from pem")
	}

	if err := mysql.RegisterTLSConfig(tlsKey, &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}); err != nil {
		log.Fatal(err)
	}

	return tlsKey
}

func (db *Database) CreateUser(input repository.UserInput) error {
	tx, err := db.sess.Begin()
	defer tx.RollbackUnlessCommitted()
	if err != nil {
		return err
	}

	var el educationLevel
	if err := tx.Select("*").
		From(tableEducationLevels).
		Where("name = ?", input.Level).LoadOne(&el); err != nil {
		return err
	}

	u := user{
		ID:               input.ID,
		Email:            input.Email,
		EducationLevelID: el.ID,
		Grade:            input.Grade,
		AcceptEmail:      input.AcceptEmail,
	}
	if _, err := tx.InsertInto(tableUsers).Columns(u.columns()...).Record(&u).Exec(); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *Database) FindScholarships(option repository.FilterOption) ([]model.Scholarship, map[int64]bool, error) {
	scholarships := make([]model.Scholarship, 0)

	tx, err := db.sess.Begin()
	defer tx.RollbackUnlessCommitted()
	if err != nil {
		return nil, nil, err
	}

	type scholarshipWithTarget struct {
		scholarship
		Target string `db:"target"`
	}
	flat := make([]scholarshipWithTarget, 0)
	stmt := tx.Select(fmt.Sprintf("%s.*, %s.name AS target", tableScholarships, tableEducationLevels)).
		From(tableScholarships).
		Join(tableScholarshipTargets, fmt.Sprintf("%s.id = %s.scholarship_id", tableScholarships, tableScholarshipTargets)).
		Join(tableEducationLevels, fmt.Sprintf("%s.id = %s.education_level_id", tableEducationLevels, tableScholarshipTargets))

	if option.Target != "" {
		stmt = stmt.Where(fmt.Sprintf("%s.name = ?", tableEducationLevels), option.Target)
	}

	if option.Type != "" {
		stmt = stmt.Where(fmt.Sprintf("%s.type_detail LIKE ?", tableScholarships), fmt.Sprintf("[%s]%%", option.Type))
	}

	if _, err := stmt.Load(&flat); err != nil {
		return nil, nil, err
	}

	scholarshipMap := make(map[int64]*model.Scholarship)
	favoriteMap := make(map[int64]bool)
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
				Deadline:       f.Deadline.Format("2006-01-02"),
				DeadlineDetail: f.DeadlineDetail,
				ContactPoint:   f.ContactPoint,
				Remark:         f.Remark,
				PostingDate:    f.PostingDate.Format("2006-01-02"), // 文字列に変換
			}
			scholarshipMap[f.ID] = s
		}
		s.Targets = append(s.Targets, f.Target)
		favoriteMap[f.ID] = false
	}

	for _, s := range scholarshipMap {
		scholarships = append(scholarships, *s)
	}

	sort.Slice(scholarships, func(i, j int) bool {
		return scholarships[i].Deadline < scholarships[j].Deadline
	})

	if option.UserID != nil {
		ufs := make([]userFavorite, 0)
		if _, err := tx.Select("*").
			From(tableUserFavorites).
			Where("user_id = ?", option.UserID).
			Load(&ufs); err != nil {
			return nil, nil, err
		}
		for _, uf := range ufs {
			favoriteMap[uf.ScholarshipID] = true
		}
	}

	return scholarships, favoriteMap, tx.Commit()
}

func (db *Database) UserFavoriteAction(input repository.UserFavoriteInput) error {
	tx, err := db.sess.Begin()
	defer tx.RollbackUnlessCommitted()
	if err != nil {
		return err
	}

	switch input.Mode {
	case "REGISTER":
		if _, err := tx.InsertInto(tableUserFavorites).
			Pair("user_id", input.UserID).
			Pair("scholarship_id", input.ScholarshipID).
			Exec(); err != nil {
			return err
		}
	case "DELETE":
		if _, err := tx.DeleteFrom(tableUserFavorites).
			Where("user_id = ? AND scholarship_id = ?", input.UserID, input.ScholarshipID).
			Exec(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("error: invalid action mode (%s)", input.Mode)
	}

	return tx.Commit()
}

func (db *Database) FindUsersToNotifyForUpcomingDeadlines() (map[string][]model.Scholarship, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	deadlineDate := todayStart.AddDate(0, 0, 7)

	results := make([]UserScholarshipNotification, 0)
	tx, err := db.sess.Begin()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Select(
		fmt.Sprintf("%s.id AS user_id", tableUsers),
		fmt.Sprintf("%s.name AS user_name", tableUsers),
		fmt.Sprintf("%s.email AS user_email", tableUsers),
		fmt.Sprintf("%s.*", tableScholarships),
	).
		From(tableScholarships).
		Join(tableUserFavorites, fmt.Sprintf("%s.id = %s.scholarship_id", tableScholarships, tableUserFavorites)).
		Join(tableUsers, fmt.Sprintf("%s.id = %s.user_id", tableUsers, tableUserFavorites)).
		Where(fmt.Sprintf("%s.deadline = ?", tableScholarships), deadlineDate).
		Load(&results); err != nil {
		return nil, err
	}

	// key: email address, value: scholarship info
	userScholarships := make(map[string][]model.Scholarship)
	for _, info := range results {
		if _, exist := userScholarships[info.UserEmail]; !exist {
			userScholarships[info.UserEmail] = make([]model.Scholarship, 0)
		}
		scholarship := model.Scholarship{
			ID:             info.scholarship.ID,
			Name:           info.scholarship.Name,
			Address:        info.scholarship.Address,
			TargetDetail:   info.scholarship.TargetDetail,
			AmountDetail:   info.scholarship.AmountDetail,
			TypeDetail:     info.scholarship.TypeDetail,
			CapacityDetail: info.scholarship.CapacityDetail,
			Deadline:       info.scholarship.Deadline.Format("2006-01-02"),
			DeadlineDetail: info.scholarship.DeadlineDetail,
			ContactPoint:   info.scholarship.ContactPoint,
			Remark:         info.scholarship.Remark,
		}
		userScholarships[info.UserEmail] = append(userScholarships[info.UserEmail], scholarship)
	}

	return userScholarships, nil
}
