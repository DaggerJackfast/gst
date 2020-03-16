package main

type SessionService interface {
	CreateSession(session *Session) error
	UpdateSession(session *Session) error
	DeleteSession(session *Session) error
}

type sessionService struct {
	sessionRepo SessionRepository
}

func NewSessionService(sessionRepo SessionRepository) SessionService {
	return &sessionService{
		sessionRepo:sessionRepo,
	}
}

func (service *sessionService) CreateSession(session *Session) error {
	err := service.sessionRepo.Store(session)
	if err != nil{
		return err
	}
	return nil
}

func (service *sessionService) UpdateSession(session *Session) error {
	err := service.sessionRepo.Update(session)
	if err != nil{
		return err
	}
	return nil
}

func (service *sessionService) DeleteSession(session *Session) error {
	err := service.sessionRepo.Delete(session.Id)
	if err != nil {
		return err
	}
	return nil
}