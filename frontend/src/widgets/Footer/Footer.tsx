import { useTranslation } from 'react-i18next';
import './footer.css';

export function Footer() {
  const { t } = useTranslation();

  return (
    <footer className="site-footer">
      <div>
        <span>{t('copyright')}</span>
        <span aria-hidden="true"> · · · </span>
        <span>{t('footer')}</span>
      </div>
    </footer>
  );
}
