package main

import (
	"database/sql"
	"time"
)

type SessionRepository interface {
	//Find(userId uint64) (*Session, error)
	Update(session *Session) error
	Store(session *Session) error
	Delete(id uint64) error
}

type sessionRepository struct {
	db sql.DB
}

func NewSessionRepository(db sql.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (repo *sessionRepository) Store(session *Session) error {
	nowTime := time.Now()
	session.UpdatedAt = nowTime
	session.CreatedAt = nowTime
	err := repo.db.QueryRow(`
		insert into sessions(user_id, refresh_token, user_agent, fingerprint, ip, expired_in, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8) returning id`,
		session.User.Id, session.RefreshToken, session.UserAgent, session.FingerPrint,
		session.Ip, session.ExpiredIn, session.CreatedAt, session.UpdatedAt,
	).Scan(&session.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *sessionRepository) Update(session *Session) error {
	session.UpdatedAt = time.Now()
	_, err := repo.db.Exec(`
		update sessions set user_id=$2, refresh_token=$3, user_agent=$4, fingerprint=$5, 
		                    ip=$6, expired_in=$7, updated_at=$8
						where id=$1`,
		session.Id, session.User.Id, session.RefreshToken, session.UserAgent,
		session.FingerPrint, session.Ip, session.ExpiredIn, session.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (repo *sessionRepository) Delete(id uint64) error {
	_, err := repo.db.Exec("delete from sessions where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
