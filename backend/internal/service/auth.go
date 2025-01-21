package service

import (
	"context"
	"fmt"
	"time"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	"fitness-trainer/internal/utils"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) Login(ctx context.Context, email string, password string) (domain.Tokens, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Login")
	defer span.Finish()

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Errorf("error getting user by email: %v", err)
		return domain.Tokens{}, domain.ErrInvalidArgument
	}

	err = utils.VerifyPassword(user.Password, password)
	if err != nil {
		logger.Errorf("error verifying password: %v", err)
		return domain.Tokens{}, domain.ErrInvalidArgument
	}

	tokens, err := s.jwtProvider.GeneratePair(ctx, user.ID, domain.NewID(), time.Now())
	if err != nil {
		logger.Errorf("error generating pair: %v", err)
		return domain.Tokens{}, domain.ErrInternal
	}

	_, err = s.sessionRepository.CreateSession(
		ctx,
		domain.NewSession(
			user.ID,
			time.Time{},
			tokens.RefreshToken,
		),
	)
	if err != nil {
		logger.Errorf("error creating session: %v", err)
		return domain.Tokens{}, fmt.Errorf("error creating session: %w", err)
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, tokens domain.Tokens) (domain.Tokens, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Refresh")
	defer span.Finish()

	refreshToken, err := s.sessionRepository.GetSessionByToken(ctx, tokens.RefreshToken)
	if err != nil {
		logger.Errorf("error getting refresh token: %v", err)
		return domain.Tokens{}, domain.ErrInvalidArgument
	}

	if !refreshToken.ExpiredAt.IsZero() && refreshToken.ExpiredAt.Before(time.Now()) {
		logger.Errorf("refresh token is expired")
		return domain.Tokens{}, domain.ErrInvalidArgument
	}

	if err := s.jwtProvider.VerifyPair(ctx, refreshToken.UserID, tokens, time.Now()); err != nil {
		logger.Errorf("error verifying token: %v", err)
		return domain.Tokens{}, domain.ErrInvalidArgument
	}

	err = s.sessionRepository.SetSessionExpired(ctx, refreshToken.ID, time.Now())
	if err != nil {
		logger.Errorf("error setting expired session: %v", err)
		return domain.Tokens{}, domain.ErrInternal
	}

	newTokens, err := s.jwtProvider.GeneratePair(ctx, refreshToken.UserID, domain.NewID(), time.Now())
	if err != nil {
		logger.Errorf("error generating pair: %v", err)
		return domain.Tokens{}, domain.ErrInternal
	}

	_, err = s.sessionRepository.CreateSession(
		ctx,
		domain.NewSession(
			refreshToken.UserID,
			time.Time{},
			newTokens.RefreshToken,
		),
	)
	if err != nil {
		logger.Errorf("error creating session: %v", err)
		return domain.Tokens{}, domain.ErrInternal
	}

	return newTokens, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Logout")
	defer span.Finish()

	session, err := s.sessionRepository.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		logger.Errorf("error getting session: %v", err)
		return domain.ErrInvalidArgument
	}

	err = s.sessionRepository.SetSessionExpired(ctx, session.ID, time.Now())
	if err != nil {
		logger.Errorf("error setting expired session: %v", err)
		return domain.ErrInternal
	}

	return nil
}

func (s *Service) ParseToken(ctx context.Context, token string) (domain.ID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ParseToken")
	defer span.Finish()

	userID, err := s.jwtProvider.ParseToken(ctx, token)
	if err != nil {
		logger.Errorf("error parsing token: %v", err)
		return domain.ID{}, fmt.Errorf("error parsing token: %w", err)
	}

	return userID, nil
}
