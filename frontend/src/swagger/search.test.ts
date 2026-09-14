import { createElement, type ReactElement } from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { fromJS, type Map, type List } from 'immutable';
import { createFilterContainer, matches, opsFilter, searchPlugin, type Operation, type TaggedOperations } from './search';

const taggedOps = fromJS({
  Index: {
    tagDetails: {},
    operations: [
      { path: '/api/index', method: 'get', operation: { summary: 'List indexes', tags: ['Index'] } },
      { path: '/api/index', method: 'put', operation: { summary: 'Create index', description: 'Creates a new index with mappings', tags: ['Index'] } },
    ],
  },
  Search: {
    tagDetails: {},
    operations: [
      { path: '/api/{index}/_search', method: 'post', operation: { summary: 'Search V1', tags: ['Search'] } },
      { path: '/es/{index}/_search', method: 'post', operation: { tags: ['Search'] } },
    ],
  },
}) as unknown as Map<string, Map<string, List<Map<string, unknown>>>>;

const paths = (result: unknown) => {
  const out: Record<string, string[]> = {};
  (result as Map<string, Map<string, List<Map<string, string>>>>).forEach((tag, name) => {
    out[name] = tag.get('operations')!.map(op => `${op.get('method')} ${op.get('path')}`).toArray();
  });
  return out;
};

describe('opsFilter', () => {
  it.each([
    ['path', 'es/', { Search: ['post /es/{index}/_search'] }],
    ['method, case-insensitive', 'PUT', { Index: ['put /api/index'] }],
    ['summary', 'list', { Index: ['get /api/index'] }],
    ['description', 'mappings', { Index: ['put /api/index'] }],
    ['tag', 'search', { Search: ['post /api/{index}/_search', 'post /es/{index}/_search'] }],
    ['all words must match', 'post api', { Search: ['post /api/{index}/_search'] }],
    ['no match', 'nothing-here', {}],
    ['blank keeps everything', '  ', { Index: ['get /api/index', 'put /api/index'], Search: ['post /api/{index}/_search', 'post /es/{index}/_search'] }],
  ])('%s', (_name, phrase, expected) => {
    expect(paths(opsFilter(taggedOps as unknown as TaggedOperations, phrase))).toEqual(expected);
  });

  it('preserves tag order', () => {
    expect(Object.keys(paths(opsFilter(taggedOps as unknown as TaggedOperations, 'index')))).toEqual(['Index', 'Search']);
  });
});

describe('matches', () => {
  it('ignores operations without summary/description/tags', () => {
    const op = fromJS({ path: '/healthz', method: 'get', operation: {} }) as unknown as Operation;
    expect(matches(op, 'health')).toBe(true);
    expect(matches(op, 'version')).toBe(false);
  });
});

describe('FilterContainer', () => {
  const Col = ({ children }: { children?: ReactElement }) => createElement('div', null, children);
  const props = (filter: string | boolean | null, status = 'success') => ({
    specSelectors: { loadingStatus: () => status },
    layoutSelectors: { currentFilter: () => filter },
    layoutActions: { updateFilter: vi.fn() },
    getComponent: () => Col,
  });

  it('is wired into the plugin', () => {
    const plugin = searchPlugin({ React: { createElement } });
    expect(plugin.fn.opsFilter).toBe(opsFilter);
    expect(typeof plugin.components.FilterContainer).toBe('function');
  });

  it('renders a search box and reports input', async () => {
    const FilterContainer = createFilterContainer({ createElement });
    const p = props(true);
    render(createElement(FilterContainer, p));
    const input = screen.getByRole('searchbox', { name: 'Search operations' });
    expect(input).toHaveValue('');
    await userEvent.type(input, 'x');
    expect(p.layoutActions.updateFilter).toHaveBeenCalledWith('x');
  });

  it('shows the current phrase and disables while loading', () => {
    const FilterContainer = createFilterContainer({ createElement });
    render(createElement(FilterContainer, props('index', 'loading')));
    const input = screen.getByRole('searchbox');
    expect(input).toHaveValue('index');
    expect(input).toBeDisabled();
  });

  it('renders nothing when filtering is off', () => {
    const FilterContainer = createFilterContainer({ createElement });
    const { container } = render(createElement(FilterContainer, props(false)));
    expect(container).toBeEmptyDOMElement();
  });
});
