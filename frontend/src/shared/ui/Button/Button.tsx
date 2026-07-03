import type { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'outline' | 'ghost';
  size?: 'md' | 'sm';
  iconOnly?: boolean;
  children: ReactNode;
}

export function Button({ variant = 'primary', size = 'md', iconOnly, className = '', children, ...props }: ButtonProps) {
  return <button className={`button button-${variant} button-${size} ${iconOnly ? 'button-icon' : ''} ${className}`} {...props}>{children}</button>;
}
