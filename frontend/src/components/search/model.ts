export type Period = 'Minutes' | 'Hours' | 'Days' | 'Weeks' | 'Months';
export interface TimeRange {
  tab: 'relative' | 'absolute';
  startDate: string; startTime: string; endDate: string; endTime: string;
  selectedRelativePeriod: Period; selectedRelativeValue: number; selectedFullTime: boolean;
}
export interface QueryData { query: string; time: TimeRange }
export const initialTime: TimeRange = {
  tab: 'relative', startDate: '', startTime: '', endDate: '', endTime: '',
  selectedRelativePeriod: 'Minutes', selectedRelativeValue: 30, selectedFullTime: false,
};
export interface Hit { _id: string; _index?: string; _score?: number; '@timestamp'?: string; _source: Record<string, unknown> }
export interface SearchResponse {
  took: number;
  hits: { total: { value: number }; hits: Hit[] };
  aggregations?: { histogram?: { interval?: string; buckets: { key: string | number; doc_count: number }[] } };
}
export function timeBounds(time: TimeRange, now = new Date()) {
  const end = new Date(now);
  const start = new Date(now);
  if (time.tab === 'absolute') {
    const parse = (day: string, clock: string) => day || clock ? new Date(`${day.replace(/\//g, '-')}T${clock || '00:00'}`) : new Date(now);
    return { start: parse(time.startDate, time.startTime), end: parse(time.endDate, time.endTime) };
  }
  const value = time.selectedRelativeValue;
  switch (time.selectedRelativePeriod) {
    case 'Minutes': start.setMinutes(start.getMinutes() - value); break;
    case 'Hours': start.setHours(start.getHours() - value); break;
    case 'Days': start.setDate(start.getDate() - value); break;
    case 'Weeks': start.setDate(start.getDate() - value * 7); break;
    case 'Months': start.setMonth(start.getMonth() - value); break;
  }
  return { start, end };
}
export function histogramInterval(duration: number) {
  if (duration >= 604800000) return { calendar_interval: '1d', fixed_interval: '' };
  if (duration >= 86400000) return { calendar_interval: '1h', fixed_interval: '' };
  if (duration >= 10800000) return { calendar_interval: '', fixed_interval: '5m' };
  if (duration >= 3600000) return { calendar_interval: '1m', fixed_interval: '' };
  if (duration >= 1800000) return { calendar_interval: '', fixed_interval: '30s' };
  if (duration >= 600000) return { calendar_interval: '', fixed_interval: '10s' };
  if (duration >= 300000) return { calendar_interval: '', fixed_interval: '5s' };
  return { calendar_interval: '1s' };
}
export function buildSearch(data: QueryData, size: number, now = new Date()) {
  if (!Number.isInteger(size) || size < 1) throw new Error('Max records must be a positive integer.');
  const must: Record<string, unknown>[] = [];
  let histogram: Record<string, unknown>;
  if (data.time.selectedFullTime) {
    histogram = { auto_date_histogram: { field: '@timestamp', buckets: 100 } };
  } else {
    const { start, end } = timeBounds(data.time, now);
    if (!Number.isFinite(+start) || !Number.isFinite(+end) || start >= end) throw new Error('Select a valid start and end time.');
    must.push({ range: { '@timestamp': { gte: start.toISOString(), lt: end.toISOString(), format: '2006-01-02T15:04:05Z07:00' } } });
    histogram = { date_histogram: { field: '@timestamp', ...histogramInterval(+end - +start) } };
  }
  must.push(data.query === '' ? { match_all: {} } : { query_string: { query: data.query } });
  return { query: { bool: { must } }, sort: ['-@timestamp'], from: 0, size, aggs: { histogram } };
}
export function deepKeys(value: unknown): string[] {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return [];
  return Object.entries(value).flatMap(([key, child]) => {
    const nested = deepKeys(child);
    return nested.length ? nested.map(path => `${key}.${path}`) : [key];
  });
}
export function fieldValue(value: unknown, path: string): unknown {
  if (!value || typeof value !== 'object') return undefined;
  const object = value as Record<string, unknown>;
  if (Object.prototype.hasOwnProperty.call(object, path)) return object[path];
  const parts = path.replace(/\[(\w+)\]/g, '.$1').replace(/^\./, '').split('.');
  for (let i = parts.length - 1; i > 0; i--) {
    const key = parts.slice(0, i).join('.');
    if (Object.prototype.hasOwnProperty.call(object, key)) return fieldValue(object[key], parts.slice(i).join('.'));
  }
  return undefined;
}
export function displayValue(value: unknown): string {
  return value == null ? '' : typeof value === 'object' ? JSON.stringify(value) : String(value);
}
export function formatTimestamp(value: unknown) {
  if (value == null) return '';
  const date = new Date(String(value));
  if (!Number.isFinite(+date)) return String(value);
  const pad = (n: number, digits = 2) => String(n).padStart(digits, '0');
  const offset = -date.getTimezoneOffset();
  return `${date.toLocaleString('en-US', { month: 'short' })} ${pad(date.getDate())}, ${date.getFullYear()} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}.${pad(date.getMilliseconds(), 3)} ${offset >= 0 ? '+' : '-'}${pad(Math.floor(Math.abs(offset) / 60))}:${pad(Math.abs(offset) % 60)}`;
}
