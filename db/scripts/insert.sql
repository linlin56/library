-- =============================================================================
-- LIBRARY — PostgreSQL script
-- Native PostgreSQL inheritance:
--   - books    INHERITS items
--   - journals INHERITS items
--   - students INHERITS users
--   - teachers INHERITS users
-- =============================================================================

-- -----------------------------------------------------------------------------
-- CLEANUP (reverse dependency order)
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS teachers    CASCADE;
DROP TABLE IF EXISTS students    CASCADE;
DROP TABLE IF EXISTS users       CASCADE;
DROP TABLE IF EXISTS librarians  CASCADE;
DROP TABLE IF EXISTS copies      CASCADE;
DROP TABLE IF EXISTS books       CASCADE;
DROP TABLE IF EXISTS journals    CASCADE;
DROP TABLE IF EXISTS items       CASCADE;
DROP TABLE IF EXISTS authors     CASCADE;
DROP TABLE IF EXISTS locations   CASCADE;
DROP TABLE IF EXISTS categories  CASCADE;
DROP TABLE IF EXISTS addresses   CASCADE;
DROP TABLE IF EXISTS departments CASCADE;


-- =============================================================================
-- INDEPENDENT TABLES
-- =============================================================================

CREATE TABLE authors (
    id         SERIAL PRIMARY KEY,
    last_name  VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL
);

CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE locations (
    id           SERIAL PRIMARY KEY,
    bay_number   INT NOT NULL,
    shelf_number INT NOT NULL,
    level        INT NOT NULL,
    category_id  INT REFERENCES categories(id)
);

CREATE TABLE addresses (
    id            SERIAL PRIMARY KEY,
    city          VARCHAR(100) NOT NULL,
    postal_code   VARCHAR(10)  NOT NULL,
    street_name   VARCHAR(200) NOT NULL,
    street_number VARCHAR(10)  NOT NULL
);

CREATE TABLE departments (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE librarians (
    id         SERIAL PRIMARY KEY,
    last_name  VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    login      VARCHAR(50)  NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL
);


-- =============================================================================
-- ITEM (parent table — book and journal inherit from it)
-- =============================================================================

CREATE TABLE items (
    id      SERIAL PRIMARY KEY,
    title   VARCHAR(255) NOT NULL,
    deposit DOUBLE PRECISION NOT NULL DEFAULT 0.0
);

-- Book inherits from item
CREATE TABLE books (
    isbn      VARCHAR(20) NOT NULL,
    author_id INT         NOT NULL REFERENCES authors(id)
) INHERITS (items);

-- Journal inherits from item
CREATE TABLE journals (
    issue_number INT NOT NULL
) INHERITS (items);


-- =============================================================================
-- COPY (linked to an item, stored at a location)
-- =============================================================================

CREATE TABLE copies (
    id               SERIAL PRIMARY KEY,
    is_borrowed      BOOLEAN     NOT NULL DEFAULT FALSE,
    loan_start_date  DATE,
    loan_due_date    DATE,
    barcode          VARCHAR(50) NOT NULL UNIQUE,
    loan_period_days INT         NOT NULL DEFAULT 14,
    item_id          INT         NOT NULL,  -- logical reference (inheritance)
    location_id      INT         REFERENCES locations(id),
    borrower_id      INT                    -- logical FK to users
);


-- =============================================================================
-- USER (parent table — student and teacher inherit from it)
-- =============================================================================

CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    last_name       VARCHAR(100)     NOT NULL,
    first_name      VARCHAR(100)     NOT NULL,
    phone_number    VARCHAR(20),
    deposit_balance DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    total_deposit   DOUBLE PRECISION NOT NULL DEFAULT 20.0,
    login           VARCHAR(50)      NOT NULL UNIQUE,
    password        VARCHAR(255)     NOT NULL,
    birth_date      DATE,
    email           VARCHAR(150)     NOT NULL UNIQUE,
    address_id      INT              REFERENCES addresses(id)
);

-- Student inherits from user
CREATE TABLE students (
    study_year VARCHAR(20) NOT NULL
) INHERITS (users);

-- Teacher inherits from user
CREATE TABLE teachers (
    department_id INT NOT NULL REFERENCES departments(id)
) INHERITS (users);


-- =============================================================================
-- INSERT — TEST DATA
-- =============================================================================

-- -----------------------------------------------------------------------------
-- Categories
-- -----------------------------------------------------------------------------
INSERT INTO categories (name) VALUES
    ('Informatique'),
    ('Sciences'),
    ('Littérature'),
    ('Histoire'),
    ('Mathématiques');

-- -----------------------------------------------------------------------------
-- Locations
-- -----------------------------------------------------------------------------
INSERT INTO locations (bay_number, shelf_number, level, category_id) VALUES
    (1, 1, 1, 1),  -- Informatique
    (1, 2, 1, 1),  -- Informatique
    (2, 1, 1, 2),  -- Sciences
    (2, 2, 1, 3),  -- Littérature
    (3, 1, 2, 4),  -- Histoire
    (3, 2, 2, 5);  -- Mathématiques

-- -----------------------------------------------------------------------------
-- Authors
-- -----------------------------------------------------------------------------
INSERT INTO authors (last_name, first_name) VALUES
    ('Kernighan', 'Brian'),
    ('Ritchie',   'Dennis'),
    ('Donovan',   'Alan'),
    ('Kernighan', 'Brian'),
    ('Martin',    'Robert C.'),
    ('Zola',      'Émile'),
    ('Hugo',      'Victor'),
    ('Cormen',    'Thomas H.');

-- -----------------------------------------------------------------------------
-- Items — Books (inserted directly into books, automatically inherited in items)
-- -----------------------------------------------------------------------------
INSERT INTO books (title, deposit, isbn, author_id) VALUES
    ('The C Programming Language',       5.00, '978-0131103627', 1),
    ('The Go Programming Language',      5.00, '978-0134190440', 3),
    ('Clean Code',                       5.00, '978-0132350884', 5),
    ('Introduction to Algorithms',      10.00, '978-0262033848', 8),
    ('Germinal',                          3.00, '978-2070408597', 6),
    ('Les Misérables',                    3.00, '978-2070409228', 7);

-- -----------------------------------------------------------------------------
-- Items — Journals
-- -----------------------------------------------------------------------------
INSERT INTO journals (title, deposit, issue_number) VALUES
    ('Linux Magazine',      2.00,  258),
    ('Pour la Science',     2.50,  551),
    ('L''Histoire',         2.00,  512),
    ('Programmez!',         2.00,   95);

-- -----------------------------------------------------------------------------
-- Addresses
-- -----------------------------------------------------------------------------
INSERT INTO addresses (city, postal_code, street_name, street_number) VALUES
    ('Paris',      '75001', 'Rue de Rivoli',       '10'),
    ('Lyon',       '69001', 'Rue de la République','25'),
    ('Nantes',     '44000', 'Rue du Général de Gaulle', '7'),
    ('Bordeaux',   '33000', 'Allées de Tourny',    '3'),
    ('Marseille',  '13001', 'La Canebière',         '42'),
    ('Toulouse',   '31000', 'Rue du Taur',          '15');

-- -----------------------------------------------------------------------------
-- Departments
-- -----------------------------------------------------------------------------
INSERT INTO departments (name) VALUES
    ('Informatique'),
    ('Mathématiques'),
    ('Physique'),
    ('Langues et Lettres'),
    ('Sciences Humaines');

-- -----------------------------------------------------------------------------
-- Librarians
-- -----------------------------------------------------------------------------
INSERT INTO librarians (last_name, first_name, login, password) VALUES
    ('Dupont',  'Marie',   'mdupont',  'password'),
    ('Lambert', 'Jacques', 'jlambert', 'password');

-- -----------------------------------------------------------------------------
-- Users — Students (inserted directly into students)
-- -----------------------------------------------------------------------------
INSERT INTO students (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, study_year) VALUES
    ('Martin',   'Alice',   '0601020304', 0.00, 'amartin',   'password', '2002-03-15', 'alice.martin@etud.fr',   1, 'L3'),
    ('Bernard',  'Lucas',   '0605060708', 0.00, 'lbernard',  'password', '2001-07-22', 'lucas.bernard@etud.fr',  2, 'M1'),
    ('Petit',    'Emma',    '0609101112', 5.00, 'epetit',    'password', '2003-01-10', 'emma.petit@etud.fr',     3, 'L2'),
    ('Robert',   'Noah',    '0613141516', 0.00, 'nrobert',   'password', '2000-11-05', 'noah.robert@etud.fr',    4, 'M2'),
    ('Leroy',    'Chloé',   '0617181920', 2.50, 'cleroy',    'password', '2002-09-30', 'chloe.leroy@etud.fr',    5, 'L1');

-- -----------------------------------------------------------------------------
-- Users — Teachers (inserted directly into teachers)
-- -----------------------------------------------------------------------------
INSERT INTO teachers (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, department_id) VALUES
    ('Moreau',   'Julien',  '0621222324', 0.00, 'jmoreau',   'password', '1978-04-12', 'julien.moreau@univ.fr',  6, 1),
    ('Simon',    'Claire',  '0625262728', 0.00, 'csimon',    'password', '1985-08-19', 'claire.simon@univ.fr',   1, 2),
    ('Girard',   'Pierre',  '0629303132', 0.00, 'pgirard',   'password', '1972-12-03', 'pierre.girard@univ.fr',  2, 3);

-- -----------------------------------------------------------------------------
-- Copies
-- (item_id logically references the ids inserted into books/journals)
-- The ids are retrieved with a sub-query
-- -----------------------------------------------------------------------------
INSERT INTO copies (is_borrowed, barcode, loan_period_days, item_id, location_id) VALUES
    -- The C Programming Language (2 copies)
    (FALSE, 'EX-0001', 14, (SELECT id FROM books WHERE isbn = '978-0131103627'), 1),
    (FALSE, 'EX-0002', 14, (SELECT id FROM books WHERE isbn = '978-0131103627'), 1),
    -- The Go Programming Language
    (FALSE, 'EX-0003', 14, (SELECT id FROM books WHERE isbn = '978-0134190440'), 2),
    -- Clean Code
    (FALSE, 'EX-0004', 21, (SELECT id FROM books WHERE isbn = '978-0132350884'), 2),
    -- Introduction to Algorithms
    (FALSE, 'EX-0005', 21, (SELECT id FROM books WHERE isbn = '978-0262033848'), 6),
    -- Germinal
    (FALSE, 'EX-0006', 14, (SELECT id FROM books WHERE isbn = '978-2070408597'), 4),
    -- Les Misérables
    (FALSE, 'EX-0007', 14, (SELECT id FROM books WHERE isbn = '978-2070409228'), 4),
    -- Linux Magazine (2 copies)
    (FALSE, 'EX-0008',  7, (SELECT id FROM journals WHERE issue_number = 258), 1),
    (FALSE, 'EX-0009',  7, (SELECT id FROM journals WHERE issue_number = 258), 1),
    -- Pour la Science
    (FALSE, 'EX-0010',  7, (SELECT id FROM journals WHERE issue_number = 551), 3),
    -- L'Histoire
    (FALSE, 'EX-0011',  7, (SELECT id FROM journals WHERE issue_number = 512), 5),
    -- Programmez!
    (FALSE, 'EX-0012',  7, (SELECT id FROM journals WHERE issue_number = 95),  1);

-- -----------------------------------------------------------------------------
-- Simulated loan (Alice borrows "The Go Programming Language")
-- -----------------------------------------------------------------------------
UPDATE copies
SET
    is_borrowed     = TRUE,
    loan_start_date = CURRENT_DATE,
    loan_due_date   = CURRENT_DATE + INTERVAL '14 days',
    borrower_id     = (SELECT id FROM students WHERE login = 'amartin')
WHERE barcode = 'EX-0003';


-- =============================================================================
-- QUICK CHECKS
-- =============================================================================
SELECT 'authors'     AS table_name, COUNT(*) FROM authors
UNION ALL
SELECT 'books',       COUNT(*) FROM books
UNION ALL
SELECT 'journals',    COUNT(*) FROM journals
UNION ALL
SELECT 'copies',      COUNT(*) FROM copies
UNION ALL
SELECT 'students',    COUNT(*) FROM students
UNION ALL
SELECT 'teachers',    COUNT(*) FROM teachers
UNION ALL
SELECT 'librarians',  COUNT(*) FROM librarians
UNION ALL
SELECT 'departments', COUNT(*) FROM departments
UNION ALL
SELECT 'categories',  COUNT(*) FROM categories
UNION ALL
SELECT 'locations',   COUNT(*) FROM locations
UNION ALL
SELECT 'addresses',   COUNT(*) FROM addresses;
