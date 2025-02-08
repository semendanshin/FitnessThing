package repository

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type userEntity struct {
	ID pgtype.UUID

	Email    string
	Password string

	FirstName string
	LastName  string

	PictureProfileURL pgtype.Text `db:"picture_profile_url"`

	DateOfBirth pgtype.Timestamptz

	Weight pgtype.Float4
	Height pgtype.Float4

	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (u userEntity) toDomain() domain.User {
	return domain.User{
		Model: domain.Model{
			ID:        domain.ID(u.ID.Bytes),
			CreatedAt: timeFromPgtype(u.CreatedAt),
			UpdatedAt: timeFromPgtype(u.UpdatedAt),
		},
		Email:         u.Email,
		Password:      u.Password,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		DateOfBirth:   timeFromPgtype(u.DateOfBirth),
		Weight:        u.Weight.Float32,
		Height:        u.Height.Float32,
		ProfilePicURL: u.PictureProfileURL.String,
	}
}

func userFromDomain(user domain.User) userEntity {
	return userEntity{
		ID:                uuidToPgtype(user.ID),
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Password:          user.Password,
		DateOfBirth:       timeToPgtype(user.DateOfBirth),
		Weight:            floatToPgtype(user.Weight),
		Height:            floatToPgtype(user.Height),
		CreatedAt:         timeToPgtype(user.CreatedAt),
		UpdatedAt:         timeToPgtype(user.UpdatedAt),
		PictureProfileURL: pgtype.Text{String: user.ProfilePicURL, Valid: user.ProfilePicURL != ""},
	}
}

func (r *PGXRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetUserByEmail")
	defer span.Finish()

	const query = `
		select id, email, password, first_name, last_name, date_of_birth, height, weight, created_at, updated_at, picture_profile_url
		from users u 
		where u.email=$1;
	`

	var user userEntity

	engine := r.contextManager.GetEngineFromContext(ctx)

	err := pgxscan.Get(ctx, engine, &user, query, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		logger.Errorf("error getting user by email: %v", err)
	}

	return user.toDomain(), nil
}

func (r *PGXRepository) GetUserByID(ctx context.Context, id domain.ID) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetUserByID")
	defer span.Finish()

	const query = `
		select id, email, password, first_name, last_name, date_of_birth, height, weight, created_at, updated_at, picture_profile_url
		from users u 
		where u.id=$1;
	`

	var user userEntity

	engine := r.contextManager.GetEngineFromContext(ctx)

	err := pgxscan.Get(ctx, engine, &user, query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		logger.Errorf("error getting user by id: %v", err)
	}

	return user.toDomain(), nil
}

func (r *PGXRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateUser")
	defer span.Finish()

	const query = `
		insert into users (id, email, password, first_name, last_name, date_of_birth, height, weight, created_at, updated_at, picture_profile_url)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	userEntity := userFromDomain(user)

	engine := r.contextManager.GetEngineFromContext(ctx)

	err := pgxscan.Get(ctx, engine, &userEntity, query,
		uuidToPgtype(user.ID),
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		timeToPgtype(user.DateOfBirth),
		floatToPgtype(user.Height),
		floatToPgtype(user.Weight),
		timeToPgtype(user.CreatedAt),
		timeToPgtype(user.UpdatedAt),
		pgtype.Text{String: user.ProfilePicURL, Valid: user.ProfilePicURL != ""},
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.User{}, domain.ErrAlreadyExists
			}
		}
		logger.Errorf("error creating user: %v", err)
		return domain.User{}, err
	}

	return userEntity.toDomain(), nil
}

func (r *PGXRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UpdateUser")
	defer span.Finish()

	const query = `
		update users
		set email=$2, first_name=$3, last_name=$4, date_of_birth=$5, height=$6, weight=$7, updated_at=$8, picture_profile_url=$9
		where id=$1
		returning id, email, password, first_name, last_name, date_of_birth, height, weight, created_at, updated_at;
	`

	userEntity := userFromDomain(user)

	engine := r.contextManager.GetEngineFromContext(ctx)

	err := pgxscan.Get(
		ctx, engine, &userEntity, query,
		uuidToPgtype(user.ID),
		user.Email,
		user.FirstName,
		user.LastName,
		timeToPgtype(user.DateOfBirth),
		floatToPgtype(user.Height),
		floatToPgtype(user.Weight),
		timeToPgtype(user.UpdatedAt),
		pgtype.Text{String: user.ProfilePicURL, Valid: user.ProfilePicURL != ""},
	)
	if err != nil {
		logger.Errorf("error updating user: %v", err)
		return domain.User{}, err
	}

	return userEntity.toDomain(), nil
}
