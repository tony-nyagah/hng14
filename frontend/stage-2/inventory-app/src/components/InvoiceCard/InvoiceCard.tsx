import { Link } from 'react-router-dom';
import { format } from 'date-fns';
import { ChevronRight } from 'lucide-react';
import type { Invoice } from '../../types';
import StatusBadge from '../ui/StatusBadge/StatusBadge';
import './InvoiceCard.css';

interface InvoiceCardProps {
    invoice: Invoice;
}

export default function InvoiceCard({ invoice }: InvoiceCardProps) {
    // Format the date to something like "19 Aug 2026"
    const formattedDate = format(new Date(invoice.paymentDue), 'dd MMM yyyy');

    // Format total as GBP currency
    const formattedTotal = new Intl.NumberFormat('en-GB', {
        style: 'currency',
        currency: 'GBP',
    }).format(invoice.total);

    return (
        <Link to={`/invoice/${invoice.id}`} className="invoice-card">
            <span className="card-id"><span className="hash">#</span>{invoice.id}</span>
            <span className="card-date">Due {formattedDate}</span>
            <span className="card-name">{invoice.clientName}</span>
            <span className="card-total">{formattedTotal}</span>

            <div className="card-status-wrapper">
                <StatusBadge status={invoice.status} />
                {/* The chevron arrow only shows on desktop */}
                <ChevronRight className="card-arrow" size={20} color="#7C5DFA" />
            </div>
        </Link>
    );
}
