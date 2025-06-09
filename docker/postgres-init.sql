CREATE TABLE IF NOT EXISTS advices (
  user_id TEXT NOT NULL,
  advice TEXT NOT NULL,
  created_at TIMESTAMP,
  PRIMARY KEY (user_id)
);