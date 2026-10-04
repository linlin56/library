import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { LanguageService } from '../i18n/language.service';
import { AmountPipe, TranslatePipe } from '../i18n/translate.pipe';
import { AuthService } from '../services/auth.service';
import { Loan, LoanService } from '../services/loan.service';
import { DepositInfo, UserService } from '../services/user.service';

@Component({
  selector: 'app-my-loans',
  standalone: true,
  imports: [CommonModule, TranslatePipe, AmountPipe],
  templateUrl: './my-loans.html',
  styleUrl: './my-loans.scss',
})
// Displays the active loans of the logged-in user.
export class MyLoansComponent implements OnInit {
  loans: Loan[] = [];
  loading = true;
  error: string | null = null;
  userName = '';
  deposit: DepositInfo | null = null;
  returning: Set<number> = new Set();

  constructor(
    private loanService: LoanService,
    private authService: AuthService,
    private userService: UserService,
    private languageService: LanguageService,
  ) {}

  ngOnInit(): void {
    const u = this.authService.getLoggedInUser();
    if (!u) return;
    this.userName = `${u.first_name} ${u.last_name}`;
    this.loanService.listLoans(u.id).subscribe({
      next: (loans) => {
        this.loans = loans;
        this.loading = false;
      },
      error: () => {
        this.error = this.languageService.t('loans.errors.load');
        this.loading = false;
      },
    });
    this.userService.getDeposit(u.id).subscribe({
      next: (info) => (this.deposit = info),
    });
  }

  isOverdue(loan: Loan): boolean {
    return loan.overdue;
  }

  markReturned(loan: Loan): void {
    this.returning.add(loan.id);
    this.loanService.returnCopy(loan.id).subscribe({
      next: () => {
        this.loans = this.loans.filter((l) => l.id !== loan.id);
        this.returning.delete(loan.id);
        const u = this.authService.getLoggedInUser();
        if (u) {
          this.userService.getDeposit(u.id).subscribe({
            next: (info) => (this.deposit = info),
          });
        }
      },
      error: () => {
        this.returning.delete(loan.id);
      },
    });
  }
}
