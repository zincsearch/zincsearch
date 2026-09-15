import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useRef } from 'react';
import { expect, it, vi } from 'vitest';
import Select from '../components/Select';
import useClickOutside from './useClickOutside';

function Example({ active = true, onOutside }: { active?: boolean; onOutside: () => void }) {
  const root = useRef<HTMLDivElement>(null);
  useClickOutside(root, active, onOutside);
  return <>
    <div ref={root}>
      <Select aria-label="Inside" value="one" onValueChange={() => {}}>
        <option value="one">Inside option</option>
      </Select>
    </div>
    <Select aria-label="Outside" value="one" onValueChange={() => {}}>
      <option value="one">Outside option</option>
    </Select>
  </>;
}

it('treats only the contained Select portal as inside', async () => {
  const user = userEvent.setup();
  const onOutside = vi.fn();
  const { container } = render(<Example onOutside={onOutside} />);
  await user.click(screen.getByRole('combobox', { name: 'Inside' }));
  expect(container).not.toContainElement(screen.getByRole('listbox'));
  await user.click(screen.getByText('Inside option', { selector: '[role="option"] span' }));
  expect(onOutside).not.toHaveBeenCalled();

  await user.tab();
  expect(screen.getByRole('combobox', { name: 'Outside' })).toHaveFocus();
  await user.keyboard('{ArrowDown}');
  expect(onOutside).not.toHaveBeenCalled();
  await user.click(screen.getByText('Outside option', { selector: '[role="option"] span' }));
  expect(onOutside).toHaveBeenCalledTimes(1);
});

it('closes for ordinary outside pointers and removes inactive or unmounted listeners', () => {
  const onOutside = vi.fn();
  const { rerender, unmount } = render(<Example onOutside={onOutside} />);
  fireEvent.pointerDown(document.body);
  expect(onOutside).toHaveBeenCalledTimes(1);
  rerender(<Example active={false} onOutside={onOutside} />);
  fireEvent.pointerDown(document.body);
  expect(onOutside).toHaveBeenCalledTimes(1);
  rerender(<Example onOutside={onOutside} />);
  unmount();
  fireEvent.pointerDown(document.body);
  expect(onOutside).toHaveBeenCalledTimes(1);
});
