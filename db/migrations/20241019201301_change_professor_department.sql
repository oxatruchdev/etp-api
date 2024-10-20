DROP TABLE professor_department;

ALTER TABLE professor ADD COLUMN department_id INTEGER REFERENCES department(id);
