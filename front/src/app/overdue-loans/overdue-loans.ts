import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { LanguageService } from '../i18n/language.service';
import { TranslatePipe } from '../i18n/translate.pipe';
import { LoanService, OverdueLoan } from '../services/loan.service';

@Component({
  selector: 'app-overdue-loans',
  standalone: true,
  imports: [CommonModule, TranslatePipe],
  templateUrl: './overdue-loans.html',
  styleUrl: './overdue-loans.scss',
})
// Lists every loan whose due date has passed (librarian view).
export class OverdueLoansComponent implements OnInit {
  loans: OverdueLoan[] = [];
  loading = true;
  error: string | null = null;
  notification: string | null = null;

  constructor(
    private loanService: LoanService,
    private languageService: LanguageService,
  ) {}

  ngOnInit(): void {
    this.loanService.listOverdueLoans().subscribe({
      next: (loans) => {
        this.loans = loans;
        this.loading = false;
      },
      error: () => {
        this.error = this.languageService.t('overdue.errors.load');
        this.loading = false;
      },
    });
  }

  sendReminders(): void {
    this.notification = this.languageService.t('overdue.remindersSent');
    setTimeout(() => (this.notification = null), 4000);
  }
}
