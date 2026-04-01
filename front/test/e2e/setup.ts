/**
 * Root hooks Mocha — initialisation de l'environnement de test e2e.
 *
 * Avant chaque session de tests :
 *  1. Applique insert.sql pour remettre la BDD dans un état connu.
 *  2. Démarre le backend Go s'il n'est pas déjà lancé.
 *
 * Variables d'environnement (toutes optionnelles, valeurs par défaut = dev) :
 *   API_URL, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
 */

import { ChildProcess, spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

import { Client } from 'pg';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';

// Résolution des chemins depuis test/e2e/
const BACK_DIR = path.resolve(__dirname, '..', '..', '..', 'back');
const SQL_FILE = path.resolve(__dirname, '..', '..', '..', 'db', 'scripts', 'insert.sql');

// Paramètres BDD (mêmes valeurs que back/db/db.env)
const DB_HOST = process.env['DB_HOST'] ?? '127.0.0.1';
const DB_PORT = process.env['DB_PORT'] ?? '5435';
const DB_USER = process.env['DB_USER'] ?? 'postgres';
const DB_PASSWORD = process.env['DB_PASSWORD'] ?? 'admin_bibli';
const DB_NAME = process.env['DB_NAME'] ?? 'bibliotheque';

let backendProcess: ChildProcess | null = null;

// ─── Helpers ──────────────────────────────────────────────────────────────────

async function seedDatabase(): Promise<void> {
  console.log('\n  → Initialisation de la BDD de test (insert.sql)...');
  const sql = fs.readFileSync(SQL_FILE, 'utf-8');
  const client = new Client({
    host: DB_HOST,
    port: parseInt(DB_PORT),
    user: DB_USER,
    password: DB_PASSWORD,
    database: DB_NAME,
  });
  await client.connect();
  await client.query(sql);
  await client.end();
  console.log('  ✔ BDD initialisée avec les données de test');
}

async function isBackendUp(): Promise<boolean> {
  try {
    await fetch(`${API_URL}/api/ouvrages`, { signal: AbortSignal.timeout(1500) });
    return true;
  } catch {
    return false;
  }
}

async function waitForBackend(timeoutMs = 20000): Promise<void> {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    if (await isBackendUp()) return;
    await new Promise((r) => setTimeout(r, 500));
  }
  throw new Error(`Backend (${API_URL}) inaccessible après ${timeoutMs / 1000}s`);
}

// ─── Hooks Mocha ──────────────────────────────────────────────────────────────

export const mochaHooks = {
  async beforeAll() {
    // Si le backend tourne déjà, ne pas relancer insert.sql :
    // DROP + CREATE TABLE casserait les connexions actives du backend (→ 500 sur les queries).
    // Chaque fichier de test appellera TestContext.before() qui fera un TRUNCATE (sans DROP).
    if (await isBackendUp()) {
      console.log('  ✔ Backend déjà démarré — réutilisation (sans re-seed)\n');
      return;
    }

    // Backend non disponible : créer le schéma + données initiales, puis démarrer.
    await seedDatabase();

    console.log('  → Démarrage du backend (go run main.go)...');
    backendProcess = spawn('go', ['run', 'main.go'], {
      cwd: BACK_DIR,
      env: {
        ...process.env,
        DB_HOST,
        DB_PORT,
        DB_USER,
        DB_PASSWORD,
        DB_NAME,
      },
      stdio: 'pipe',
    });

    backendProcess.on('error', (err) => {
      throw new Error(`Impossible de lancer le backend : ${err.message}`);
    });

    await waitForBackend();
    console.log('  ✔ Backend prêt\n');
  },

  async afterAll() {
    if (backendProcess) {
      backendProcess.kill('SIGTERM');
      backendProcess = null;
    }
  },
};
