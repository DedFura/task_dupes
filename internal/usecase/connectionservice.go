package usecase

import (
	"task/internal/domain/repositories"

	"github.com/sirupsen/logrus"
)

type ConnectionChecker interface {
	CheckDupes(userID1, userID2 int64) (bool, error)
}

type connectionService struct {
	repo repositories.ConnectionRepository
	log  *logrus.Logger
}

func NewConnectionService(repo repositories.ConnectionRepository, logger *logrus.Logger) ConnectionChecker {
	return &connectionService{
		repo: repo,
		log:  logger,
	}
}

func (cs *connectionService) CheckDupes(userID1, userID2 int64) (bool, error) {
	connections1, connections2, err := cs.repo.FindConnectionsByUserIDs(userID1, userID2)
	if err != nil {
		return false, err
	}

	ipMap := make(map[string]bool)
	for _, conn := range connections1 {
		ipMap[conn.IPAddr] = true
	}

	for _, conn := range connections2 {
		if ipMap[conn.IPAddr] {
			return true, nil
		}
	}

	return false, nil
}
