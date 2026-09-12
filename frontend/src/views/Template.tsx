import { useState } from 'react';
import templateService from '../services/template';
import { useTranslation } from '../locales';
import { type Column, ConfirmDelete, LocalList, useList } from '../components/management/Common';
import SchemaEditor, {
  SchemaPreview,
  type TemplateRecord,
} from '../components/management/SchemaEditor';

export default function Template() {
  const { t } = useTranslation();
  const list = useList<TemplateRecord>(templateService.list);
  const [editor, setEditor] = useState<TemplateRecord | 'new' | null>(null);
  const [preview, setPreview] = useState<TemplateRecord | null>(null);
  const [deleting, setDeleting] = useState<TemplateRecord | null>(null);
  const columns: Column<TemplateRecord>[] = [
    { label: 'NAME', render: (row) => <button onClick={() => setPreview(row)}>{row.name}</button> },
    { label: 'PATTERNS', render: (row) => row.index_template.index_patterns.join(', ') },
    { label: 'PRIORITY', render: (row) => row.index_template.priority ?? '' },
    {
      label: 'TEMPLATE',
      render: (row) => (
        <span className='template-indicators'>
          {row.index_template.template?.mappings && <span title='Mappings'>M</span>}
          {row.index_template.template?.settings && <span title='Settings'>S</span>}
          {!row.index_template.template?.mappings && !row.index_template.template?.settings &&
            'None'}
        </span>
      ),
    },
    {
      label: 'ACTIONS',
      render: (row) => (
        <div className='toolbar'>
          <button aria-label={`Edit ${row.name}`} onClick={() => setEditor(row)}>Edit</button>
          <button aria-label={`Delete ${row.name}`} onClick={() => setDeleting(row)}>Delete</button>
        </div>
      ),
    },
  ];
  return (
    <>
      <LocalList
        {...list}
        title={t('template.header')}
        searchLabel={t('template.search')}
        addLabel={t('template.add')}
        columns={columns}
        name={(row) => row.name}
        rowKey={(row) => row.name}
        onAdd={() => setEditor('new')}
        onRetry={list.refresh}
      />
      {editor && (
        <SchemaEditor
          kind='template'
          value={editor === 'new' ? undefined : editor}
          onClose={() => setEditor(null)}
          onUpdated={() => {
            setEditor(null);
            list.refresh();
          }}
        />
      )}
      {preview && (
        <SchemaPreview
          name={preview.name}
          summary={{
            Name: preview.name,
            'Index patterns': preview.index_template.index_patterns.join(', '),
            Priority: preview.index_template.priority,
          }}
          schema={preview.index_template.template || {}}
          data={{ name: preview.name, ...preview.index_template }}
          onClose={() => setPreview(null)}
        />
      )}
      {deleting && (
        <ConfirmDelete
          names={[deleting.name]}
          onClose={() => setDeleting(null)}
          onDelete={async () => {
            await templateService.delete(deleting.name);
            list.refresh();
          }}
        />
      )}
    </>
  );
}
