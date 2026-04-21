import './StatusBadge.css';

interface StatusBadgeProps {
    status: 'paid' | 'pending' | 'draft';
}

export default function StatusBadge({ status }: StatusBadgeProps) {
    return (
        <div className={`status-badge status-${status}`}>
            <span className="status-dot"></span>
            <span className="status-text">{status.charAt(0).toUpperCase() + status.slice(1)}</span>
        </div>
    );
}
