import type { ReactNode, ButtonHTMLAttributes } from 'react';
import './Button.css';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    children: ReactNode;
    variant?: 'primary' | 'secondary' | 'danger' | 'draft' | 'edit';
    icon?: ReactNode;
}

export default function Button({
    children,
    variant = 'primary',
    icon,
    className = '',
    ...props
}: ButtonProps) {
    return (
        <button
            className={`btn btn-${variant} ${className}`}
            {...props}
        >
            {/* If an icon is provided (like the plus sign for new invoices), render it in the circle */}
            {icon && <span className="btn-icon">{icon}</span>}
            {children}
        </button>
    );
}
