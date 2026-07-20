import { useEffect, useRef, useState, type PointerEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '@/shared/ui/Button';
import { Modal } from '@/shared/ui/Modal';

const CROP_WIDTH = 320;
const CROP_HEIGHT = 480;
const OUTPUT_WIDTH = 800;
const OUTPUT_HEIGHT = 1200;

interface ImageState {
  naturalWidth: number;
  naturalHeight: number;
  minScale: number;
  scale: number;
  offsetX: number;
  offsetY: number;
}

const initialImageState: ImageState = { naturalWidth: 0, naturalHeight: 0, minScale: 1, scale: 1, offsetX: 0, offsetY: 0 };

export function CardBackCropModal({ open, imageUrl, fileName, onClose, onConfirm }: {
  open: boolean;
  imageUrl: string;
  fileName: string;
  onClose: () => void;
  onConfirm: (file: File) => void;
}) {
  const { t } = useTranslation('admin');
  const imageRef = useRef<HTMLImageElement>(null);
  const dragRef = useRef({ active: false, x: 0, y: 0, offsetX: 0, offsetY: 0 });
  const [image, setImage] = useState<ImageState>(initialImageState);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) {
      dragRef.current.active = false;
      setSaving(false);
      setImage(initialImageState);
    }
  }, [open]);

  const clamp = (next: ImageState): ImageState => {
    const scale = Math.min(Math.max(next.scale, next.minScale), next.minScale * 3);
    const minX = Math.min(0, CROP_WIDTH - next.naturalWidth * scale);
    const minY = Math.min(0, CROP_HEIGHT - next.naturalHeight * scale);
    return { ...next, scale, offsetX: Math.min(Math.max(next.offsetX, minX), 0), offsetY: Math.min(Math.max(next.offsetY, minY), 0) };
  };

  const handleLoad = () => {
    const element = imageRef.current;
    if (!element) return;
    const minScale = Math.max(CROP_WIDTH / element.naturalWidth, CROP_HEIGHT / element.naturalHeight);
    setImage(clamp({ naturalWidth: element.naturalWidth, naturalHeight: element.naturalHeight, minScale, scale: minScale, offsetX: (CROP_WIDTH - element.naturalWidth * minScale) / 2, offsetY: (CROP_HEIGHT - element.naturalHeight * minScale) / 2 }));
  };

  const handlePointerDown = (event: PointerEvent<HTMLImageElement>) => {
    event.currentTarget.setPointerCapture(event.pointerId);
    dragRef.current = { active: true, x: event.clientX, y: event.clientY, offsetX: image.offsetX, offsetY: image.offsetY };
  };

  const handlePointerMove = (event: PointerEvent<HTMLDivElement>) => {
    if (!dragRef.current.active) return;
    setImage((current) => clamp({ ...current, offsetX: dragRef.current.offsetX + event.clientX - dragRef.current.x, offsetY: dragRef.current.offsetY + event.clientY - dragRef.current.y }));
  };

  const confirm = async () => {
    const element = imageRef.current;
    if (!element || !image.naturalWidth) return;
    setSaving(true);
    const canvas = document.createElement('canvas');
    canvas.width = OUTPUT_WIDTH;
    canvas.height = OUTPUT_HEIGHT;
    const context = canvas.getContext('2d');
    if (!context) { setSaving(false); return; }
    context.drawImage(element, -image.offsetX / image.scale, -image.offsetY / image.scale, CROP_WIDTH / image.scale, CROP_HEIGHT / image.scale, 0, 0, OUTPUT_WIDTH, OUTPUT_HEIGHT);
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/png'));
    setSaving(false);
    if (!blob) return;
    const baseName = fileName.replace(/\.[^.]+$/, '') || 'card-back';
    onConfirm(new File([blob], `${baseName}-cropped.png`, { type: 'image/png' }));
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={t('assets.crop.title')}
      panelClassName="admin-crop-modal"
      footer={<><Button type="button" variant="ghost" onClick={onClose}>{t('actions.cancel')}</Button><Button type="button" disabled={saving || !image.naturalWidth} onClick={() => void confirm()}>{saving ? t('actions.saving') : t('assets.crop.confirm')}</Button></>}
    >
      <div className="admin-crop">
        <p>{t('assets.crop.copy')}</p>
        <div className="admin-crop-stage" onPointerMove={handlePointerMove} onPointerUp={() => { dragRef.current.active = false; }} onPointerCancel={() => { dragRef.current.active = false; }}>
          {imageUrl && <img ref={imageRef} src={imageUrl} alt={fileName} draggable={false} onLoad={handleLoad} onPointerDown={handlePointerDown} style={{ width: image.naturalWidth * image.scale, height: image.naturalHeight * image.scale, transform: `translate(${image.offsetX}px, ${image.offsetY}px)` }} />}
          <span aria-hidden="true" />
        </div>
        <label className="admin-crop-zoom">
          <span>{t('assets.crop.zoom')}</span>
          <input type="range" min={image.minScale} max={image.minScale * 3} step="0.01" value={image.scale} disabled={!image.naturalWidth} onChange={(event) => setImage((current) => clamp({ ...current, scale: Number(event.target.value) }))} />
        </label>
      </div>
    </Modal>
  );
}
