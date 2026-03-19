-- =============================================================================
-- BIBLIOTHÈQUE — Script SQL PostgreSQL
-- Héritage natif PostgreSQL :
--   - livres    INHERITS ouvrages
--   - revues    INHERITS ouvrages
--   - etudiants INHERITS utilisateurs
--   - enseignants INHERITS utilisateurs
-- =============================================================================

-- -----------------------------------------------------------------------------
-- NETTOYAGE (ordre inverse des dépendances)
-- -----------------------------------------------------------------------------
DROP TABLE IF EXISTS enseignants    CASCADE;
DROP TABLE IF EXISTS etudiants      CASCADE;
DROP TABLE IF EXISTS utilisateurs   CASCADE;
DROP TABLE IF EXISTS bibliothecaires CASCADE;
DROP TABLE IF EXISTS exemplaires    CASCADE;
DROP TABLE IF EXISTS livres         CASCADE;
DROP TABLE IF EXISTS revues         CASCADE;
DROP TABLE IF EXISTS ouvrages       CASCADE;
DROP TABLE IF EXISTS auteurs        CASCADE;
DROP TABLE IF EXISTS emplacements   CASCADE;
DROP TABLE IF EXISTS categories     CASCADE;
DROP TABLE IF EXISTS adresses       CASCADE;
DROP TABLE IF EXISTS departements_ecole CASCADE;


-- =============================================================================
-- TABLES INDÉPENDANTES
-- =============================================================================

CREATE TABLE auteurs (
    id     SERIAL PRIMARY KEY,
    nom    VARCHAR(100) NOT NULL,
    prenom VARCHAR(100) NOT NULL
);

CREATE TABLE categories (
    id  SERIAL PRIMARY KEY,
    nom VARCHAR(100) NOT NULL
);

CREATE TABLE emplacements (
    id               SERIAL PRIMARY KEY,
    numero_travee    INT NOT NULL,
    numero_etagere   INT NOT NULL,
    niveau           INT NOT NULL,
    categorie_id     INT REFERENCES categories(id)
);

CREATE TABLE adresses (
    id          SERIAL PRIMARY KEY,
    ville       VARCHAR(100) NOT NULL,
    code_postal VARCHAR(10)  NOT NULL,
    nom_rue     VARCHAR(200) NOT NULL,
    numero_rue  VARCHAR(10)  NOT NULL
);

CREATE TABLE departements_ecole (
    id  SERIAL PRIMARY KEY,
    nom VARCHAR(100) NOT NULL
);

CREATE TABLE bibliothecaires (
    id          SERIAL PRIMARY KEY,
    nom         VARCHAR(100) NOT NULL,
    prenom      VARCHAR(100) NOT NULL,
    login       VARCHAR(50)  NOT NULL UNIQUE,
    mot_de_passe VARCHAR(255) NOT NULL
);


-- =============================================================================
-- OUVRAGE (table parente — livre et revue héritent de celle-ci)
-- =============================================================================

CREATE TABLE ouvrages (
    id           SERIAL PRIMARY KEY,
    titre        VARCHAR(255) NOT NULL,
    caution      DOUBLE PRECISION NOT NULL DEFAULT 0.0
);

-- Livre hérite d'ouvrage
CREATE TABLE livres (
    isbn         VARCHAR(20)  NOT NULL,
    auteur_id    INT          NOT NULL REFERENCES auteurs(id)
) INHERITS (ouvrages);

-- Revue hérite d'ouvrage
CREATE TABLE revues (
    numero INT NOT NULL
) INHERITS (ouvrages);


-- =============================================================================
-- EXEMPLAIRE (lié à un ouvrage, se trouve à un emplacement)
-- =============================================================================

CREATE TABLE exemplaires (
    id                  SERIAL PRIMARY KEY,
    est_emprunte        BOOLEAN          NOT NULL DEFAULT FALSE,
    date_debut_emprunt  DATE,
    date_fin_emprunt    DATE,
    code_barre          VARCHAR(50)      NOT NULL UNIQUE,
    delai_emprunt_jours INT              NOT NULL DEFAULT 14,
    ouvrage_id          INT              NOT NULL,  -- référence logique (héritage)
    emplacement_id      INT              REFERENCES emplacements(id),
    emprunteur_id       INT              -- FK vers utilisateurs (logique)
);


-- =============================================================================
-- UTILISATEUR (table parente — etudiant et enseignant héritent de celle-ci)
-- =============================================================================

CREATE TABLE utilisateurs (
    id               SERIAL PRIMARY KEY,
    nom              VARCHAR(100)     NOT NULL,
    prenom           VARCHAR(100)     NOT NULL,
    numero_telephone VARCHAR(20),
    solde_caution    DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    login            VARCHAR(50)      NOT NULL UNIQUE,
    mot_de_passe     VARCHAR(255)     NOT NULL,
    date_de_naissance DATE,
    email            VARCHAR(150)     NOT NULL UNIQUE,
    adresse_id       INT              REFERENCES adresses(id)
);

-- Etudiant hérite d'utilisateur
CREATE TABLE etudiants (
    annee_etude VARCHAR(20) NOT NULL
) INHERITS (utilisateurs);

-- Enseignant hérite d'utilisateur
CREATE TABLE enseignants (
    departement_id INT NOT NULL REFERENCES departements_ecole(id)
) INHERITS (utilisateurs);


-- =============================================================================
-- INSERT — DONNÉES DE TEST
-- =============================================================================

-- -----------------------------------------------------------------------------
-- Catégories
-- -----------------------------------------------------------------------------
INSERT INTO categories (nom) VALUES
    ('Informatique'),
    ('Sciences'),
    ('Littérature'),
    ('Histoire'),
    ('Mathématiques');

-- -----------------------------------------------------------------------------
-- Emplacements
-- -----------------------------------------------------------------------------
INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id) VALUES
    (1, 1, 1, 1),  -- Informatique
    (1, 2, 1, 1),  -- Informatique
    (2, 1, 1, 2),  -- Sciences
    (2, 2, 1, 3),  -- Littérature
    (3, 1, 2, 4),  -- Histoire
    (3, 2, 2, 5);  -- Mathématiques

-- -----------------------------------------------------------------------------
-- Auteurs
-- -----------------------------------------------------------------------------
INSERT INTO auteurs (nom, prenom) VALUES
    ('Kernighan', 'Brian'),
    ('Ritchie',   'Dennis'),
    ('Donovan',   'Alan'),
    ('Kernighan', 'Brian'),
    ('Martin',    'Robert C.'),
    ('Zola',      'Émile'),
    ('Hugo',      'Victor'),
    ('Cormen',    'Thomas H.');

-- -----------------------------------------------------------------------------
-- Ouvrages — Livres (INSERT direct dans livres, héritage auto dans ouvrages)
-- -----------------------------------------------------------------------------
INSERT INTO livres (titre, caution, isbn, auteur_id) VALUES
    ('The C Programming Language',       5.00, '978-0131103627', 1),
    ('The Go Programming Language',      5.00, '978-0134190440', 3),
    ('Clean Code',                       5.00, '978-0132350884', 5),
    ('Introduction to Algorithms',      10.00, '978-0262033848', 8),
    ('Germinal',                          3.00, '978-2070408597', 6),
    ('Les Misérables',                    3.00, '978-2070409228', 7);

-- -----------------------------------------------------------------------------
-- Ouvrages — Revues
-- -----------------------------------------------------------------------------
INSERT INTO revues (titre, caution, numero) VALUES
    ('Linux Magazine',      2.00,  258),
    ('Pour la Science',     2.50,  551),
    ('L''Histoire',         2.00,  512),
    ('Programmez!',         2.00,   95);

-- -----------------------------------------------------------------------------
-- Adresses
-- -----------------------------------------------------------------------------
INSERT INTO adresses (ville, code_postal, nom_rue, numero_rue) VALUES
    ('Paris',      '75001', 'Rue de Rivoli',       '10'),
    ('Lyon',       '69001', 'Rue de la République','25'),
    ('Nantes',     '44000', 'Rue du Général de Gaulle', '7'),
    ('Bordeaux',   '33000', 'Allées de Tourny',    '3'),
    ('Marseille',  '13001', 'La Canebière',         '42'),
    ('Toulouse',   '31000', 'Rue du Taur',          '15');

-- -----------------------------------------------------------------------------
-- Départements
-- -----------------------------------------------------------------------------
INSERT INTO departements_ecole (nom) VALUES
    ('Informatique'),
    ('Mathématiques'),
    ('Physique'),
    ('Langues et Lettres'),
    ('Sciences Humaines');

-- -----------------------------------------------------------------------------
-- Bibliothécaires
-- -----------------------------------------------------------------------------
INSERT INTO bibliothecaires (nom, prenom, login, mot_de_passe) VALUES
    ('Dupont',  'Marie',   'mdupont',  'mdp'),
    ('Lambert', 'Jacques', 'jlambert', 'mdp');

-- -----------------------------------------------------------------------------
-- Utilisateurs — Etudiants (INSERT direct dans etudiants)
-- -----------------------------------------------------------------------------
INSERT INTO etudiants (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id, annee_etude) VALUES
    ('Martin',   'Alice',   '0601020304', 0.00, 'amartin',   'mdp', '2002-03-15', 'alice.martin@etud.fr',   1, 'L3'),
    ('Bernard',  'Lucas',   '0605060708', 0.00, 'lbernard',  'mdp', '2001-07-22', 'lucas.bernard@etud.fr',  2, 'M1'),
    ('Petit',    'Emma',    '0609101112', 5.00, 'epetit',    'mdp', '2003-01-10', 'emma.petit@etud.fr',     3, 'L2'),
    ('Robert',   'Noah',    '0613141516', 0.00, 'nrobert',   'mdp', '2000-11-05', 'noah.robert@etud.fr',    4, 'M2'),
    ('Leroy',    'Chloé',   '0617181920', 2.50, 'cleroy',    'mdp', '2002-09-30', 'chloe.leroy@etud.fr',    5, 'L1');

-- -----------------------------------------------------------------------------
-- Utilisateurs — Enseignants (INSERT direct dans enseignants)
-- -----------------------------------------------------------------------------
INSERT INTO enseignants (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id, departement_id) VALUES
    ('Moreau',   'Julien',  '0621222324', 0.00, 'jmoreau',   'mdp', '1978-04-12', 'julien.moreau@univ.fr',  6, 1),
    ('Simon',    'Claire',  '0625262728', 0.00, 'csimon',    'mdp', '1985-08-19', 'claire.simon@univ.fr',   1, 2),
    ('Girard',   'Pierre',  '0629303132', 0.00, 'pgirard',   'mdp', '1972-12-03', 'pierre.girard@univ.fr',  2, 3);

-- -----------------------------------------------------------------------------
-- Exemplaires
-- (ouvrage_id référence logiquement les ids insérés dans livres/revues)
-- On récupère les ids via sous-requête sur le titre
-- -----------------------------------------------------------------------------
INSERT INTO exemplaires (est_emprunte, code_barre, delai_emprunt_jours, ouvrage_id, emplacement_id) VALUES
    -- The C Programming Language (2 exemplaires)
    (FALSE, 'EX-0001', 14, (SELECT id FROM livres WHERE isbn = '978-0131103627'), 1),
    (FALSE, 'EX-0002', 14, (SELECT id FROM livres WHERE isbn = '978-0131103627'), 1),
    -- The Go Programming Language
    (FALSE, 'EX-0003', 14, (SELECT id FROM livres WHERE isbn = '978-0134190440'), 2),
    -- Clean Code
    (FALSE, 'EX-0004', 21, (SELECT id FROM livres WHERE isbn = '978-0132350884'), 2),
    -- Introduction to Algorithms
    (FALSE, 'EX-0005', 21, (SELECT id FROM livres WHERE isbn = '978-0262033848'), 6),
    -- Germinal
    (FALSE, 'EX-0006', 14, (SELECT id FROM livres WHERE isbn = '978-2070408597'), 4),
    -- Les Misérables
    (FALSE, 'EX-0007', 14, (SELECT id FROM livres WHERE isbn = '978-2070409228'), 4),
    -- Linux Magazine (2 exemplaires)
    (FALSE, 'EX-0008',  7, (SELECT id FROM revues WHERE numero = 258), 1),
    (FALSE, 'EX-0009',  7, (SELECT id FROM revues WHERE numero = 258), 1),
    -- Pour la Science
    (FALSE, 'EX-0010',  7, (SELECT id FROM revues WHERE numero = 551), 3),
    -- L'Histoire
    (FALSE, 'EX-0011',  7, (SELECT id FROM revues WHERE numero = 512), 5),
    -- Programmez!
    (FALSE, 'EX-0012',  7, (SELECT id FROM revues WHERE numero = 95),  1);

-- -----------------------------------------------------------------------------
-- Simulation d'un emprunt (Alice emprunte "The Go Programming Language")
-- -----------------------------------------------------------------------------
UPDATE exemplaires
SET
    est_emprunte       = TRUE,
    date_debut_emprunt = CURRENT_DATE,
    date_fin_emprunt   = CURRENT_DATE + INTERVAL '14 days',
    emprunteur_id      = (SELECT id FROM etudiants WHERE login = 'amartin')
WHERE code_barre = 'EX-0003';


-- =============================================================================
-- VÉRIFICATIONS RAPIDES
-- =============================================================================
SELECT 'auteurs'          AS table_name, COUNT(*) FROM auteurs
UNION ALL
SELECT 'livres',           COUNT(*) FROM livres
UNION ALL
SELECT 'revues',           COUNT(*) FROM revues
UNION ALL
SELECT 'exemplaires',      COUNT(*) FROM exemplaires
UNION ALL
SELECT 'etudiants',        COUNT(*) FROM etudiants
UNION ALL
SELECT 'enseignants',      COUNT(*) FROM enseignants
UNION ALL
SELECT 'bibliothecaires',  COUNT(*) FROM bibliothecaires
UNION ALL
SELECT 'departements',     COUNT(*) FROM departements_ecole
UNION ALL
SELECT 'categories',       COUNT(*) FROM categories
UNION ALL
SELECT 'emplacements',     COUNT(*) FROM emplacements
UNION ALL
SELECT 'adresses',         COUNT(*) FROM adresses;