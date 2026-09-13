import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import Histogram, { bucketLabel } from './Histogram';

afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });
const histogram = { buckets: [{ key: new Date(2026, 1, 15, 12, 30, 45).getTime(), doc_count: 7 }] };
describe('Histogram', () => {
  it.each([['1s', '12:30:45'], ['1m', '12:30'], ['1h', '02-15 12:30'], ['1d', '2026-02-15']] as const)('formats %s buckets', (interval, label) => {
    expect(bucketLabel(histogram.buckets[0].key, interval)).toBe(label);
  });
  it('zooms and resets while exposing bucket counts', () => {
    render(<Histogram histogram={histogram} interval="1s" />);
    expect(screen.getByText('12:30:45: 7')).toBeInTheDocument();
    fireEvent.click(screen.getByLabelText('Zoom in histogram'));
    expect(screen.getByRole('img')).toHaveStyle({ width: '200%' });
    fireEvent.click(screen.getByRole('button', { name: 'Reset zoom' }));
    expect(screen.getByRole('img')).toHaveStyle({ width: '100%' });
  });
  it.each(['CSV', 'SVG', 'PNG'])('exports %s', async format => {
    const create = vi.fn(() => 'blob:histogram');
    vi.stubGlobal('URL', { createObjectURL: create, revokeObjectURL: vi.fn() });
    vi.stubGlobal('Image', class { src = ''; decode = vi.fn().mockResolvedValue(undefined); });
    vi.spyOn(HTMLCanvasElement.prototype, 'getContext').mockReturnValue({ fillStyle: '', fillRect: vi.fn(), drawImage: vi.fn() } as unknown as CanvasRenderingContext2D);
    vi.spyOn(HTMLCanvasElement.prototype, 'toBlob').mockImplementation(callback => callback(new Blob(['png'], { type: 'image/png' })));
    const click = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {});
    render(<Histogram histogram={histogram} interval="1s" />);
    fireEvent.click(screen.getByRole('button', { name: `Download ${format}` }));
    await waitFor(() => expect(click).toHaveBeenCalledOnce());
    expect((click.mock.instances[0] as HTMLAnchorElement).download).toBe(`search-summary.${format.toLowerCase()}`);
    expect(create.mock.calls.length).toBeGreaterThan(0);
    await waitFor(() => expect(URL.revokeObjectURL).toHaveBeenCalledTimes(format === 'PNG' ? 2 : 1));
  });
  it('reports failed export without breaking the chart', async () => {
    vi.stubGlobal('URL', { createObjectURL: () => { throw new Error('unavailable'); } });
    render(<Histogram histogram={histogram} interval="1s" />);
    fireEvent.click(screen.getByRole('button', { name: 'Download SVG' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Unable to export histogram');
    expect(screen.getByRole('img')).toBeInTheDocument();
  });
});
