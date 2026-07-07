import { Badge, Chip } from '@/shared/ui/Chip';
import { Button } from '@/shared/ui/Button';
import { Field, TextareaField } from '@/shared/ui/Field';
import { StarIcon } from '@/shared/icons';
import { CodeCaption, DemoRow, UiKitSection } from './UiKitSection';

const colors = [
  ['Gold', '#C9A96E', '#C9A96E'],
  ['Gold Light', '#E8D5A3', '#E8D5A3'],
  ['Gold Dim', '#8A6F47', '#8A6F47'],
  ['Ink', '#0A0A0F', '#0A0A0F'],
  ['Ink Mid', '#13131C', '#13131C'],
  ['Ink Soft', '#1C1C2A', '#1C1C2A'],
  ['Smoke', '#2A2A3D', '#2A2A3D'],
  ['Mist', '#E8E2D9', '#E8E2D9'],
  ['Mist Dim', '#9B9189', '#9B9189'],
  ['Cream', '#F5F0E8', '#F5F0E8'],
];

export function FoundationSections() {
  return (
    <>
      <UiKitSection
        id="foundations"
        number="01"
        title="Основы — Палитра"
        description="Базовые токены цвета. Золотая ось акцентов поверх многослойных чернильных поверхностей."
      >
        <div className="swatch-grid">
          {colors.map(([name, hex, value]) => (
            <div className="swatch" key={name}>
              <i style={{ background: value }} />
              <b>{name}</b>
              <code>{hex}</code>
            </div>
          ))}
        </div>
      </UiKitSection>
      <UiKitSection
        id="typography"
        number="02"
        title="Типографика"
        description="Три гарнитуры: Cinzel — заголовки и метки, Cormorant Garamond — основной курсивный текст, Raleway — интерфейс."
      >
        <div className="type-row">
          <small>Cinzel · 500 · letter-spacing .08em</small>
          <b className="type-cinzel title">Верховная Жрица</b>
        </div>
        <div className="type-row">
          <small>Cinzel · 300 · 20px</small>
          <b className="type-cinzel light">Свет и тьма</b>
        </div>
        <div className="type-row">
          <small>Cinzel · 600 · .38em · uppercase</small>
          <b className="type-cinzel label">Карта дня</b>
        </div>
        <div className="type-row">
          <small>Cormorant Garamond · italic 300 · 24px</small>
          <b className="type-corm italic">Ответы, которые вы ищете, живут не снаружи.</b>
        </div>
        <div className="type-row">
          <small>Cormorant Garamond · 300 · 20px</small>
          <b className="type-corm body">
            Полумесяц у её ног напоминает: всё движется циклами, и то, что скрыто, выйдет на свет.
          </b>
        </div>
        <div className="type-row">
          <small>Raleway · 400 · .25em · uppercase</small>
          <b className="type-raleway">Развёрнутая интерпретация</b>
        </div>
      </UiKitSection>
      <UiKitSection
        id="buttons"
        number="03"
        title="Кнопки"
        description="Золотая первичная, контурная, призрачная — плюс компактный и иконочный варианты."
      >
        <DemoRow>
          <Button>Вытянуть карту</Button>
          <Button variant="outline">Подробнее</Button>
          <Button variant="ghost">Пропустить</Button>
          <Button variant="outline">
            <StarIcon width="14" height="14" />
            Сохранить расклад
          </Button>
        </DemoRow>
        <CodeCaption>.btn-primary · .btn-outline · .btn-ghost · .btn-ic</CodeCaption>
        <DemoRow>
          <Button size="sm">Малая</Button>
          <Button size="sm" variant="outline">
            Малая контурная
          </Button>
        </DemoRow>
        <CodeCaption>+ .btn-sm</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="chips"
        number="04"
        title="Чипы & Бейджи"
        description="Ключевые слова карты и статусные метки."
      >
        <DemoRow>
          <Chip>Интуиция</Chip>
          <Chip>Тайное знание</Chip>
          <Chip solid>Внутренний голос</Chip>
          <Chip>Цикличность</Chip>
        </DemoRow>
        <CodeCaption>.chip · .chip-solid</CodeCaption>
        <div className="kit-ornament-gap" />
        <DemoRow>
          <Badge dot>Прямое положение</Badge>
          <Badge variant="soft">Старший Аркан</Badge>
          <Badge variant="dim">Перевёрнутая</Badge>
        </DemoRow>
        <CodeCaption>.badge-gold · .badge-soft · .badge-dim</CodeCaption>
      </UiKitSection>
      <UiKitSection
        id="forms"
        number="05"
        title="Поля ввода"
        description="Текстовые поля с засечным вводом и золотым фокус-кольцом."
      >
        <div className="field-demo">
          <Field label="Ваше имя" placeholder="Мария А." />
          <Field label="Дата рождения" placeholder="28 мая 1994" />
          <TextareaField
            label="Вопрос к картам"
            placeholder="О чём вы хотите спросить колоду сегодня?"
          />
        </div>
        <CodeCaption>.field › input / textarea</CodeCaption>
      </UiKitSection>
    </>
  );
}
