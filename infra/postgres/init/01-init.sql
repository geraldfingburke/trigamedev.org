-- Initial shared database setup
CREATE SCHEMA IF NOT EXISTS trigamedev;
CREATE SCHEMA IF NOT EXISTS shared;

CREATE TABLE IF NOT EXISTS shared.service_registry (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) UNIQUE NOT NULL,
    runtime VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO shared.service_registry (service_name, runtime)
VALUES
    ('api-go', 'Go 1.24'),
    ('api-fastapi', 'Python FastAPI')
ON CONFLICT (service_name) DO NOTHING;
