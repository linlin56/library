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
    caution      DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    emplacement_id INT REFERENCES emplacements(id)
);

-- Livre hérite d'ouvrage
CREATE TABLE livres (
    isbn         VARCHAR(20)  NOT NULL,
    auteur_id    INT          NOT NULL REFERENCES auteurs(id)
) INHERITS (ouvrages);

-- Revue hérite d'ouvrage
CREATE TABLE revues (
    numero        INT  NOT NULL,
    date_parution DATE NOT NULL
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
    delai_emprunt_jours INT              NOT NULL DEFAULT 15,
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
    solde_caution    DOUBLE PRECISION NOT NULL DEFAULT 20.0,
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
    ('Littérature'),
    ('Sciences'),
    ('Consommation');

-- -----------------------------------------------------------------------------
-- Emplacements
-- -----------------------------------------------------------------------------
INSERT INTO emplacements (numero_travee, numero_etagere, niveau, categorie_id) VALUES
    (1, 1, 1, 1),  -- Littérature
    (2, 1, 1, 2),  -- Sciences
    (2, 2, 1, 3);  -- Consommation

-- -----------------------------------------------------------------------------
-- Auteurs
-- -----------------------------------------------------------------------------
INSERT INTO auteurs (nom, prenom) VALUES
    ('Tolkien',  'J.R.R.'),
    ('Anonyme',  'Collectif');

-- -----------------------------------------------------------------------------
-- Livres
-- -----------------------------------------------------------------------------
INSERT INTO livres (titre, caution, isbn, auteur_id, emplacement_id) VALUES
    ('Le Seigneur des anneaux — La Communauté de l''anneau', 5.00, '978-2267011296', 1, 1),
    ('Le Seigneur des anneaux — Les Deux Tours',             5.00, '978-2267011555', 1, 1),
    ('Le Seigneur des anneaux — Le Retour du roi',           5.00, '978-2267011777', 1, 1),
    ('Les Mille et Une Nuits',                               3.00, '978-2070379583', 2, 1);

-- -----------------------------------------------------------------------------
-- Revues
-- -----------------------------------------------------------------------------
INSERT INTO revues (titre, caution, numero, date_parution, emplacement_id) VALUES
    ('60 millions de consommateurs', 2.00, 581, '2024-01-01', 3),
    ('60 millions de consommateurs', 2.00, 582, '2024-02-01', 3),
    ('Science & Vie',                2.00, 1265, '2024-03-01', 2);

-- -----------------------------------------------------------------------------
-- Adresses
-- -----------------------------------------------------------------------------
INSERT INTO adresses (ville, code_postal, nom_rue, numero_rue) VALUES
    ('Paris',   '75001', 'Rue de Rivoli',              '10'),  -- Jean Martin
    ('Lyon',    '69001', 'Rue de la République',       '25'),  -- Marie Dupont
    ('Nantes',  '44000', 'Rue du Général de Gaulle',   '7');   -- Pierre Roche

-- -----------------------------------------------------------------------------
-- Départements
-- -----------------------------------------------------------------------------
INSERT INTO departements_ecole (nom) VALUES
    ('Mathématiques'),
    ('Informatique'),
    ('Physique');

-- -----------------------------------------------------------------------------
-- Bibliothécaire
-- -----------------------------------------------------------------------------
INSERT INTO bibliothecaires (nom, prenom, login, mot_de_passe) VALUES
    ('Petit', 'Valérie', 'vpetit', 'mdp');

-- -----------------------------------------------------------------------------
-- Particulier — insertion directe dans utilisateurs (ONLY)
-- -----------------------------------------------------------------------------
INSERT INTO utilisateurs (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id) VALUES
    ('Martin', 'Jean', '0601020304', 20.00, 'jmartin', 'mdp', '1985-06-15', 'jean.martin@email.fr', 1);

-- -----------------------------------------------------------------------------
-- Étudiant
-- -----------------------------------------------------------------------------
INSERT INTO etudiants (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id, annee_etude) VALUES
    ('Dupont', 'Marie', '0607080910', 20.00, 'mdupont', 'mdp', '2002-03-20', 'marie.dupont@etud.fr', 2, 'L3');

-- -----------------------------------------------------------------------------
-- Enseignant
-- -----------------------------------------------------------------------------
INSERT INTO enseignants (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email, adresse_id, departement_id) VALUES
    ('Roche', 'Pierre', '0611121314', 20.00, 'proche', 'mdp', '1975-11-08', 'pierre.roche@univ.fr', 3, 1);

-- -----------------------------------------------------------------------------
-- Exemplaires (1 par ouvrage)
-- -----------------------------------------------------------------------------
INSERT INTO exemplaires (est_emprunte, code_barre, delai_emprunt_jours, ouvrage_id) VALUES
    -- La Communauté de l'anneau
    (FALSE, 'EX-0001', 15, (SELECT id FROM livres WHERE isbn = '978-2267011296')),
    -- Les Deux Tours
    (FALSE, 'EX-0002', 15, (SELECT id FROM livres WHERE isbn = '978-2267011555')),
    -- Le Retour du roi
    (FALSE, 'EX-0003', 15, (SELECT id FROM livres WHERE isbn = '978-2267011777')),
    -- Les Mille et Une Nuits
    (FALSE, 'EX-0004', 15, (SELECT id FROM livres WHERE isbn = '978-2070379583')),
    -- 60 millions de consommateurs n°581
    (FALSE, 'EX-0005',  7, (SELECT id FROM revues WHERE numero = 581)),
    -- 60 millions de consommateurs n°582
    (FALSE, 'EX-0006',  7, (SELECT id FROM revues WHERE numero = 582)),
    -- Science & Vie n°1265
    (FALSE, 'EX-0007',  7, (SELECT id FROM revues WHERE numero = 1265));

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
SELECT 'utilisateurs',     COUNT(*) FROM ONLY utilisateurs
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