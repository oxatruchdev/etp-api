-- Seed data for country table
INSERT INTO country (name, abbreviation, additional_fields, flag_code) VALUES
('Brazil', 'BR', '{"continent": "South America", "population": 212559417}', 'br'),
('Mexico', 'MX', '{"continent": "North America", "population": 128932753}', 'mx'),
('Argentina', 'AR', '{"continent": "South America", "population": 45195774}', 'ar'),
('Colombia', 'CO', '{"continent": "South America", "population": 50882891}', 'co'),
('Chile', 'CL', '{"continent": "South America", "population": 19116201}', 'cl'),
('Peru', 'PE', '{"continent": "South America", "population": 32971854}', 'pe'),
('Uruguay', 'UY', '{"continent": "South America", "population": 3473730}', 'uy'),
('Ecuador', 'EC', '{"continent": "South America", "population": 17643054}', 'ec');

-- Seed data for school table
INSERT INTO school (name, abbreviation, metadata, country_id) VALUES
('Universidade de São Paulo', 'USP', '{"founded": 1934, "type": "Public"}', 1),
('Universidade Estadual de Campinas', 'UNICAMP', '{"founded": 1966, "type": "Public"}', 1),
('Universidade Federal do Rio de Janeiro', 'UFRJ', '{"founded": 1920, "type": "Public"}', 1),
('Universidad Nacional Autónoma de México', 'UNAM', '{"founded": 1551, "type": "Public"}', 2),
('Instituto Tecnológico y de Estudios Superiores de Monterrey', 'ITESM', '{"founded": 1943, "type": "Private"}', 2),
('Universidad de Guadalajara', 'UDG', '{"founded": 1792, "type": "Public"}', 2),
('Universidad de Buenos Aires', 'UBA', '{"founded": 1821, "type": "Public"}', 3),
('Universidad Nacional de La Plata', 'UNLP', '{"founded": 1905, "type": "Public"}', 3),
('Universidad Nacional de Colombia', 'UNAL', '{"founded": 1867, "type": "Public"}', 4),
('Universidad de los Andes', 'Uniandes', '{"founded": 1948, "type": "Private"}', 4),
('Pontificia Universidad Católica de Chile', 'UC', '{"founded": 1888, "type": "Private"}', 5),
('Universidad de Chile', 'UCh', '{"founded": 1842, "type": "Public"}', 5),
('Universidad Nacional Mayor de San Marcos', 'UNMSM', '{"founded": 1551, "type": "Public"}', 6),
('Pontificia Universidad Católica del Perú', 'PUCP', '{"founded": 1917, "type": "Private"}', 6),
('Universidad de la República', 'UdelaR', '{"founded": 1849, "type": "Public"}', 7),
('Universidad ORT Uruguay', 'ORT', '{"founded": 1942, "type": "Private"}', 7),
('Universidad San Francisco de Quito', 'USFQ', '{"founded": 1988, "type": "Private"}', 8),
('Escuela Politécnica Nacional', 'EPN', '{"founded": 1869, "type": "Public"}', 8);

-- Seed data for school_rating table
INSERT INTO school_rating (rating, comment, is_approved, approval_count, school_id) VALUES
(5, 'Excelente institución con una amplia gama de programas académicos.', true, 25, 1),
(4, 'Gran ambiente de investigación y oportunidades de desarrollo.', true, 20, 2),
(5, 'Programas de alta calidad y profesores reconocidos internacionalmente.', true, 22, 3),
(4, 'Gran universidad con una rica historia y fuerte en investigación.', true, 18, 4),
(5, 'Innovación y excelencia en educación tecnológica.', true, 24, 5),
(4, 'Amplia oferta académica y buena infraestructura.', true, 19, 6),
(5, 'Educación de alta calidad y ambiente estudiantil vibrante.', true, 23, 7),
(4, 'Excelente formación académica y oportunidades de investigación.', true, 21, 8),
(5, 'Universidad líder con gran impacto en la educación colombiana.', true, 26, 9),
(4, 'Excelente calidad de enseñanza y enfoque internacional.', true, 20, 10),
(5, 'Prestigiosa institución con programas de clase mundial.', true, 27, 11),
(4, 'Fuerte tradición académica y excelencia en investigación.', true, 22, 12),
(5, 'Universidad histórica con gran impacto en la educación peruana.', true, 23, 13),
(4, 'Formación integral y compromiso con la excelencia académica.', true, 19, 14),
(5, 'Principal institución de educación superior en Uruguay.', true, 24, 15),
(4, 'Enfoque innovador en educación tecnológica y emprendimiento.', true, 18, 16),
(5, 'Universidad líder en Ecuador con enfoque liberal en artes.', true, 21, 17),
(4, 'Excelencia en ingeniería y ciencias aplicadas.', true, 20, 18);

-- Seed data for department table
INSERT INTO department (name, code, school_id) VALUES
('Engenharia de Computação', 'EC', 1),
('Física', 'FIS', 1),
('Medicina', 'MED', 1),
('Economia', 'ECO', 1),
('Biologia', 'BIO', 2),
('Química', 'QUI', 2),
('Matemática', 'MAT', 2),
('Engenharia Elétrica', 'EE', 3),
('Direito', 'DIR', 3),
('Psicologia', 'PSI', 3),
('Ingeniería Civil', 'IC', 4),
('Filosofía', 'FIL', 4),
('Arquitectura', 'ARQ', 4),
('Administración', 'ADM', 5),
('Mecatrónica', 'MEC', 5),
('Ciencias de la Comunicación', 'COM', 6),
('Sociología', 'SOC', 6),
('Medicina', 'MED', 7),
('Ingeniería Informática', 'INF', 7),
('Ciencias Políticas', 'CP', 8),
('Agronomía', 'AGR', 8),
('Ingeniería Ambiental', 'IA', 9),
('Artes Plásticas', 'ART', 9),
('Economía', 'ECO', 10),
('Ingeniería Industrial', 'II', 10),
('Medicina Veterinaria', 'VET', 11),
('Astronomía', 'AST', 11),
('Derecho', 'DER', 12),
('Geología', 'GEO', 12),
('Lingüística', 'LIN', 13),
('Farmacia y Bioquímica', 'FAR', 13),
('Antropología', 'ANT', 14),
('Ingeniería de Sistemas', 'IS', 14),
('Odontología', 'ODO', 15),
('Ciencias de la Educación', 'EDU', 15),
('Diseño Gráfico', 'DG', 16),
('Biotecnología', 'BIO', 16),
('Relaciones Internacionales', 'RI', 17),
('Gastronomía', 'GAS', 17),
('Ingeniería en Petróleos', 'IP', 18),
('Matemática Aplicada', 'MA', 18);

-- Seed data for course table
INSERT INTO course (name, code, credits, department_id, school_id) VALUES
('Algoritmos e Estruturas de Dados', 'EC101', 4, 1, 1),
('Mecânica Quântica', 'FIS301', 3, 2, 1),
('Anatomia Humana', 'MED201', 5, 3, 1),
('Microeconomia', 'ECO202', 4, 4, 1),
('Genética Molecular', 'BIO202', 4, 5, 2),
('Química Orgânica', 'QUI301', 3, 6, 2),
('Cálculo Avançado', 'MAT401', 4, 7, 2),
('Circuitos Elétricos', 'EE201', 4, 8, 3),
('Direito Constitucional', 'DIR301', 5, 9, 3),
('Psicologia Cognitiva', 'PSI401', 3, 10, 3),
('Resistencia de Materiales', 'IC301', 4, 11, 4),
('Ética y Filosofía Política', 'FIL201', 3, 12, 4),
('Diseño Arquitectónico', 'ARQ301', 5, 13, 4),
('Gestión Estratégica', 'ADM401', 4, 14, 5),
('Robótica', 'MEC301', 4, 15, 5),
('Teoría de la Comunicación', 'COM201', 3, 16, 6),
('Métodos de Investigación Social', 'SOC301', 4, 17, 6),
('Fisiología Humana', 'MED301', 5, 18, 7),
('Inteligencia Artificial', 'INF401', 4, 19, 7),
('Sistemas Políticos Comparados', 'CP301', 3, 20, 8),
('Producción Agrícola Sostenible', 'AGR401', 4, 21, 8),
('Gestión de Residuos', 'IA301', 4, 22, 9),
('Escultura Contemporánea', 'ART401', 3, 23, 9),
('Macroeconomía Avanzada', 'ECO401', 4, 24, 10),
('Logística y Cadena de Suministro', 'II301', 4, 25, 10),
('Cirugía Veterinaria', 'VET401', 5, 26, 11),
('Astrofísica Estelar', 'AST301', 4, 27, 11),
('Derecho Internacional', 'DER401', 4, 28, 12),
('Tectónica de Placas', 'GEO301', 3, 29, 12),
('Semántica y Pragmática', 'LIN301', 3, 30, 13),
('Farmacología Clínica', 'FAR401', 5, 31, 13),
('Antropología Urbana', 'ANT301', 3, 32, 14),
('Desarrollo de Software', 'IS401', 4, 33, 14),
('Ortodoncia Avanzada', 'ODO401', 5, 34, 15),
('Tecnología Educativa', 'EDU301', 3, 35, 15),
('Diseño de Interfaces', 'DG301', 4, 36, 16),
('Ingeniería Genética', 'BIO401', 4, 37, 16),
('Diplomacia y Negociación', 'RI301', 3, 38, 17),
('Cocina Internacional', 'GAS301', 4, 39, 17),
('Perforación y Producción', 'IP401', 5, 40, 18),
('Optimización y Control', 'MA301', 4, 41, 18);

-- Seed data for professor table
INSERT INTO professor (first_name, last_name, full_name, school_id, department_id) VALUES
('João', 'Silva', 'João Silva', 1, 1),          -- EC
('Ana', 'Rodrigues', 'Ana Rodrigues', 1, 2),    -- FIS
('Carlos', 'Oliveira', 'Carlos Oliveira', 1, 3), -- MED
('Mariana', 'Santos', 'Mariana Santos', 2, 5),   -- BIO
('Pedro', 'Ferreira', 'Pedro Ferreira', 2, 6),   -- QUI
('Luísa', 'Carvalho', 'Luísa Carvalho', 2, 7),  -- MAT
('Ricardo', 'Almeida', 'Ricardo Almeida', 3, 8), -- EE
('Beatriz', 'Lima', 'Beatriz Lima', 3, 9),       -- DIR
('María', 'González', 'María González', 4, 11),   -- IC
('Alejandro', 'López', 'Alejandro López', 4, 12), -- FIL
('Sofía', 'Ramírez', 'Sofía Ramírez', 4, 13),    -- ARQ
('Javier', 'Hernández', 'Javier Hernández', 5, 14), -- ADM
('Isabella', 'Martínez', 'Isabella Martínez', 5, 15), -- MEC
('Diego', 'Torres', 'Diego Torres', 6, 16),       -- COM
('Valentina', 'Flores', 'Valentina Flores', 6, 17), -- SOC
('Martín', 'Rodríguez', 'Martín Rodríguez', 7, 18), -- MED
('Lucía', 'Fernández', 'Lucía Fernández', 7, 19),   -- INF
('Mateo', 'Gómez', 'Mateo Gómez', 7, 18),          -- MED
('Camila', 'Pérez', 'Camila Pérez', 8, 20),        -- CP
('Sebastián', 'Díaz', 'Sebastián Díaz', 8, 21),    -- AGR
('Andrés', 'Vargas', 'Andrés Vargas', 9, 22),      -- IA
('Daniela', 'Morales', 'Daniela Morales', 9, 23),  -- ART
('Felipe', 'Castro', 'Felipe Castro', 10, 24),      -- ECO
('Catalina', 'Rojas', 'Catalina Rojas', 10, 25),   -- II
('Joaquín', 'Muñoz', 'Joaquín Muñoz', 11, 26),     -- VET
('Antonia', 'Sepúlveda', 'Antonia Sepúlveda', 11, 27), -- AST
('Gabriel', 'Contreras', 'Gabriel Contreras', 12, 28), -- DER
('Isidora', 'Pinto', 'Isidora Pinto', 12, 29),     -- GEO
('Miguel', 'Chávez', 'Miguel Chávez', 13, 30),      -- LIN
('Valeria', 'Mendoza', 'Valeria Mendoza', 13, 31),  -- FAR
('Rodrigo', 'Quispe', 'Rodrigo Quispe', 14, 32),    -- ANT
('Adriana', 'Huamán', 'Adriana Huamán', 14, 33),    -- IS
('Gonzalo', 'Silva', 'Gonzalo Silva', 15, 34),      -- ODO
('Cecilia', 'Fernández', 'Cecilia Fernández', 15, 35), -- EDU
('Martín', 'Sosa', 'Martín Sosa', 16, 36),         -- DG
('Laura', 'Giménez', 'Laura Giménez', 16, 37),     -- BIO
('Andrés', 'Benítez', 'Andrés Benítez', 17, 38),   -- RI
('Gabriela', 'Espinoza', 'Gabriela Espinoza', 17, 39), -- GAS
('Santiago', 'Córdova', 'Santiago Córdova', 18, 40), -- IP
('Camila', 'Andrade', 'Camila Andrade', 18, 41);    -- MA

-- Seed data for professor_rating table
INSERT INTO professor_rating (rating, difficulty, comment, would_take_again, mandatory_attendance, grade, textbook_required, is_approved, approvals_count, professor_id, course_id) VALUES
(5, 4, 'Excelente professor, aulas muito dinâmicas.', true, true, 'A', true, true, 30, 1, 1),
(4, 4, 'Boa professora, explica bem conceitos complexos.', true, false, 'B+', true, true, 25, 2, 2),
(5, 4, 'Professor muito dedicado, sempre disponível para tirar dúvidas.', true, true, 'A-', false, true, 28, 3, 3),
(4, 4, 'Aulas interessantes, mas às vezes um pouco confusas.', true, true, 'B', true, true, 22, 4, 4),
(5, 4, 'Excelente pedagogo, torna a matéria fascinante.', true, false, 'A', true, true, 32, 5, 5),
(4, 4, 'Professora competente, mas um pouco rígida nas avaliações.', true, true, 'B+', false, true, 26, 6, 6),
(5, 4, 'Ótimo professor, relaciona bem a teoria com a prática.', true, false, 'A-', true, true, 29, 7, 7),
(4, 4, 'Boa professora, mas às vezes falta clareza nas explicações.', true, true, 'B', true, true, 24, 8, 8),
(5, 4, 'Excelente maestra, clases muy dinámicas e interesantes.', true, true, 'A', true, true, 31, 9, 9),
(4, 4, 'Buen profesor, exigente pero justo en las evaluaciones.', true, false, 'B+', true, true, 27, 10, 10),
(5, 4, 'Profesora muy dedicada, siempre dispuesta a ayudar.', true, true, 'A-', false, true, 30, 11, 11),
(4, 4, 'Clases interesantes, aunque a veces un poco teóricas.', true, true, 'B', true, true, 23, 12, 12),
(5, 4, 'Excelente profesor, hace que temas complejos sean comprensibles.', true, false, 'A', true, true, 33, 13, 13),
(4, 3, 'Buen docente, pero podría mejorar en la organización de las clases.', true, true, 'B+', false, true, 25, 14, 14),
(5, 3, 'Profesora excepcional, muy apasionada por su materia.', true, false, 'A-', true, true, 28, 15, 15),
(4, 3, 'Profesor competente, aunque a veces le falta paciencia.', true, true, 'B', true, true, 22, 16, 16),
(5, 3, 'Excelente docente, clases muy prácticas y útiles.', true, true, 'A', false, true, 31, 17, 17),
(4, 3, 'Buen profesor, pero las evaluaciones son muy difíciles.', true, false, 'B+', true, true, 26, 18, 18),
(5, 3, 'Profesora extraordinaria, muy clara en sus explicaciones.', true, true, 'A-', true, true, 29, 19, 19),
(4, 5, 'Docente competente, aunque podría ser más dinámico en clase.', true, true, 'B', false, true, 24, 20, 20),
(5, 5, 'Excelente profesor, muy actualizado en su campo.', true, false, 'A', true, true, 32, 21, 21),
(4, 5, 'Buena profesora, pero a veces va muy rápido con el material.', true, true, 'B+', true, true, 27, 22, 22),
(5, 5, 'Profesor excepcional, muy dedicado a sus estudiantes.', true, false, 'A-', false, true, 30, 23, 23),
(4, 5, 'Docente competente, aunque podría mejorar en feedback.', true, true, 'B', true, true, 25, 24, 24),
(5, 5, 'Excelente profesor, clases muy interactivas y enriquecedoras.', true, true, 'A', true, true, 33, 25, 25),
(4, 5, 'Buena profesora, pero a veces falta claridad en las instrucciones.', true, false, 'B+', false, true, 26, 26, 26),
(5, 1, 'Profesor muy competente, excelente en resolver dudas.', true, true, 'A-', true, true, 29, 27, 27),
(4, 1, 'Docente dedicada, aunque podría variar más las actividades.', true, true, 'B', true, true, 23, 28, 28),
(5, 1, 'Excelente profesor, muy inspirador y motivador.', true, false, 'A', false, true, 31, 29, 29),
(4, 1, 'Buena profesora, pero las lecturas asignadas son muy extensas.', true, true, 'B+', true, true, 25, 30, 30),
(5, 1, 'Profesor excepcional, muy accesible fuera de clase.', true, false, 'A-', true, true, 28, 31, 31),
(4, 1, 'Docente competente, aunque a veces le falta organización.', true, true, 'B', false, true, 24, 32, 32),
(5, 1, 'Excelente profesora, muy clara y paciente en sus explicaciones.', true, true, 'A', true, true, 32, 33, 33),
(4, 2, 'Buen profesor, pero podría mejorar en la retroalimentación de trabajos.', true, false, 'B+', true, true, 26, 34, 34),
(5, 2, 'Profesora extraordinaria, hace que la materia sea muy interesante.', true, true, 'A-', false, true, 29, 35, 35),
(4, 2, 'Docente competente, aunque a veces las clases son un poco monótonas.', true, true, 'B', true, true, 23, 36, 36),
(5, 2, 'Excelente profesor, muy actualizado y con experiencia práctica.', true, false, 'A', true, true, 31, 37, 37),
(4, 2, 'Buena profesora, pero las evaluaciones son muy exigentes.', true, true, 'B+', false, true, 25, 38, 38),
(5, 2, 'Profesor excepcional, muy dedicado y apasionado por su materia.', true, false, 'A-', true, true, 30, 39, 39),
(4, 2, 'Docente competente, aunque podría mejorar en la gestión del tiempo en clase.', true, true, 'B', true, true, 24, 40, 40);

-- Seed data for professor_course table
INSERT INTO professor_course (professor_id, course_id)
SELECT professor_id, course_id FROM professor_rating;

-- Seed data for role table
INSERT INTO role (name, display_name) VALUES
('student', 'Estudiante'),
('professor', 'Profesor'),
('moderator', 'Moderador'),
('school_admin', 'Administrador de Escuela'),
('admin', 'Administrador');

-- Seed data for tag table in Spanish

INSERT INTO tag (name) VALUES
('Excelente Enseñanza'),
('Calificador Estricto'),
('Servicial'),
('Explicaciones Claras'),
('Clases Dinámicas'),
('Participación Obligatoria'),
('Créditos Extra Disponibles'),
('Proyectos Grupales'),
('Evaluación Justa'),
('Asistencia Obligatoria'),
('Plazos Flexibles'),
('Da Buenas Retroalimentaciones'),
('Inspirador'),
('Exámenes Difíciles'),
('Tareas Exigentes'),
('Accesible Fuera de Clase'),
('Exámenes Pesados'),
('Conocimiento Práctico'),
('Basado en Teoría'),
('Usa Ejemplos del Mundo Real');

-- Seed data for profesor_rating_tag table

INSERT INTO professor_rating_tag (professor_rating_id, tag_id) VALUES
(1, 1),  -- Excelente Enseñanza
(1, 4),  -- Explicaciones Claras
(1, 9),  -- Evaluación Justa
(2, 2),  -- Calificador Estricto
(2, 6),  -- Participación Obligatoria
(2, 13), -- Inspirador
(3, 8),  -- Proyectos Grupales
(3, 3),  -- Servicial
(3, 12), -- Da Buenas Retroalimentaciones
(4, 7),  -- Créditos Extra Disponibles
(4, 14), -- Exámenes Difíciles
(4, 17), -- Exámenes Pesados
(5, 5),  -- Clases Dinámicas
(5, 11), -- Plazos Flexibles
(5, 16), -- Accesible Fuera de Clase
(6, 10), -- Asistencia Obligatoria
(6, 18), -- Conocimiento Práctico
(7, 19), -- Basado en Teoría
(7, 20), -- Usa Ejemplos del Mundo Real
(7, 9),  -- Evaluación Justa
(8, 3),  -- Servicial
(8, 5),  -- Clases Dinámicas
(8, 8),  -- Proyectos Grupales
(9, 2),  -- Calificador Estricto
(9, 14), -- Exámenes Difíciles
(9, 17), -- Exámenes Pesados
(10, 6), -- Participación Obligatoria
(10, 11),-- Plazos Flexibles
(10, 13),-- Inspirador
(11, 12),-- Da Buenas Retroalimentaciones
(11, 18),-- Conocimiento Práctico
(11, 20),-- Usa Ejemplos del Mundo Real
(12, 1), -- Excelente Enseñanza
(12, 4), -- Explicaciones Claras
(12, 9); -- Evaluación Justa

