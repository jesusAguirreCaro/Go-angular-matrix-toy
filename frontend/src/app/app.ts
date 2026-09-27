import { Component, inject, signal } from '@angular/core';
import { MatrixService } from './services/matrix.service';

const MIN_SIZE = 1;
const MAX_SIZE = 10;

function emptyMatrix(n: number): (number | null)[][] {
  return Array.from({ length: n }, () => Array<number | null>(n).fill(null));
}

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  private readonly matrixService = inject(MatrixService);

  readonly minSize = MIN_SIZE;
  readonly maxSize = MAX_SIZE;

  readonly size = signal(3);
  readonly matrix = signal<(number | null)[][]>(emptyMatrix(3));
  readonly result = signal<number | null>(null);
  readonly error = signal<string | null>(null);
  readonly loading = signal(false);

  onSizeChange(raw: string): void {
    const parsed = Math.floor(Number(raw));
    if (!Number.isFinite(parsed)) return;

    const n = Math.min(MAX_SIZE, Math.max(MIN_SIZE, parsed));
    this.size.set(n);
    this.matrix.set(emptyMatrix(n));
    this.resetOutput();
  }

  setCell(row: number, col: number, raw: string): void {
    const trimmed = raw.trim();
    const value = trimmed === '' ? null : Number(trimmed);

    this.matrix.update((m) => {
      const copy = m.map((r) => [...r]);
      copy[row][col] = value !== null && Number.isFinite(value) ? value : null;
      return copy;
    });
  }

  randomize(): void {
    const n = this.size();
    this.matrix.set(
      Array.from({ length: n }, () =>
        Array.from({ length: n }, () => Math.floor(Math.random() * 19) - 9),
      ),
    );
    this.resetOutput();
  }

  clear(): void {
    this.matrix.set(emptyMatrix(this.size()));
    this.resetOutput();
  }

  fillEmptyWithZero(): void {
    this.matrix.update((m) =>
      m.map((row) => row.map((v) => (v === null ? 0 : v))),
    );
  }

  calculate(): void {
    this.loading.set(true);
    this.error.set(null);
    this.result.set(null);

    // Empty cells are treated as 0 when sending to the server.
    const payload = this.matrix().map((row) => row.map((v) => v ?? 0));

    this.matrixService.determinant(payload).subscribe({
      next: (res) => {
        this.result.set(res.determinant);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err?.error?.error ?? 'Could not reach the server.');
        this.loading.set(false);
      },
    });
  }

  private resetOutput(): void {
    this.result.set(null);
    this.error.set(null);
  }
}
