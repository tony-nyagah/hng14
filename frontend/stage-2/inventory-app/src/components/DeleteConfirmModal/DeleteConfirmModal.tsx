import Button from '../ui/Button/Button';
import './DeleteConfirmModal.css';

interface DeleteConfirmModalProps {
    invoiceId: string;
    onCancel: () => void;
    onConfirm: () => void;
}

export default function DeleteConfirmModal({ invoiceId, onCancel, onConfirm }: DeleteConfirmModalProps) {
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>Confirm Deletion</h2>
                <p>Are you sure you want to delete invoice #{invoiceId}? This action cannot be undone.</p>
                <div className="modal-actions">
                    <Button variant="edit" onClick={onCancel}>Cancel</Button>
                    <Button variant="danger" onClick={onConfirm}>Delete</Button>
                </div>
            </div>
        </div>
    );
}