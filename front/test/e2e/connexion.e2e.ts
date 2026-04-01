/**
 * Tests e2e — Écran Connexion
 *
 * Ces tests appellent directement l'API REST du backend (Go/Gin).
 * Ils nécessitent que le serveur et la base de données soient démarrés.
 *
 * Variables d'environnement :
 *   API_URL  (défaut : http://localhost:8080)
 */

import { expect } from 'chai';
import { before, describe, it } from 'mocha';

const API_URL = process.env['API_URL'] ?? 'http://localhost:8080';
const ENDPOINT = `${API_URL}/api/connexion`;

// ─── Vérification de disponibilité du backend ─────────────────────────────────

async function backendDisponible(): Promise<boolean> {
  try {
    await fetch(`${API_URL}/api/ouvrages`, { signal: AbortSignal.timeout(3000) });
    return true;
  } catch {
    return false;
  }
}

before(async function () {
  if (!(await backendDisponible())) {
    console.warn(`\n  ⚠ Backend inaccessible (${API_URL}) — tests e2e ignorés.\n`);
    this.skip();
  }
});

// ─── Cas nominal — étudiant ───────────────────────────────────────────────────

describe('[E2E] Connexion — étudiant valide (amartin / mdp)', () => {
  let status: number;
  let body: Record<string, unknown>;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: 'amartin', mot_de_passe: 'mdp' }),
    });
    status = res.status;
    body = (await res.json()) as Record<string, unknown>;
  });

  it('retourne HTTP 200', () => expect(status).to.equal(200));
  it('retourne le nom "Martin"', () => expect(body['nom']).to.equal('Martin'));
  it('retourne le prénom "Alice"', () => expect(body['prenom']).to.equal('Alice'));
  it('retourne le rôle "utilisateur"', () => expect(body['role']).to.equal('utilisateur'));
  it('retourne un id numérique positif', () =>
    expect(body['id']).to.be.a('number').and.greaterThan(0));
  it('retourne un message de confirmation', () =>
    expect(body['message']).to.be.a('string').and.not.empty);
});

// ─── Cas nominal — bibliothécaire ────────────────────────────────────────────

describe('[E2E] Connexion — bibliothécaire valide (mdupont / mdp)', () => {
  let status: number;
  let body: Record<string, unknown>;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: 'mdupont', mot_de_passe: 'mdp' }),
    });
    status = res.status;
    body = (await res.json()) as Record<string, unknown>;
  });

  it('retourne HTTP 200', () => expect(status).to.equal(200));
  it('retourne le rôle "bibliothecaire"', () => expect(body['role']).to.equal('bibliothecaire'));
  it('retourne le nom "Dupont"', () => expect(body['nom']).to.equal('Dupont'));
  it('retourne le prénom "Marie"', () => expect(body['prenom']).to.equal('Marie'));
});

// ─── Mauvais mot de passe ─────────────────────────────────────────────────────

describe('[E2E] Connexion — mauvais mot de passe', () => {
  let status: number;
  let body: Record<string, unknown>;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: 'amartin', mot_de_passe: 'mauvais' }),
    });
    status = res.status;
    body = (await res.json()) as Record<string, unknown>;
  });

  it('retourne HTTP 401', () => expect(status).to.equal(401));
  it("retourne un champ 'erreur'", () => expect(body['erreur']).to.be.a('string').and.not.empty);
});

// ─── Login inexistant ─────────────────────────────────────────────────────────

describe('[E2E] Connexion — login inexistant', () => {
  let status: number;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: 'utilisateur_inconnu', mot_de_passe: 'mdp' }),
    });
    status = res.status;
  });

  it('retourne HTTP 401', () => expect(status).to.equal(401));
});

// ─── Body incomplet — champs manquants ────────────────────────────────────────

describe('[E2E] Connexion — body sans mot de passe', () => {
  let status: number;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: 'amartin' }),
    });
    status = res.status;
  });

  it('retourne HTTP 400', () => expect(status).to.equal(400));
});

describe('[E2E] Connexion — body vide', () => {
  let status: number;

  before(async () => {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({}),
    });
    status = res.status;
  });

  it('retourne HTTP 400', () => expect(status).to.equal(400));
});
