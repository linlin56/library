import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LanguageService } from '../i18n/language.service';
import { AmountPipe, TranslatePipe } from '../i18n/translate.pipe';
import { Loan, LoanService } from '../services/loan.service';
import { DepositInfo, UserService, UserSummary } from '../services/user.service';

type Step = 'search' | 'select' | 'loans';

@Component({
  selector: 'app-user-loans',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslatePipe, AmountPipe],
  templateUrl: './user-loans.html',
  styleUrl: './user-loans.scss',
})
// Librarian view: find a user, then view their loans and manage their deposit.
export class UserLoansComponent {
  // Search form
  lastName = '';
  firstName = '';
  postalCode = '';
  phoneNumber = '';

  // State
  step: Step = 'search';
  loading = false;
  error: string | null = null;

  // Results
  users: UserSummary[] = [];
  selectedUser: UserSummary | null = null;
  loans: Loan[] = [];
  deposit: DepositInfo | null = null;

  // Deposit editing
  editingDeposit = false;
  newTotalDeposit: number | null = null;
  savingDeposit = false;
  depositError: string | null = null;
  depositSaved = false;

  constructor(
    private userService: UserService,
    private loanService: LoanService,
    private languageService: LanguageService,
  ) {}

  canSearch(): boolean {
    return !!(
      this.lastName.trim() ||
      this.firstName.trim() ||
      this.postalCode.trim() ||
      this.phoneNumber.trim()
    );
  }

  searchUser(): void {
    if (!this.canSearch()) {
      this.error = this.languageService.t('librarian.errors.criteria');
      return;
    }

    this.loading = true;
    this.error = null;
    this.users = [];
    this.selectedUser = null;

    this.userService
      .searchUsers(
        this.lastName.trim(),
        this.firstName.trim(),
        this.postalCode.trim(),
        this.phoneNumber.trim(),
      )
      .subscribe({
        next: (users) => {
          this.users = users;
          this.loading = false;
          if (users.length === 0) {
            this.error = this.languageService.t('librarian.errors.noUser');
          } else {
            this.step = 'select';
          }
        },
        error: () => {
          this.error = this.languageService.t('librarian.errors.search');
          this.loading = false;
        },
      });
  }

  selectUser(user: UserSummary): void {
    this.selectedUser = user;
    this.loading = true;
    this.loans = [];
    this.deposit = null;
    this.editingDeposit = false;
    this.depositError = null;
    this.depositSaved = false;

    this.loanService.listLoans(user.id).subscribe({
      next: (loans) => {
        this.loans = loans;
        this.step = 'loans';
        this.loading = false;
      },
      error: () => {
        this.error = this.languageService.t('librarian.errors.loadLoans');
        this.loading = false;
      },
    });
    this.userService.getDeposit(user.id).subscribe({
      next: (info) => (this.deposit = info),
    });
  }

  openDepositEditing(): void {
    this.newTotalDeposit = this.deposit?.total_deposit ?? null;
    this.editingDeposit = true;
    this.depositError = null;
    this.depositSaved = false;
  }

  cancelDepositEditing(): void {
    this.editingDeposit = false;
    this.depositError = null;
  }

  saveDeposit(): void {
    if (!this.selectedUser || this.newTotalDeposit === null || this.newTotalDeposit < 0) {
      this.depositError = this.languageService.t('librarian.errors.invalidValue');
      return;
    }
    const borrowedAmount = Math.max(
      0,
      (this.deposit?.total_deposit ?? 0) - (this.deposit?.deposit_balance ?? 0),
    );
    if (this.newTotalDeposit < borrowedAmount) {
      this.depositError = this.languageService.t('librarian.errors.minValue', {
        amount: this.languageService.formatAmount(borrowedAmount),
      });
      return;
    }
    const userId = this.selectedUser.id;
    this.savingDeposit = true;
    this.depositError = null;
    this.userService.updateTotalDeposit(userId, this.newTotalDeposit).subscribe({
      next: (fresh) => {
        this.deposit = fresh;
        this.editingDeposit = false;
        this.depositSaved = true;
        this.savingDeposit = false;
      },
      error: (err) => {
        this.depositError =
          err.error?.error ?? this.languageService.t('librarian.errors.updateDeposit');
        this.savingDeposit = false;
      },
    });
  }

  backToSearch(): void {
    this.step = 'search';
    this.selectedUser = null;
    this.loans = [];
    this.users = [];
    this.error = null;
    this.lastName = '';
    this.firstName = '';
    this.postalCode = '';
    this.phoneNumber = '';
  }

  formatRole(role: string): string {
    return this.languageService.roleLabel(role);
  }
}
