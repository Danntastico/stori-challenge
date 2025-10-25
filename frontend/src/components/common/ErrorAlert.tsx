interface ErrorAlertProps {
  message: string;
  onRetry?: () => void;
  className?: string;
}

export default function ErrorAlert({ message, onRetry, className = '' }: ErrorAlertProps) {
  return (
    <div className={`card bg-red-50 border border-red-200 ${className}`}>
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-3">
          <span className="text-2xl">⚠️</span>
          <div>
            <p className="text-red-800 font-medium">{message}</p>
            {onRetry && (
              <button
                onClick={onRetry}
                className="mt-2 text-sm text-red-700 hover:text-red-900 underline"
              >
                Try again
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
