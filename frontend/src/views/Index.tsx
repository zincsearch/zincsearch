import { useEffect, useState } from 'react';
import indexService from '../services/index';
import { useTranslation } from '../locales';
import { ConfirmDelete, ErrorMessage, Pagination } from '../components/management/Common';
import SchemaEditor, { type Schema, SchemaPreview } from '../components/management/SchemaEditor';

type IndexRecord = Schema & {
  name: string;
  shard_num: number;
  storage_type: string;
  stats: { doc_num: number; storage_size: number; wal_size?: number };
};
export function storageSize(bytes: number) {
  const unit = bytes > 1024 ** 3 ? 3 : bytes > 1024 ** 2 ? 2 : 1;
  return `${(bytes / 1024 ** unit).toFixed(2)} ${['', 'KB', 'MB', 'GB'][unit]}`;
}
export default function Index() {
  const { t } = useTranslation();
  const [rows, setRows] = useState<IndexRecord[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(20);
  const [sort, setSort] = useState('name');
  const [descending, setDescending] = useState(false);
  const [filter, setFilter] = useState('');
  const [revision, setRevision] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selected, setSelected] = useState<string[]>([]);
  const [adding, setAdding] = useState(false);
  const [preview, setPreview] = useState<IndexRecord | null>(null);
  const [deleting, setDeleting] = useState<string[] | null>(null);
  const refresh = () => setRevision((n) => n + 1);
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError('');
    indexService.list(page, size, sort, descending, filter).then((res) => {
      if (!active) return;
      const last = Math.max(1, Math.ceil(res.data.page.total / size));
      setTotal(res.data.page.total);
      if (page > last) {
        setPage(last);
        return;
      }
      setRows(res.data.list);
    }).catch(() => {
      if (active) setError('Unable to load indexes. Please try again.');
    }).finally(() => {
      if (active) setLoading(false);
    });
    return () => {
      active = false;
    };
  }, [page, size, sort, descending, filter, revision]);
  const columns = ['name', 'doc_num', 'shard_num', 'storage_size', 'storage_type'];
  return (
    <>
    <dl className="metric-grid" aria-label="Index overview" aria-busy={loading}>
      {[
        ['Indexes', total, false],
        ['Documents', rows.reduce((sum, row) => sum + row.stats.doc_num, 0).toLocaleString(), true],
        ['Storage', storageSize(rows.reduce((sum, row) => sum + row.stats.storage_size, 0)), true],
        ['Shards', rows.reduce((sum, row) => sum + row.shard_num, 0).toLocaleString(), true],
      ].map(([label, value, pageOnly]) => <div className="card metric-card" key={String(label)}>
        <dt>{label}{pageOnly && <small> · this page</small>}</dt>
        <dd>{loading || error ? '—' : value}</dd>
      </div>)}
    </dl>
    <section className='card data-panel'>
      <div className='toolbar justify-between'>
        <h1>{t('index.header')}</h1>
        <input
          aria-label={t('index.search')}
          placeholder={t('index.search')}
          value={filter}
          onChange={(e) => {
            setFilter(e.target.value);
            setPage(1);
            setSelected([]);
          }}
        />
        <button onClick={refresh}>Refresh</button>
        <button className='primary' onClick={() => setAdding(true)}>{t('index.add')}</button>
        <button
          onClick={() => {
            if (selected.length) setDeleting(selected);
            else setError('Please select index for deletion');
          }}
        >
          {t('index.delete')}
        </button>
      </div>
      <ErrorMessage message={error} />
      {loading && <p role='status'>Loading…</p>}
      <div className='overflow-x-auto'>
        <table>
          <thead>
            <tr>
              <th>
                <input
                  type='checkbox'
                  aria-label='Select all indexes on this page'
                  checked={rows.length > 0 && rows.every((r) => selected.includes(r.name))}
                  onChange={(e) =>
                    setSelected((current) =>
                      e.target.checked
                        ? [...new Set([...current, ...rows.map((r) => r.name)])]
                        : current.filter((name) => !rows.some((r) => r.name === name))
                    )}
                />
              </th>
              <th>#</th>
              {columns.map((column) => (
                <th
                  key={column}
                  aria-sort={sort === column ? descending ? 'descending' : 'ascending' : 'none'}
                >
                  <button
                    onClick={() => {
                      setSort(column);
                      setDescending(sort === column ? !descending : false);
                      setPage(1);
                    }}
                  >
                    {column.toUpperCase()}
                  </button>
                </th>
              ))}
              <th>ACTIONS</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr key={row.name}>
                <td>
                  <input
                    type='checkbox'
                    aria-label={`Select ${row.name}`}
                    checked={selected.includes(row.name)}
                    onChange={(e) =>
                      setSelected((current) =>
                        e.target.checked
                          ? [...current, row.name]
                          : current.filter((name) => name !== row.name)
                      )}
                  />
                </td>
                <td>{(page - 1) * size + i + 1}</td>
                <td>
                  <button onClick={() => setPreview(row)}>{row.name}</button>
                </td>
                <td>{row.stats.doc_num}</td>
                <td>{row.shard_num}</td>
                <td>{storageSize(row.stats.storage_size)}</td>
                <td>{row.storage_type}</td>
                <td>
                  <div className='toolbar'>
                    <button
                      aria-label={`Preview ${row.name}`}
                      onClick={() => setPreview(row)}
                    >
                      Preview
                    </button>
                    <button
                      aria-label={`Delete ${row.name}`}
                      onClick={() => setDeleting([row.name])}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!loading && !rows.length && <p>No results</p>}
      <Pagination
        page={page}
        size={size}
        total={total}
        onPage={setPage}
        onSize={(n) => {
          setSize(n);
          setPage(1);
        }}
      />
      {adding && (
        <SchemaEditor
          kind='index'
          onClose={() => setAdding(false)}
          onUpdated={() => {
            setAdding(false);
            refresh();
          }}
        />
      )}
      {preview && (
        <SchemaPreview
          name={preview.name}
          summary={{
            Name: preview.name,
            'Docs Count': preview.stats.doc_num,
            'Shards Num': preview.shard_num,
            'Storage Size': storageSize(preview.stats.storage_size),
            'Storage Type': preview.storage_type,
            'WAL Entries': preview.stats.wal_size || 0,
          }}
          schema={preview}
          data={{
            name: preview.name,
            doc_num: preview.stats.doc_num,
            shard_num: preview.shard_num,
            storage_type: preview.storage_type,
            storage_size: storageSize(preview.stats.storage_size),
            wal_size: preview.stats.wal_size || 0,
            settings: preview.settings || {},
            mappings: preview.mappings || {},
          }}
          onClose={() => setPreview(null)}
        />
      )}
      {deleting && (
        <ConfirmDelete
          names={deleting}
          onClose={() => setDeleting(null)}
          onDelete={async () => {
            await indexService.delete(deleting.join(','));
            setSelected((current) => current.filter((name) => !deleting.includes(name)));
            refresh();
          }}
        />
      )}
    </section>
    </>
  );
}
