import { useState } from 'react';
import { Avatar, UserPill } from '@/shared/ui/Avatar';
import { Checkbox, RadioCard, SegmentedControl, Toggle } from '@/shared/ui/Controls';
import { NavBar, type NavItem } from '@/shared/ui/NavBar';
import { Divider, Ornament } from '@/shared/ui/Ornament';
import { MoonPhases } from '@/shared/icons';
import { CodeCaption, DemoRow, UiKitSection } from './UiKitSection';

export function ControlSections() {
  const [notifications, setNotifications] = useState(true);
  const [sound, setSound] = useState(false);
  const [checked, setChecked] = useState(true);
  const [history, setHistory] = useState(false);
  const [segment, setSegment] = useState<'direct' | 'reversed'>('direct');
  const [spread, setSpread] = useState('three');
  const [nav, setNav] = useState('daily');
  const navItems: NavItem[] = [
    { id: 'daily', label: 'Карта дня' },
    { id: 'spreads', label: 'Расклады' },
    { id: 'horoscope', label: 'Гороскоп' },
    { id: 'learn', label: 'Обучение' },
  ];
  return (
    <>
      <UiKitSection
        id="controls"
        number="06"
        title="Переключатели"
        description="Интерактивные тоглы, чекбоксы и сегмент-контрол — нажмите, чтобы проверить состояние."
      >
        <div className="control-stack">
          <label className="control-label">
            <Toggle
              checked={notifications}
              onChange={setNotifications}
              label="Ежедневное уведомление о карте"
            />
            Ежедневное уведомление о карте
          </label>
          <label className="control-label">
            <Toggle checked={sound} onChange={setSound} label="Звук перетасовки" />
            Звук перетасовки
          </label>
        </div>
        <CodeCaption>.tgl / .tgl.on</CodeCaption>
        <div className="kit-ornament-gap" />
        <div className="control-stack">
          <Checkbox checked={checked} onChange={setChecked}>
            Учитывать перевёрнутые карты
          </Checkbox>
          <Checkbox checked={history} onChange={setHistory}>
            Сохранять историю раскладов
          </Checkbox>
        </div>
        <CodeCaption>.cbx / .cbx.on</CodeCaption>
        <div className="kit-ornament-gap" />
        <DemoRow>
          <SegmentedControl
            value={segment}
            options={[
              { value: 'direct', label: 'Прямое' },
              { value: 'reversed', label: 'Перевёрнутое' },
            ]}
            onChange={setSegment}
            label="Положение карты"
          />
        </DemoRow>
        <CodeCaption>.seg › .seg-opt.on</CodeCaption>
        <div className="radio-row" role="radiogroup">
          <RadioCard checked={spread === 'three'} onClick={() => setSpread('three')}>
            Три карты · прошлое · настоящее · будущее
          </RadioCard>
          <RadioCard checked={spread === 'cross'} onClick={() => setSpread('cross')}>
            Кельтский крест · 10 карт
          </RadioCard>
          <RadioCard checked={spread === 'one'} onClick={() => setSpread('one')}>
            Одна карта дня
          </RadioCard>
        </div>
        <CodeCaption>.radio › .rad.on</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="avatars"
        number="07"
        title="Аватары & User Pill"
        description="Золотая градиентная заливка с инициалами, в трёх размерах и в составе плашки пользователя."
      >
        <DemoRow className="kit-avatar-row">
          <Avatar initials="МА" size="sm" />
          <Avatar initials="МА" />
          <Avatar initials="МА" size="lg" ring />
        </DemoRow>
        <CodeCaption>.av-sm · .av-md · .av-lg · .av-ring</CodeCaption>
        <div className="kit-ornament-gap" />
        <DemoRow className="kit-user-pill-row">
          <UserPill initials="МА" name="Мария А." sign="Скорпион · ♀ Весы" />
        </DemoRow>
        <CodeCaption>.upill</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="nav"
        number="08"
        title="Навигация"
        description="Горизонтальное меню с золотым подчёркиванием на ховере и активном пункте."
      >
        <div className="nav-demo">
          <NavBar items={navItems} active={nav} onSelect={(item) => setNav(item.id)} />
        </div>
        <CodeCaption>.navbar › a.on</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="ornaments"
        number="09"
        title="Орнаменты & Разделители"
        description="Астральные разделители: звезда, солнце, фазы луны."
      >
        <div className="ornament-demo">
          <Ornament />
          <Ornament variant="cross" />
          <MoonPhases className="kit-moon-phases" />
          <CodeCaption>.orn (divider) · фазы луны</CodeCaption>
          <Divider>Развёрнутая интерпретация</Divider>
          <div className="kit-divider-right">
            <span>Развёрнутая интерпретация</span>
          </div>
          <CodeCaption>.divlbl (по центру) · .divlbl-r (метка справа)</CodeCaption>
        </div>
      </UiKitSection>
    </>
  );
}
