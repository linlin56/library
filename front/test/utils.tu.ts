import { expect } from 'chai';
import { describe, it } from 'mocha';

import { buildFiltresParams, estBibliothecaire, parseUtilisateur } from '../src/utils';

// ─── buildFiltresParams ───────────────────────────────────────────────────────

describe('buildFiltresParams', () => {
  it('retourne un objet vide pour des filtres vides', () => {
    expect(buildFiltresParams({})).to.deep.equal({});
  });

  it('inclut le titre quand il est renseigné', () => {
    expect(buildFiltresParams({ titre: 'Le Petit Prince' })).to.deep.equal({
      titre: 'Le Petit Prince',
    });
  });

  it('supprime les espaces en début/fin de titre', () => {
    expect(buildFiltresParams({ titre: '  Dune  ' })).to.deep.equal({ titre: 'Dune' });
  });

  it("exclut le titre s'il ne contient que des espaces", () => {
    expect(buildFiltresParams({ titre: '   ' })).to.deep.equal({});
  });

  it('inclut disponible=true quand la valeur est true', () => {
    expect(buildFiltresParams({ disponible: true })).to.deep.equal({ disponible: 'true' });
  });

  it("n'inclut pas disponible quand la valeur est false", () => {
    expect(buildFiltresParams({ disponible: false })).to.deep.equal({});
  });

  it('mappe codeBarre vers code_barre', () => {
    expect(buildFiltresParams({ codeBarre: '0012345' })).to.deep.equal({
      code_barre: '0012345',
    });
  });

  it('mappe codeRevue vers code_revue', () => {
    expect(buildFiltresParams({ codeRevue: 'REV42' })).to.deep.equal({ code_revue: 'REV42' });
  });

  it('inclut tous les champs renseignés simultanément', () => {
    expect(
      buildFiltresParams({
        titre: 'Le Monde',
        auteur: 'Dupont',
        isbn: '978-3-16-148410-0',
        disponible: true,
      }),
    ).to.deep.equal({
      titre: 'Le Monde',
      auteur: 'Dupont',
      isbn: '978-3-16-148410-0',
      disponible: 'true',
    });
  });
});

// ─── parseUtilisateur ─────────────────────────────────────────────────────────

describe('parseUtilisateur', () => {
  it('retourne null pour un JSON invalide', () => {
    expect(parseUtilisateur('not-json')).to.be.null;
  });

  it('retourne null pour une chaîne vide', () => {
    expect(parseUtilisateur('')).to.be.null;
  });

  it('désérialise un utilisateur valide', () => {
    const user = { id: 1, nom: 'Durand', prenom: 'Alice', role: 'etudiant', message: 'ok' };
    expect(parseUtilisateur(JSON.stringify(user))).to.deep.equal(user);
  });
});

// ─── estBibliothecaire ────────────────────────────────────────────────────────

describe('estBibliothecaire', () => {
  it('retourne false pour null', () => {
    expect(estBibliothecaire(null)).to.be.false;
  });

  it("retourne false pour un utilisateur avec le rôle 'etudiant'", () => {
    expect(estBibliothecaire({ id: 1, nom: 'A', prenom: 'B', role: 'etudiant', message: '' })).to.be
      .false;
  });

  it("retourne false pour un utilisateur avec le rôle 'enseignant'", () => {
    expect(estBibliothecaire({ id: 2, nom: 'C', prenom: 'D', role: 'enseignant', message: '' })).to
      .be.false;
  });

  it("retourne true pour un utilisateur avec le rôle 'bibliothecaire'", () => {
    expect(estBibliothecaire({ id: 3, nom: 'E', prenom: 'F', role: 'bibliothecaire', message: '' }))
      .to.be.true;
  });
});
