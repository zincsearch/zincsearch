import { useEffect, useRef, useState } from 'react';
import { useTranslation } from '../../locales';
import DateTime from './DateTime';
import Select from '../Select';
import SyntaxGuide from './SyntaxGuide';
import type { QueryData, TimeRange } from './model';

export const refreshTimes = [['Off', 0], ['5s', 5], ['10s', 10], ['15s', 15], ['30s', 30], ['1m', 60], ['5m', 300], ['15m', 900], ['30m', 1800], ['1h', 3600], ['2h', 7200], ['1d', 86400]] as const;
export default function SearchBar({ value, onChange, onSearch, onRefresh, loading }: {
  value: QueryData; onChange: (value: QueryData) => void; onSearch: (value: QueryData) => void; onRefresh: () => void; loading: boolean;
}) {
  const { t } = useTranslation();
  const [refresh, setRefresh] = useState(0);
  const refreshRef = useRef(onRefresh);
  const submitRef = useRef(onSearch);
  const valueRef = useRef(value);
  const focusedQuery = useRef(value.query);
  refreshRef.current = onRefresh;
  submitRef.current = onSearch;
  valueRef.current = value;
  const timeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  useEffect(() => {
    if (!refresh) return;
    const timer = setInterval(() => refreshRef.current(), refresh * 1000);
    return () => clearInterval(timer);
  }, [refresh]);
  useEffect(() => () => clearTimeout(timeout.current), []);
  const changeTime = (time: TimeRange) => {
    const next = { ...value, time };
    onChange(next);
    clearTimeout(timeout.current);
    if (time.startDate && time.endDate) timeout.current = setTimeout(() => submitRef.current(valueRef.current), 1000);
  };
  return <form className="toolbar search-toolbar flex-wrap" onSubmit={event => { event.preventDefault(); clearTimeout(timeout.current); onSearch(value); }}>
    <input className="min-w-48 flex-1" type="search" data-cy="search-bar-input" aria-label={t('search.typeSearch')} placeholder={t('search.typeSearch')} value={value.query} onChange={event => onChange({ ...value, query: event.target.value })} onFocus={() => { focusedQuery.current = value.query; }} onBlur={() => { if (!loading && focusedQuery.current !== value.query) onSearch(value); }} />
    <SyntaxGuide /><DateTime value={value.time} onChange={changeTime} />
    <button className="primary" data-cy="search-bar-refresh-button" type="submit" disabled={loading}>{loading ? 'Searching…' : 'Search'}</button>
    <label>Auto refresh<Select data-cy="search-bar-button-dropdown" value={refresh} onValueChange={value => setRefresh(Number(value))}>{refreshTimes.map(([label, seconds]) => <option value={seconds} key={seconds}>{label}</option>)}</Select></label>
  </form>;
}
