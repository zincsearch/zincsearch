import { useRef, useState } from 'react';
import type { SearchResponse } from './model';

export function bucketLabel(key: string | number, interval: string) {
  const date = new Date(key);
  if (!Number.isFinite(+date)) return String(key);
  const pad = (n: number) => String(n).padStart(2, '0');
  const day = `${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
  const clock = `${pad(date.getHours())}:${pad(date.getMinutes())}`;
  if (interval.includes('d')) return `${date.getFullYear()}-${day}`;
  if (interval.includes('h') || interval === '5m') return `${day} ${clock}`;
  return interval.includes('s') ? `${clock}:${pad(date.getSeconds())}` : clock;
}
export default function Histogram({ histogram, interval }: { histogram: NonNullable<SearchResponse['aggregations']>['histogram']; interval: string }) {
  const [zoom, setZoom] = useState(1);
  const chart = useRef<SVGSVGElement>(null);
  const [exportError, setExportError] = useState('');
  const buckets = histogram?.buckets || [];
  const maximum = Math.max(1, ...buckets.map(bucket => Number(bucket.doc_count)));
  const label = (key: string | number) => bucketLabel(key, histogram?.interval || interval);
  async function download(format: 'csv' | 'svg' | 'png') {
    setExportError('');
    try {
      let blob: Blob;
      if (format === 'csv') {
        const csv = ['Timestamp,Count', ...buckets.map(bucket => `${new Date(bucket.key).toISOString()},${Number(bucket.doc_count)}`)].join('\n');
        blob = new Blob([csv], { type: 'text/csv' });
      } else {
        const svg = chart.current!.cloneNode(true) as SVGSVGElement;
        svg.setAttribute('width', String(800 * zoom));
        svg.setAttribute('height', '170');
        svg.removeAttribute('style');
        blob = new Blob([new XMLSerializer().serializeToString(svg)], { type: 'image/svg+xml' });
        if (format === 'png') {
          const url = URL.createObjectURL(blob);
          try {
            const image = new Image();
            image.src = url;
            await image.decode();
            const canvas = document.createElement('canvas');
            canvas.width = 800 * zoom;
            canvas.height = 170;
            const context = canvas.getContext('2d');
            if (!context) throw new Error('Canvas unavailable');
            context.fillStyle = 'white';
            context.fillRect(0, 0, canvas.width, canvas.height);
            context.drawImage(image, 0, 0);
            blob = await new Promise<Blob>((resolve, reject) => canvas.toBlob(value => value ? resolve(value) : reject(new Error('PNG export failed')), 'image/png'));
          } finally { URL.revokeObjectURL(url); }
        }
      }
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `search-summary.${format}`;
      link.click();
      setTimeout(() => URL.revokeObjectURL(url), 0);
    } catch {
      setExportError('Unable to export histogram. Please try another format.');
    }
  }
  return <section aria-label="Search histogram">
    {exportError && <p className="error" role="alert">{exportError}</p>}
    {buckets.length ? <>
      <div className="toolbar"><button type="button" aria-label="Zoom in histogram" onClick={() => setZoom(Math.min(8, zoom * 2))}>+</button><button type="button" aria-label="Zoom out histogram" onClick={() => setZoom(Math.max(1, zoom / 2))}>−</button><button type="button" onClick={() => setZoom(1)}>Reset zoom</button>{(['csv', 'svg', 'png'] as const).map(format => <button type="button" key={format} onClick={() => { void download(format); }}>Download {format.toUpperCase()}</button>)}</div>
      <div className="overflow-x-auto"><svg ref={chart} role="img" aria-label="Document counts over time" viewBox={`0 0 ${800 * zoom} 170`} style={{ width: `${zoom * 100}%`, minWidth: 300, height: 170 }}>
        {buckets.map((bucket, index) => {
          const width = 760 * zoom / buckets.length;
          const height = Number(bucket.doc_count) / maximum * 125;
          return <g key={`${bucket.key}-${index}`}><rect x={30 + index * width} y={135 - height} width={width * 0.9} height={height} fill="#26A69A"><title>{label(bucket.key)}: {bucket.doc_count}</title></rect>{index % Math.max(1, Math.ceil(buckets.length / (6 * zoom))) === 0 && <text x={30 + index * width} y="155" fontSize="10" fill="currentColor">{label(bucket.key)}</text>}</g>;
        })}<text x="0" y="12" fontSize="10" fill="currentColor">{maximum}</text>
      </svg></div>
    </> : <p>No data available</p>}
  </section>;
}
