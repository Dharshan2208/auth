import type { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export function Input({ label, error, className = '', ...props }: InputProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-semibold text-black mb-1.5">
          {label}
        </label>
      )}
      <input
        className={`
          w-full px-4 py-2.5 bg-white border-2 border-black
          text-black placeholder:text-gray-400
          transition-all duration-150
          focus:outline-none focus:-translate-y-0.5 focus:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]
          ${error ? 'border-red-500 focus:shadow-red-500/50' : 'border-black'}
          ${className}
        `}
        {...props}
      />
      {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
    </div>
  );
}
