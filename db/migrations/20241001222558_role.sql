CREATE TABLE IF NOT EXISTS role (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL default now(),
  updated_at TIMESTAMPTZ NOT NULL default now()
);

ALTER TABLE "user" ADD COLUMN role_id integer REFERENCES role(id);
