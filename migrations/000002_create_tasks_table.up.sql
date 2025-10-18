-- Create task_status enum type
CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');
-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'todo',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
