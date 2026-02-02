import { Link } from 'react-router-dom'
import { useLocale } from '../i18n/LocaleContext'
import styles from './Landing.module.css'

export function Landing() {
  const { t, locale, setLocale } = useLocale()

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <Link to="/" className={styles.logo}>
          {t.appName}
        </Link>
        <nav className={styles.nav}>
          <button
            type="button"
            className={styles.lang}
            onClick={() => setLocale(locale === 'ru' ? 'en' : 'ru')}
          >
            {locale === 'ru' ? 'EN' : 'RU'}
          </button>
          <Link to="/login">{t.login}</Link>
          <Link to="/register" className={styles.cta}>
            {t.register}
          </Link>
        </nav>
      </header>

      <section className={styles.hero}>
        <h1 className={styles.heroTitle}>{t.landingHero}</h1>
        <p className={styles.heroSub}>{t.landingSub}</p>
        <Link to="/register" className={styles.heroButton}>
          {t.getStarted}
        </Link>
      </section>

      <section className={styles.steps}>
        <h2 className={styles.stepsTitle}>
          {locale === 'ru' ? 'Как сделать звонок в 3 шага' : 'How to call in 3 steps'}
        </h2>
        <div className={styles.stepsGrid}>
          <div className={styles.step}>
            <div className={styles.stepNum}>1</div>
            <h3>{t.step1}</h3>
            <p>{t.step1Desc}</p>
          </div>
          <div className={styles.step}>
            <div className={styles.stepNum}>2</div>
            <h3>{t.step2}</h3>
            <p>{t.step2Desc}</p>
          </div>
          <div className={styles.step}>
            <div className={styles.stepNum}>3</div>
            <h3>{t.step3}</h3>
            <p>{t.step3Desc}</p>
          </div>
        </div>
      </section>

      <section className={styles.features}>
        <h2 className={styles.featuresTitle}>
          {locale === 'ru' ? 'Почему выбирают нас?' : 'Why choose us?'}
        </h2>
        <div className={styles.featuresGrid}>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>🌐</span>
            <h3>{locale === 'ru' ? 'Звоните откуда угодно' : 'Call from anywhere'}</h3>
            <p>
              {locale === 'ru'
                ? 'Международные звонки прямо из браузера без установки приложений'
                : 'International calls from your browser, no app installation'}
            </p>
          </div>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>💳</span>
            <h3>{locale === 'ru' ? 'Платите за минуты' : 'Pay per minute'}</h3>
            <p>
              {locale === 'ru'
                ? 'Никаких подписок. Платите только за использованные минуты'
                : 'No subscriptions. Pay only for minutes used'}
            </p>
          </div>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>📱</span>
            <h3>{locale === 'ru' ? 'Без задержек' : 'Without delays'}</h3>
            <p>
              {locale === 'ru'
                ? 'Использует современные технологии'
                : 'Uses modern technologies'}
            </p>
          </div>
        </div>
      </section>

      <section className={styles.ctaSection}>
        <h2>{locale === 'ru' ? 'Начните звонить сейчас' : 'Start calling now'}</h2>
        <Link to="/register" className={styles.ctaButton}>
          {t.getStarted}
        </Link>
      </section>

      <footer className={styles.footer}>
        <p>© 2025 {t.appName}</p>
      </footer>
    </div>
  )
}
