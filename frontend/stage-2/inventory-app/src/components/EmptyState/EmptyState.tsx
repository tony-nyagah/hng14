import './EmptyState.css';

export default function EmptyState() {
    return (
        <div className="empty-state">
            <div className="empty-state-img">
                <img src="/illustration-empty.svg" alt="" onError={(e) => e.currentTarget.style.display = 'none'} />
            </div>
            <h2>There is nothing here</h2>
            <p>Create an invoice by clicking the <strong>New Invoice</strong> button and get started</p>
        </div>
    );
}
