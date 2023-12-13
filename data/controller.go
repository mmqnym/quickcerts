package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	cfg "QuickCertS/configs"
	"QuickCertS/model"
	"QuickCertS/utils"
)

var db *sql.DB = nil

// Connect to the specified database.
func ConnectDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_CONFIG.HOST, cfg.DB_CONFIG.PORT, cfg.DB_CONFIG.USER, cfg.DB_CONFIG.PWD, cfg.DB_CONFIG.DB_NAME)
	var err error
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		db = nil
		return errors.New("failed to connect the database")
	}

	err = db.Ping()

	if err != nil {
		db = nil
		return errors.New("failed to access the database")
	}

	return nil
}

// Disconnect from the database.
func DisconnectDB() error {
	if db == nil {
		return errors.New("currently not connecting the database")
	}

	err := db.Close()
	db = nil
	return err
}

// Add a new S/N into the database.
func AddNewSN(sn string) error {
	if db == nil {
		return errors.New("currently not connecting the database")
	}

	stmt, err := db.Prepare("INSERT INTO certs (sn, key, note) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(sn, sql.NullString{}, sql.NullString{})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("the s/n already exists")
		}
		return err
	}

	return err
}

// Add new S/N(s) into the database.
func AddNewSNs(snList []string) error {
	if db == nil {
		return errors.New("currently not connecting the database")
	}

	var valuesStrings []string

	for _, sn := range snList {
		valuesStrings = append(valuesStrings, fmt.Sprintf("('%s', NULL, NULL)", sn))
	}

	query := fmt.Sprintf("INSERT INTO certs (sn, key, note) VALUES %s;", strings.Join(valuesStrings, ", "))

	_, err := db.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("some s/ns already exist")
		}
		return err
	}

	return nil
}

// Check if the given S/N exists in the database.
func IsSNExist(sn string) (bool, error) {
	if db == nil {
		return false, errors.New("currently not connecting the database")
	}

	query := "SELECT EXISTS (SELECT 1 FROM certs WHERE sn = $1)"

	var exists bool
	err := db.QueryRow(query, sn).Scan(&exists)
	if !exists {
		return false, errors.New("the s/n does not exist")
	}

	if err != nil {
		return false, err
	}

	return exists, nil
}

// Bind the given serial number to the key. (Update the key field corresponding to the given S/N.)
func BindSNWithKey(sn string, key string) error {
	if db == nil {
		return errors.New("currently not connecting the database")
	}

	stmt, err := db.Prepare(`
		UPDATE certs
		SET key = $1 
		WHERE sn = $2 
		AND (key IS NULL OR key = $1)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(key, sn)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("the s/n does not exist or has already been used")
	}

	return err
}

// Get the remaining trial period for the given key.
//
// If the key is not found, allow for temporary permit application.
func GetTemporaryPermitExpiredTime(key string) (int64, error) {
	if db == nil {
		return 0, errors.New("currently not connecting the database")
	}

	var expiration time.Time

	query := "SELECT expiration FROM temporary_permits WHERE key = $1"
	err := db.QueryRow(query, key).Scan(&expiration)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("allowed new key: %s", key)
	} else if err != nil {
		return 0, err
	}

	durationLeft := (expiration.Unix()) - time.Now().Unix()

	if durationLeft < 0 {
		return 0, nil
	}
	return durationLeft, nil
}

// Providing temporary usage rights to trial clients.
func AddTemporaryPermit(key string) (int64, error) {
	if db == nil {
		return 0, errors.New("currently not connecting the database")
	}

	stmt, err := db.Prepare("INSERT INTO temporary_permits (key, expiration) VALUES ($1, $2)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	timeUnit, err := utils.TimeUnitStrToTimeDuration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT)
	if err != nil {
		return 0, err
	}

	expiration := time.Now().Add(time.Duration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME) * timeUnit)
	_, err = stmt.Exec(key, expiration)

	if err != nil {
		return 0, err
	}

	timeLeft := expiration.Unix() - time.Now().Unix()

	return timeLeft, nil
}

// Get all certificate records in the database.
func GetAllCerts() ([]model.Cert, error) {
	if db == nil {
		return nil, errors.New("currently not connecting the database")
	}

	query := "SELECT * FROM certs"

	rows, err := db.Query(query)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, nil
	}

	defer rows.Close()

	var certs []model.Cert

	for rows.Next() {
		var cert model.Cert
		var tmpKey sql.NullString
		var tmpNote sql.NullString
		if err := rows.Scan(&cert.SerialNumber, &tmpKey, &tmpNote); err != nil {
			return nil, err
		}

		cert.Key = tmpKey.String
		cert.Note = tmpNote.String
		certs = append(certs, cert)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return certs, nil
}

// Get avaliable S/N in the database.
func GetAvaliableSN() ([]string, error) {
	if db == nil {
		return nil, errors.New("currently not connecting the database")
	}

	query := "SELECT sn FROM certs where key is NULL"

	rows, err := db.Query(query)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, nil
	}

	defer rows.Close()

	var res []string

	for rows.Next() {
		var sn string
		if err := rows.Scan(&sn); err != nil {
			return nil, err
		}

		res = append(res, sn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Update the note field corresponding to the given S/N.
func UpdateCertNote(sn string, note string) error {
	if db == nil {
		return errors.New("currently not connecting the database")
	}

	stmt, err := db.Prepare("UPDATE certs SET note = $1 WHERE sn = $2")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(note, sn)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("the s/n does not exist")
	}

	return nil
}
