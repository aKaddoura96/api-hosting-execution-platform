-- Make code_path nullable to allow API creation before code upload
ALTER TABLE apis ALTER COLUMN code_path DROP NOT NULL;
