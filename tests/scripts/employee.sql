CREATE TABLE IF NOT EXISTS role
(
    id         BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employee
(
    id         BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       TEXT                        NOT NULL,
    created_at TIMESTAMPTZ                 NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ                 NOT NULL DEFAULT NOW(),
    role_id    BIGINT REFERENCES role (id) NOT NULL
);