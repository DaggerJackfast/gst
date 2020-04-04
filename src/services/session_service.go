package services

import (
	"github.com/DaggerJackfast/gst/src/domains"
	"github.com/DaggerJackfast/gst/src/repositories"
)

type SessionService interface {
	CreateSession(session *domains.Session) error
	UpdateSession(session *domains.Session) error
	DeleteSession(session *domains.Session) error
	GetSession(user *domains.User, token string) (*domains.Session, error)
}

type sessionService struct {
	sessionRepo repositories.SessionRepository
}

func NewSessionService(sessionRepo repositories.SessionRepository) SessionService {
	return &sessionService{
		sessionRepo:sessionRepo,
	}
}

func (service *sessionService) CreateSession(session *domains.Session) error {
	err := service.sessionRepo.Store(session)
	if err != nil{
		return err
	}
	return nil
}

func (service *sessionService) UpdateSession(session *domains.Session) error {
	err := service.sessionRepo.Update(session)
	if err != nil{
		return err
	}
	return nil
}

func (service *sessionService) DeleteSession(session *domains.Session) error {
	err := service.sessionRepo.Delete(session.Id)
	if err != nil {
		return err
	}
	return nil
}

func (service *sessionService) GetSession(user *domains.User, token string) (*domains.Session, error){
	session, err := service.sessionRepo.Find(user, token)
	if err != nil {
		return nil, err
	}
	return session, nil
}