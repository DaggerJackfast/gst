package main

import (
	"database/sql"
	"time"
)

type UserProfileTokenRepository interface {
	FindUserTokenByStatus(user *User, tokenType string) (*UserProfileToken, error)
	Store(token *UserProfileToken) error
	Update(token *UserProfileToken) error
}

type userProfileTokenRepository struct {
	db sql.DB
}

func NewUserProfileTokenRepository(db sql.DB) UserProfileTokenRepository {
	return &userProfileTokenRepository{
		db: db,
	}
}

func (repo *userProfileTokenRepository) FindUserTokenByStatus(user *User, tokenType string) (*UserProfileToken, error) {
	row := repo.db.QueryRow(`select * from user_profile_tokens where user_id=$1 
                                    and token_type= $2 and is_active=true order by created_at`,
		user.Id, tokenType)
	var token UserProfileToken
	var userId uint64
	err := row.Scan(&token.Id, &userId, &token.ProfileToken, &token.TokenType, &token.IsActive, &token.ExpiredIn, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, err
	}
	token.User = user
	return &token, nil
}

func (repo *userProfileTokenRepository) Store(token *UserProfileToken) error {
	nowTime := time.Now()
	token.UpdatedAt = nowTime
	token.CreatedAt = nowTime
	err := repo.db.QueryRow(`
		insert into user_profile_tokens(user_id, profile_token, token_type, is_active, expired_in, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7) returning id`,
		token.User.Id, token.ProfileToken, token.TokenType,
		token.IsActive, token.ExpiredIn, token.CreatedAt,
		token.UpdatedAt,
	).Scan(&token.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userProfileTokenRepository) Update(token *UserProfileToken) error {
	token.UpdatedAt = time.Now()
	_, err := repo.db.Exec(`update user_profile_tokens set user_id=$2, profile_token=$3, token_type=$4,
									is_active=$5, expired_in=$6, updated_at=$7 where id=$1`,
		token.Id, token.User.Id, token.ProfileToken, token.TokenType, token.IsActive,
		token.ExpiredIn, token.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
