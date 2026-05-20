interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  text?: string;
}

const sizeStyles = {
  sm: 'h-4 w-4 border-2',
  md: 'h-8 w-8 border-3',
  lg: 'h-12 w-12 border-4',
};

export function LoadingSpinner({ size = 'md', text }: LoadingSpinnerProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-8">
      <div
        className={`
          ${sizeStyles[size]}
          rounded-full border-black border-t-transparent
          animate-spin
        `}
      />
      {text && <p className="text-sm font-medium text-gray-500">{text}</p>}
    </div>
  );
}
