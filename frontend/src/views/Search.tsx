import { useEffect, useRef, useState } from 'react';
import searchService from '../services/search';
import SearchBar from '../components/search/SearchBar';
import IndexList from '../components/search/IndexList';
import SearchResult from '../components/search/SearchResult';
import { buildSearch, deepKeys, histogramInterval, initialTime, timeBounds } from '../components/search/model';
import type { QueryData, SearchResponse } from '../components/search/model';

export default function Search() {
  const [index, setIndex] = useState('');
  const [columns, setColumns] = useState<string[]>([]);
  const [fields, setFields] = useState<string[]>([]);
  const [draft, setDraft] = useState<QueryData>({ query: '', time: { ...initialTime } });
  const [submitted, setSubmitted] = useState(draft);
  const [result, setResult] = useState<SearchResponse | null>(null);
  const [maxRecords, setMaxRecords] = useState(100);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [interval, setInterval] = useState('30s');
  const requestId = useRef(0);
  const busy = useRef(false);
  useEffect(() => {
    busy.current = false;
    setLoading(false);
    return () => { requestId.current++; };
  }, []);

  async function search(data: QueryData, name = index, indexChanged = false) {
    if (busy.current && !indexChanged) return;
    const id = ++requestId.current;
    busy.current = true;
    setLoading(true);
    setError('');
    try {
      const query = buildSearch(data, maxRecords);
      // The unchanged service declares string, but the Vue caller sends an object to axios.
      const response = await searchService.search({ index: name, query: query as unknown as string });
      if (id !== requestId.current) return;
      const next: SearchResponse = response.data;
      setResult(next);
      setSubmitted(data);
      setFields(previous => [...new Set([...previous, ...next.hits.hits.flatMap(hit => deepKeys(hit._source))])]);
      const bounds = timeBounds(data.time);
      const histogram = histogramInterval(+bounds.end - +bounds.start);
      setInterval(data.time.selectedFullTime ? '1s' : histogram.fixed_interval || histogram.calendar_interval);
    } catch (cause) {
      if (id === requestId.current) setError(cause instanceof Error && !('response' in cause) ? cause.message : 'Search failed. Please try again.');
    } finally {
      if (id === requestId.current) { busy.current = false; setLoading(false); }
    }
  }
  const selectIndex = (name: string) => {
    if (name === index) return;
    setIndex(name);
    setColumns([]);
    setFields(['_id', '_index', '_score']);
    setResult(null);
    const data = { ...draft, query: '' };
    setDraft(data);
    void search(data, name, true);
  };
  return <section aria-label="Search workspace">
    <SearchBar value={draft} onChange={setDraft} onSearch={data => { void search(data); }} onRefresh={() => { void search(submitted); }} loading={loading} />
    {error && <p role="alert" className="error">{error}</p>}
    <div className="flex flex-col md:flex-row gap-4 mt-4">
      <IndexList index={index} fields={fields} columns={columns} onIndex={selectIndex} onColumns={setColumns} />
      <SearchResult key={index} result={result} columns={columns} query={submitted.query} interval={interval} loading={loading} maxRecords={maxRecords} onMaxRecords={setMaxRecords} />
    </div>
  </section>;
}
