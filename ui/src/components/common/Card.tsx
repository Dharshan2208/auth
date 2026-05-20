import type { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  hover?: boolean;
  accent?: 'yellow' | 'blue' | 'green' | 'pink' | 'none';
}

const accentBorders: Record<string, string> = {
  yellow: 'border-l-yellow-400',
  blue: 'border-l-blue-500',
  green: 'border-l-green-400',
  pink: 'border-l-pink-400',
  none: 'border-l-black',
};

export function Card({
  children,
  className = '',
  hover = false,
  accent = 'none',
}: CardProps) {
  return (
    <div
      className={`
        bg-white border-2 border-black
        shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]
        ${accentBorders[accent]}
        ${hover ? 'transition-all duration-200 hover:-translate-y-1 hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]' : ''}
        ${className}
      `}
    >
      {children}
    </div>
  );
}
