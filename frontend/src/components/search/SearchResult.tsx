import { Fragment, useState } from 'react';
import { useTranslation } from '../../locales';
import HighLight from './HighLight';
import Select from '../Select';
import Histogram from './Histogram';
import { displayValue, fieldValue, formatTimestamp } from './model';
import type { Hit, SearchResponse } from './model';

export default function SearchResult({ result, columns, query, interval, loading, maxRecords, onMaxRecords }: {
  result: SearchResponse | null; columns: string[]; query: string; interval: string; loading: boolean; maxRecords: number; onMaxRecords: (value: number) => void;
}) {
  const { t } = useTranslation();
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [expanded, setExpanded] = useState<string[]>([]);
  const [mode, setMode] = useState<'table' | 'json'>('table');
  const [sort, setSort] = useState({ field: '', descending: false });
  const rows = result?.hits.hits || [];
  const fields = ['@timestamp', ...(columns.length ? columns.filter(field => field !== '@timestamp') : ['_source'])];
  const cell = (row: Hit, field: string): unknown => field === '_source' ? JSON.stringify(row._source) : field === '@timestamp' ? row['@timestamp'] : ['_id', '_index', '_score'].includes(field) ? row[field as '_id' | '_index' | '_score'] : fieldValue(row._source, field);
  const sorted = sort.field ? [...rows].sort((a, b) => {
    const left = cell(a, sort.field), right = cell(b, sort.field);
    const order = typeof left === 'number' && typeof right === 'number' ? left - right : displayValue(left).localeCompare(displayValue(right));
    return sort.descending ? -order : order;
  }) : rows;
  const pages = perPage ? Math.max(1, Math.ceil(rows.length / perPage)) : 1;
  const current = Math.min(page, pages);
  const start = perPage ? (current - 1) * perPage : 0;
  const visible = perPage ? sorted.slice(start, start + perPage) : sorted;
  return <section className="card search-results flex-1 min-w-0" data-cy="search-result-area" aria-label={t('search.searchResult')} aria-busy={loading}>
    <h2>{t('search.searchResult')}</h2>
    {result && <p role="status">Found {result.hits.total.value.toLocaleString()} hits in {result.took} ms</p>}
    <div className={`result-controls ${result?.aggregations?.histogram?.buckets?.length ? '' : 'result-controls-empty'}`}>
    <Histogram histogram={result?.aggregations?.histogram} interval={interval} />
    <div className="toolbar result-view-toggle"><button type="button" aria-pressed={mode === 'table'} onClick={() => setMode('table')}>Table</button><button type="button" aria-pressed={mode === 'json'} onClick={() => setMode('json')}>JSON</button></div>
    </div>
    {mode === 'json' ? <pre className="overflow-auto whitespace-pre-wrap break-all" aria-label="Results JSON"><HighLight content={JSON.stringify(result || {}, null, 2)} queryString={query} /></pre> : <div className="overflow-x-auto"><table className="w-full"><thead><tr><th aria-label="Expand document" />{fields.map(field => <th key={field} aria-sort={sort.field === field ? sort.descending ? 'descending' : 'ascending' : 'none'}><button type="button" onClick={() => { setSort({ field, descending: sort.field === field && !sort.descending }); setPage(1); }}>{field === '@timestamp' ? t('search.timestamp') : field}</button></th>)}</tr></thead><tbody>
      {visible.map((row, index) => {
        const key = `${row._index || ''}:${row._id}:${start + index}`;
        const open = expanded.includes(key);
        return <Fragment key={key}><tr><td><button type="button" aria-label={`${open ? 'Collapse' : 'Expand'} document ${row._id}`} aria-expanded={open} onClick={() => setExpanded(open ? expanded.filter(item => item !== key) : [...expanded, key])}>{open ? '−' : '+'}</button></td>{fields.map(field => <td className="align-top break-all" key={field}>{field === '@timestamp' ? formatTimestamp(cell(row, field)) : <div className="line-clamp-5"><HighLight content={displayValue(cell(row, field))} queryString={query} /></div>}</td>)}</tr>
          {open && <tr><td colSpan={fields.length + 1}><pre className="whitespace-pre-wrap break-all"><HighLight content={JSON.stringify(row, null, 2)} queryString={query} /></pre></td></tr>}</Fragment>;
      })}
      {!rows.length && <tr><td colSpan={fields.length + 1}>{t('search.noResult')}</td></tr>}
    </tbody></table></div>}
    <div className="toolbar flex-wrap mt-3">
      <label>{t('search.maxRecords')}<input className="w-28" type="number" min="1" value={Number.isNaN(maxRecords) ? '' : maxRecords} onChange={event => onMaxRecords(event.target.value === '' ? NaN : Number(event.target.value))} /></label>
      <label>Records per page:<Select value={perPage} onValueChange={value => { setPerPage(Number(value)); setPage(1); }}>{[5, 10, 20, 50, 100, 0].map(value => <option key={value} value={value}>{value || 'All'}</option>)}</Select></label>
      <span>{rows.length ? start + 1 : 0}-{Math.min(start + visible.length, rows.length)} of {rows.length}</span>
      <button type="button" aria-label="First page" disabled={current === 1} onClick={() => setPage(1)}>«</button>
      <button type="button" aria-label="Previous page" disabled={current === 1} onClick={() => setPage(current - 1)}>‹</button>
      <button type="button" aria-label="Next page" disabled={current === pages} onClick={() => setPage(current + 1)}>›</button>
      <button type="button" aria-label="Last page" disabled={current === pages} onClick={() => setPage(pages)}>»</button>
    </div>
  </section>;
}
