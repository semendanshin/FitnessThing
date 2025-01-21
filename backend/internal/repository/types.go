package repository

import (
	"fitness-trainer/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func timeToPgtype(t time.Time) pgtype.Timestamptz {
	var valid bool
	if !t.IsZero() {
		valid = true
	}

	return pgtype.Timestamptz{Time: t, Valid: valid}
}

func floatToPgtype(f float32) pgtype.Float4 {
	var valid bool
	if f != 0 {
		valid = true
	}

	return pgtype.Float4{Float32: f, Valid: valid}
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
