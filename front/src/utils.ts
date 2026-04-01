/**
 * Fonctions utilitaires pures (sans dépendance Angular) utilisées par les services.
 * Elles peuvent être testées avec Mocha directement sous Node.js.
 */

export interface FiltresRecherche {
  titre?: string;
  auteur?: string;
  isbn?: string;
  codeBarre?: string;
  codeRevue?: string;
  disponible?: boolean;
}

export interface UtilisateurConnecte {
  id: number;
  nom: string;
  prenom: string;
  role: string;
  message: string;
}

/**
 * Convertit un objet FiltresRecherche en dictionnaire de paramètres URL.
 * Les champs vides ou composés uniquement d'espaces sont exclus.
 */
export function buildFiltresParams(filtres: FiltresRecherche): Record<string, string> {
  const params: Record<string, string> = {};
  if (filtres.titre?.trim()) params['titre'] = filtres.titre.trim();
  if (filtres.auteur?.trim()) params['auteur'] = filtres.auteur.trim();
  if (filtres.isbn?.trim()) params['isbn'] = filtres.isbn.trim();
  if (filtres.codeBarre?.trim()) params['code_barre'] = filtres.codeBarre.trim();
  if (filtres.codeRevue?.trim()) params['code_revue'] = filtres.codeRevue.trim();
  if (filtres.disponible) params['disponible'] = 'true';
  return params;
}

/**
 * Désérialise un JSON brut en UtilisateurConnecte.
 * Retourne null si le JSON est invalide.
 */
export function parseUtilisateur(raw: string): UtilisateurConnecte | null {
  try {
    return JSON.parse(raw) as UtilisateurConnecte;
  } catch {
    return null;
  }
}

/**
 * Indique si l'utilisateur connecté a le rôle bibliothécaire.
 */
export function estBibliothecaire(user: UtilisateurConnecte | null): boolean {
  return user?.role === 'bibliothecaire';
}
