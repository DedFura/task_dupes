package repositories

import (
	"database/sql"
	"fmt"
	"task/configs"
	"task/internal/domain/models"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type ConnectionRepository interface {
	FindConnectionsByUserIDs(userID1, userID2 int64) ([]models.Connection, []models.Connection, error)
	FindConnectionsByUserID(userID int64) ([]models.Connection, error)
}

type connectionRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewConnectionRepository(config *configs.Config, logger *logrus.Logger) (ConnectionRepository, error) {
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User,
		config.Database.Password, config.Database.DBName, config.Database.SSLMode)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	return &connectionRepository{
		db:  db,
		log: logger}, nil
}

func (cr *connectionRepository) FindConnectionsByUserIDs(userID1, userID2 int64) ([]models.Connection, []models.Connection, error) {
	query := `
	SELECT DISTINCT ip_addr FROM conn_log
	WHERE user_id IN ($1, $2);
`
	rows, err := cr.db.Query(query, userID1, userID2)
	if err != nil {
		cr.log.Error("Failed to execute query:", err)
		return nil, nil, err
	}
	defer rows.Close()

	var connections1 []models.Connection
	var connections2 []models.Connection

	for rows.Next() {
		var ipAddr string
		if err := rows.Scan(&ipAddr); err != nil {
			cr.log.Error("Failed to scan row:", err)
			return nil, nil, err
		}

		connection := models.Connection{
			UserID: userID1,
			IPAddr: ipAddr,
		}
		connections1 = append(connections1, connection)

		connection.UserID = userID2
		connections2 = append(connections2, connection)
	}
	if err := rows.Err(); err != nil {
		cr.log.Error("Error occurred while iterating through rows:", err)
		return nil, nil, err
	}

	return connections1, connections2, nil
}

func (cr *connectionRepository) FindConnectionsByUserID(userID int64) ([]models.Connection, error) {
	query := `
	SELECT DISTINCT ip_addr FROM conn_log
	WHERE user_id = $1;
`
	rows, err := cr.db.Query(query, userID)
	if err != nil {
		cr.log.Error("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var connections []models.Connection

	for rows.Next() {
		var ipAddr string
		if err := rows.Scan(&ipAddr); err != nil {
			cr.log.Error("Failed to scan row:", err)
			return nil, err
		}

		connection := models.Connection{
			UserID: userID,
			IPAddr: ipAddr,
		}
		connections = append(connections, connection)
	}
	if err := rows.Err(); err != nil {
		cr.log.Error("Error occurred while iterating through rows:", err)
		return nil, err
	}

	return connections, nil
}
