import { type ReactNode, useEffect, useId, useRef, useState } from 'react';
import Select from '../Select';

export function Modal(
  { title, onClose, children }: { title: string; onClose: () => void; children: ReactNode },
) {
  const id = useId();
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const previous = document.activeElement as HTMLElement | null;
    ref.current?.focus();
    return () => previous?.focus();
  }, []);
  return (
    <div className='dialog-backdrop'>
      <div
        ref={ref}
        tabIndex={-1}
        role='dialog'
        aria-modal='true'
        aria-labelledby={id}
        className='dialog w-full max-w-4xl max-h-[90vh] overflow-auto'
        onKeyDown={(event) => {
          if (event.key === 'Escape') onClose();
          if (event.key === 'Tab') {
            const elements = ref.current?.querySelectorAll<HTMLElement>(
              'button:not(:disabled), input:not(:disabled), select:not(:disabled), textarea:not(:disabled), [tabindex="0"]',
            );
            if (!elements?.length) return;
            const first = elements[0], last = elements[elements.length - 1];
            if (
              event.shiftKey &&
              (document.activeElement === first || document.activeElement === ref.current)
            ) {
              event.preventDefault();
              last.focus();
            } else if (
              !event.shiftKey &&
              (document.activeElement === last || document.activeElement === ref.current)
            ) {
              event.preventDefault();
              first.focus();
            }
          }
        }}
      >
        <div className='toolbar justify-between'>
          <h2 id={id}>{title}</h2>
          <button type='button' onClick={onClose} aria-label='Close'>×</button>
        </div>
        {children}
      </div>
    </div>
  );
}

export function ErrorMessage({ message }: { message: string }) {
  return message ? <p className='error' role='alert'>{message}</p> : null;
}

export function Pagination(
  { page, size, total, onPage, onSize }: {
    page: number;
    size: number;
    total: number;
    onPage: (page: number) => void;
    onSize: (size: number) => void;
  },
) {
  const pages = Math.max(1, Math.ceil(total / size));
  return (
    <div className='toolbar mt-4'>
      <label>
        Rows per page{' '}
        <Select value={size} onValueChange={(value) => onSize(Number(value))}>
          {[5, 10, 20, 50, 100, 500, 1000].map((n) => <option key={n}>{n}</option>)}
        </Select>
      </label>
      <button disabled={page <= 1} onClick={() => onPage(page - 1)}>Previous</button>
      <span>Page {page} of {pages} · {total} rows</span>
      <button disabled={page >= pages} onClick={() => onPage(page + 1)}>Next</button>
    </div>
  );
}

export function ConfirmDelete(
  { names, onDelete, onClose }: {
    names: string[];
    onDelete: () => Promise<void>;
    onClose: () => void;
  },
) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  return (
    <Modal
      title='Confirm deletion'
      onClose={() => {
        if (!busy) onClose();
      }}
    >
      <p>You are about to delete:</p>
      <ul>{names.map((name) => <li key={name}>{name}</li>)}</ul>
      <ErrorMessage message={error} />
      <div className='toolbar'>
        <button disabled={busy} onClick={onClose}>Cancel</button>
        <button
          className='primary'
          disabled={busy}
          onClick={async () => {
            setBusy(true);
            setError('');
            try {
              await onDelete();
              onClose();
            } catch {
              setError('Deletion failed. Please try again.');
            } finally {
              setBusy(false);
            }
          }}
        >
          Delete
        </button>
      </div>
    </Modal>
  );
}

export type Column<T> = { label: string; render: (row: T) => ReactNode };
export function LocalList<T>({
  title,
  searchLabel,
  addLabel,
  rows,
  name,
  rowKey,
  columns,
  loading,
  error,
  onAdd,
  onRetry,
}: {
  title: string;
  searchLabel: string;
  addLabel: string;
  rows: T[];
  name: (row: T) => string;
  rowKey: (row: T) => string;
  columns: Column<T>[];
  loading: boolean;
  error: string;
  onAdd: () => void;
  onRetry: () => void;
}) {
  const [filter, setFilter] = useState('');
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(20);
  const filtered = rows.filter((row) => name(row).toLowerCase().includes(filter.toLowerCase()));
  const current = Math.min(page, Math.max(1, Math.ceil(filtered.length / size)));
  return (
    <section className='card data-panel'>
      <div className='toolbar justify-between'>
        <h1>{title}</h1>
        <input
          aria-label={searchLabel}
          placeholder={searchLabel}
          value={filter}
          onChange={(e) => {
            setFilter(e.target.value);
            setPage(1);
          }}
        />
        <button className='primary' onClick={onAdd}>{addLabel}</button>
      </div>
      <ErrorMessage message={error} />
      {error && <button onClick={onRetry}>Retry</button>}
      {loading && <p role='status'>Loading…</p>}
      <div className='overflow-x-auto'>
        <table>
          <thead>
            <tr>
              <th>#</th>
              {columns.map((c) => <th key={c.label}>{c.label}</th>)}
            </tr>
          </thead>
          <tbody>
            {filtered.slice((current - 1) * size, current * size).map((row, i) => (
              <tr key={rowKey(row)}>
                <td>{(current - 1) * size + i + 1}</td>
                {columns.map((c) => <td key={c.label}>{c.render(row)}</td>)}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!loading && !filtered.length && <p>No results</p>}
      <Pagination
        page={current}
        size={size}
        total={filtered.length}
        onPage={setPage}
        onSize={(n) => {
          setSize(n);
          setPage(1);
        }}
      />
    </section>
  );
}

export function useList<T>(load: () => Promise<{ data: T[] }>) {
  const [rows, setRows] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [revision, setRevision] = useState(0);
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError('');
    load().then((res) => {
      if (active) setRows(res.data);
    }).catch(() => {
      if (active) setError('Unable to load data. Please try again.');
    }).finally(() => {
      if (active) setLoading(false);
    });
    return () => {
      active = false;
    };
  }, [load, revision]);
  return { rows, loading, error, refresh: () => setRevision((n) => n + 1) };
}

export function formatDate(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}
