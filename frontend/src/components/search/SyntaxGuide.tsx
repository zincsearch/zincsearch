import { useRef, useState } from 'react';
import { useTranslation } from '../../locales';
import useClickOutside from '../../utils/useClickOutside';

const examples = [
  ['Search for Gold', 'Gold'],
  ['Search for City with Paris', 'City:Paris'],
  ['Search for City Not Paris', '-City:Paris'],
  ['Search for City with Paris or Gold', 'City:Paris Gold'],
  ['Search for City with Paris and Gold', '+City:Paris +Gold'],
  ['Search for City with Paris but not Gold', '+City:Paris -Gold'],
  ['Search for Medal Gold and Year>2000', '+Medal:Gold +Year:>2000'],
  ['Search text starting with par . e.g. Paris, Part, Paramater', 'par*'],
];
export default function SyntaxGuide() {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const root = useRef<HTMLDetailsElement>(null);
  useClickOutside(root, open, () => setOpen(false));
  return <details className="search-popover-control" ref={root} open={open}><summary data-cy="syntax-guide-button" onClick={event => { event.preventDefault(); setOpen(!open); }}>{t('search.syntaxGuide')}</summary>
    <div hidden={!open} className="card search-popover syntax-guide-panel" role="region" aria-label={t('search.syntaxGuide')}>
      {examples.map(([label, query]) => <div className="mb-3" key={query}><p>{label}</p><code>{query}</code></div>)}
      <button type="button" onClick={() => setOpen(false)}>Close</button>
    </div>
  </details>;
}
