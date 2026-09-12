import { useEffect, useState } from 'react';
import about from '../services/about';
import { useTranslation } from '../locales';

export default function About() {
  const { t } = useTranslation();
  const [version, setVersion] = useState<Record<string, string>>({});
  const [error, setError] = useState('');
  useEffect(() => {
    let active = true;
    about.get().then(({ data }) => { if (active) setVersion(data); }).catch(error => {
      if (active) setError(error instanceof Error ? error.message : 'Unable to load version');
    });
    return () => { active = false; };
  }, []);
  return <section className="card space-y-6"><h1>Zinc Search</h1><p>{t('about.introduction')}</p>
    {error && <p className="error" role="alert">{error}</p>}
    <dl className="grid grid-cols-[auto_1fr] gap-x-8 gap-y-4">
      {Object.entries({ version: 'Version', build: 'Build', commit_hash: 'CommitHash', branch: 'Branch', build_date: 'BuildDate' }).map(([key, label]) => <div className="contents" key={key}><dt className="font-medium">{label}</dt><dd>{version[key] || '—'}</dd></div>)}
    </dl></section>;
}
