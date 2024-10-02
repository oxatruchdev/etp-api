-- Seed data for country table
INSERT INTO country (name, abbreviation, additional_fields) VALUES
('Brazil', 'BR', '{"continent": "South America", "population": 212559417}'),
('Mexico', 'MX', '{"continent": "North America", "population": 128932753}'),
('Argentina', 'AR', '{"continent": "South America", "population": 45195774}');

-- Seed data for school table
INSERT INTO school (name, abbreviation, metadata, country_id) VALUES
('Universidade de São Paulo', 'USP', '{"founded": 1934, "type": "Public"}', 1),
('Universidad Nacional Autónoma de México', 'UNAM', '{"founded": 1551, "type": "Public"}', 2),
('Universidad de Buenos Aires', 'UBA', '{"founded": 1821, "type": "Public"}', 3);

-- Seed data for school_rating table
INSERT INTO school_rating (rating, comment, is_approved, approval_count, school_id) VALUES
(5, 'Excelente institución con una amplia gama de programas académicos.', true, 15, 1),
(4, 'Gran universidad con una rica historia y fuerte en investigación.', true, 12, 2),
(5, 'Educación de alta calidad y ambiente estudiantil vibrante.', true, 18, 3);

-- Seed data for department table
INSERT INTO department (name, code, school_id) VALUES
('Engenharia de Computação', 'EC', 1),
('Física', 'FIS', 2),
('Medicina', 'MED', 3);

-- Seed data for course table
INSERT INTO course (name, code, credits, department_id, school_id) VALUES
('Algoritmos e Estruturas de Dados', 'EC101', 4, 1, 1),
('Mecánica Cuántica', 'FIS301', 3, 2, 2),
('Anatomía Humana', 'MED201', 5, 3, 3);

-- Seed data for professor table
INSERT INTO professor (first_name, last_name, school_id) VALUES
('João', 'Silva', 1),
('María', 'González', 2),
('Carlos', 'Fernández', 3);

-- Seed data for professor_rating table
INSERT INTO professor_rating (rating, comment, would_take_again, mandatory_attendance, grade, textbook_required, is_approved, approvals_count, professor_id, course_id) VALUES
(5, 'Professor muito dedicado e aulas interessantes.', true, true, 'A', true, true, 20, 1, 1),
(4, 'Clases desafiantes pero muy instructivas.', true, false, 'B+', true, true, 15, 2, 2),
(5, 'Excelente pedagogo, explica conceptos complejos con claridad.', true, true, 'A-', true, true, 18, 3, 3);

-- Seed data for professor_course table
INSERT INTO professor_course (professor_id, course_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- Seed data for professor_department table
INSERT INTO professor_department (professor_id, department_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- Seed data for role table
INSERT INTO role (name) VALUES
('Student'),
('Professor'),
('Moderator'),
('School Administrator'),
('Admin');
