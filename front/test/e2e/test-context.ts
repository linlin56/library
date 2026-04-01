/**
 * TestContext — réinitialisation de la BDD entre fichiers de test.
 *
 * Usage dans chaque fichier .e2e.ts :
 *
 *   const ctx = new TestContext();
 *
 *   describe('Suite ...', () => {
 *     before(async () => {
 *       await ctx.before();               // vide toutes les tables, réinitialise les séquences
 *       await ctx.exec(`INSERT ...`);     // insère les données propres à ce fichier
 *     });
 *     after(async () => await ctx.after());
 *     // ... describes imbriqués ...
 *   });
 */

import { Client } from 'pg';

export class TestContext {
  private readonly client: Client;

  constructor() {
    this.client = new Client({
      host: process.env['DB_HOST'] ?? '127.0.0.1',
      port: parseInt(process.env['DB_PORT'] ?? '5435'),
      user: process.env['DB_USER'] ?? 'postgres',
      password: process.env['DB_PASSWORD'] ?? 'admin_bibli',
      database: process.env['DB_NAME'] ?? 'bibliotheque',
    });
  }

  /**
   * Ouvre la connexion, applique les migrations de schéma manquantes,
   * puis vide toutes les tables avec RESTART IDENTITY CASCADE.
   * Utilise TRUNCATE (pas DROP) pour ne pas invalider les connexions du backend.
   */
  async before(): Promise<void> {
    await this.client.connect();

    // Colonnes requises par le backend mais absentes du schéma initial
    await this.client.query(
      `ALTER TABLE ouvrages ADD COLUMN IF NOT EXISTS emplacement_id INT REFERENCES emplacements(id)`,
    );
    await this.client.query(
      `ALTER TABLE revues ADD COLUMN IF NOT EXISTS date_parution DATE`,
    );

    await this.client.query(`
      TRUNCATE
        enseignants, etudiants, ONLY utilisateurs,
        bibliothecaires,
        exemplaires,
        livres, revues, ONLY ouvrages,
        auteurs, adresses,
        emplacements, categories,
        departements_ecole
      RESTART IDENTITY CASCADE
    `);
  }

  /** Ferme la connexion. */
  async after(): Promise<void> {
    await this.client.end();
  }

  /** Exécute un ou plusieurs INSERT (séparés par ;) dans le contexte du test. */
  async exec(sql: string): Promise<void> {
    const statements = sql
      .split(';')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);
    for (const stmt of statements) {
      await this.client.query(stmt);
    }
  }
}
