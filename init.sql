-- Create the task_manager database safely using DO block
DO $$
BEGIN
  CREATE DATABASE task_manager;
EXCEPTION WHEN duplicate_database THEN
  NULL;
END
$$;
