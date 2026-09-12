import { describe, expect, it } from 'vitest';
import { buildSearch, deepKeys, fieldValue, histogramInterval, initialTime, timeBounds } from './model';

describe('search request parity', () => {
  const now = new Date('2026-02-15T12:00:00Z');
  it('builds the default match-all, timestamp bounds, sort and record limit', () => {
    expect(buildSearch({ query: '', time: initialTime }, 100, now)).toEqual({
      query: { bool: { must: [{ range: { '@timestamp': { gte: '2026-02-15T11:30:00.000Z', lt: now.toISOString(), format: '2006-01-02T15:04:05Z07:00' } } }, { match_all: {} }] } },
      sort: ['-@timestamp'], from: 0, size: 100,
      aggs: { histogram: { date_histogram: { field: '@timestamp', calendar_interval: '', fixed_interval: '30s' } } },
    });
  });
  it.each(['Gold', 'City:Paris', '-City:Paris', 'City:Paris Gold', '+City:Paris +Gold', '+City:Paris -Gold', '+Medal:Gold +Year:>2000', 'par*', '"New York"'])('preserves query-string syntax %s', query => {
    const request = buildSearch({ query, time: { ...initialTime, selectedFullTime: true } }, 50, now);
    expect(request.query.bool.must).toEqual([{ query_string: { query } }]);
    expect(request.aggs.histogram).toEqual({ auto_date_histogram: { field: '@timestamp', buckets: 100 } });
  });
  it.each([[60000, '1s'], [300000, '5s'], [600000, '10s'], [1800000, '30s'], [3600000, '1m'], [10800000, '5m'], [86400000, '1h'], [604800000, '1d']] as const)('selects histogram interval at %i ms', (duration, expected) => {
    const interval = histogramInterval(duration);
    expect(interval.fixed_interval || interval.calendar_interval).toBe(expected);
  });
  it('supports weeks and calendar months', () => {
    expect(timeBounds({ ...initialTime, selectedRelativePeriod: 'Weeks', selectedRelativeValue: 2 }, now).start).toEqual(new Date('2026-02-01T12:00:00Z'));
    expect(timeBounds({ ...initialTime, selectedRelativePeriod: 'Months', selectedRelativeValue: 1 }, now).start).toEqual(new Date('2026-01-15T12:00:00Z'));
  });
  it('uses local absolute times, validates ranges and max records', () => {
    const time = { ...initialTime, tab: 'absolute' as const, startDate: '2026-02-10', startTime: '10:30', endDate: '2026-02-11', endTime: '12:00' };
    expect(timeBounds(time).start).toEqual(new Date(2026, 1, 10, 10, 30));
    expect(() => buildSearch({ query: '', time: { ...time, endDate: '2026-02-01' } }, 100)).toThrow('valid start');
    expect(() => buildSearch({ query: '', time }, 0)).toThrow('positive integer');
  });
  it('discovers nested fields and resolves literal dotted keys and arrays', () => {
    const source = { 'a.b': { c: 'literal' }, a: { b: { c: 'nested' } }, list: [{ name: 'first' }], empty: null };
    expect(deepKeys(source)).toEqual(['a.b.c', 'a.b.c', 'list', 'empty']);
    expect(fieldValue(source, 'a.b.c')).toBe('literal');
    expect(fieldValue(source, 'list[0].name')).toBe('first');
    expect(fieldValue(source, 'a.b.missing')).toBeUndefined();
    expect(fieldValue(source, 'empty.child')).toBeUndefined();
  });
});
