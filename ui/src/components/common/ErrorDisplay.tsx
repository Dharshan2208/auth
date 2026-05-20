interface ErrorDisplayProps {
  message: string;
  onDismiss?: () => void;
}

export function ErrorDisplay({ message, onDismiss }: ErrorDisplayProps) {
  return (
    <div className="flex items-start gap-3 bg-red-50 border-2 border-red-500 p-4 shadow-[4px_4px_0px_0px_rgba(239,68,68,1)]">
      <div className="flex-1">
        <p className="text-sm font-semibold text-red-700">Error</p>
        <p className="text-sm text-red-600 mt-0.5">{message}</p>
      </div>
      {onDismiss && (
        <button
          onClick={onDismiss}
          className="text-red-500 hover:text-red-700 font-bold text-lg leading-none"
        >
          ×
        </button>
      )}
    </div>
  );
}
