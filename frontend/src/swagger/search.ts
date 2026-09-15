import type { ComponentType, ReactElement } from 'react';

// Structural subset of the Immutable.js collections Swagger UI passes to fn.opsFilter:
// OrderedMap<tag, Map{ tagDetails, operations: List<Map{ path, method, operation }> }>.
export interface Operation {
  get(key: 'method' | 'path'): string;
  getIn(path: string[]): unknown;
}
interface Operations {
  filter(fn: (op: Operation) => boolean): Operations;
  size: number;
}
interface TagObject {
  get(key: 'operations'): Operations;
  update(key: 'operations', fn: (ops: Operations) => Operations): TagObject;
}
export interface TaggedOperations {
  map(fn: (tag: TagObject) => TagObject): TaggedOperations;
  filter(fn: (tag: TagObject) => boolean): TaggedOperations;
}

export function operationText(op: Operation): string {
  const tags = op.getIn(['operation', 'tags']) as { join(sep: string): string } | undefined;
  return [
    op.get('method'),
    op.get('path'),
    op.getIn(['operation', 'summary']),
    op.getIn(['operation', 'description']),
    tags?.join(' '),
  ].filter(Boolean).join(' ').toLowerCase();
}

// Every whitespace-separated word of phrase must occur somewhere in the operation.
export function matches(op: Operation, phrase: string): boolean {
  const text = operationText(op);
  return phrase.toLowerCase().split(/\s+/).filter(Boolean).every(word => text.includes(word));
}

export function opsFilter(taggedOps: TaggedOperations, phrase: string): TaggedOperations {
  return taggedOps
    .map(tag => tag.update('operations', ops => ops.filter(op => matches(op, phrase))))
    .filter(tag => tag.get('operations').size > 0);
}

export interface FilterContainerProps {
  specSelectors: { loadingStatus(): string };
  layoutSelectors: { currentFilter(): string | boolean | null };
  layoutActions: { updateFilter(value: string): void };
  getComponent(name: string): ComponentType<{ className?: string; mobile?: number; children?: ReactElement }>;
}

// Replaces Swagger UI's FilterContainer, whose placeholder says "Filter by tag".
export function createFilterContainer(React: { createElement: typeof import('react').createElement }) {
  return function FilterContainer({ specSelectors, layoutSelectors, layoutActions, getComponent }: FilterContainerProps) {
    const filter = layoutSelectors.currentFilter();
    if (filter === false || filter === null) return null;
    const Col = getComponent('Col');
    return React.createElement('div', { className: 'filter-container' },
      React.createElement(Col, { className: 'filter wrapper', mobile: 12 },
        React.createElement('input', {
          className: 'operation-filter-input',
          type: 'search',
          placeholder: 'Search by method, path, summary or tag…',
          'aria-label': 'Search operations',
          value: typeof filter === 'string' ? filter : '',
          disabled: specSelectors.loadingStatus() === 'loading',
          onChange: (e: { target: { value: string } }) => layoutActions.updateFilter(e.target.value),
        })));
  };
}

export const searchPlugin = (system: { React: { createElement: typeof import('react').createElement } }) => ({
  fn: { opsFilter },
  components: { FilterContainer: createFilterContainer(system.React) },
});
