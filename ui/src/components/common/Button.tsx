import type { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  children: ReactNode;
  loading?: boolean;
}

const variantStyles: Record<string, string> = {
  primary:
    'bg-blue-500 text-white border-2 border-black hover:bg-blue-600 active:bg-blue-700',
  secondary:
    'bg-yellow-400 text-black border-2 border-black hover:bg-yellow-300 active:bg-yellow-500',
  outline:
    'bg-white text-black border-2 border-black hover:bg-gray-100 active:bg-gray-200',
  ghost:
    'bg-transparent text-black border-2 border-transparent hover:bg-gray-100',
  danger:
    'bg-red-500 text-white border-2 border-black hover:bg-red-600 active:bg-red-700',
};

export function Button({
  variant = 'primary',
  size = 'md',
  children,
  loading,
  disabled,
  className = '',
  ...props
}: ButtonProps) {
  const sizeStyles = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-5 py-2.5 text-base',
    lg: 'px-8 py-3.5 text-lg',
  };

  return (
    <button
      className={`
        inline-flex items-center justify-center gap-2 font-semibold
        transition-all duration-150 ease-in-out
        disabled:opacity-50 disabled:cursor-not-allowed
        active:translate-y-0.5
        ${variantStyles[variant]}
        ${sizeStyles[size]}
        ${className}
      `}
      disabled={disabled || loading}
      {...props}
    >
      {loading && (
        <svg
          className="animate-spin h-4 w-4"
          viewBox="0 0 24 24"
          fill="none"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
          />
        </svg>
      )}
      {children}
    </button>
  );
}
