import { useEffect, useState } from 'react';
import indexService from '../../services/index';
import Select from '../Select';
import { useTranslation } from '../../locales';

export default function IndexList({ index, fields, columns, onIndex, onColumns }: {
  index: string; fields: string[]; columns: string[]; onIndex: (name: string) => void; onColumns: (columns: string[]) => void;
}) {
  const { t } = useTranslation();
  const [filter, setFilter] = useState('');
  const [fieldFilter, setFieldFilter] = useState('');
  const [options, setOptions] = useState<string[]>([]);
  const [error, setError] = useState('');
  useEffect(() => {
    let active = true;
    setError('');
    indexService.nameList(filter).then(response => { if (active) setOptions(response.data || []); }).catch(() => { if (active) setError('Unable to load indexes.'); });
    return () => { active = false; };
  }, [filter]);
  return <aside className="card search-fields w-full md:w-60 shrink-0" aria-label="Indexes and fields">
    <label>{t('search.selectIndex')}<input aria-label="Filter indexes" value={filter} onChange={event => setFilter(event.target.value)} /></label>
    <Select className="w-full" data-cy="index-dropdown" aria-label={t('search.selectIndex')} value={index} onValueChange={onIndex}>
      <option value="">{t('search.selectIndex')}</option>
      {[...new Set([...(index ? [index] : []), ...options])].map(name => <option key={name}>{name}</option>)}
    </Select>
    {!options.length && <p>{t('search.noResult')}</p>}
    {error && <p role="alert" className="error">{error}</p>}
    <input className="w-full mt-3" data-cy="index-field-search-input" aria-label={t('search.searchField')} placeholder={t('search.searchField')} value={fieldFilter} onChange={event => setFieldFilter(event.target.value)} />
    <div className="max-h-[70vh] overflow-auto">{fields.filter(field => field.toLowerCase().includes(fieldFilter.toLowerCase())).map(field => <label className="flex gap-2 py-1 break-all" key={field}><input type="checkbox" checked={columns.includes(field)} onChange={() => onColumns(columns.includes(field) ? columns.filter(column => column !== field) : [...columns, field])} />{field}</label>)}</div>
  </aside>;
}
