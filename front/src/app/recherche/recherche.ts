import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ConnexionService } from '../services/connexion.service';
import { EmpruntService, ExemplaireDisponible, PreviewEmprunt } from '../services/emprunt.service';
import { FiltresRecherche, Ouvrage, RechercheService } from '../services/recherche.service';

type EtapeModal = 'saisie' | 'preview' | 'succes' | 'echec';

@Component({
  selector: 'app-recherche',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './recherche.html',
  styleUrl: './recherche.scss',
})
// Recherche d'ouvrages par filtres avec modal d'emprunt intégrée.
export class RechercheComponent implements OnInit {
  // Champs recherche
  titre = '';
  auteur = '';
  isbn = '';
  codeRevue = '';
  codeBarre = '';
  disponible = false;
  rechercheAvanceeOuverte = false;

  // Etat recherche
  resultats: Ouvrage[] = [];
  enCours = false;
  erreur: string | null = null;
  aRecherche = false;
  utilisateurId: number | null = null;
  estConnecte = false;

  // Etat modal emprunt
  modalVisible = false;
  etapeModal: EtapeModal = 'saisie';
  ouvrageSelectionne: Ouvrage | null = null;
  codeBareModal = '';
  preview: PreviewEmprunt | null = null;
  erreurModal: string | null = null;
  enCoursModal = false;
  exemplairesDisponibles: ExemplaireDisponible[] = [];
  chargementExemplaires = false;

  constructor(
    private rechercheService: RechercheService,
    private empruntService: EmpruntService,
    private connexionService: ConnexionService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    const u = this.connexionService.getUtilisateurConnecte();
    this.estConnecte = u !== null;
    this.utilisateurId = u?.id ?? null;
  }

  // ---- Recherche ----

  rechercher(): void {
    const auMoinsUnChamp =
      this.titre.trim() ||
      this.auteur.trim() ||
      this.isbn.trim() ||
      this.codeRevue.trim() ||
      this.codeBarre.trim();

    if (!auMoinsUnChamp) {
      this.erreur = 'Veuillez renseigner au moins un champ.';
      return;
    }

    this.enCours = true;
    this.erreur = null;
    this.resultats = [];
    this.aRecherche = false;

    const filtres: FiltresRecherche = {
      titre: this.titre,
      auteur: this.auteur,
      isbn: this.isbn,
      codeRevue: this.codeRevue,
      codeBarre: this.codeBarre,
      disponible: this.disponible,
    };

    this.rechercheService.rechercher(filtres).subscribe({
      next: (ouvrages) => {
        this.resultats = ouvrages ?? [];
        this.aRecherche = true;
        this.enCours = false;
      },
      error: () => {
        this.erreur = 'Erreur lors de la recherche. Veuillez reessayer.';
        this.resultats = [];
        this.aRecherche = true;
        this.enCours = false;
      },
    });
  }

  effacer(): void {
    this.titre = '';
    this.auteur = '';
    this.isbn = '';
    this.codeRevue = '';
    this.codeBarre = '';
    this.disponible = false;
    this.resultats = [];
    this.erreur = null;
    this.aRecherche = false;
    this.enCours = false;
  }

  getTypeBadge(ouvrage: Ouvrage): string {
    return ouvrage.type === 'livre' ? 'Livre' : 'Revue';
  }

  getTypeBadgeClass(ouvrage: Ouvrage): string {
    return ouvrage.type === 'livre' ? 'bg-info' : 'bg-success';
  }

  // ---- Modal emprunt ----

  ouvrirModal(ouvrage: Ouvrage): void {
    this.ouvrageSelectionne = ouvrage;
    this.etapeModal = 'saisie';
    this.codeBareModal = '';
    this.preview = null;
    this.erreurModal = null;
    this.enCoursModal = false;
    this.exemplairesDisponibles = [];
    this.chargementExemplaires = true;
    this.modalVisible = true;
    this.empruntService.getExemplairesDisponibles(ouvrage.id).subscribe({
      next: (ex) => {
        this.exemplairesDisponibles = ex;
        this.chargementExemplaires = false;
      },
      error: () => {
        this.erreurModal = 'Impossible de charger les exemplaires.';
        this.chargementExemplaires = false;
      },
    });
  }

  fermerModal(): void {
    this.modalVisible = false;
  }

  verifierEmprunt(): void {
    if (!this.codeBareModal.trim() || !this.utilisateurId) return;

    this.enCoursModal = true;
    this.erreurModal = null;

    this.empruntService.verifier(this.utilisateurId, this.codeBareModal.trim()).subscribe({
      next: (preview) => {
        this.preview = preview;
        this.etapeModal = 'preview';
        this.enCoursModal = false;
      },
      error: (err) => {
        this.erreurModal = err.error?.erreur ?? 'Impossible de verifier cet emprunt.';
        this.enCoursModal = false;
      },
    });
  }

  confirmerEmprunt(): void {
    if (!this.utilisateurId || !this.codeBareModal) return;

    this.enCoursModal = true;
    this.empruntService.emprunter(this.utilisateurId, this.codeBareModal.trim()).subscribe({
      next: () => {
        this.etapeModal = 'succes';
        this.enCoursModal = false;
      },
      error: (err) => {
        this.erreurModal = err.error?.erreur ?? "Erreur lors de l'enregistrement.";
        this.etapeModal = 'echec';
        this.enCoursModal = false;
      },
    });
  }

  refuserEmprunt(): void {
    this.etapeModal = 'echec';
    this.erreurModal = null;
  }

  voirMesEmprunts(): void {
    this.fermerModal();
    this.router.navigate(['/mes-emprunts']);
  }
}
