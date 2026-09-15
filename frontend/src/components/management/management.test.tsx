import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import indexService from '../../services/index';
import templateService from '../../services/template';
import userService from '../../services/user';
import roleService from '../../services/role';
import permissionService from '../../services/permission';
import Index from '../../views/Index';
import Template from '../../views/Template';
import User from '../../views/User';
import Role from '../../views/Role';
import SchemaEditor from './SchemaEditor';
import AccountEditor from './AccountEditor';

vi.mock('../../locales', async (importOriginal) => ({
  ...await importOriginal<typeof import('../../locales')>(),
  useTranslation: () => ({ t: (key: string) => key }),
}));
vi.mock(
  '../../services/index',
  () => ({ default: { list: vi.fn(), update: vi.fn(), delete: vi.fn() } }),
);
vi.mock(
  '../../services/template',
  () => ({ default: { list: vi.fn(), update: vi.fn(), delete: vi.fn() } }),
);
vi.mock(
  '../../services/user',
  () => ({ default: { list: vi.fn(), update: vi.fn(), delete: vi.fn() } }),
);
vi.mock(
  '../../services/role',
  () => ({ default: { list: vi.fn(), update: vi.fn(), delete: vi.fn() } }),
);
vi.mock('../../services/permission', () => ({ default: { list: vi.fn() } }));
const response = (data: unknown) => ({ data } as Awaited<ReturnType<typeof userService.list>>);
const change = (label: string, value: string) => {
  const control = screen.getByLabelText(label);
  if (control.getAttribute('role') === 'combobox') {
    fireEvent.click(control);
    const option = screen.getAllByRole('option').find(option => option.dataset.value === value);
    expect(option).toBeDefined();
    fireEvent.click(option!);
  } else fireEvent.change(control, { target: { value } });
};
const click = (name: string) => fireEvent.click(screen.getByRole('button', { name }));
afterEach(cleanup);
beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(roleService.list).mockResolvedValue(
    response([{ _id: 'reader', name: 'Reader', permission: ['search'] }]),
  );
  vi.mocked(permissionService.list).mockResolvedValue(response(['search', 'index']));
  vi.mocked(userService.list).mockResolvedValue(
    response([{ _id: 'alice', name: 'Alice', role: 'reader' }]),
  );
  vi.mocked(templateService.list).mockResolvedValue(
    response([{
      name: 'logs',
      index_template: {
        index_patterns: ['logs-*'],
        priority: 0,
        template: { settings: { analysis: {} }, mappings: {} },
      },
    }]),
  );
  vi.mocked(indexService.list).mockResolvedValue(
    response({
      page: { total: 21 },
      list: [{
        name: 'logs',
        shard_num: 2,
        storage_type: 'disk',
        stats: { doc_num: 9, storage_size: 2048, wal_size: 3 },
        settings: { analysis: {} },
        mappings: {},
      }],
    }),
  );
  for (const service of [indexService, templateService, userService, roleService]) {
    vi.mocked(service.update).mockResolvedValue(response({}));
    vi.mocked(service.delete).mockResolvedValue(response({}));
  }
});

it('uses independent, labeled password toggles in the user editor', async () => {
  render(<AccountEditor kind='user' onClose={vi.fn()} onUpdated={vi.fn()} />);
  await waitFor(() => expect(screen.getByRole('button', { name: 'Save User' })).toBeEnabled());
  const fields = ['user.password', 'user.repassword'].map(label => screen.getByLabelText(label));
  for (const field of fields) {
    expect(field.closest('label')).toBeNull();
    expect(field).toHaveAttribute('autoComplete', 'new-password');
    fireEvent.change(field, { target: { value: 'example1' } });
    const toggle = screen.getAllByRole('button', { name: 'passwordInput.show' }).find(button => button.getAttribute('aria-controls') === field.id)!;
    expect(toggle.closest('label')).toBeNull();
    fireEvent.click(toggle);
    expect(field).toHaveAttribute('type', 'text');
    expect(field).toHaveValue('example1');
    for (const other of fields.filter(input => input !== field)) expect(other).toHaveAttribute('type', 'password');
    fireEvent.click(screen.getByRole('button', { name: 'passwordInput.hide' }));
    expect(field).toHaveAttribute('type', 'password');
  }
  expect(userService.update).not.toHaveBeenCalled();
});

describe('schema editor parity', () => {
  it.each(['index', 'template'] as const)('validates %s JSON immediately and recovers across steps and formats', (kind) => {
    render(<SchemaEditor kind={kind} onClose={vi.fn()} onUpdated={vi.fn()} />);
    change(kind === 'index' ? 'Index Name' : 'Template Name', 'logs');
    if (kind === 'template') change('Index Patterns', 'logs-*');
    click('Continue');
    const settingsLabel = kind === 'index' ? 'Settings JSON' : 'Index settings JSON';
    change(settingsLabel, '{');
    expect(screen.getByRole('alert')).toHaveTextContent('Invalid JSON format.');
    expect(screen.getByLabelText(settingsLabel)).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByRole('button', { name: 'Continue' })).toBeDisabled();
    click('Back');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    click('Continue');
    expect(screen.getByRole('alert')).toHaveTextContent('Invalid JSON format.');
    for (const invalid of ['[]', 'null', '42']) {
      change(settingsLabel, invalid);
      expect(screen.getByRole('alert')).toHaveTextContent('JSON must be an object.');
    }
    change(settingsLabel, '{}');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.getByLabelText(settingsLabel)).toHaveAttribute('aria-invalid', 'false');
    click('Continue');
    change('Mappings JSON', '{"properties":');
    expect(screen.getByRole('alert')).toHaveTextContent('Invalid JSON format.');
    expect(screen.getByRole('button', { name: 'Continue' })).toBeDisabled();
    change('Mappings format', 'sql');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    change('CREATE TABLE SQL', 'CREATE TABLE logs (message TEXT);');
    click('Continue');
    expect(screen.getByLabelText('Review JSON')).toHaveTextContent('CREATE TABLE logs');
    click('Back');
    change('Mappings format', 'json');
    expect(screen.getByRole('alert')).toHaveTextContent('Invalid JSON format.');
    change('Mappings JSON', '{"properties":{"message":{"type":"text"}}}');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    click('Continue');
    expect(screen.getByLabelText('Review JSON')).toHaveTextContent('"message"');
    expect(indexService.update).not.toHaveBeenCalled();
    expect(templateService.update).not.toHaveBeenCalled();
  });
  it('checks geo_point and vector mapping properties before review', () => {
    render(<SchemaEditor kind='index' onClose={vi.fn()} onUpdated={vi.fn()} />);
    change('Index Name', 'places');
    click('Continue');
    click('Continue');
    change('Mappings JSON', '{"properties":{"embedding":{"type":"vector","dims":0}}}');
    expect(screen.getByRole('alert')).toHaveTextContent('dims must be a positive integer');
    expect(screen.getByRole('button', { name: 'Continue' })).toBeDisabled();
    change('Mappings JSON', '{"properties":{"location":{"type":"geo_point","dims":2}}}');
    expect(screen.getByRole('alert')).toHaveTextContent('only valid for vector');
    change('Mappings JSON', '{"properties":{"location":{"type":"geo_pt"}}}');
    expect(screen.getByRole('alert')).toHaveTextContent('unsupported type "geo_pt"');
    change('Mappings JSON', '{"properties":{"location":{"type":"geo_point"},"embedding":{"type":"vector","dims":3}}}');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    click('Continue');
    expect(screen.getByLabelText('Review JSON')).toHaveTextContent('"geo_point"');
    expect(screen.getByLabelText('Review JSON')).toHaveTextContent('"dims": 3');
  });
});

describe('index management', () => {
  it('converts the styled pagination choice to a numeric server page size', async () => {
    render(<Index />);
    await screen.findByText('logs');
    change('Rows per page', '50');
    await waitFor(() => expect(indexService.list).toHaveBeenLastCalledWith(1, 50, 'name', false, ''));
  });
  it('shows real index totals and explicitly page-scoped resource metrics', async () => {
    render(<Index />);
    const overview = screen.getByLabelText('Index overview');
    expect(overview).toHaveAttribute('aria-busy', 'true');
    expect(within(overview).getAllByText('—')).toHaveLength(4);
    await screen.findByText('logs');
    expect(overview).toHaveAttribute('aria-busy', 'false');
    expect(within(overview).getByText('21')).toBeInTheDocument();
    expect(within(overview).getByText('9')).toBeInTheDocument();
    expect(within(overview).getByText('2.00 KB')).toBeInTheDocument();
    expect(within(overview).getByText('2')).toBeInTheDocument();
    expect(within(overview).getAllByText('· this page')).toHaveLength(3);
    vi.mocked(indexService.list).mockRejectedValueOnce(new Error('offline'));
    click('Refresh');
    await screen.findByRole('alert');
    expect(within(overview).getAllByText('—')).toHaveLength(4);
  });
  it('uses server pagination, sorting and filtering, previews metadata, and confirms bulk deletion', async () => {
    render(<Index />);
    await screen.findByText('logs');
    expect(indexService.list).toHaveBeenLastCalledWith(1, 20, 'name', false, '');
    click('Next');
    await waitFor(() =>
      expect(indexService.list).toHaveBeenLastCalledWith(2, 20, 'name', false, '')
    );
    click('DOC_NUM');
    await waitFor(() =>
      expect(indexService.list).toHaveBeenLastCalledWith(1, 20, 'doc_num', false, '')
    );
    change('index.search', 'logs');
    await waitFor(() =>
      expect(indexService.list).toHaveBeenLastCalledWith(1, 20, 'doc_num', false, 'logs')
    );
    click('Preview logs');
    expect(within(screen.getByRole('dialog')).getByText('WAL Entries')).toBeTruthy();
    fireEvent.click(screen.getByRole('tab', { name: 'Settings' }));
    expect(screen.getByRole('tabpanel').textContent).toContain('analysis');
    click('Close');
    fireEvent.click(screen.getByLabelText('Select logs'));
    click('index.delete');
    expect(indexService.delete).not.toHaveBeenCalled();
    click('Delete');
    await waitFor(() => expect(indexService.delete).toHaveBeenCalledWith('logs'));
  });
  it('validates JSON before creating an index and submits parsed objects', async () => {
    const updated = vi.fn();
    render(<SchemaEditor kind='index' onClose={vi.fn()} onUpdated={updated} />);
    click('Continue');
    expect(screen.getByRole('alert').textContent).toContain('name is required');
    change('Index Name', 'events');
    change('Shard num (optional)', '2');
    click('Continue');
    change('Settings JSON', '{broken');
    click('Continue');
    expect(screen.getByRole('alert').textContent).toContain('Invalid JSON');
    expect(indexService.update).not.toHaveBeenCalled();
    change('Settings JSON', '{"analysis":{}}');
    click('Continue');
    change('Mappings JSON', '{"properties":{"message":{"type":"text"}}}');
    click('Continue');
    click('Save Index');
    await waitFor(() => expect(updated).toHaveBeenCalled());
    expect(indexService.update).toHaveBeenCalledWith({
      name: 'events',
      storage_type: 'disk',
      shard_num: 2,
      settings: { analysis: {} },
      mappings: { properties: { message: { type: 'text' } } },
    });
  });
  it('offers disk, s3, minio, gcs and oss storage and submits the selection', async () => {
    const updated = vi.fn();
    render(<SchemaEditor kind='index' onClose={vi.fn()} onUpdated={updated} />);
    const select = screen.getByLabelText('Storage Type');
    expect(select).toHaveValue('disk');
    fireEvent.click(select);
    expect(screen.getAllByRole('option').map((o) => o.dataset.value)).toEqual(['disk', 's3', 'minio', 'gcs', 'oss']);
    fireEvent.click(select);
    expect(screen.queryByText(/ZINC_S3_BUCKET/)).toBeNull();
    change('Index Name', 'events');
    // the hint renders inside the label, so later picks go through the captured select
    const pick = (value: string) => {
      fireEvent.click(select);
      fireEvent.click(screen.getAllByRole('option').find((o) => o.dataset.value === value)!);
    };
    pick('gcs');
    expect(screen.getByText(/ZINC_GCS_BUCKET/)).toBeTruthy();
    pick('oss');
    expect(screen.getByText(/ZINC_OSS_BUCKET/)).toBeTruthy();
    pick('minio');
    expect(screen.getByText(/ZINC_S3_BUCKET/)).toBeTruthy();
    click('Continue');
    click('Continue');
    click('Continue');
    expect(JSON.parse(screen.getByLabelText('Review JSON').textContent!).storage_type).toBe('minio');
    click('Save Index');
    await waitFor(() => expect(updated).toHaveBeenCalled());
    expect(indexService.update).toHaveBeenCalledWith({
      name: 'events',
      storage_type: 'minio',
      settings: {},
      mappings: {},
    });
  });
});

describe('template management', () => {
  it('edits immutable names, patterns, priority and schema without mutating list data', async () => {
    render(<Template />);
    await screen.findByText('logs');
    click('Edit logs');
    expect((screen.getByLabelText('Template Name') as HTMLInputElement).disabled).toBe(true);
    change('Index Patterns', 'logs-*, events-*, logs-*');
    change('Priority (optional)', '0');
    click('Continue');
    change('Index settings JSON', '{"number_of_shards":2}');
    click('Continue');
    click('Continue');
    click('Save Template');
    await waitFor(() =>
      expect(templateService.update).toHaveBeenCalledWith({
        name: 'logs',
        index_patterns: ['logs-*', 'events-*'],
        priority: 0,
        template: { settings: { number_of_shards: 2 }, mappings: {} },
      })
    );
  });
  it('filters by name, previews settings and cancels deletion', async () => {
    render(<Template />);
    await screen.findByText('logs');
    change('template.search', 'missing');
    expect(screen.queryByText('logs')).toBeNull();
    change('template.search', 'LOG');
    click('logs');
    fireEvent.click(screen.getByRole('tab', { name: 'Settings' }));
    expect(screen.getByRole('tabpanel').textContent).toContain('analysis');
    click('Close');
    click('Delete logs');
    click('Cancel');
    expect(templateService.delete).not.toHaveBeenCalled();
  });
});

describe.each(['index', 'template'] as const)('%s mappings formats', (kind) => {
  const sql = 'CREATE TABLE events (message TEXT, count INT);';
  const json = '{"properties":{"message":{"type":"text"}}}';
  const settings = { analysis: {} };
  const service = kind === 'index' ? indexService : templateService;
  const save = kind === 'index' ? 'Save Index' : 'Save Template';

  function start() {
    const updated = vi.fn();
    render(<SchemaEditor kind={kind} onClose={vi.fn()} onUpdated={updated} />);
    change(kind === 'index' ? 'Index Name' : 'Template Name', 'events');
    if (kind === 'template') change('Index Patterns', 'events-*');
    click('Continue');
    change(kind === 'index' ? 'Settings JSON' : 'Index settings JSON', JSON.stringify(settings));
    click('Continue');
    return updated;
  }

  function expected(mappings: Record<string, unknown>) {
    return kind === 'index'
      ? { name: 'events', storage_type: 'disk', settings, mappings }
      : { name: 'events', index_patterns: ['events-*'], template: { settings, mappings } };
  }

  it('preserves both drafts across switching and submits only SQL with JSON settings', async () => {
    const updated = start();
    expect(screen.getByLabelText('Mappings format')).toHaveValue('json');
    change('Mappings JSON', json);
    change('Mappings format', 'sql');
    change('CREATE TABLE SQL', sql);
    change('Mappings format', 'json');
    expect(screen.getByLabelText('Mappings JSON')).toHaveValue(json);
    change('Mappings format', 'sql');
    expect(screen.getByLabelText('CREATE TABLE SQL')).toHaveValue(sql);
    click('Back');
    expect(screen.getByLabelText(kind === 'index' ? 'Settings JSON' : 'Index settings JSON'))
      .toHaveValue(JSON.stringify(settings));
    click('Continue');
    expect(screen.getByLabelText('CREATE TABLE SQL')).toHaveValue(sql);
    click('Continue');
    expect(JSON.parse(screen.getByLabelText('Review JSON').textContent!)).toEqual(expected({ sql }));
    click(save);
    await waitFor(() => expect(updated).toHaveBeenCalledOnce());
    expect(service.update).toHaveBeenCalledWith(expected({ sql }));
  });

  it('submits JSON properties unchanged after switching from SQL', async () => {
    start();
    change('Mappings JSON', json);
    change('Mappings format', 'sql');
    change('CREATE TABLE SQL', sql);
    change('Mappings format', 'json');
    click('Continue');
    click(save);
    await waitFor(() => expect(service.update).toHaveBeenCalledWith(expected(JSON.parse(json))));
  });

  it('rejects blank SQL locally, including the JSON object alternative', () => {
    start();
    change('Mappings format', 'sql');
    for (const blank of ['', ' \n\t ']) {
      change('CREATE TABLE SQL', blank);
      click('Continue');
      expect(screen.getByRole('alert')).toHaveTextContent('CREATE TABLE SQL is required');
    }
    change('Mappings format', 'json');
    change('Mappings JSON', '{"sql":"  "}');
    click('Continue');
    expect(screen.getByRole('alert')).toHaveTextContent('CREATE TABLE SQL is required');
    expect(service.update).not.toHaveBeenCalled();
  });

  it('accepts the SQL object in JSON mode and displays server errors with retry', async () => {
    const updated = start();
    vi.mocked(service.update).mockRejectedValueOnce({
      isAxiosError: true,
      response: { data: { error: 'Unsupported SQL column type' }, status: 400 },
    });
    change('Mappings JSON', JSON.stringify({ sql }));
    click('Continue');
    click(save);
    expect(await screen.findByRole('alert')).toHaveTextContent('Unsupported SQL column type');
    expect(updated).not.toHaveBeenCalled();
    expect(screen.getByRole('button', { name: save })).toBeEnabled();
    expect(service.update).toHaveBeenCalledWith(expected({ sql }));
    click(save);
    await waitFor(() => expect(updated).toHaveBeenCalledOnce());
  });

  it('keeps SQL input available after a network failure', async () => {
    start();
    vi.mocked(service.update).mockRejectedValueOnce(new Error('Network Error'));
    change('Mappings format', 'sql');
    change('CREATE TABLE SQL', sql);
    click('Continue');
    click(save);
    expect(await screen.findByRole('alert')).toHaveTextContent(`Unable to save ${kind}`);
    click('Back');
    expect(screen.getByLabelText('CREATE TABLE SQL')).toHaveValue(sql);
  });
});

it('initializes SQL mode when editing a template with SQL mappings', async () => {
  const sql = 'CREATE TABLE events (message TEXT);';
  render(<SchemaEditor
    kind='template'
    value={{ name: 'events', index_template: { index_patterns: ['events-*'], template: { mappings: { sql } } } }}
    onClose={vi.fn()}
    onUpdated={vi.fn()}
  />);
  click('Continue');
  click('Continue');
  expect(screen.getByLabelText('Mappings format')).toHaveValue('sql');
  expect(screen.getByLabelText('CREATE TABLE SQL')).toHaveValue(sql);
  click('Continue');
  click('Save Template');
  await waitFor(() => expect(templateService.update).toHaveBeenCalledWith({
    name: 'events', index_patterns: ['events-*'], template: { settings: {}, mappings: { sql } },
  }));
});

describe('users and roles', () => {
  it('enforces password rules and confirmation before saving a new user', async () => {
    const updated = vi.fn();
    render(<AccountEditor kind='user' onClose={vi.fn()} onUpdated={updated} />);
    await waitFor(() => expect(screen.getByRole('combobox', { name: 'user.role' })).toBeEnabled());
    change('user.id', 'bobby');
    change('user.name', 'Bobby');
    change('user.role', 'reader');
    change('user.password', 'abcdefgh');
    change('user.repassword', 'abcdefgh');
    click('Save User');
    expect(screen.getByRole('alert').textContent).toContain('digit');
    change('user.password', 'abcdefg1');
    click('Save User');
    expect(screen.getByRole('alert').textContent).toContain('should match');
    change('user.repassword', 'abcdefg1');
    click('Save User');
    await waitFor(() => expect(updated).toHaveBeenCalled());
    expect(userService.update).toHaveBeenCalledWith({
      _id: 'bobby',
      name: 'Bobby',
      role: 'reader',
      password: 'abcdefg1',
      confirmPassword: 'abcdefg1',
    });
  });
  it('updates users with a blank unchanged password and immutable ID', async () => {
    render(<User />);
    await screen.findByText('Alice');
    click('Edit alice');
    await waitFor(() => expect(screen.getByRole('combobox', { name: 'user.role' })).toBeEnabled());
    expect((screen.getByLabelText('user.id') as HTMLInputElement).disabled).toBe(true);
    change('user.name', 'Alice Updated');
    click('Save User');
    await waitFor(() =>
      expect(userService.update).toHaveBeenCalledWith({
        _id: 'alice',
        name: 'Alice Updated',
        role: 'reader',
        password: '',
        confirmPassword: '',
      })
    );
  });
  it('loads permissions and saves changed role selections', async () => {
    render(<Role />);
    await screen.findByText('Reader');
    click('Edit reader');
    await screen.findByLabelText('search');
    expect((screen.getByLabelText('search') as HTMLInputElement).checked).toBe(true);
    fireEvent.click(screen.getByLabelText('search'));
    fireEvent.click(screen.getByLabelText('index'));
    click('Save Role');
    await waitFor(() =>
      expect(roleService.update).toHaveBeenCalledWith({
        _id: 'reader',
        name: 'Reader',
        permission: ['index'],
      })
    );
  });
  it('keeps failed deletion open and supports retry', async () => {
    vi.mocked(userService.delete).mockRejectedValueOnce(new Error('unavailable'));
    render(<User />);
    await screen.findByText('Alice');
    click('Delete alice');
    click('Delete');
    expect((await screen.findByRole('alert')).textContent).toContain('Deletion failed');
    click('Delete');
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    expect(userService.delete).toHaveBeenCalledTimes(2);
  });
  it('paginates and filters local lists', async () => {
    vi.mocked(userService.list).mockResolvedValue(
      response(Array.from({ length: 21 }, (_, i) => ({ _id: `user${i}`, name: `Person ${i}` }))),
    );
    render(<User />);
    await screen.findByText('Person 0');
    expect(screen.queryByText('Person 20')).toBeNull();
    click('Next');
    expect(screen.getByText('Person 20')).toBeTruthy();
    change('user.search', 'PERSON 1');
    expect(screen.getByText('Person 1')).toBeTruthy();
    expect(screen.queryByText('Person 20')).toBeNull();
  });
});
