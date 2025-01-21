package mappers

import (
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"
)

func TokensToProto(tokens domain.Tokens) *desc.TokensPair {
	return &desc.TokensPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func ProtoToTokens(pb *desc.TokensPair) domain.Tokens {
	return domain.Tokens{
		AccessToken:  pb.AccessToken,
		RefreshToken: pb.RefreshToken,
	}
}
