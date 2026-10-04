import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { LanguageService } from '../i18n/language.service';
import { TranslatePipe } from '../i18n/translate.pipe';
import { CreateUserResponse, CreateUserService, Department } from '../services/create-user.service';

@Component({
  selector: 'app-create-user',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslatePipe],
  templateUrl: './create-user.html',
  styleUrl: './create-user.scss',
})
// Reactive form to create a user, with conditional fields depending on the status.
export class CreateUserComponent implements OnInit {
  readonly form: FormGroup;
  departments: Department[] = [];
  loading = false;
  error: string | null = null;
  existingAccount: string | null = null;
  result: CreateUserResponse | null = null;

  readonly studyYears = ['L1', 'L2', 'L3', 'M1', 'M2'];

  constructor(
    private fb: FormBuilder,
    private service: CreateUserService,
    private languageService: LanguageService,
  ) {
    this.form = this.fb.group({
      last_name: ['', Validators.required],
      first_name: ['', Validators.required],
      phone_number: ['', [Validators.required, Validators.pattern(/^0[1-9][0-9]{8}$/)]],
      birth_date: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      status: ['student', Validators.required],
      study_year: ['L3', Validators.required],
      department_id: [null as number | null],
      street_number: ['', Validators.required],
      street_name: ['', Validators.required],
      postal_code: ['', Validators.required],
      city: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.service.getDepartments().subscribe({
      next: (departments) => (this.departments = departments),
    });

    this.form.get('status')?.valueChanges.subscribe((status: string) => {
      const yearCtrl = this.form.get('study_year')!;
      const deptCtrl = this.form.get('department_id')!;
      if (status === 'student') {
        yearCtrl.setValidators(Validators.required);
        deptCtrl.clearValidators();
        deptCtrl.setValue(null);
      } else if (status === 'teacher') {
        deptCtrl.setValidators(Validators.required);
        yearCtrl.clearValidators();
        yearCtrl.setValue(null);
      } else {
        // individual: no conditional field
        yearCtrl.clearValidators();
        yearCtrl.setValue(null);
        deptCtrl.clearValidators();
        deptCtrl.setValue(null);
      }
      yearCtrl.updateValueAndValidity();
      deptCtrl.updateValueAndValidity();
    });
  }

  get status(): string {
    return this.form.get('status')?.value as string;
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.loading = true;
    this.error = null;
    this.existingAccount = null;
    this.result = null;

    const val = this.form.value;
    const payload = {
      last_name: val.last_name,
      first_name: val.first_name,
      phone_number: val.phone_number as string,
      birth_date: val.birth_date as string,
      email: val.email,
      status: val.status,
      study_year: val.status === 'student' ? val.study_year : undefined,
      department_id: val.status === 'teacher' ? Number(val.department_id) : undefined,
      street_number: val.street_number as string,
      street_name: val.street_name as string,
      postal_code: val.postal_code as string,
      city: val.city as string,
    };

    this.service.createUser(payload).subscribe({
      next: (res) => {
        this.result = res;
        this.form.reset({ status: 'student', study_year: 'L3' });
        this.loading = false;
      },
      error: (err: HttpErrorResponse) => {
        this.loading = false;
        if (err.status === 409) {
          this.existingAccount = (err.error as { existing_login: string }).existing_login;
        } else {
          this.error =
            (err.error as { error?: string })?.error ??
            this.languageService.t('createUser.errors.generic');
        }
      },
    });
  }

  newForm(): void {
    this.result = null;
    this.existingAccount = null;
    this.error = null;
  }
}
