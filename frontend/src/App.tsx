import { Activity, useState } from 'react';
import { BrowserRouter, Link, NavLink, Navigate, Route, Routes, useLocation } from 'react-router';
import { setCredentials, useAuth } from './auth';
import { languages, useTranslation, type Locale } from './locales';
import { apiEndpoint } from './services/http';
import Login from './views/Login';
import Search from './views/Search';
import Index from './views/Index';
import Template from './views/Template';
import User from './views/User';
import Role from './views/Role';
import About from './views/About';
import Icon from './components/Icon';
import Select from './components/Select';
import ThemeSelect from './components/ThemeSelect';
import Account from './components/Account';

const pages = ['search', 'index', 'template', 'user', 'role', 'about'] as const;

const retainedPages = { '/search': Search, '/index': Index, '/template': Template };
function Console() {
  const user = useAuth();
  const location = useLocation();
  const { t, locale, changeLanguage } = useTranslation();
  const [drawer, setDrawer] = useState(true);
  const [visited, setVisited] = useState<string[]>([]);
  const [editingAccount, setEditingAccount] = useState(false);
  const path = location.pathname.replace(/\/$/, '') || '/';
  if (path in retainedPages && !visited.includes(path)) setVisited([...visited, path]);
  if (!user) return <Navigate to="/login" state={{ from: location.pathname + location.search }} replace />;
  return <div className={`console-shell ${drawer ? '' : 'rail-collapsed'}`}>
    <a className="skip-link" href="#workspace">Skip to content</a>
    {drawer && <aside className="icon-rail" id="console-navigation">
      <button className="workspace-avatar" type="button" aria-label={t('menu.account')} title={`${user.name} · ${t('menu.account')}`} onClick={() => setEditingAccount(true)}><Icon name="account" /></button>
      <nav aria-label="Workspace shortcuts" className="rail-pages">
        {pages.map(name => <NavLink key={name} to={`/${name}`} aria-label={`${t(`menu.${name}`)} workspace`} title={t(`menu.${name}`)} className={({ isActive }) => `rail-link ${isActive ? 'active' : ''}`}><Icon name={name} /></NavLink>)}
      </nav>
      <div className="rail-utilities">
        <a className="rail-link" href={`${apiEndpoint}/swagger/index.html`} target="_blank" rel="noreferrer" aria-label={t('menu.openapi')} title={t('menu.openapi')}><Icon name="openapi" /></a>
        <a className="rail-link" href="https://zincsearch-docs.zinc.dev" target="_blank" rel="noreferrer" aria-label={t('menu.documentation')} title={t('menu.documentation')}><Icon name="documentation" /></a>
        <button className="rail-link" type="button" aria-label={t('menu.signOut')} title={t('menu.signOut')} onClick={() => setCredentials(null)}><Icon name="signOut" /></button>
      </div>
    </aside>}
    <div className="console-body">
      <header className="console-header">
        <button className="menu-toggle" type="button" aria-label="Menu" aria-expanded={drawer} aria-controls="console-navigation" onClick={() => setDrawer(!drawer)}><Icon name="menu" /></button>
        <h1>{pages.some(name => path === `/${name}`) ? t(`menu.${path.slice(1)}`) : t('menu.zincSearch')}</h1>
        <div className="header-account">
          <ThemeSelect />
          <Select displayValue={locale.toUpperCase()} title={languages[locale]} className="language-select" aria-label={t('menu.language')} value={locale} onValueChange={value => changeLanguage(value as Locale)}>{Object.entries(languages).map(([value, label]) => <option key={value} value={value}>{label}</option>)}</Select>
          {!drawer && <button type="button" aria-label={t('menu.signOut')} onClick={() => setCredentials(null)}><Icon name="signOut" /></button>}
          <button className="account-name" type="button" title={`${user.name} · ${t('menu.account')}`} onClick={() => setEditingAccount(true)}>{user.name}</button>
        </div>
      </header>
      {editingAccount && <Account user={user} onClose={() => setEditingAccount(false)} />}
      <nav aria-label="Main navigation" className="section-tabs">
        {pages.map(name => <NavLink key={name} to={`/${name}`} className={({ isActive }) => `section-tab ${isActive ? 'active' : ''}`}>{t(`menu.${name}`)}</NavLink>)}
      </nav>
      <main id="workspace" tabIndex={-1} className="workspace">
        {Object.entries(retainedPages).filter(([key]) => visited.includes(key)).map(([key, Page]) => <Activity key={key} mode={path === key ? 'visible' : 'hidden'}><Page /></Activity>)}
        <Routes>
          <Route index element={<Navigate to="/search" replace />} />
          {Object.keys(retainedPages).map(key => <Route key={key} path={key} element={null} />)}
          <Route path="/user" element={<User />} /><Route path="/role" element={<Role />} /><Route path="/about" element={<About />} />
          <Route path="*" element={<section className="card"><h1>404</h1><p>Sorry, nothing here…</p><Link to="/search">Go Home</Link></section>} />
        </Routes>
      </main>
    </div>
  </div>;
}
export function AppRoutes() {
  return <Routes><Route path="/login" element={<Login />} /><Route path="/*" element={<Console />} /></Routes>;
}
export default function App() {
  return <BrowserRouter basename={import.meta.env.BASE_URL}><AppRoutes /></BrowserRouter>;
}
