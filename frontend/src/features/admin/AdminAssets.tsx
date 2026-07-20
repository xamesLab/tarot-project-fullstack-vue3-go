import { useState, type FormEvent } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import { resolveApiAssetUrl } from '@/api/client';
import type { CardBackDto, DeckDto, PaginationDto } from '@/api/types';
import { usePreferencesStore } from '@/features/preferences/store';
import { Button } from '@/shared/ui/Button';
import { SegmentedControl } from '@/shared/ui/Controls';
import { Field, TextareaField } from '@/shared/ui/Field';
import { Modal } from '@/shared/ui/Modal';
import { AdminState } from './AdminDashboard';
import { CardBackCropModal } from './CardBackCropModal';
import { adminKeys, getAdminError } from './shared';
import { readDeckMetaFromZip } from './zipMeta';

type AssetScope = 'active' | 'archived';
type AssetAction = { type: 'archive' | 'restore' | 'delete'; asset: 'deck' | 'back'; id: string; name: string };

const emptyDeckForm = { nameEn: '', nameRu: '', descriptionEn: '', descriptionRu: '', author: '', year: '', tags: '' };
const emptyBackForm = { nameEn: '', nameRu: '', descriptionEn: '', descriptionRu: '', tags: '' };

export function AdminAssets() {
  const { t } = useTranslation('admin');
  const language = usePreferencesStore((state) => state.language);
  const queryClient = useQueryClient();
  const [scope, setScope] = useState<AssetScope>('active');
  const [deckPage, setDeckPage] = useState(1);
  const [backPage, setBackPage] = useState(1);
  const [deckForm, setDeckForm] = useState(emptyDeckForm);
  const [backForm, setBackForm] = useState(emptyBackForm);
  const [zipFile, setZipFile] = useState<File | null>(null);
  const [backFile, setBackFile] = useState<File | null>(null);
  const [zipKey, setZipKey] = useState(0);
  const [backKey, setBackKey] = useState(0);
  const [metaMessage, setMetaMessage] = useState('');
  const [crop, setCrop] = useState({ open: false, url: '', name: '' });
  const [pendingAction, setPendingAction] = useState<AssetAction | null>(null);
  const limit = 12;
  const deckQuery = useQuery({ queryKey: adminKeys.decks(scope, deckPage, limit), queryFn: () => api.admin.getDecks({ scope, page: deckPage, limit }) });
  const backQuery = useQuery({ queryKey: adminKeys.cardBacks(scope, backPage, limit), queryFn: () => api.admin.getCardBacks({ scope, page: backPage, limit }) });
  const localize = (value: { ru: string; en: string }) => value[language] || value.en || value.ru;

  const invalidateAssets = async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: ['admin', 'decks'] }),
      queryClient.invalidateQueries({ queryKey: ['admin', 'card-backs'] }),
      queryClient.invalidateQueries({ queryKey: adminKeys.dashboard }),
      queryClient.invalidateQueries({ queryKey: adminKeys.settings }),
    ]);
  };

  const deckCreate = useMutation({
    mutationFn: (payload: FormData) => api.admin.createDeck(payload),
    onSuccess: async () => { setDeckForm(emptyDeckForm); setZipFile(null); setZipKey((value) => value + 1); setMetaMessage(t('assets.deckCreated')); await invalidateAssets(); },
  });
  const backCreate = useMutation({
    mutationFn: (payload: FormData) => api.admin.createCardBack(payload),
    onSuccess: async () => { setBackForm(emptyBackForm); setBackFile(null); setBackKey((value) => value + 1); await invalidateAssets(); },
  });
  const assetMutation = useMutation({
    mutationFn: async (action: AssetAction) => {
      if (action.asset === 'deck') {
        if (action.type === 'archive') return api.admin.archiveDeck(action.id);
        if (action.type === 'restore') return api.admin.restoreDeck(action.id);
        return api.admin.deleteDeck(action.id);
      }
      if (action.type === 'archive') return api.admin.archiveCardBack(action.id);
      if (action.type === 'restore') return api.admin.restoreCardBack(action.id);
      return api.admin.deleteCardBack(action.id);
    },
    onSuccess: async () => { setPendingAction(null); await invalidateAssets(); },
  });

  const submitDeck = (event: FormEvent) => {
    event.preventDefault();
    if (!zipFile) return;
    const payload = new FormData();
    payload.append('name_en', deckForm.nameEn.trim()); payload.append('name_ru', deckForm.nameRu.trim());
    payload.append('description_en', deckForm.descriptionEn.trim()); payload.append('description_ru', deckForm.descriptionRu.trim());
    payload.append('author', deckForm.author.trim()); payload.append('year', deckForm.year.trim()); payload.append('tags', deckForm.tags.trim()); payload.append('cards_zip', zipFile);
    deckCreate.mutate(payload);
  };
  const submitBack = (event: FormEvent) => {
    event.preventDefault();
    if (!backFile) return;
    const payload = new FormData();
    payload.append('name_en', backForm.nameEn.trim()); payload.append('name_ru', backForm.nameRu.trim());
    payload.append('description_en', backForm.descriptionEn.trim()); payload.append('description_ru', backForm.descriptionRu.trim()); payload.append('tags', backForm.tags.trim()); payload.append('back_image', backFile);
    backCreate.mutate(payload);
  };
  const closeCrop = () => {
    if (crop.url) URL.revokeObjectURL(crop.url);
    setCrop({ open: false, url: '', name: '' });
    if (!backFile) setBackKey((value) => value + 1);
  };

  return (
    <div className="admin-stack">
      <section className="admin-upload-grid">
        <form className="admin-surface admin-upload-card" onSubmit={submitDeck}>
          <div className="admin-section-head"><div><span>{t('assets.uploadEyebrow')}</span><h2>{t('assets.deckUploadTitle')}</h2><p>{t('assets.deckUploadCopy')}</p></div></div>
          <div className="admin-form-grid"><Field required label={t('assets.nameRu')} value={deckForm.nameRu} onChange={(event) => setDeckForm((current) => ({ ...current, nameRu: event.target.value }))} /><Field required label={t('assets.nameEn')} value={deckForm.nameEn} onChange={(event) => setDeckForm((current) => ({ ...current, nameEn: event.target.value }))} /><TextareaField id="admin-deck-description-ru" label={t('assets.descriptionRu')} value={deckForm.descriptionRu} onChange={(event) => setDeckForm((current) => ({ ...current, descriptionRu: event.target.value }))} /><TextareaField id="admin-deck-description-en" label={t('assets.descriptionEn')} value={deckForm.descriptionEn} onChange={(event) => setDeckForm((current) => ({ ...current, descriptionEn: event.target.value }))} /><Field label={t('assets.author')} value={deckForm.author} onChange={(event) => setDeckForm((current) => ({ ...current, author: event.target.value }))} /><Field label={t('assets.year')} type="number" min="0" value={deckForm.year} onChange={(event) => setDeckForm((current) => ({ ...current, year: event.target.value }))} /></div>
          <Field label={t('assets.tags')} placeholder={t('assets.tagsPlaceholder')} value={deckForm.tags} onChange={(event) => setDeckForm((current) => ({ ...current, tags: event.target.value }))} />
          <Field key={zipKey} required label={t('assets.zip')} type="file" accept=".zip,application/zip" onChange={(event) => { const file = event.target.files?.[0] || null; setZipFile(file); setMetaMessage(''); if (file) void readDeckMetaFromZip(file).then((meta) => { if (!meta) { setMetaMessage(t('assets.metaMissing')); return; } setDeckForm(meta); setMetaMessage(t('assets.metaLoaded')); }).catch(() => setMetaMessage(t('assets.metaError'))); }} />
          {metaMessage && <p className="admin-form-hint">{metaMessage}</p>}
          {deckCreate.isError && <p className="admin-feedback error" role="alert">{getAdminError(deckCreate.error, t('states.error'))}</p>}
          <Button type="submit" disabled={!zipFile || deckCreate.isPending}>{deckCreate.isPending ? t('actions.uploading') : t('assets.uploadDeck')}</Button>
        </form>

        <form className="admin-surface admin-upload-card" onSubmit={submitBack}>
          <div className="admin-section-head"><div><span>{t('assets.uploadEyebrow')}</span><h2>{t('assets.backUploadTitle')}</h2><p>{t('assets.backUploadCopy')}</p></div></div>
          <div className="admin-form-grid"><Field required label={t('assets.nameRu')} value={backForm.nameRu} onChange={(event) => setBackForm((current) => ({ ...current, nameRu: event.target.value }))} /><Field required label={t('assets.nameEn')} value={backForm.nameEn} onChange={(event) => setBackForm((current) => ({ ...current, nameEn: event.target.value }))} /><TextareaField id="admin-back-description-ru" label={t('assets.descriptionRu')} value={backForm.descriptionRu} onChange={(event) => setBackForm((current) => ({ ...current, descriptionRu: event.target.value }))} /><TextareaField id="admin-back-description-en" label={t('assets.descriptionEn')} value={backForm.descriptionEn} onChange={(event) => setBackForm((current) => ({ ...current, descriptionEn: event.target.value }))} /></div>
          <Field label={t('assets.tags')} placeholder={t('assets.tagsPlaceholder')} value={backForm.tags} onChange={(event) => setBackForm((current) => ({ ...current, tags: event.target.value }))} />
          <Field key={backKey} required label={t('assets.image')} type="file" accept="image/png,image/jpeg,image/webp" onChange={(event) => { const file = event.target.files?.[0]; if (!file) return; setBackFile(null); setCrop({ open: true, url: URL.createObjectURL(file), name: file.name }); }} />
          {backFile && <p className="admin-selected-file">{t('assets.selectedFile', { name: backFile.name })}</p>}
          {backCreate.isError && <p className="admin-feedback error" role="alert">{getAdminError(backCreate.error, t('states.error'))}</p>}
          {backCreate.isSuccess && <p className="admin-feedback success" role="status">{t('assets.backCreated')}</p>}
          <Button type="submit" disabled={!backFile || backCreate.isPending}>{backCreate.isPending ? t('actions.uploading') : t('assets.uploadBack')}</Button>
        </form>
      </section>

      <section className="admin-surface admin-catalog">
        <div className="admin-catalog-head"><div><span>{t('assets.catalogEyebrow')}</span><h2>{scope === 'active' ? t('assets.catalogTitle') : t('assets.archiveTitle')}</h2><p>{scope === 'active' ? t('assets.catalogCopy') : t('assets.archiveCopy')}</p></div><SegmentedControl value={scope} options={[{ value: 'active', label: t('assets.catalog') }, { value: 'archived', label: t('assets.archive') }]} onChange={(value) => { setScope(value); setDeckPage(1); setBackPage(1); }} label={t('assets.scopeLabel')} /></div>
        <AssetGroup title={t('assets.decksTitle')} loading={deckQuery.isPending} error={deckQuery.isError ? getAdminError(deckQuery.error, t('states.error')) : ''} empty={!deckQuery.data?.decks.length} emptyText={t('assets.emptyDecks')} pagination={deckQuery.data?.pagination} page={deckPage} setPage={setDeckPage}>
          {deckQuery.data?.decks.map((deck) => <DeckCard key={deck.id} deck={deck} name={localize(deck.name)} language={language} busy={assetMutation.isPending} onAction={(type) => type === 'restore' ? assetMutation.mutate({ type, asset: 'deck', id: deck.id, name: localize(deck.name) }) : setPendingAction({ type, asset: 'deck', id: deck.id, name: localize(deck.name) })} />)}
        </AssetGroup>
        <AssetGroup title={t('assets.backsTitle')} loading={backQuery.isPending} error={backQuery.isError ? getAdminError(backQuery.error, t('states.error')) : ''} empty={!backQuery.data?.card_backs.length} emptyText={t('assets.emptyBacks')} pagination={backQuery.data?.pagination} page={backPage} setPage={setBackPage}>
          {backQuery.data?.card_backs.map((back) => <BackCard key={back.id} back={back} name={localize(back.name)} language={language} busy={assetMutation.isPending} onAction={(type) => type === 'restore' ? assetMutation.mutate({ type, asset: 'back', id: back.id, name: localize(back.name) }) : setPendingAction({ type, asset: 'back', id: back.id, name: localize(back.name) })} />)}
        </AssetGroup>
        {assetMutation.isError && <p className="admin-feedback error" role="alert">{getAdminError(assetMutation.error, t('states.error'))}</p>}
      </section>

      <CardBackCropModal open={crop.open} imageUrl={crop.url} fileName={crop.name} onClose={closeCrop} onConfirm={(file) => { setBackFile(file); closeCrop(); }} />
      <Modal open={Boolean(pendingAction)} onClose={() => setPendingAction(null)} title={pendingAction?.type === 'delete' ? t('assets.confirmDeleteTitle') : t('assets.confirmArchiveTitle')} footer={<><Button variant="ghost" onClick={() => setPendingAction(null)}>{t('actions.cancel')}</Button><Button disabled={assetMutation.isPending} onClick={() => pendingAction && assetMutation.mutate(pendingAction)}>{pendingAction?.type === 'delete' ? t('actions.delete') : t('actions.archive')}</Button></>}>
        <p className="admin-confirm-copy">{pendingAction?.type === 'delete' ? t('assets.confirmDelete', { name: pendingAction.name }) : t('assets.confirmArchive', { name: pendingAction?.name })}</p>
      </Modal>
    </div>
  );
}

function AssetGroup({ title, loading, error, empty, emptyText, pagination, page, setPage, children }: { title: string; loading: boolean; error: string; empty: boolean; emptyText: string; pagination?: PaginationDto; page: number; setPage: (page: number) => void; children: React.ReactNode }) {
  const { t } = useTranslation('admin');
  return <section className="admin-asset-group"><h3>{title}</h3>{loading ? <AdminState>{t('states.loading')}</AdminState> : error ? <AdminState error>{error}</AdminState> : empty ? <div className="admin-empty">{emptyText}</div> : <div className="admin-asset-list">{children}</div>}{pagination && <div className="admin-pagination"><Button variant="ghost" size="sm" disabled={page <= 1 || loading} onClick={() => setPage(page - 1)}>{t('actions.previous')}</Button><span>{t('pagination', { page: pagination.page, total: pagination.total_pages })}</span><Button variant="ghost" size="sm" disabled={page >= pagination.total_pages || loading} onClick={() => setPage(page + 1)}>{t('actions.next')}</Button></div>}</section>;
}

function DeckCard({ deck, name, language, busy, onAction }: { deck: DeckDto; name: string; language: 'ru' | 'en'; busy: boolean; onAction: (type: AssetAction['type']) => void }) {
  const { t } = useTranslation('admin');
  const description = deck.description[language] || deck.description.en || deck.description.ru;
  return <article className="admin-asset-card"><div className="admin-deck-preview">{deck.back_image_path ? <img src={resolveApiAssetUrl(deck.back_image_path)} alt="" /> : <span>✦</span>}</div><div className="admin-asset-copy"><div className="admin-asset-title"><div><h4>{name}</h4><p>{description}</p></div><span className={`admin-status admin-status-${deck.status}`}>{t(`statuses.${deck.status}`, { defaultValue: deck.status })}</span></div><dl><div><dt>{t('assets.author')}</dt><dd>{deck.author || '—'}</dd></div><div><dt>{t('assets.year')}</dt><dd>{deck.year || '—'}</dd></div><div><dt>{t('assets.cards')}</dt><dd>{deck.card_count}</dd></div><div><dt>{t('assets.tags')}</dt><dd>{deck.tags?.join(', ') || '—'}</dd></div></dl><AssetActions archived={deck.status === 'archived'} busy={busy} onAction={onAction} /></div></article>;
}

function BackCard({ back, name, language, busy, onAction }: { back: CardBackDto; name: string; language: 'ru' | 'en'; busy: boolean; onAction: (type: AssetAction['type']) => void }) {
  const { t } = useTranslation('admin');
  return <article className="admin-asset-card"><div className="admin-back-preview">{back.image_path ? <img src={resolveApiAssetUrl(back.image_path)} alt={name} /> : <span>✦</span>}</div><div className="admin-asset-copy"><div className="admin-asset-title"><div><h4>{name}</h4><p>{back.description[language] || back.description.en || back.description.ru}</p></div><span className={`admin-status admin-status-${back.status}`}>{t(`statuses.${back.status}`, { defaultValue: back.status })}</span></div><dl><div><dt>{t('assets.tags')}</dt><dd>{back.tags?.join(', ') || '—'}</dd></div><div><dt>{t('assets.created')}</dt><dd>{back.created_at ? new Date(back.created_at).toLocaleDateString(language === 'ru' ? 'ru-RU' : 'en-US') : '—'}</dd></div></dl><AssetActions archived={back.status === 'archived'} busy={busy} onAction={onAction} /></div></article>;
}

function AssetActions({ archived, busy, onAction }: { archived: boolean; busy: boolean; onAction: (type: AssetAction['type']) => void }) {
  const { t } = useTranslation('admin');
  return <div className="admin-asset-actions">{archived ? <><Button variant="outline" size="sm" disabled={busy} onClick={() => onAction('restore')}>{t('actions.restore')}</Button><Button variant="ghost" size="sm" disabled={busy} onClick={() => onAction('delete')}>{t('actions.deletePermanent')}</Button></> : <Button variant="outline" size="sm" disabled={busy} onClick={() => onAction('archive')}>{t('actions.archive')}</Button>}</div>;
}
