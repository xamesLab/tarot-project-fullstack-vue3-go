import { forwardRef, type InputHTMLAttributes, type TextareaHTMLAttributes } from 'react';

interface FieldProps extends InputHTMLAttributes<HTMLInputElement> { label: string; error?: string; }
export const Field = forwardRef<HTMLInputElement, FieldProps>(({ label, error, id, ...props }, ref) => { const fieldId = id ?? props.name; return <label className="field" htmlFor={fieldId}><span className="field-label">{label}</span><input ref={ref} id={fieldId} className="field-control" aria-invalid={Boolean(error)} aria-describedby={error ? `${fieldId}-error` : undefined} {...props}/>{error && <span id={`${fieldId}-error`} className="field-error">{error}</span>}</label>; });
Field.displayName = 'Field';

export function TextareaField({ label, id, ...props }: TextareaHTMLAttributes<HTMLTextAreaElement> & { label: string }) { return <label className="field" htmlFor={id}><span className="field-label">{label}</span><textarea id={id} className="field-control" {...props}/></label>; }
