import { useState } from 'react';
import { LogoSigil, MoonPhases, StarIcon, ZodiacWheel } from '@/shared/icons';
import { Button } from '@/shared/ui/Button';
import { AspectBlock, CornerPanel, Tooltip } from '@/shared/ui/Content';
import { Logo } from '@/shared/ui/Logo';
import { Modal } from '@/shared/ui/Modal';
import { TarotCard } from '@/shared/ui/TarotCard';
import { DemoRow, UiKitSection } from './UiKitSection';

export function ShowcaseSections() {
  const [modalOpen, setModalOpen] = useState(false);
  return <><UiKitSection id="tarot" number="10" title="Tarot card" description="Та же карта, что используется на Home; здесь включён hover tilt."><div className="tarot-demo"><TarotCard roman="II" name="The High Priestess" tilt/></div></UiKitSection><UiKitSection id="aspects" number="11" title="Aspect blocks"><div className="aspect-demo"><AspectBlock type="love" title="Любовь">Скрытые чувства требуют тишины.</AspectBlock><AspectBlock type="career" title="Карьера">Сначала соберите все данные.</AspectBlock><AspectBlock type="health" title="Здоровье">Прислушайтесь к телу.</AspectBlock><AspectBlock type="spirit" title="Духовность">Сновидения открывают путь.</AspectBlock></div></UiKitSection><UiKitSection id="panels" number="12" title="Corner-ornament panel"><CornerPanel title="Послание дня"><p>Ответ появляется там, где заканчивается шум.</p><p><em>Доверьтесь своей интуиции.</em></p></CornerPanel></UiKitSection><UiKitSection id="tooltips" number="13" title="Tooltips & Hints"><DemoRow><Tooltip label="Тихая подсказка"><button className="info-dot">i</button></Tooltip><Tooltip below label={<><b>Старший Аркан</b><br/>Карта глубинной интуиции.</>}><Button variant="outline">Rich popover</Button></Tooltip></DemoRow></UiKitSection><UiKitSection id="modal" number="14" title="Modal"><Button onClick={() => setModalOpen(true)}>Открыть модальное окно</Button><Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Полный расклад дня" footer={<Button variant="outline" onClick={() => setModalOpen(false)}>Принять послание</Button>}><p className="modal-sample">Вселенная приглашает вас замолчать и прислушаться.</p></Modal></UiKitSection><UiKitSection id="icons" number="15" title="Icons & Decor"><div className="icon-grid"><LogoSigil/><StarIcon/><MoonPhases/><ZodiacWheel/></div></UiKitSection><UiKitSection id="logo" number="16" title="Logo"><div className="logo-grid"><div><Logo size="lg" tagline="Карты & Судьба"/></div><div><Logo size="md" tagline="Карты & Судьба"/></div><div><Logo size="sm"/></div></div></UiKitSection></>;
}
