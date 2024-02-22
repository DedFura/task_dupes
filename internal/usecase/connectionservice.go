package usecase

import (
	"task/internal/domain/models"
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
	ch1 := make(chan []models.Connection)
	ch2 := make(chan []models.Connection)
	errCh := make(chan error)

	go func() {
		connections, err := cs.repo.FindConnectionsByUserID(userID1)
		if err != nil {
			errCh <- err
			return
		}
		ch1 <- connections
	}()

	go func() {
		connections, err := cs.repo.FindConnectionsByUserID(userID2)
		if err != nil {
			errCh <- err
			return
		}
		ch2 <- connections
	}()

	connections1 := <-ch1
	connections2 := <-ch2

	select {
	case err := <-errCh:
		return false, err
	default:
	}

	ipMap := make(map[string]struct{})
	for _, conn := range connections1 {
		ipMap[conn.IPAddr] = struct{}{}
	}

	for _, conn := range connections2 {
		if _, found := ipMap[conn.IPAddr]; found {
			return true, nil
		}
	}

	return false, nil
}
