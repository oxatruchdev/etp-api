-- Country table
CREATE TABLE IF NOT EXISTS country (
  id SERIAL PRIMARY KEY,
  name character varying(255) NOT NULL,
  abbreviation character varying(255) NOT NULL,
  additional_fields jsonb NULL,

  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now()
);

-- School table
CREATE TABLE IF NOT EXISTS school (
  id SERIAL PRIMARY KEY,
  name character varying(255) NOT NULL,
  abbreviation character varying(255) NOT NULL,
  metadata jsonb NULL,

  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  -- relations
  country_id integer REFERENCES country(id)
);

-- School ratings table
CREATE TABLE IF NOT EXISTS school_rating (
  id SERIAL PRIMARY KEY,
  rating integer NOT NULL,
  comment text NOT NULL,
  is_approved boolean NOT NULL default false,
  approval_count integer NOT NULL default 0,
  updated_count integer NOT NULL default 0,

  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  -- relations
  school_id integer REFERENCES school(id)
);

-- Department table
CREATE TABLE IF NOT EXISTS department (
  id SERIAL PRIMARY KEY,
  name character varying(255) NOT NULL,
  code character varying(255) NOT NULL,

  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  school_id integer REFERENCES school(id)
);

-- Course table
CREATE TABLE IF NOT EXISTS course (
  id SERIAL PRIMARY KEY,
  name character varying(255) NOT NULL,
  code character varying(255) NOT NULL,
  credits integer NOT NULL,
  
  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  department_id integer REFERENCES department(id),
  school_id integer REFERENCES school(id)
);

-- Professor table
CREATE TABLE IF NOT EXISTS professor (
  id SERIAL PRIMARY KEY,
  first_name character varying(255) NOT NULL,
  last_name character varying(255) NOT NULL,

  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  school_id integer REFERENCES school(id)
);

-- Professor Ratings table
CREATE TABLE IF NOT EXISTS professor_rating (
  id SERIAL PRIMARY KEY,

  -- rating fields
  rating integer NOT NULL,
  comment text NOT NULL,
  would_take_again boolean NOT NULL,
  mandatory_attendance boolean NOT NULL,
  grade character varying(255) NOT NULL,
  textbook_required boolean NOT NULL,

  is_approved boolean NOT NULL default false,
  approvals_count integer NOT NULL default 0,
  updated_count integer NOT NULL default 0,

  -- date
  created_at timestamp with time zone NOT NULL default now(),
  updated_at timestamp with time zone NOT NULL default now(),

  -- relations
  professor_id integer REFERENCES professor(id),
  course_id integer REFERENCES course(id)
);

-- many-to-many relations

-- professor_course table
CREATE TABLE IF NOT EXISTS professor_course (
  id SERIAL PRIMARY KEY,
  professor_id integer REFERENCES professor(id),
  course_id integer REFERENCES course(id)
);

-- professor_department table
CREATE TABLE IF NOT EXISTS professor_department (
  id SERIAL PRIMARY KEY,
  professor_id integer REFERENCES professor(id),
  department_id integer REFERENCES department(id)
);
