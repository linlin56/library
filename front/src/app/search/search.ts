import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { LanguageService } from '../i18n/language.service';
import { AmountPipe, TranslatePipe } from '../i18n/translate.pipe';
import { AuthService } from '../services/auth.service';
import { AvailableCopy, CopyDetail, LoanPreview, LoanService } from '../services/loan.service';
import { Item, SearchFilters, SearchService } from '../services/search.service';

type ModalStep = 'input' | 'preview' | 'success' | 'failure';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslatePipe, AmountPipe],
  templateUrl: './search.html',
  styleUrl: './search.scss',
})
// Item search with filters, with a built-in loan modal.
export class SearchComponent implements OnInit {
  // Search fields
  title = '';
  author = '';
  isbn = '';
  issueNumber = '';
  barcode = '';
  available = false;
  advancedSearchOpen = false;

  // Search state
  results: Item[] = [];
  loading = false;
  error: string | null = null;
  hasSearched = false;
  userId: number | null = null;
  isLoggedIn = false;
  isLibrarian = false;

  // Loan modal state
  modalVisible = false;
  modalStep: ModalStep = 'input';
  selectedItem: Item | null = null;
  modalBarcode = '';
  preview: LoanPreview | null = null;
  modalError: string | null = null;
  modalLoading = false;
  availableCopies: AvailableCopy[] = [];
  loadingCopies = false;

  // Copy editing modal state
  copiesModalVisible = false;
  allCopies: CopyDetail[] = [];
  newBarcode = '';
  loadingAllCopies = false;
  copiesError: string | null = null;
  creatingCopy = false;
  copyCreated = false;

  constructor(
    private searchService: SearchService,
    private loanService: LoanService,
    private authService: AuthService,
    private languageService: LanguageService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    const u = this.authService.getLoggedInUser();
    this.isLoggedIn = u !== null;
    this.userId = u?.id ?? null;
    this.isLibrarian = u?.role === 'librarian';
  }

  // ---- Search ----

  search(): void {
    const atLeastOneField =
      this.title.trim() ||
      this.author.trim() ||
      this.isbn.trim() ||
      this.issueNumber.trim() ||
      this.barcode.trim();

    if (!atLeastOneField) {
      this.error = this.languageService.t('search.errors.atLeastOne');
      return;
    }

    this.loading = true;
    this.error = null;
    this.results = [];
    this.hasSearched = false;

    const filters: SearchFilters = {
      title: this.title,
      author: this.author,
      isbn: this.isbn,
      issueNumber: this.issueNumber,
      barcode: this.barcode,
      available: this.available,
    };

    this.searchService.search(filters).subscribe({
      next: (items) => {
        this.results = items ?? [];
        this.hasSearched = true;
        this.loading = false;
      },
      error: () => {
        this.error = this.languageService.t('search.errors.search');
        this.results = [];
        this.hasSearched = true;
        this.loading = false;
      },
    });
  }

  clear(): void {
    this.title = '';
    this.author = '';
    this.isbn = '';
    this.issueNumber = '';
    this.barcode = '';
    this.available = false;
    this.results = [];
    this.error = null;
    this.hasSearched = false;
    this.loading = false;
  }

  getTypeBadge(item: Item): string {
    return this.languageService.t(
      item.type === 'book' ? 'search.types.book' : 'search.types.journal',
    );
  }

  getTypeBadgeClass(item: Item): string {
    return item.type === 'book' ? 'bg-info' : 'bg-success';
  }

  // ---- Loan modal ----

  openModal(item: Item): void {
    this.selectedItem = item;
    this.modalStep = 'input';
    this.modalBarcode = '';
    this.preview = null;
    this.modalError = null;
    this.modalLoading = false;
    this.availableCopies = [];
    this.loadingCopies = true;
    this.modalVisible = true;
    this.loanService.getAvailableCopies(item.id).subscribe({
      next: (copies) => {
        this.availableCopies = copies;
        this.loadingCopies = false;
      },
      error: () => {
        this.modalError = this.languageService.t('search.errors.loadCopies');
        this.loadingCopies = false;
      },
    });
  }

  closeModal(): void {
    this.modalVisible = false;
  }

  checkLoan(): void {
    if (!this.modalBarcode.trim() || !this.userId) return;

    this.modalLoading = true;
    this.modalError = null;

    this.loanService.check(this.userId, this.modalBarcode.trim()).subscribe({
      next: (preview) => {
        this.preview = preview;
        this.modalStep = 'preview';
        this.modalLoading = false;
      },
      error: (err) => {
        this.modalError = err.error?.error ?? this.languageService.t('search.errors.verify');
        this.modalLoading = false;
      },
    });
  }

  confirmLoan(): void {
    if (!this.userId || !this.modalBarcode) return;

    this.modalLoading = true;
    this.loanService.borrow(this.userId, this.modalBarcode.trim()).subscribe({
      next: () => {
        this.modalStep = 'success';
        this.modalLoading = false;
      },
      error: (err) => {
        this.modalError = err.error?.error ?? this.languageService.t('search.errors.save');
        this.modalStep = 'failure';
        this.modalLoading = false;
      },
    });
  }

  declineLoan(): void {
    this.modalStep = 'failure';
    this.modalError = null;
  }

  viewMyLoans(): void {
    this.closeModal();
    this.router.navigate(['/my-loans']);
  }

  // ---- Copy editing modal ----

  openCopiesModal(item: Item): void {
    this.selectedItem = item;
    this.allCopies = [];
    this.newBarcode = '';
    this.copiesError = null;
    this.creatingCopy = false;
    this.copyCreated = false;
    this.loadingAllCopies = true;
    this.copiesModalVisible = true;
    this.loanService.getAllCopies(item.id).subscribe({
      next: (copies) => {
        this.allCopies = copies;
        this.loadingAllCopies = false;
      },
      error: () => {
        this.copiesError = this.languageService.t('search.errors.loadCopies');
        this.loadingAllCopies = false;
      },
    });
  }

  closeCopiesModal(): void {
    this.copiesModalVisible = false;
    this.selectedItem = null;
  }

  createCopy(): void {
    if (!this.newBarcode.trim() || !this.selectedItem) return;
    this.creatingCopy = true;
    this.copiesError = null;
    this.copyCreated = false;
    this.loanService.createCopy(this.selectedItem.id, this.newBarcode.trim()).subscribe({
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
}
