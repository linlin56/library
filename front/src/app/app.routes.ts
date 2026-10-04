import { Routes } from '@angular/router';
import { CreateItemComponent } from './create-item/create-item';
import { CreateUserComponent } from './create-user/create-user';
import { authGuard } from './guards/auth.guard';
import { guestGuard } from './guards/guest.guard';
import { librarianGuard } from './guards/librarian.guard';
import { HomeComponent } from './home/home';
import { LoginComponent } from './login/login';
import { MyLoansComponent } from './my-loans/my-loans';
import { OverdueLoansComponent } from './overdue-loans/overdue-loans';
import { SearchComponent } from './search/search';
import { UserLoansComponent } from './user-loans/user-loans';

// Application routes and the guards that protect them.
export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent, canActivate: [guestGuard] },
  {
    path: 'create-user',
    component: CreateUserComponent,
    canActivate: [librarianGuard],
  },
  {
    path: 'create-item',
    component: CreateItemComponent,
    canActivate: [librarianGuard],
  },
  { path: 'search', component: SearchComponent },
  { path: 'my-loans', component: MyLoansComponent, canActivate: [authGuard] },
  {
    path: 'user-loans',
    component: UserLoansComponent,
    canActivate: [librarianGuard],
  },
  {
    path: 'overdue-loans',
    component: OverdueLoansComponent,
    canActivate: [librarianGuard],
  },
  { path: '**', redirectTo: '' },
];
