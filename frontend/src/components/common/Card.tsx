import { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  hover?: boolean;
  gradient?: boolean;
}

export default function Card({ children, className = '', hover = false, gradient = false }: CardProps) {
  const baseClasses = 'card';
  const hoverClasses = hover ? 'hover:shadow-lg transition-shadow duration-200' : '';
  const gradientClasses = gradient
    ? 'bg-linear-to-br from-purple-50 to-blue-50 border-2 border-dashed border-purple-300'
    : '';

  return (
    <div className={`${baseClasses} ${hoverClasses} ${gradientClasses} ${className}`.trim()}>
      {children}
    </div>
  );
}

