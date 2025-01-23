-- +goose Up
-- +goose StatementBegin
-- Create users table
CREATE TABLE users
(
    id            UUID PRIMARY KEY      ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW() ,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    first_name    VARCHAR(255),
    last_name     VARCHAR(255),
    date_of_birth TIMESTAMPTZ,
    height        FLOAT,
    weight        FLOAT
);

-- Create muscle_groups table
CREATE TABLE muscle_groups
(
    id   UUID PRIMARY KEY ,
    name VARCHAR(255) NOT NULL
);

-- Create exercises table
CREATE TABLE exercises
(
    id          UUID PRIMARY KEY      ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    video_url   VARCHAR(255)
);

-- Create helper table
CREATE TABLE exercise_muscle_groups
(
    exercise_id     UUID NOT NULL,
    muscle_group_id UUID NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id),
    FOREIGN KEY (muscle_group_id) REFERENCES muscle_groups (id),
    PRIMARY KEY (exercise_id, muscle_group_id)
);

-- Create routines table
CREATE TABLE routines
(
    id          UUID PRIMARY KEY      ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    user_id     UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create exercise_instances table
CREATE TABLE exercise_instances
(
    id          UUID PRIMARY KEY     ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    exercise_id UUID        NOT NULL,
    routine_id  UUID        NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id),
    FOREIGN KEY (routine_id) REFERENCES routines (id)
);

-- Create sets table
CREATE TABLE sets
(
    id                   UUID PRIMARY KEY     ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    exercise_instance_id UUID        NOT NULL,
    set_type             VARCHAR(50) NOT NULL,
    reps                 INT,
    weight               FLOAT,
    time INTERVAL,
    FOREIGN KEY (exercise_instance_id) REFERENCES exercise_instances (id)
);

-- Create workouts table
CREATE TABLE workouts
(
    id          UUID PRIMARY KEY     ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id     UUID        NOT NULL,
    routine_id  UUID        NOT NULL,
    notes       TEXT,
    rating      INT,
    finished_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (routine_id) REFERENCES routines (id)
);

-- Create exercise_logs table
CREATE TABLE exercise_logs
(
    id           UUID PRIMARY KEY     ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    workout_id   UUID        NOT NULL,
    exercise_id  UUID        NOT NULL,
    notes        TEXT,
    power_rating INT,
    FOREIGN KEY (workout_id) REFERENCES workouts (id),
    FOREIGN KEY (exercise_id) REFERENCES exercises (id)
);

-- Create set_logs table
CREATE TABLE set_logs
(
    id              UUID PRIMARY KEY     ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    exercise_log_id UUID        NOT NULL,
    reps            INT,
    weight          FLOAT,
    time INTERVAL,
    FOREIGN KEY (exercise_log_id) REFERENCES exercise_logs (id)
);

-- Table for refresh tokens
CREATE TABLE sessions
(
    id         UUID PRIMARY KEY ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id    UUID        NOT NULL,
    token      TEXT        NOT NULL,
    expired_at TIMESTAMPTZ      DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create indexes
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_exercises_name ON exercises (name);
CREATE INDEX idx_routines_user_id ON routines (user_id);
CREATE INDEX idx_exercise_instances_exercise_id ON exercise_instances (exercise_id);
CREATE INDEX idx_sets_exercise_instance_id ON sets (exercise_instance_id);
CREATE INDEX idx_workouts_user_id ON workouts (user_id);
CREATE INDEX idx_workouts_routine_id ON workouts (routine_id);
CREATE INDEX idx_exercise_logs_workout_id ON exercise_logs (workout_id);
CREATE INDEX idx_exercise_logs_exercise_id ON exercise_logs (exercise_id);
CREATE INDEX idx_set_logs_exercise_log_id ON set_logs (exercise_log_id);
CREATE INDEX idx_workouts_finished_at ON workouts (finished_at);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop set_logs index
DROP INDEX IF EXISTS idx_set_logs_exercise_log_id;

-- Drop exercise_logs index
DROP INDEX IF EXISTS idx_exercise_logs_workout_id;
DROP INDEX IF EXISTS idx_exercise_logs_exercise_id;

-- Drop workouts index
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workouts_routine_id;
DROP INDEX IF EXISTS idx_workouts_finished_at;

-- Drop sets index
DROP INDEX IF EXISTS idx_sets_exercise_instance_id;

-- Drop exercise_instances index
DROP INDEX IF EXISTS idx_exercise_instances_exercise_id;
DROP INDEX IF EXISTS idx_exercise_instances_routine_id;

-- Drop routines index
DROP INDEX IF EXISTS idx_routines_user_id;

-- Drop exercise_muscle_groups indexes
DROP INDEX IF EXISTS idx_exercise_muscle_groups_exercise_id;
DROP INDEX IF EXISTS idx_exercise_muscle_groups_muscle_group_id;

-- Drop exercises index
DROP INDEX IF EXISTS idx_exercises_name;

-- Drop users index
DROP INDEX IF EXISTS idx_users_email;

-- Drop set_logs table
DROP TABLE IF EXISTS set_logs;

-- Drop exercise_logs table
DROP TABLE IF EXISTS exercise_logs;

-- Drop workouts table
DROP TABLE IF EXISTS workouts;

-- Drop sets table
DROP TABLE IF EXISTS sets;

-- Drop exercise_instances table
DROP TABLE IF EXISTS exercise_instances;

-- Drop routines table
DROP TABLE IF EXISTS routines;

-- Drop exercise_muscle_groups table
DROP TABLE IF EXISTS exercise_muscle_groups;

-- Drop exercises table
DROP TABLE IF EXISTS exercises;

-- Drop muscle_groups table
DROP TABLE IF EXISTS muscle_groups;

-- Drop users table
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
