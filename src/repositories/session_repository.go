package repositories

import (
	"database/sql"
	"github.com/DaggerJackfast/gst/src/domains"
	"time"
)

type SessionRepository interface {
	Find(user *domains.User, refreshToken string) (*domains.Session, error)
	Update(session *domains.Session) error
	Store(session *domains.Session) error
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

func (repo *sessionRepository) Store(session *domains.Session) error {
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

func (repo *sessionRepository) Update(session *domains.Session) error {
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

func (repo *sessionRepository) Find(user *domains.User, refreshToken string) (*domains.Session, error) {
	row := repo.db.QueryRow("select * from sessions where refresh_token=$1 and user_id=$2", refreshToken, user.Id)
	var session domains.Session
	var userId int
	err := row.Scan(&session.Id, &userId, &session.RefreshToken, &session.UserAgent,
		&session.FingerPrint, &session.Ip, &session.ExpiredIn,
		&session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
