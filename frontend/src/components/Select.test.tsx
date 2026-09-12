import { act, fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, it, vi } from 'vitest';
import { useState } from 'react';
import Select from './Select';
import { Modal } from './management/Common';

function Example({ disabled = false }: { disabled?: boolean }) {
  const [value, setValue] = useState('');
  return <><label>Size<Select value={value} onValueChange={setValue} disabled={disabled}>
    <option value="">Choose size</option>
    <option value={5}>Five</option>
    <option value={10} disabled>Ten</option>
    <option value={0}>All</option>
    <option>Another</option>
  </Select></label><button type="button">Next</button></>;
}
afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

it('uses a centered, fixed-size chevron independent of font metrics', () => {
  render(<Example />);
  const trigger = screen.getByRole('combobox', { name: 'Size' });
  const arrow = trigger.querySelector('svg.ui-select-arrow');
  expect(arrow).toHaveAttribute('width', '12');
  expect(arrow).toHaveAttribute('height', '12');
  expect(arrow).toHaveAttribute('viewBox', '0 0 12 12');
  expect(arrow).toHaveAttribute('aria-hidden', 'true');
  expect(arrow).toHaveAttribute('focusable', 'false');
  expect(arrow?.querySelector('path')).toHaveAttribute('d', 'm3 4.5 3 3 3-3');
  expect(trigger).not.toHaveTextContent('⌄');
});

it('portals styled options, preserves numeric/empty/implicit values, and returns focus without submitting', async () => {
  const user = userEvent.setup();
  const submit = vi.fn(event => event.preventDefault());
  const { container } = render(<form onSubmit={submit}><Example /></form>);
  const trigger = screen.getByRole('combobox', { name: 'Size' });
  expect(container.querySelector('select')).toBeNull();
  await user.click(trigger);
  const list = screen.getByRole('listbox', { name: 'Size' });
  expect(list.parentElement).toBe(document.body);
  expect(trigger).toHaveAttribute('aria-controls', list.id);
  expect(screen.getByRole('option', { name: 'Choose size' })).toHaveAttribute('aria-selected', 'true');
  await user.click(screen.getByRole('option', { name: 'All' }));
  expect(trigger).toHaveValue('0');
  expect(trigger).toHaveAccessibleDescription('All');
  expect(trigger).toHaveFocus();
  expect(screen.queryByRole('listbox')).toBeNull();
  await user.click(trigger);
  await user.click(screen.getByRole('option', { name: 'Another' }));
  expect(trigger).toHaveValue('Another');
  await user.click(trigger);
  await user.click(screen.getByRole('option', { name: 'Choose size' }));
  expect(trigger).toHaveValue('');
  expect(submit).not.toHaveBeenCalled();
});

it('supports arrows, Home/End, typeahead, Enter/Space, Escape cancellation and Tab', async () => {
  const user = userEvent.setup();
  render(<Example />);
  const trigger = screen.getByRole('combobox', { name: 'Size' });
  await user.tab();
  await user.keyboard('{ArrowDown}{ArrowDown}{ArrowDown}');
  expect(trigger).toHaveAttribute('aria-activedescendant', screen.getByRole('option', { name: 'All' }).id);
  expect(trigger).toHaveValue('');
  await user.keyboard('{Escape}');
  expect(trigger).toHaveFocus();
  expect(trigger).toHaveValue('');
  await user.keyboard(' {End}{Enter}');
  expect(trigger).toHaveValue('Another');
  await user.keyboard('{ArrowUp}{Home}{ArrowDown} ');
  expect(trigger).toHaveValue('5');
  await user.keyboard('a{Enter}');
  expect(trigger).toHaveValue('0');
  await user.click(trigger);
  await user.keyboard('an{Enter}');
  expect(trigger).toHaveValue('Another');
  await user.click(trigger);
  await user.tab();
  expect(screen.getByRole('button', { name: 'Next' })).toHaveFocus();
  expect(screen.queryByRole('listbox')).toBeNull();
});

it('ignores disabled controls/options and closes on outside interaction or becoming disabled', async () => {
  const user = userEvent.setup();
  const { rerender } = render(<Example />);
  const trigger = screen.getByRole('combobox');
  await user.click(trigger);
  await user.click(screen.getByRole('option', { name: 'Ten' }));
  expect(trigger).toHaveValue('');
  expect(screen.getByRole('listbox')).toBeInTheDocument();
  await user.click(screen.getByRole('button', { name: 'Next' }));
  expect(screen.queryByRole('listbox')).toBeNull();
  await user.click(trigger);
  rerender(<Example disabled />);
  expect(screen.queryByRole('listbox')).toBeNull();
  await user.click(trigger);
  expect(screen.queryByRole('listbox')).toBeNull();
  rerender(<Example />);
  expect(screen.queryByRole('listbox')).toBeNull();
});

it('handles empty choices, unchanged values and asynchronously replaced options', async () => {
  const user = userEvent.setup();
  const change = vi.fn();
  const { rerender } = render(<Select aria-label="Async" value="saved" onValueChange={change}>{[]}</Select>);
  const trigger = screen.getByRole('combobox');
  await user.click(trigger);
  await user.keyboard('{End}{Enter}');
  expect(change).not.toHaveBeenCalled();
  rerender(<Select aria-label="Async" value="saved" onValueChange={change}><option value="saved">Saved</option></Select>);
  await user.click(screen.getByRole('option', { name: 'Saved' }));
  expect(change).not.toHaveBeenCalled();
  expect(trigger).toHaveTextContent('Saved');
});

it('keeps modal focus and lets Escape close only the dropdown first', async () => {
  const user = userEvent.setup();
  const close = vi.fn();
  render(<Modal title="Edit" onClose={close}><Example /></Modal>);
  const trigger = screen.getByRole('combobox');
  await user.click(trigger);
  await user.keyboard('{Escape}');
  expect(close).not.toHaveBeenCalled();
  expect(trigger).toHaveFocus();
  await user.keyboard('{Escape}');
  expect(close).toHaveBeenCalledOnce();
});

it.each([
  { top: 4, left: 290, expectedTop: 48, expectedLeft: 152, maxHeight: 144 },
  { top: 156, left: -20, expectedTop: 8, expectedLeft: 8, maxHeight: 144 },
])('clamps viewport edges and flips in a short viewport: $top/$left', ({ top, left, expectedTop, expectedLeft, maxHeight }) => {
  vi.stubGlobal('innerWidth', 320);
  vi.stubGlobal('innerHeight', 200);
  vi.spyOn(HTMLElement.prototype, 'scrollHeight', 'get').mockReturnValue(400);
  const { unmount } = render(<Example />);
  const trigger = screen.getByRole('combobox');
  vi.spyOn(trigger, 'getBoundingClientRect').mockReturnValue({ top, bottom: top + 40, left, right: left + 120, width: 120, height: 40, x: left, y: top, toJSON() {} });
  fireEvent.click(trigger);
  expect(screen.getByRole('listbox')).toHaveStyle({ top: `${expectedTop}px`, left: `${expectedLeft}px`, maxHeight: `${maxHeight}px`, width: '160px' });
  vi.stubGlobal('innerWidth', 140);
  fireEvent.resize(window);
  expect(screen.getByRole('listbox')).toHaveStyle({ left: '8px', width: '124px' });
  vi.stubGlobal('innerHeight', 600);
  fireEvent.scroll(window);
  expect(screen.getByRole('listbox')).toHaveStyle({ top: `${top + 44}px`, maxHeight: '288px' });
  unmount();
  expect(screen.queryByRole('listbox')).toBeNull();
});

it.each([-400, 1000])('keeps the menu bounded when scrolling moves the trigger offscreen (%s)', (top) => {
  vi.stubGlobal('innerWidth', 320);
  vi.stubGlobal('innerHeight', 200);
  vi.spyOn(HTMLElement.prototype, 'scrollHeight', 'get').mockReturnValue(400);
  render(<Example />);
  const trigger = screen.getByRole('combobox');
  vi.spyOn(trigger, 'getBoundingClientRect').mockReturnValue({ top, bottom: top + 40, left: 0, right: 120, width: 120, height: 40, x: 0, y: top, toJSON() {} });
  fireEvent.click(trigger);
  expect(screen.getByRole('listbox')).toHaveStyle({ top: '8px', maxHeight: '184px' });
});

it('scrolls active options within the bounded menu instead of scrolling the page', async () => {
  const user = userEvent.setup();
  render(<Example />);
  await user.click(screen.getByRole('combobox'));
  const list = screen.getByRole('listbox');
  const last = screen.getByRole('option', { name: 'Another' });
  Object.defineProperties(list, { clientHeight: { value: 80 } });
  Object.defineProperties(last, { offsetTop: { value: 160 }, offsetHeight: { value: 40 } });
  await user.keyboard('{End}');
  expect(list.scrollTop).toBe(120);
  await user.keyboard('{Home}');
  expect(list.scrollTop).toBe(0);
});

it('observes size changes and cleans up the observer', () => {
  let resize!: ResizeObserverCallback;
  const disconnect = vi.fn();
  vi.stubGlobal('ResizeObserver', class {
    constructor(callback: ResizeObserverCallback) { resize = callback; }
    observe() {}
    disconnect = disconnect;
  });
  const { unmount } = render(<Example />);
  fireEvent.click(screen.getByRole('combobox'));
  act(() => resize([], {} as ResizeObserver));
  unmount();
  expect(disconnect).toHaveBeenCalledOnce();
});
