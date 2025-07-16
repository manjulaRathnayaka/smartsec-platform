import React from 'react';

const LoadingSpinner = ({ size = 'md', color = 'primary' }) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12'
  };

  const colorClasses = {
    primary: 'text-primary-600',
    secondary: 'text-secondary-600',
    white: 'text-white',
    gray: 'text-gray-600'
  };

  return (
    <div className="flex justify-center items-center">
      <div
        className={`animate-spin rounded-full border-2 border-gray-200 border-t-current ${sizeClasses[size]} ${colorClasses[color]}`}
      />
    </div>
  );
};

const LoadingCard = ({ title = 'Loading...', description }) => {
  return (
    <div className="bg-white rounded-lg shadow p-6 text-center">
      <LoadingSpinner size="lg" />
      <h3 className="mt-4 text-lg font-medium text-gray-900">{title}</h3>
      {description && (
        <p className="mt-2 text-sm text-gray-600">{description}</p>
      )}
    </div>
  );
};

const LoadingPage = ({ title = 'Loading...', description = 'Please wait while we load your data.' }) => {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center">
        <LoadingSpinner size="xl" />
        <h2 className="mt-6 text-2xl font-bold text-gray-900">{title}</h2>
        <p className="mt-2 text-gray-600">{description}</p>
      </div>
    </div>
  );
};

const LoadingTable = ({ rows = 5, columns = 4 }) => {
  return (
    <div className="animate-pulse">
      <div className="space-y-3">
        {Array(rows).fill(0).map((_, i) => (
          <div key={i} className="flex space-x-4">
            {Array(columns).fill(0).map((_, j) => (
              <div key={j} className="flex-1 h-4 bg-gray-300 rounded"></div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

const LoadingButton = ({ children, loading = false, disabled = false, ...props }) => {
  return (
    <button
      {...props}
      disabled={loading || disabled}
      className={`${props.className || ''} ${loading || disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
    >
      {loading ? (
        <span className="flex items-center justify-center">
          <LoadingSpinner size="sm" color="white" />
          <span className="ml-2">Loading...</span>
        </span>
      ) : (
        children
      )}
    </button>
  );
};

export {
  LoadingSpinner,
  LoadingCard,
  LoadingPage,
  LoadingTable,
  LoadingButton
};

export default LoadingSpinner;
