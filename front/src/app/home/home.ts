import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { TranslatePipe } from '../i18n/translate.pipe';
import { AuthService } from '../services/auth.service';
import { HomeService, HomeStats } from '../services/home.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink, TranslatePipe],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
// Home page: displays the library statistics.
export class HomeComponent implements OnInit {
  stats: HomeStats | null = null;
  error = false;
  isLoggedIn = false;

  constructor(
    private homeService: HomeService,
    private authService: AuthService,
  ) {}

  ngOnInit(): void {
    this.isLoggedIn = this.authService.isLoggedIn();
    this.homeService.getStats().subscribe({
      next: (data) => (this.stats = data),
      error: () => (this.error = true),
    });
  }
}
