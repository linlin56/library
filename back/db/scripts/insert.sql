-- This script is run when the project starts (if the database is empty)
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
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    deposit     DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    location_id INT REFERENCES locations(id)
);

-- Book inherits from item
CREATE TABLE books (
    isbn      VARCHAR(20) NOT NULL,
    author_id INT         NOT NULL REFERENCES authors(id)
) INHERITS (items);

-- Journal inherits from item
CREATE TABLE journals (
    issue_number     INT  NOT NULL,
    publication_date DATE NOT NULL
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
    loan_period_days INT         NOT NULL DEFAULT 15,
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
    deposit_balance DOUBLE PRECISION NOT NULL DEFAULT 20.0,
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
    ('Littérature'),
    ('Sciences'),
    ('Consommation');

-- -----------------------------------------------------------------------------
-- Locations
-- -----------------------------------------------------------------------------
INSERT INTO locations (bay_number, shelf_number, level, category_id) VALUES
    (1, 1, 1, 1),  -- Littérature
    (2, 1, 1, 2),  -- Sciences
    (2, 2, 1, 3);  -- Consommation

-- -----------------------------------------------------------------------------
-- Authors
-- -----------------------------------------------------------------------------
INSERT INTO authors (last_name, first_name) VALUES
    ('Tolkien',  'J.R.R.'),
    ('Anonyme',  'Collectif');

-- -----------------------------------------------------------------------------
-- Books
-- -----------------------------------------------------------------------------
INSERT INTO books (title, deposit, isbn, author_id, location_id) VALUES
    ('Le Seigneur des anneaux — La Communauté de l''anneau', 5.00, '978-2267011296', 1, 1),
    ('Le Seigneur des anneaux — Les Deux Tours',             5.00, '978-2267011555', 1, 1),
    ('Le Seigneur des anneaux — Le Retour du roi',           5.00, '978-2267011777', 1, 1),
    ('Les Mille et Une Nuits',                               3.00, '978-2070379583', 2, 1);

-- -----------------------------------------------------------------------------
-- Journals
-- -----------------------------------------------------------------------------
INSERT INTO journals (title, deposit, issue_number, publication_date, location_id) VALUES
    ('60 millions de consommateurs', 2.00, 581, '2024-01-01', 3),
    ('60 millions de consommateurs', 2.00, 582, '2024-02-01', 3),
    ('Science & Vie',                2.00, 1265, '2024-03-01', 2);

-- -----------------------------------------------------------------------------
-- Addresses
-- -----------------------------------------------------------------------------
INSERT INTO addresses (city, postal_code, street_name, street_number) VALUES
    ('Paris',   '75001', 'Rue de Rivoli',              '10'),  -- Jean Martin
    ('Lyon',    '69001', 'Rue de la République',       '25'),  -- Marie Dupont
    ('Nantes',  '44000', 'Rue du Général de Gaulle',   '7');   -- Pierre Roche

-- -----------------------------------------------------------------------------
-- Departments
-- -----------------------------------------------------------------------------
INSERT INTO departments (name) VALUES
    ('Mathématiques'),
    ('Informatique'),
    ('Physique');

-- -----------------------------------------------------------------------------
-- Librarian
-- -----------------------------------------------------------------------------
INSERT INTO librarians (last_name, first_name, login, password) VALUES
    ('Petit', 'Valérie', 'vpetit', 'password');

-- -----------------------------------------------------------------------------
-- Individual — inserted directly into users (ONLY)
-- -----------------------------------------------------------------------------
INSERT INTO users (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id) VALUES
    ('Martin', 'Jean', '0601020304', 20.00, 'jmartin', 'password', '1985-06-15', 'jean.martin@email.fr', 1);

-- -----------------------------------------------------------------------------
-- Student
-- -----------------------------------------------------------------------------
INSERT INTO students (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, study_year) VALUES
    ('Dupont', 'Marie', '0607080910', 20.00, 'mdupont', 'password', '2002-03-20', 'marie.dupont@etud.fr', 2, 'L3');

-- -----------------------------------------------------------------------------
-- Teacher
-- -----------------------------------------------------------------------------
INSERT INTO teachers (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, department_id) VALUES
    ('Roche', 'Pierre', '0611121314', 20.00, 'proche', 'password', '1975-11-08', 'pierre.roche@univ.fr', 3, 1);

-- -----------------------------------------------------------------------------
-- Copies (1 per item)
-- -----------------------------------------------------------------------------
INSERT INTO copies (is_borrowed, barcode, loan_period_days, item_id) VALUES
    -- La Communauté de l'anneau
    (FALSE, 'EX-0001', 15, (SELECT id FROM books WHERE isbn = '978-2267011296')),
    -- Les Deux Tours
    (FALSE, 'EX-0002', 15, (SELECT id FROM books WHERE isbn = '978-2267011555')),
    -- Le Retour du roi
    (FALSE, 'EX-0003', 15, (SELECT id FROM books WHERE isbn = '978-2267011777')),
    -- Les Mille et Une Nuits — copy 1
    (FALSE, 'EX-0004', 15, (SELECT id FROM books WHERE isbn = '978-2070379583')),
    -- Les Mille et Une Nuits — copy 2
    (FALSE, 'EX-0008', 15, (SELECT id FROM books WHERE isbn = '978-2070379583')),
    -- 60 millions de consommateurs no. 581
    (FALSE, 'EX-0005',  7, (SELECT id FROM journals WHERE issue_number = 581)),
    -- 60 millions de consommateurs no. 582
    (FALSE, 'EX-0006',  7, (SELECT id FROM journals WHERE issue_number = 582)),
    -- Science & Vie no. 1265
    (FALSE, 'EX-0007',  7, (SELECT id FROM journals WHERE issue_number = 1265));

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
SELECT 'users',       COUNT(*) FROM ONLY users
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
