import { useState } from 'react';
import userService from '../../services/user';
import roleService from '../../services/role';
import { useTranslation } from '../../locales';
import AccountEditor, { type Account } from './AccountEditor';
import { type Column, ConfirmDelete, formatDate, LocalList, useList } from './Common';

export default function Accounts({ kind }: { kind: 'user' | 'role' }) {
  const { t } = useTranslation();
  const service = kind === 'user' ? userService : roleService;
  const list = useList<Account>(service.list);
  const [editor, setEditor] = useState<Account | 'new' | null>(null);
  const [deleting, setDeleting] = useState<Account | null>(null);
  const columns: Column<Account>[] = [
    { label: 'ID', render: (row) => row._id },
    { label: 'NAME', render: (row) => row.name || row._id },
    ...(kind === 'user' ? [{ label: 'ROLE', render: (row: Account) => row.role }] : []),
    { label: 'CREATED', render: (row) => formatDate(row.created_at) },
    { label: 'UPDATED', render: (row) => formatDate(row.updated_at) },
    {
      label: 'ACTIONS',
      render: (row) => (
        <div className='toolbar'>
          <button aria-label={`Edit ${row._id}`} onClick={() => setEditor(row)}>Edit</button>
          <button aria-label={`Delete ${row._id}`} onClick={() => setDeleting(row)}>Delete</button>
        </div>
      ),
    },
  ];
  return (
    <>
      <LocalList
        {...list}
        title={t(`${kind}.header`)}
        searchLabel={t(`${kind}.search`)}
        addLabel={t(`${kind}.add`)}
        name={(row) => row.name || row._id}
        rowKey={(row) => row._id}
        columns={columns}
        onAdd={() => setEditor('new')}
        onRetry={list.refresh}
      />
      {editor && (
        <AccountEditor
          kind={kind}
          value={editor === 'new' ? undefined : editor}
          onClose={() => setEditor(null)}
          onUpdated={() => {
            setEditor(null);
            list.refresh();
          }}
        />
      )}
      {deleting && (
        <ConfirmDelete
          names={[deleting._id]}
          onClose={() => setDeleting(null)}
          onDelete={async () => {
            await service.delete(deleting._id);
            list.refresh();
          }}
        />
      )}
    </>
  );
}
