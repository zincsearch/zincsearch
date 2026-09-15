import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import Search from '../../views/Search';
import { Activity } from 'react';
import SearchBar from './SearchBar';
import HighLight from './HighLight';
import { initialTime } from './model';
import type { SearchResponse } from './model';

const mocks = vi.hoisted(() => ({ search: vi.fn(), nameList: vi.fn() }));
vi.mock('../../services/search', () => ({ default: { search: mocks.search } }));
vi.mock('../../services/index', () => ({ default: { nameList: mocks.nameList } }));
vi.mock('../../locales', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const response = (name = 'logs', count = 23): SearchResponse => ({
  took: 4, hits: { total: { value: 1234 }, hits: Array.from({ length: count }, (_, i) => ({ _id: String(i), _index: name, '@timestamp': '2026-02-15T12:00:00Z', _source: { message: `Gold ${i}`, nested: { city: 'Paris' }, tags: ['one', 'two'] } })) },
  aggregations: { histogram: { interval: '1h', buckets: [{ key: '2026-02-15T12:00:00Z', doc_count: count }] } },
});
async function selectIndex(name = 'logs') {
  fireEvent.click(screen.getByRole('combobox', { name: 'search.selectIndex' }));
  fireEvent.click(await screen.findByRole('option', { name }));
  await screen.findByText('Found 1,234 hits in 4 ms');
}
beforeEach(() => {
  mocks.nameList.mockReset().mockResolvedValue({ data: ['logs', 'other'] });
  mocks.search.mockReset().mockResolvedValue({ data: response() });
});
afterEach(() => { cleanup(); vi.useRealTimers(); });

describe('Search migration', () => {
  it('selects and filters indexes, discovers columns, expands JSON and paginates locally', async () => {
    render(<Search />);
    await selectIndex();
    expect(mocks.search).toHaveBeenCalledWith(expect.objectContaining({ index: 'logs', query: expect.objectContaining({ size: 100, from: 0 }) }));
    expect(screen.getByRole('img', { name: 'Document counts over time' })).toBeInTheDocument();
    expect(screen.getByText('1-20 of 23')).toBeInTheDocument();
    fireEvent.click(screen.getByLabelText('Next page'));
    expect(screen.getByText('21-23 of 23')).toBeInTheDocument();
    expect(mocks.search).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByLabelText('Records per page:'));
    fireEvent.click(screen.getByRole('option', { name: 'All' }));
    expect(screen.getByText('1-23 of 23')).toBeInTheDocument();
    fireEvent.click(screen.getByLabelText('Expand document 0'));
    expect(screen.getByLabelText('Collapse document 0')).toHaveAttribute('aria-expanded', 'true');
    fireEvent.click(screen.getByRole('checkbox', { name: 'nested.city' }));
    expect(screen.getByRole('columnheader', { name: 'nested.city' })).toBeInTheDocument();
    expect(screen.getAllByText('Paris').length).toBeGreaterThan(0);
    expect(mocks.search).toHaveBeenCalledTimes(1);
    fireEvent.change(screen.getByLabelText('search.searchField'), { target: { value: 'NESTED' } });
    expect(screen.queryByRole('checkbox', { name: 'message' })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'JSON' }));
    expect(screen.getByLabelText('Results JSON')).toHaveTextContent('aggregations');
    fireEvent.change(screen.getByLabelText('Filter indexes'), { target: { value: 'oth' } });
    await waitFor(() => expect(mocks.nameList).toHaveBeenLastCalledWith('oth'));
  });
  it('submits query syntax, full time, limits, highlights and resets on index change', async () => {
    const { container } = render(<Search />);
    await selectIndex();
    fireEvent.change(screen.getByLabelText('search.typeSearch'), { target: { value: 'City:Paris Gold' } });
    fireEvent.change(screen.getByLabelText('search.maxRecords'), { target: { value: '50' } });
    fireEvent.click(screen.getByRole('button', { name: '30 Minutes' }));
    fireEvent.click(screen.getByLabelText('FullTime'));
    fireEvent.click(screen.getByRole('button', { name: 'Close' }));
    fireEvent.click(screen.getByRole('button', { name: 'Search' }));
    await waitFor(() => expect(mocks.search).toHaveBeenCalledTimes(2));
    expect(mocks.search.mock.calls[1][0].query).toMatchObject({ size: 50, query: { bool: { must: [{ query_string: { query: 'City:Paris Gold' } }] } }, aggs: { histogram: { auto_date_histogram: { buckets: 100 } } } });
    await waitFor(() => expect(container.querySelector('mark')).toHaveTextContent('Gold'));
    fireEvent.click(screen.getByRole('checkbox', { name: 'message' }));
    await selectIndex('other');
    expect(screen.getByLabelText('search.typeSearch')).toHaveValue('');
    expect(screen.getByRole('checkbox', { name: 'message' })).not.toBeChecked();
    expect(screen.getByRole('columnheader', { name: '_source' })).toBeInTheDocument();
    const source = screen.getAllByRole('cell').find(cell => cell.textContent?.includes('"message":"Gold 0"'));
    expect(source).toHaveTextContent(JSON.stringify({ message: 'Gold 0', nested: { city: 'Paris' }, tags: ['one', 'two'] }));
    expect(source).not.toHaveTextContent('"_id"');
  });
  it('recovers after a failed request and clears stale histogram data', async () => {
    render(<Search />);
    await selectIndex();
    mocks.search.mockRejectedValueOnce(new Error('Network unavailable'));
    fireEvent.click(screen.getByRole('button', { name: 'Search' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Network unavailable');
    mocks.search.mockResolvedValueOnce({ data: { took: 0, hits: { total: { value: 0 }, hits: [] } } });
    fireEvent.click(screen.getByRole('button', { name: 'Search' }));
    await screen.findByText('Found 0 hits in 0 ms');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.queryByRole('img')).not.toBeInTheDocument();
    expect(screen.getByText('0-0 of 0')).toBeInTheDocument();
    expect(screen.getByLabelText('Next page')).toBeDisabled();
  });
  it('can search again after hiding an in-flight search', async () => {
    mocks.search.mockImplementationOnce(() => new Promise(() => {}));
    const { rerender } = render(<Activity mode="visible"><Search /></Activity>);
    fireEvent.click(screen.getByRole('combobox', { name: 'search.selectIndex' }));
    fireEvent.click(await screen.findByRole('option', { name: 'logs' }));
    rerender(<Activity mode="hidden"><Search /></Activity>);
    rerender(<Activity mode="visible"><Search /></Activity>);
    fireEvent.click(screen.getByRole('button', { name: 'Search' }));
    await waitFor(() => expect(mocks.search).toHaveBeenCalledTimes(2));
    expect(await screen.findByText('1-20 of 23')).toBeInTheDocument();
  });
  it('ignores late responses when switching indexes', async () => {
    let resolveOld!: (value: { data: SearchResponse }) => void;
    mocks.search.mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    render(<Search />);
    fireEvent.click(screen.getByRole('combobox', { name: 'search.selectIndex' }));
    fireEvent.click(await screen.findByRole('option', { name: 'logs' }));
    mocks.search.mockResolvedValueOnce({ data: response('other', 1) });
    await selectIndex('other');
    await act(async () => resolveOld({ data: response('logs') }));
    expect(screen.getByText('1-1 of 1')).toBeInTheDocument();
    expect(within(screen.getByRole('table')).queryByText(/"_index":"logs"/)).not.toBeInTheDocument();
  });
  it('cleans up refresh timers and refreshes the latest callback', () => {
    vi.useFakeTimers();
    const refresh = vi.fn(), nextRefresh = vi.fn();
    const props = { value: { query: '', time: initialTime }, onChange: vi.fn(), onSearch: vi.fn(), onRefresh: refresh, loading: false };
    const { rerender, unmount } = render(<SearchBar {...props} />);
    fireEvent.click(screen.getByLabelText('Auto refresh'));
    fireEvent.click(screen.getByRole('option', { name: '5s' }));
    act(() => vi.advanceTimersByTime(5000));
    expect(refresh).toHaveBeenCalledTimes(1);
    rerender(<SearchBar {...props} onRefresh={nextRefresh} />);
    act(() => vi.advanceTimersByTime(5000));
    expect(nextRefresh).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByLabelText('Auto refresh'));
    fireEvent.click(screen.getByRole('option', { name: 'Off' }));
    act(() => vi.advanceTimersByTime(5000));
    expect(nextRefresh).toHaveBeenCalledTimes(1);
    unmount();
    expect(vi.getTimerCount()).toBe(0);
  });
  it('supports relative presets, custom periods, absolute debounce and syntax help', async () => {
    render(<Search />);
    await selectIndex();
    fireEvent.click(screen.getByRole('button', { name: '30 Minutes' }));
    fireEvent.click(screen.getByRole('button', { name: '2 Weeks' }));
    expect(screen.getAllByRole('button', { name: '2 Weeks' })).toHaveLength(2);
    fireEvent.change(screen.getByLabelText('Relative value'), { target: { value: '3' } });
    fireEvent.click(screen.getByLabelText('Relative period'));
    fireEvent.click(screen.getByRole('option', { name: 'Months' }));
    expect(screen.getAllByRole('button', { name: '3 Months' })).toHaveLength(2);
    fireEvent.click(screen.getByRole('tab', { name: 'absolute' }));
    vi.useFakeTimers();
    fireEvent.change(screen.getByLabelText('Start Date'), { target: { value: '2026-02-01' } });
    fireEvent.change(screen.getByLabelText('Start Time'), { target: { value: '10:00' } });
    fireEvent.change(screen.getByLabelText('End Date'), { target: { value: '2026-02-02' } });
    fireEvent.change(screen.getByLabelText('End Time'), { target: { value: '11:00' } });
    await act(async () => vi.advanceTimersByTime(999));
    expect(mocks.search).toHaveBeenCalledTimes(1);
    await act(async () => vi.advanceTimersByTime(1));
    expect(mocks.search).toHaveBeenCalledTimes(2);
    expect(mocks.search.mock.calls[1][0].query.query.bool.must[0].range['@timestamp']).toMatchObject({ gte: new Date(2026, 1, 1, 10).toISOString(), lt: new Date(2026, 1, 2, 11).toISOString() });
    fireEvent.click(screen.getByRole('button', { name: 'Close' }));
    fireEvent.click(screen.getByText('search.syntaxGuide'));
    expect(screen.getByText('+Medal:Gold +Year:>2000')).toBeVisible();
  });
  it('closes the time range and syntax guide popovers when pointing outside', async () => {
    render(<Search />);
    await selectIndex();
    fireEvent.click(screen.getByRole('button', { name: '30 Minutes' }));
    fireEvent.pointerDown(screen.getByLabelText('Relative value'));
    expect(screen.getByRole('tab', { name: 'relative' })).toBeInTheDocument();
    fireEvent.click(screen.getByLabelText('Relative period'));
    fireEvent.pointerDown(screen.getByRole('option', { name: 'Hours' }));
    fireEvent.click(screen.getByRole('option', { name: 'Hours' }));
    expect(screen.getByRole('tab', { name: 'relative' })).toBeInTheDocument();
    fireEvent.pointerDown(document.body);
    expect(screen.queryByRole('tab', { name: 'relative' })).toBeNull();
    fireEvent.click(screen.getByText('search.syntaxGuide'));
    expect(screen.getByText('+Medal:Gold +Year:>2000')).toBeVisible();
    fireEvent.pointerDown(document.body);
    expect(screen.getByText('+Medal:Gold +Year:>2000')).not.toBeVisible();
  });
  it('provides an explicit close action for syntax help without submitting a search', () => {
    const onSearch = vi.fn();
    render(<SearchBar value={{ query: '', time: initialTime }} onChange={vi.fn()} onSearch={onSearch} onRefresh={vi.fn()} loading={false} />);
    fireEvent.click(screen.getByText('search.syntaxGuide'));
    const guide = screen.getByRole('region', { name: 'search.syntaxGuide' });
    expect(guide).toHaveClass('search-popover', 'syntax-guide-panel');
    expect(within(guide).getByText('+Medal:Gold +Year:>2000')).toBeVisible();
    fireEvent.click(within(guide).getByRole('button', { name: 'Close' }));
    expect(guide).not.toBeVisible();
    expect(onSearch).not.toHaveBeenCalled();
  });
  it('keeps absolute ranges and preset controls in responsive containers', () => {
    const time = { ...initialTime, tab: 'absolute' as const, startDate: '2026-02-01', startTime: '10:00', endDate: '2026-02-02', endTime: '11:00' };
    const props = { onChange: vi.fn(), onSearch: vi.fn(), onRefresh: vi.fn(), loading: false };
    const { rerender } = render(<SearchBar {...props} value={{ query: '', time }} />);
    const trigger = screen.getByRole('button', { name: '2026-02-01 10:00 - 2026-02-02 11:00' });
    expect(trigger.parentElement).toHaveClass('search-popover-control');
    fireEvent.click(trigger);
    const panel = screen.getByRole('region', { name: 'Time range' });
    expect(panel).toHaveClass('search-popover', 'time-range-panel');
    for (const name of ['Start Date', 'End Date', 'Start Time', 'End Time']) {
      expect(within(panel).getByLabelText(name).parentElement?.parentElement).toHaveClass('time-range-fields');
    }
    rerender(<SearchBar {...props} value={{ query: '', time: initialTime }} />);
    expect(within(panel).getByRole('button', { name: '45 Minutes' }).parentElement).toHaveClass('time-range-presets');
    fireEvent.click(within(panel).getByRole('button', { name: 'Close' }));
    expect(screen.queryByRole('region', { name: 'Time range' })).not.toBeInTheDocument();
    expect(props.onSearch).not.toHaveBeenCalled();
  });
  it('does not search when an unchanged query loses focus', async () => {
    render(<Search />);
    await selectIndex();
    fireEvent.focus(screen.getByLabelText('search.typeSearch'));
    fireEvent.blur(screen.getByLabelText('search.typeSearch'));
    expect(mocks.search).toHaveBeenCalledTimes(1);
  });
  it('renders escaped text and highlights phrases, prefixes, operators and Unicode', () => {
    const { container, rerender } = render(<HighLight content={'<img src=x> New York Paris 中文 a.b'} queryString={'"New York" +City:Par* 中文 a.b'} />);
    expect(container.querySelector('img')).toBeNull();
    expect(Array.from(container.querySelectorAll('mark')).map(mark => mark.textContent)).toEqual(['New York', 'Par', '中', '文', 'a.b']);
    rerender(<HighLight content="Paris" queryString="" />);
    expect(container.querySelector('mark')).toBeNull();
  });
});
