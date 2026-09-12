import { Fragment } from 'react';

export function getKeywords(query: string): string[] {
  const terms = query.match(/(?:[^\s"]+|"[^"]*")+/g) || [];
  return terms.flatMap(term => {
    const value = term.replace(/^[+-]/, '').replace(/^[^:"]+:/, '').replace(/^[*"]+|[*"]+$/g, '');
    return value.match(/[\u0001-\u007e\uff60-\uff9f]+|[^\u0001-\u007e\uff60-\uff9f]/gu) || [];
  }).filter(term => term.trim() !== '');
}
export default function HighLight({ content, queryString = '' }: { content: string; queryString?: string }) {
  const keywords = [...new Set(getKeywords(queryString))].sort((a, b) => b.length - a.length);
  if (!keywords.length) return <>{content}</>;
  const pattern = new RegExp(`(${keywords.map(word => word.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')).join('|')})`, 'g');
  return <>{content.split(pattern).map((text, index) => <Fragment key={index}>{index % 2 ? <mark className="highlight" style={{ backgroundColor: 'rgb(255, 213, 0)', color: '#111' }}>{text}</mark> : text}</Fragment>)}</>;
}
