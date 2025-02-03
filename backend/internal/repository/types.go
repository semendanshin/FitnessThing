package repository

import (
	"fitness-trainer/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func timeToPgtype(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: !t.IsZero()}
}

func floatToPgtype(f float32) pgtype.Float4 {
	return pgtype.Float4{Float32: f, Valid: f != 0}
}

func uuidToPgtype(id domain.ID) pgtype.UUID {
	return pgtype.UUID{Bytes: uuid.UUID(id), Valid: id != domain.ID{}}
}

func durationFromPgtype(d pgtype.Interval) time.Duration {
	return time.Duration(d.Microseconds) * time.Microsecond
}

func intervalToPgtype(d time.Duration) pgtype.Interval {
	return pgtype.Interval{Microseconds: int64(d / time.Microsecond), Valid: d != 0}
}

func timeFromPgtype(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}

	return t.Time
}

func uuidsToPgtype(ids []domain.ID) []pgtype.UUID {
	result := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		result = append(result, uuidToPgtype(id))
	}

	return result
}
