import { useState } from 'react';
import Select from '../Select';
import { isAxiosError } from 'axios';
import indexService from '../../services/index';
import templateService from '../../services/template';
import { ErrorMessage, Modal } from './Common';

export type Schema = { settings?: Record<string, unknown>; mappings?: Record<string, unknown> };
export type TemplateRecord = {
  name: string;
  index_template: { index_patterns: string[]; priority?: number; template?: Schema };
};
export function parseObject(text: string): Record<string, unknown> {
  let result: unknown;
  try {
    result = JSON.parse(text);
  } catch {
    throw new Error('Invalid JSON format.');
  }
  if (!result || typeof result !== 'object' || Array.isArray(result)) {
    throw new Error('JSON must be an object.');
  }
  return result as Record<string, unknown>;
}

export default function SchemaEditor(
  { kind, value, onClose, onUpdated }: {
    kind: 'index' | 'template';
    value?: TemplateRecord;
    onClose: () => void;
    onUpdated: () => void;
  },
) {
  const [step, setStep] = useState(0);
  const [name, setName] = useState(value?.name || '');
  const [patterns, setPatterns] = useState(value?.index_template.index_patterns.join(', ') || '');
  const [number, setNumber] = useState(value?.index_template.priority?.toString() || '');
  const [settings, setSettings] = useState(
    JSON.stringify(value?.index_template.template?.settings || {}, null, 2),
  );
  const [mappings, setMappings] = useState(
    JSON.stringify(value?.index_template.template?.mappings || {}, null, 2),
  );
  const initialSQL = value?.index_template.template?.mappings?.sql;
  const [format, setFormat] = useState(typeof initialSQL === 'string' ? 'sql' : 'json');
  const [sql, setSQL] = useState(typeof initialSQL === 'string' ? initialSQL : '');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  let jsonError = '';
  if (step === 1 || (step === 2 && format === 'json')) {
    try {
      parseObject(step === 1 ? settings : mappings);
    } catch (err) {
      jsonError = (err as Error).message;
    }
  }
  const steps = [
    'Logistics',
    kind === 'index' ? 'Settings' : 'Index settings',
    'Mappings',
    'Review',
  ];
  function mappingObject(): Record<string, unknown> {
    const result = format === 'sql' ? { sql } : parseObject(mappings);
    if (typeof result.sql === 'string' && !result.sql.trim()) {
      throw new Error('CREATE TABLE SQL is required.');
    }
    return result;
  }
  function payload() {
    const schema = { settings: parseObject(settings), mappings: mappingObject() };
    return kind === 'index'
      ? {
        name,
        storage_type: 'disk',
        ...(number !== '' ? { shard_num: Number(number) } : {}),
        ...schema,
      }
      : {
        name,
        index_patterns: [...new Set(patterns.split(/[\n,]+/).map((p) => p.trim()).filter(Boolean))],
        ...(number !== '' ? { priority: Number(number) } : {}),
        template: schema,
      };
  }
  async function next() {
    setError('');
    if (step === 0) {
      if (!name.trim()) {
        setError(`${kind === 'index' ? 'Index' : 'Template'} name is required`);
        return;
      }
      if (kind === 'template' && !patterns.split(/[\n,]+/).some((p) => p.trim())) {
        setError('Index patterns are required');
        return;
      }
      if (
        number !== '' &&
        (!Number.isInteger(Number(number)) || Number(number) < (kind === 'index' ? 1 : 0))
      ) {
        setError(
          kind === 'index'
            ? 'Shard num must be a positive integer'
            : 'Priority must be a non-negative integer',
        );
        return;
      }
    }
    if (step === 1 || step === 2) {
      try {
        if (step === 1) parseObject(settings);
        else mappingObject();
      } catch (err) {
        setError((err as Error).message);
        return;
      }
    }
    if (step < 3) {
      setStep(step + 1);
      return;
    }
    setBusy(true);
    try {
      await (kind === 'index' ? indexService : templateService).update(payload());
      onUpdated();
    } catch (err) {
      const detail = isAxiosError(err) ? err.response?.data?.error : undefined;
      setError(typeof detail === 'string' && detail
        ? detail
        : `Unable to save ${kind}. Please try again.`);
    } finally {
      setBusy(false);
    }
  }
  return (
    <Modal
      title={`${value ? 'Update' : 'Add'} ${kind}`}
      onClose={() => {
        if (!busy) onClose();
      }}
    >
      <ol className='toolbar'>
        {steps.map((label, i) => (
          <li
            key={label}
            aria-current={step === i ? 'step' : undefined}
            className={step === i ? 'font-bold' : ''}
          >
            {i + 1}. {label}
          </li>
        ))}
      </ol>
      <form
        className='grid gap-4'
        onSubmit={(e) => {
          e.preventDefault();
          void next();
        }}
      >
        {step === 0 && (
          <>
            <label>
              {kind === 'index' ? 'Index Name' : 'Template Name'}
              <input
                className='block w-full'
                value={name}
                disabled={!!value}
                onChange={(e) => setName(e.target.value)}
              />
            </label>
            {kind === 'index'
              ? (
                <label>
                  Storage Type<Select className='block w-full' value='disk' onValueChange={() => {}}>
                    <option>disk</option>
                  </Select>
                </label>
              )
              : (
                <label>
                  Index Patterns<textarea
                    className='block w-full'
                    placeholder='Comma-separated patterns, e.g. logs-*, events-*'
                    value={patterns}
                    onChange={(e) => setPatterns(e.target.value)}
                  />
                </label>
              )}
            <label>
              {kind === 'index' ? 'Shard num (optional)' : 'Priority (optional)'}
              <input
                className='block w-full'
                value={number}
                onChange={(e) => setNumber(e.target.value)}
              />
            </label>
          </>
        )}
        {step === 2 && (
          <label>
            Mappings format
            <Select
              className='block w-full'
              value={format}
              onValueChange={(value) => {
                setFormat(value);
                setError('');
              }}
            >
              <option value='json'>JSON</option>
              <option value='sql'>MySQL SQL</option>
            </Select>
          </label>
        )}
        {step === 2 && format === 'sql' && (
          <>
            <label>
              CREATE TABLE SQL
              <textarea
                className='block w-full font-mono min-h-80'
                spellCheck={false}
                placeholder='CREATE TABLE events (message TEXT);'
                value={sql}
                onChange={(e) => {
                  setSQL(e.target.value);
                  setError('');
                }}
              />
            </label>
            <p>MySQL CREATE TABLE defines the schema only; no SQL is executed. SQL replaces JSON properties. Settings remain JSON.</p>
          </>
        )}
        {(step === 1 || (step === 2 && format === 'json')) && (
          <>
            <label>
              {steps[step]} JSON<textarea
                className='block w-full font-mono min-h-80'
                spellCheck={false}
                aria-invalid={!!jsonError}
                value={step === 1 ? settings : mappings}
                onChange={(e) => {
                  (step === 1 ? setSettings : setMappings)(e.target.value);
                  setError('');
                }}
              />
            </label>
            <p>
              Use JSON format:{' '}
              <code>
                {step === 1
                  ? '{"analysis":{"analyzer":{"default":{"type":"standard"}}}}'
                  : '{"properties":{"content":{"type":"text"}}}'}
              </code>
            </p>
          </>
        )}
        {step === 3 && (
          <pre
            className='overflow-auto p-4'
            aria-label='Review JSON'
          >{JSON.stringify(payload(), null, 2)}</pre>
        )}
        <ErrorMessage message={jsonError || error} />
        <div className='toolbar'>
          {step > 0 && (
            <button
              type='button'
              disabled={busy}
              onClick={() => {
                setError('');
                setStep(step - 1);
              }}
            >
              Back
            </button>
          )}
          <button type='submit' className='primary' disabled={busy || !!jsonError}>
            {step === 3 ? `Save ${kind === 'index' ? 'Index' : 'Template'}` : 'Continue'}
          </button>
        </div>
      </form>
    </Modal>
  );
}

export function SchemaPreview(
  { name, summary, schema, data, onClose }: {
    name: string;
    summary: Record<string, string | number | undefined>;
    schema: Schema;
    data: unknown;
    onClose: () => void;
  },
) {
  const [tab, setTab] = useState('Summary');
  return (
    <Modal title={name} onClose={onClose}>
      <div className='toolbar' role='tablist' aria-label='Preview sections'>
        {['Summary', 'Settings', 'Mappings', 'Preview'].map((label) => (
          <button
            role='tab'
            aria-selected={tab === label}
            key={label}
            onClick={() => setTab(label)}
          >
            {label}
          </button>
        ))}
      </div>
      <div role='tabpanel'>
        {tab === 'Summary'
          ? (
            <dl className='grid grid-cols-2 gap-4 p-4'>
              {Object.entries(summary).map(([key, val]) => (
                <div key={key}>
                  <dt className='font-bold'>{key}</dt>
                  <dd>{val}</dd>
                </div>
              ))}
            </dl>
          )
          : (
            <pre className='overflow-auto p-4'>{JSON.stringify(tab === 'Settings' ? schema.settings || {} : tab === 'Mappings' ? schema.mappings || {} : data, null, 2)}</pre>
          )}
      </div>
    </Modal>
  );
}
