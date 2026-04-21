import { useState, useRef, useEffect } from 'react';
import { Plus, ChevronDown } from 'lucide-react';
import { useInvoices } from '../../context/InvoiceContext';
import InvoiceCard from '../../components/InvoiceCard/InvoiceCard';
import Button from '../../components/ui/Button/Button';
import EmptyState from '../../components/EmptyState/EmptyState';
import './InvoiceListPage.css';
import InvoiceForm from '../../components/InvoiceForm/InvoiceForm';

export default function InvoiceListPage() {
    const { invoices } = useInvoices();
    const [filterOpen, setFilterOpen] = useState(false);
    const [isFormOpen, setIsFormOpen] = useState(false);
    const filterRef = useRef<HTMLDivElement>(null);

    // This array holds the selected statuses ('draft', 'pending', 'paid')
    const [filters, setFilters] = useState<string[]>([]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (filterRef.current && !filterRef.current.contains(event.target as Node)) {
                setFilterOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    // Toggle filter on/off
    const handleFilterChange = (status: string) => {
        if (filters.includes(status)) {
            setFilters(filters.filter(f => f !== status));
        } else {
            setFilters([...filters, status]);
        }
    };

    // If no filters are selected, show all. Otherwise, filter them.
    const filteredInvoices = filters.length > 0
        ? invoices.filter(inv => filters.includes(inv.status))
        : invoices;

    return (
        <div className="page-container">
            <header className="page-header">
                <div className="header-info">
                    <h1>Invoices</h1>
                    <p>
                        <span className="hide-mobile">{filteredInvoices.length === 1 ? 'There is' : 'There are'} </span>
                        {filteredInvoices.length}
                        <span className="hide-mobile"> {filteredInvoices.length === 1 ? 'total invoice' : 'total invoices'}</span>
                        <span className="show-mobile-only" style={{display: 'none'}}> invoices</span>
                    </p>
                </div>

                <div className="header-actions">
                    {/* Wrap the filter toggle and dropdown in a container */}
                    <div className="filter-dropdown-container" ref={filterRef}>
                        <div className="filter-wrapper" onClick={() => setFilterOpen(!filterOpen)}>
                            <span className="filter-text">Filter <span className="hide-mobile">by status</span></span>
                            <ChevronDown
                                size={14}
                                color="#7C5DFA"
                                style={{ transform: filterOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s' }}
                            />
                        </div>

                        {/* The Dropdown Menu */}
                        {filterOpen && (
                            <div className="filter-dropdown">
                                {['draft', 'pending', 'paid'].map(status => (
                                    <label key={status} className="filter-option">
                                        <input
                                            type="checkbox"
                                            checked={filters.includes(status)}
                                            onChange={() => handleFilterChange(status)}
                                        />
                                        {/* Custom Checkbox UI */}
                                        <span className="checkmark"></span>
                                        <span className="status-label">{status.charAt(0).toUpperCase() + status.slice(1)}</span>
                                    </label>
                                ))}
                            </div>
                        )}
                    </div>

                    <Button icon={<Plus size={16} color="#7C5DFA" strokeWidth={3} />} onClick={() => setIsFormOpen(true)}>
                        <span className="hide-mobile">New</span> Invoice
                    </Button>
                </div>
            </header>

            <div className="invoice-list">
                {filteredInvoices.length === 0 ? (
                    <EmptyState />
                ) : (
                    filteredInvoices.map(invoice => (
                        <InvoiceCard key={invoice.id} invoice={invoice} />
                    ))
                )}
            </div>

            <InvoiceForm
                isOpen={isFormOpen}
                onClose={() => setIsFormOpen(false)}
            />
        </div>
    );
}
