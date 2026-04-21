import { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ChevronLeft } from 'lucide-react';
import { useInvoices } from '../../context/InvoiceContext';
import { format } from 'date-fns';
import Button from '../../components/ui/Button/Button';
import StatusBadge from '../../components/ui/StatusBadge/StatusBadge';
import DeleteConfirmModal from '../../components/DeleteConfirmModal/DeleteConfirmModal';
import './InvoiceDetailPage.css';
import InvoiceForm from '../../components/InvoiceForm/InvoiceForm';


export default function InvoiceDetailPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { invoices, markAsPaid, deleteInvoice } = useInvoices();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isFormOpen, setIsFormOpen] = useState(false);


    // Find the specific invoice from our data based on the URL ID
    const invoice = invoices.find(inv => inv.id === id);

    if (!invoice) {
        return (
            <div className="page-container">
                <Link to="/" className="go-back"><ChevronLeft size={16} color="#7C5DFA" strokeWidth={3} /> Go back</Link>
                <p style={{ marginTop: '24px' }}>Invoice not found!</p>
            </div>
        );
    }

    const handleDelete = () => {
        deleteInvoice(invoice.id);
        navigate('/');
    };

    const handleMarkAsPaid = () => {
        markAsPaid(invoice.id);
    };

    // Helper functions for clean formatting
    const fmtDate = (dateStr: string) => format(new Date(dateStr), 'dd MMM yyyy');
    const fmtMoney = (num: number) => new Intl.NumberFormat('en-GB', { style: 'currency', currency: 'GBP' }).format(num);

    return (
        <div className="page-container">
            <Link to="/" className="go-back">
                <ChevronLeft size={16} color="#7C5DFA" strokeWidth={3} />
                Go back
            </Link>

            {/* Header section (Status and Desktop Actions) */}
            <div className="detail-header">
                <div className="status-section">
                    <span className="status-label-text">Status</span>
                    <StatusBadge status={invoice.status} />
                </div>

                <div className="action-section hide-mobile">
                    <Button variant="edit" onClick={() => setIsFormOpen(true)}>Edit</Button>
                    <Button variant="danger" onClick={() => setIsModalOpen(true)}>Delete</Button>
                    {invoice.status !== 'paid' && (
                        <Button variant="primary" onClick={handleMarkAsPaid}>Mark as Paid</Button>
                    )}
                </div>
            </div>

            {/* Main Invoice Card Body */}
            <div className="detail-body">
                <div className="detail-top">
                    <div className="id-desc">
                        <h2 className="invoice-id"><span className="hash">#</span>{invoice.id}</h2>
                        <p className="invoice-desc">{invoice.description}</p>
                    </div>
                    <div className="sender-address">
                        <p>{invoice.senderAddress.street}</p>
                        <p>{invoice.senderAddress.city}</p>
                        <p>{invoice.senderAddress.postCode}</p>
                        <p>{invoice.senderAddress.country}</p>
                    </div>
                </div>

                <div className="detail-middle">
                    <div className="dates-section">
                        <div className="date-block">
                            <p className="field-label">Invoice Date</p>
                            <h3 className="field-value">{fmtDate(invoice.createdAt)}</h3>
                        </div>
                        <div className="date-block">
                            <p className="field-label">Payment Due</p>
                            <h3 className="field-value">{fmtDate(invoice.paymentDue)}</h3>
                        </div>
                    </div>

                    <div className="bill-to-section">
                        <p className="field-label">Bill To</p>
                        <h3 className="field-value">{invoice.clientName}</h3>
                        <div className="client-address">
                            <p>{invoice.clientAddress.street}</p>
                            <p>{invoice.clientAddress.city}</p>
                            <p>{invoice.clientAddress.postCode}</p>
                            <p>{invoice.clientAddress.country}</p>
                        </div>
                    </div>

                    <div className="sent-to-section">
                        <p className="field-label">Sent to</p>
                        <h3 className="field-value">{invoice.clientEmail}</h3>
                    </div>
                </div>

                {/* Items Table area */}
                <div className="detail-items">
                    <table className="items-table">
                        <thead>
                            <tr>
                                <th className="th-name">Item Name</th>
                                <th className="th-qty">QTY.</th>
                                <th className="th-price">Price</th>
                                <th className="th-total">Total</th>
                            </tr>
                        </thead>
                        <tbody>
                            {invoice.items.map((item, idx) => (
                                <tr key={idx}>
                                    <td className="item-name">{item.name}</td>
                                    <td className="item-qty">{item.quantity}</td>
                                    <td className="item-price">{fmtMoney(item.price)}</td>
                                    <td className="item-total">{fmtMoney(item.total)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    <div className="items-footer">
                        <p>Amount Due</p>
                        <h2>{fmtMoney(invoice.total)}</h2>
                    </div>
                </div>
            </div>

            {/* Mobile action bar - sticks to the bottom on mobile */}
            <div className="action-section show-mobile">
                <Button variant="edit" onClick={() => setIsFormOpen(true)}>Edit</Button>
                <Button variant="danger" onClick={() => setIsModalOpen(true)}>Delete</Button>
                {invoice.status !== 'paid' && (
                    <Button variant="primary" onClick={handleMarkAsPaid}>Mark as Paid</Button>
                )}
            </div>

            {isModalOpen && (
                <DeleteConfirmModal
                    invoiceId={invoice.id}
                    onCancel={() => setIsModalOpen(false)}
                    onConfirm={handleDelete}
                />
            )}

            <InvoiceForm
                isOpen={isFormOpen}
                onClose={() => setIsFormOpen(false)}
                invoiceToEdit={invoice}
            />

        </div>
    );
}
