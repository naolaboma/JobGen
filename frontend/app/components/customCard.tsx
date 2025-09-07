import type { ReactNode } from "react";

interface CardProps {
  children: ReactNode;
  className?: string;
}

interface CardHeaderProps {
  children: ReactNode;
  className?: string;
}

interface CardContentProps {
  children: ReactNode;
  className?: string;
}

export function Card({ children, className = "" }: CardProps) {
  return (
    <div
      className={`bg-white border border-gray-200 rounded-xl overflow-hidden ${className}`}
    >
      {children}
    </div>
  );
}

export function CardHeader({ children, className = "" }: CardHeaderProps) {
  return <div className={className}>{children}</div>;
}

export function CardContent({ children, className = "" }: CardContentProps) {
  return <div className={className}>{children}</div>;
}
