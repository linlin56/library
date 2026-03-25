import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  CreerUtilisateurResponse,
  CreerUtilisateurService,
  Departement,
} from '../services/creer-utilisateur.service';

@Component({
  selector: 'app-creer-utilisateur',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './creer-utilisateur.html',
  styleUrl: './creer-utilisateur.scss',
})
// Formulaire réactif de création d'un utilisateur avec champs conditionnels selon le statut.
export class CreerUtilisateurComponent implements OnInit {
  readonly formulaire: FormGroup;
  departements: Departement[] = [];
  enCours = false;
  erreur: string | null = null;
  compteExistant: string | null = null;
  resultat: CreerUtilisateurResponse | null = null;

  readonly anneesEtude = ['L1', 'L2', 'L3', 'M1', 'M2'];

  constructor(
    private fb: FormBuilder,
    private service: CreerUtilisateurService,
  ) {
    this.formulaire = this.fb.group({
      nom: ['', Validators.required],
      prenom: ['', Validators.required],
      numero_telephone: ['', [Validators.required, Validators.pattern(/^0[1-9][0-9]{8}$/)]],
      date_naissance: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      statut: ['etudiant', Validators.required],
      annee_etude: ['L3', Validators.required],
      departement_id: [null as number | null],
      numero_rue: ['', Validators.required],
      nom_rue: ['', Validators.required],
      code_postal: ['', Validators.required],
      ville: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.service.getDepartements().subscribe({
      next: (deps) => (this.departements = deps),
    });

    this.formulaire.get('statut')?.valueChanges.subscribe((statut: string) => {
      const anneeCtrl = this.formulaire.get('annee_etude')!;
      const deptCtrl = this.formulaire.get('departement_id')!;
      if (statut === 'etudiant') {
        anneeCtrl.setValidators(Validators.required);
        deptCtrl.clearValidators();
        deptCtrl.setValue(null);
      } else if (statut === 'enseignant') {
        deptCtrl.setValidators(Validators.required);
        anneeCtrl.clearValidators();
        anneeCtrl.setValue(null);
      } else {
        // particulier : aucun champ conditionnel
        anneeCtrl.clearValidators();
        anneeCtrl.setValue(null);
        deptCtrl.clearValidators();
        deptCtrl.setValue(null);
      }
      anneeCtrl.updateValueAndValidity();
      deptCtrl.updateValueAndValidity();
    });
  }

  get statut(): string {
    return this.formulaire.get('statut')?.value as string;
  }

  soumettre(): void {
    if (this.formulaire.invalid) {
      this.formulaire.markAllAsTouched();
      return;
    }

    this.enCours = true;
    this.erreur = null;
    this.compteExistant = null;
    this.resultat = null;

    const val = this.formulaire.value;
    const payload = {
      nom: val.nom,
      prenom: val.prenom,
      numero_telephone: val.numero_telephone as string,
      date_naissance: val.date_naissance as string,
      email: val.email,
      statut: val.statut,
      annee_etude: val.statut === 'etudiant' ? val.annee_etude : undefined,
      departement_id: val.statut === 'enseignant' ? Number(val.departement_id) : undefined,
      numero_rue: val.numero_rue as string,
      nom_rue: val.nom_rue as string,
      code_postal: val.code_postal as string,
      ville: val.ville as string,
    };

    this.service.creerUtilisateur(payload).subscribe({
      next: (res) => {
        this.resultat = res;
        this.formulaire.reset({ statut: 'etudiant', annee_etude: 'L3' });
        this.enCours = false;
      },
      error: (err: HttpErrorResponse) => {
        this.enCours = false;
        if (err.status === 409) {
          this.compteExistant = (err.error as { login_existant: string }).login_existant;
        } else {
          this.erreur =
            (err.error as { erreur?: string })?.erreur ??
            'Une erreur est survenue. Veuillez réessayer.';
        }
      },
    });
  }

  nouveauFormulaire(): void {
    this.resultat = null;
    this.compteExistant = null;
    this.erreur = null;
  }
}
