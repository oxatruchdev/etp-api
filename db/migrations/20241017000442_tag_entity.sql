CREATE TABLE IF NOT EXISTS tag (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,

  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ default now()
);

CREATE TABLE IF NOT EXISTS professor_rating_tag(
  id SERIAL PRIMARY KEY,
  tag_id INTEGER NOT NULL REFERENCES tag(id),
  professor_rating_id INTEGER NOT NULL REFERENCES professor_rating(id),
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ default now()
);

CREATE TABLE IF NOT EXISTS school_rating_tag(
  id SERIAL PRIMARY KEY,
  tag_id INTEGER NOT NULL REFERENCES tag(id),
  school_rating_id INTEGER NOT NULL REFERENCES school_rating(id),
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ default now()
);
