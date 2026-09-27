import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface DeterminantResponse {
  determinant: number;
}

@Injectable({ providedIn: 'root' })
export class MatrixService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = 'http://localhost:8080/api';

  determinant(matrix: number[][]): Observable<DeterminantResponse> {
    return this.http.post<DeterminantResponse>(`${this.baseUrl}/determinant`, { matrix });
  }
}
