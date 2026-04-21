import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import type { Invoice } from '../types';
import { dummyInvoices } from '../data';

interface InvoiceContextType {
    invoices: Invoice[];
    markAsPaid: (id: string) => void;
    deleteInvoice: (id: string) => void;
    addInvoice: (invoice: Invoice) => void;
    updateInvoice: (invoice: Invoice) => void;
}

const InvoiceContext = createContext<InvoiceContextType | undefined>(undefined);

export function InvoiceProvider({ children }: { children: ReactNode }) {
    const [invoices, setInvoices] = useState<Invoice[]>(() => {
        const saved = localStorage.getItem('invoices');
        if (saved) return JSON.parse(saved);
        return dummyInvoices;
    });

    useEffect(() => {
        localStorage.setItem('invoices', JSON.stringify(invoices));
    }, [invoices]);

    const markAsPaid = (id: string) => {
        setInvoices(invoices.map(inv =>
            inv.id === id ? { ...inv, status: 'paid' } : inv
        ));
    };

    const deleteInvoice = (id: string) => {
        setInvoices(invoices.filter(inv => inv.id !== id));
    };

    const addInvoice = (invoice: Invoice) => {
        setInvoices([...invoices, invoice]);
    };

    const updateInvoice = (updatedInvoice: Invoice) => {
        setInvoices(invoices.map(inv => inv.id === updatedInvoice.id ? updatedInvoice : inv));
    };

    return (
        <InvoiceContext.Provider value={{ invoices, markAsPaid, deleteInvoice, addInvoice, updateInvoice }}>
            {children}
        </InvoiceContext.Provider>
    );
}

export const useInvoices = () => {
    const context = useContext(InvoiceContext);
    if (context === undefined) {
        throw new Error('useInvoices must be used within an InvoiceProvider');
    }
    return context;
};
