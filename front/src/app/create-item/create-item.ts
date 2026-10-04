import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LanguageService } from '../i18n/language.service';
import { TranslatePipe } from '../i18n/translate.pipe';
import { Author, ItemService, Location } from '../services/item.service';
import { CopyDetail, LoanService } from '../services/loan.service';

type Step = 'form' | 'copies';

@Component({
  selector: 'app-create-item',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslatePipe],
  templateUrl: './create-item.html',
  styleUrl: './create-item.scss',
})
// Form to create a book or a journal, with author autocompletion (librarian view).
export class CreateItemComponent implements OnInit {
  step: Step = 'form';

  type: 'book' | 'journal' = 'book';

  // Shared fields
  title = '';
  deposit: number | null = null;

  // Book fields — author autocomplete
  isbn = '';
  authorQuery = ''; // text typed in the field
  authorSuggestions: Author[] = []; // filtered suggestions
  selectedAuthor: Author | null = null; // existing author picked
  suggestionsOpen = false;
  newAuthorVisible = false; // first/last name fields of a new author
  authorLastName = '';
  authorFirstName = '';

  // Journal fields
  issueNumber: number | null = null;
  publicationDate = '';

  // Form state
  authors: Author[] = [];
  locations: Location[] = [];
  locationId: number | null = null;
  loading = false;
  error: string | null = null;
  success: string | null = null;

  // Step 2 — copies
  createdItemTitle = '';
  itemId: number | null = null;
  allCopies: CopyDetail[] = [];
  loadingCopies = false;
  newBarcode = '';
  creatingCopy = false;
  copyCreated = false;
  copiesError: string | null = null;

  constructor(
    private service: ItemService,
    private loanService: LoanService,
    private languageService: LanguageService,
  ) {}

  ngOnInit(): void {
    this.service.getAuthors().subscribe({
      next: (authors) => (this.authors = authors),
    });
    this.service.getLocations().subscribe({
      next: (locations) => (this.locations = locations),
    });
  }

  changeType(): void {
    this.error = null;
    this.success = null;
  }

  filterAuthors(): void {
    this.selectedAuthor = null;
    this.newAuthorVisible = false;
    const q = this.authorQuery.trim().toLowerCase();
    if (!q) {
      this.authorSuggestions = [];
      this.suggestionsOpen = false;
      return;
    }
    this.authorSuggestions = this.authors.filter(
      (a) => a.last_name.toLowerCase().includes(q) || a.first_name.toLowerCase().includes(q),
    );
    this.suggestionsOpen = true;
  }

  selectAuthor(author: Author): void {
    this.selectedAuthor = author;
    this.authorQuery = `${author.first_name} ${author.last_name}`;
    this.authorSuggestions = [];
    this.suggestionsOpen = false;
    this.newAuthorVisible = false;
  }

  chooseNewAuthor(): void {
    this.selectedAuthor = null;
    this.authorSuggestions = [];
    this.suggestionsOpen = false;
    this.newAuthorVisible = true;
    this.authorLastName = '';
    this.authorFirstName = '';
  }

  closeSuggestions(): void {
    // delay so that a click on a suggestion is handled first
    setTimeout(() => (this.suggestionsOpen = false), 150);
  }

  submit(): void {
    this.error = null;
    this.success = null;

    if (!this.title.trim()) {
      this.error = this.languageService.t('createItem.errors.titleRequired');
      return;
    }
    if (this.deposit === null || this.deposit < 0) {
      this.error = this.languageService.t('createItem.errors.depositInvalid');
      return;
    }

    this.loading = true;

    if (this.type === 'journal') {
      if (!this.issueNumber || this.issueNumber <= 0) {
        this.error = this.languageService.t('createItem.errors.numberRequired');
        this.loading = false;
        return;
      }
      if (!this.publicationDate) {
        this.error = this.languageService.t('createItem.errors.dateRequired');
        this.loading = false;
        return;
      }
      if (!this.locationId) {
        this.error = this.languageService.t('createItem.errors.locationRequired');
        this.loading = false;
        return;
      }
      this.service
        .createJournal({
          title: this.title.trim(),
          deposit: this.deposit,
          issue_number: this.issueNumber,
          publication_date: this.publicationDate,
          location_id: this.locationId,
        })
        .subscribe({
          next: (res) => {
            this.loading = false;
            this.goToCopiesStep(res.id, this.title.trim());
          },
          error: (err) => {
            this.error =
              err.error?.error ?? this.languageService.t('createItem.errors.journalCreate');
            this.loading = false;
          },
        });
    } else {
      if (!this.isbn.trim()) {
        this.error = this.languageService.t('createItem.errors.isbnRequired');
        this.loading = false;
        return;
      }
      if (!this.locationId) {
        this.error = this.languageService.t('createItem.errors.locationRequired');
        this.loading = false;
        return;
      }
      const payload: Parameters<ItemService['createBook']>[0] = {
        title: this.title.trim(),
        deposit: this.deposit,
        isbn: this.isbn.trim(),
        location_id: this.locationId,
      };
      if (this.selectedAuthor) {
        payload.author_id = this.selectedAuthor.id;
      } else if (this.newAuthorVisible) {
        if (!this.authorLastName.trim() || !this.authorFirstName.trim()) {
          this.error = this.languageService.t('createItem.errors.authorNamesRequired');
          this.loading = false;
          return;
        }
        payload.author_last_name = this.authorLastName.trim();
        payload.author_first_name = this.authorFirstName.trim();
      } else {
        this.error = this.languageService.t('createItem.errors.authorRequired');
        this.loading = false;
        return;
      }
      this.service.createBook(payload).subscribe({
        next: (res) => {
          this.loading = false;
          const wasNewAuthor = this.newAuthorVisible;
          if (wasNewAuthor) {
            this.service.getAuthors().subscribe({ next: (a) => (this.authors = a) });
          }
          this.goToCopiesStep(res.id, this.title.trim());
        },
        error: (err) => {
          this.error = err.error?.error ?? this.languageService.t('createItem.errors.bookCreate');
          this.loading = false;
        },
      });
    }
  }

  private goToCopiesStep(id: number, title: string): void {
    this.itemId = id;
    this.createdItemTitle = title;
    this.allCopies = [];
    this.newBarcode = '';
    this.creatingCopy = false;
    this.copyCreated = false;
    this.copiesError = null;
    this.loadingCopies = true;
    this.step = 'copies';
    this.loanService.getAllCopies(id).subscribe({
      next: (copies) => {
        this.allCopies = copies;
        this.loadingCopies = false;
      },
      error: () => {
        this.loadingCopies = false;
      },
    });
  }

  createCopy(): void {
    if (!this.newBarcode.trim() || !this.itemId) return;
    this.creatingCopy = true;
    this.copiesError = null;
    this.copyCreated = false;
    this.loanService.createCopy(this.itemId, this.newBarcode.trim()).subscribe({
      next: (res) => {
        this.allCopies = [
          ...this.allCopies,
          { id: res.id, barcode: this.newBarcode.trim(), is_borrowed: false },
        ];
        this.newBarcode = '';
        this.copyCreated = true;
        this.creatingCopy = false;
      },
      error: (err) => {
        this.copiesError = err.error?.error ?? this.languageService.t('copies.createError');
        this.creatingCopy = false;
      },
    });
  }

  createAnotherItem(): void {
    this.step = 'form';
    this.resetFields();
  }

  private resetFields(): void {
    this.title = '';
    this.deposit = null;
    this.isbn = '';
    this.authorQuery = '';
    this.selectedAuthor = null;
    this.authorSuggestions = [];
    this.suggestionsOpen = false;
    this.newAuthorVisible = false;
    this.authorLastName = '';
    this.authorFirstName = '';
    this.issueNumber = null;
    this.publicationDate = '';
    this.locationId = null;
    this.error = null;
    this.success = null;
  }
}
