-- +migrate Up
CREATE TYPE tasks_status AS ENUM ('pending', 'in_progress', 'completed');
CREATE TYPE tasks_priority AS ENUM ('low', 'medium', 'high');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status tasks_status DEFAULT 'pending',
    priority tasks_priority DEFAULT 'medium',
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS tasks;