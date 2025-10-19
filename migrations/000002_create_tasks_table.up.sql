-- Create task_status enum type
-- Idempotent: Uses IF NOT EXISTS to allow safe re-runs
CREATE TYPE IF NOT EXISTS task_status AS ENUM ('todo', 'in_progress', 'done');

-- Create tasks table
-- Idempotent: Uses IF NOT EXISTS to allow safe re-runs
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'todo',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

