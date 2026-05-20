import type { ReactNode } from 'react';

interface BadgeProps {
  children: ReactNode;
  variant?: 'yellow' | 'blue' | 'green' | 'pink' | 'black';
}

const variantStyles: Record<string, string> = {
  yellow: 'bg-yellow-400 text-black',
  blue: 'bg-blue-500 text-white',
  green: 'bg-green-300 text-black',
  pink: 'bg-pink-300 text-black',
  black: 'bg-black text-white',
};

export function Badge({ children, variant = 'yellow' }: BadgeProps) {
  return (
    <span
      className={`
        inline-flex items-center gap-1.5 px-3 py-1
        text-xs font-bold uppercase tracking-wider
        border-2 border-black
        ${variantStyles[variant]}
      `}
    >
      {children}
    </span>
  );
}
