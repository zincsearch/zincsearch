import { useTranslation } from '../../locales';

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
  return <details className="relative"><summary data-cy="syntax-guide-button">{t('search.syntaxGuide')}</summary>
    <div className="card absolute z-20 w-96 max-w-[90vw]">{examples.map(([label, query]) => <div className="mb-3" key={query}><p>{label}</p><code>{query}</code></div>)}</div>
  </details>;
}
