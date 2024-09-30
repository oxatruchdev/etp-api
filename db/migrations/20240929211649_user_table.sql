CREATE TABLE IF NOT EXISTS "user" (
  id SERIAL PRIMARY KEY,
  email character varying(255) NOT NULL UNIQUE,
  password character varying(255) NOT NULL,
  name character varying(255) NULL,

  created_at TIMESTAMPTZ NOT NULL default now(),
  updated_at TIMESTAMPTZ NOT NULL default now(),

  -- relations
  school_id integer REFERENCES school(id)
);

ALTER TABLE professor_rating ADD COLUMN user_id integer REFERENCES "user"(id);
ALTER TABLE school_rating ADD COLUMN user_id integer REFERENCES "user"(id);
