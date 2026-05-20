interface SuccessDisplayProps {
  message: string;
  onDismiss?: () => void;
}

export function SuccessDisplay({ message, onDismiss }: SuccessDisplayProps) {
  return (
    <div className="flex items-start gap-3 bg-green-50 border-2 border-green-500 p-4 shadow-[4px_4px_0px_0px_rgba(34,197,94,1)]">
      <div className="flex-1">
        <p className="text-sm font-semibold text-green-700">Success</p>
        <p className="text-sm text-green-600 mt-0.5">{message}</p>
      </div>
      {onDismiss && (
        <button
          onClick={onDismiss}
          className="text-green-500 hover:text-green-700 font-bold text-lg leading-none"
        >
          ×
        </button>
      )}
    </div>
  );
}
