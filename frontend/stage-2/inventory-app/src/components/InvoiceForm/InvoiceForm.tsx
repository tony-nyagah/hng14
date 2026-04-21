import { useEffect } from 'react';
import { useForm, useFieldArray, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { v4 as uuidv4 } from 'uuid';
import { Trash } from 'lucide-react';
import { useInvoices } from '../../context/InvoiceContext';
import { invoiceSchema } from '../../schema';
import type { InvoiceFormData } from '../../schema';
import Button from '../ui/Button/Button';
import './InvoiceForm.css';

interface InvoiceFormProps {
    isOpen: boolean;
    onClose: () => void;
    invoiceToEdit?: any; // The invoice object if editing
}

export default function InvoiceForm({ isOpen, onClose, invoiceToEdit }: InvoiceFormProps) {
    const { addInvoice, updateInvoice } = useInvoices();

    // Set up React Hook Form with Zod validation
    const { register, control, handleSubmit, reset, getValues, formState: { errors } } = useForm<InvoiceFormData>({
        resolver: zodResolver(invoiceSchema),
        defaultValues: {
            senderAddress: { street: '', city: '', postCode: '', country: '' },
            clientName: '',
            clientEmail: '',
            clientAddress: { street: '', city: '', postCode: '', country: '' },
            createdAt: new Date().toISOString().split('T')[0],
            paymentTerms: 30,
            description: '',
            items: [],
            status: 'pending'
        }
    });

    // Manage dynamic invoice items array
    const { fields, append, remove } = useFieldArray({
        control,
        name: 'items'
    });

    // Watch items to dynamically calculate totals
    const watchedItems = useWatch({ control, name: 'items' }) || [];

    // Load data if editing, or clear if opening a new form
    useEffect(() => {
        if (invoiceToEdit) {
            reset(invoiceToEdit);
        } else {
            reset();
        }
    }, [invoiceToEdit, isOpen, reset]);

    const createFinalInvoice = (data: any, status: 'draft' | 'pending' | 'paid' = 'pending') => {
        const safeItems = Array.isArray(data.items) ? data.items : [];
        const invoiceTotal = safeItems.reduce((acc: number, item: any) => acc + ((item.quantity || 0) * (item.price || 0)), 0);
        const formattedItems = safeItems.map((item: any) => ({
            ...item,
            total: (item.quantity || 0) * (item.price || 0)
        }));

        const createdAtDate = data.createdAt ? new Date(data.createdAt) : new Date();
        const paymentTerms = parseInt(data.paymentTerms) || 30;
        const paymentDue = new Date(createdAtDate.getTime() + paymentTerms * 24 * 60 * 60 * 1000).toISOString().split('T')[0];

        return {
            ...data,
            id: invoiceToEdit ? invoiceToEdit.id : uuidv4().slice(0, 6).toUpperCase(),
            items: formattedItems,
            total: invoiceTotal,
            paymentDue,
            status: status
        };
    };

    // Handle successful validation
    const onSubmit = (data: InvoiceFormData) => {
        const finalInvoice = createFinalInvoice(data, invoiceToEdit ? invoiceToEdit.status : 'pending');
        if (invoiceToEdit) updateInvoice(finalInvoice);
        else addInvoice(finalInvoice);
        onClose();
    };

    const handleSaveAsDraft = () => {
        // Bypass strict validation for drafts
        const currentData = getValues();
        const finalInvoice = createFinalInvoice(currentData, 'draft');
        
        if (invoiceToEdit) updateInvoice(finalInvoice);
        else addInvoice(finalInvoice);
        
        onClose();
    };

    return (
        <>
            <div className={`form-overlay ${isOpen ? 'open' : ''}`} onClick={onClose} />

            <div className={`form-drawer ${isOpen ? 'open' : ''}`}>
                <div className="form-header">
                    <h2>{invoiceToEdit ? `Edit #${invoiceToEdit.id}` : 'New Invoice'}</h2>
                </div>

                <div className="form-content">
                    <form id="invoice-form" onSubmit={handleSubmit(onSubmit)}>

                        {/* Bill From */}
                        <section>
                            <h4 className="section-title">Bill From</h4>
                            <div className="form-group">
                                <label>Street Address</label>
                                <input type="text" {...register('senderAddress.street')} />
                                {errors.senderAddress?.street && <span className="error">{errors.senderAddress.street.message}</span>}
                            </div>

                            <div className="form-row-3">
                                <div className="form-group">
                                    <label>City</label>
                                    <input type="text" {...register('senderAddress.city')} />
                                    {errors.senderAddress?.city && <span className="error">{errors.senderAddress.city.message}</span>}
                                </div>
                                <div className="form-group">
                                    <label>Post Code</label>
                                    <input type="text" {...register('senderAddress.postCode')} />
                                    {errors.senderAddress?.postCode && <span className="error">{errors.senderAddress.postCode.message}</span>}
                                </div>
                                <div className="form-group">
                                    <label>Country</label>
                                    <input type="text" {...register('senderAddress.country')} />
                                    {errors.senderAddress?.country && <span className="error">{errors.senderAddress.country.message}</span>}
                                </div>
                            </div>
                        </section>

                        {/* Bill To */}
                        <section>
                            <h4 className="section-title">Bill To</h4>
                            <div className="form-group">
                                <label>Client's Name</label>
                                <input type="text" {...register('clientName')} />
                                {errors.clientName && <span className="error">{errors.clientName.message}</span>}
                            </div>
                            <div className="form-group">
                                <label>Client's Email</label>
                                <input type="email" placeholder="e.g. email@example.com" {...register('clientEmail')} />
                                {errors.clientEmail && <span className="error">{errors.clientEmail.message}</span>}
                            </div>
                            <div className="form-group">
                                <label>Street Address</label>
                                <input type="text" {...register('clientAddress.street')} />
                                {errors.clientAddress?.street && <span className="error">{errors.clientAddress.street.message}</span>}
                            </div>

                            <div className="form-row-3">
                                <div className="form-group">
                                    <label>City</label>
                                    <input type="text" {...register('clientAddress.city')} />
                                    {errors.clientAddress?.city && <span className="error">{errors.clientAddress.city.message}</span>}
                                </div>
                                <div className="form-group">
                                    <label>Post Code</label>
                                    <input type="text" {...register('clientAddress.postCode')} />
                                    {errors.clientAddress?.postCode && <span className="error">{errors.clientAddress.postCode.message}</span>}
                                </div>
                                <div className="form-group">
                                    <label>Country</label>
                                    <input type="text" {...register('clientAddress.country')} />
                                    {errors.clientAddress?.country && <span className="error">{errors.clientAddress.country.message}</span>}
                                </div>
                            </div>
                        </section>

                        {/* Invoice Details */}
                        <section>
                            <div className="form-row-2">
                                <div className="form-group">
                                    <label>Invoice Date</label>
                                    <input type="date" {...register('createdAt')} />
                                    {errors.createdAt && <span className="error">{errors.createdAt.message}</span>}
                                </div>
                                <div className="form-group">
                                    <label>Payment Terms</label>
                                    <select {...register('paymentTerms')}>
                                        <option value={1}>Net 1 Day</option>
                                        <option value={7}>Net 7 Days</option>
                                        <option value={14}>Net 14 Days</option>
                                        <option value={30}>Net 30 Days</option>
                                    </select>
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Project Description</label>
                                <input type="text" placeholder="e.g. Graphic Design Service" {...register('description')} />
                                {errors.description && <span className="error">{errors.description.message}</span>}
                            </div>
                        </section>

                        {/* Item List */}
                        <section className="item-list-section">
                            <h3 className="item-list-title">Item List</h3>
                            {errors.items && !Array.isArray(errors.items) && <span className="error" style={{ position: 'relative', display: 'block', marginBottom: '16px' }}>{errors.items.message}</span>}

                            {fields.map((field, index) => {
                                const qty = watchedItems[index]?.quantity || 0;
                                const price = watchedItems[index]?.price || 0;
                                const total = (qty * price).toFixed(2);
                                
                                return (
                                    <div key={field.id} className="item-row">
                                        <div className="form-group item-name">
                                            <label className="hide-desktop">Item Name</label>
                                            <input type="text" {...register(`items.${index}.name`)} />
                                            {errors.items?.[index]?.name && <span className="error">{errors.items[index]?.name?.message}</span>}
                                        </div>
                                        <div className="form-group item-qty">
                                            <label className="hide-desktop">Qty.</label>
                                            <input type="number" {...register(`items.${index}.quantity`)} />
                                        </div>
                                        <div className="form-group item-price">
                                            <label className="hide-desktop">Price</label>
                                            <input type="number" step="0.01" {...register(`items.${index}.price`)} />
                                        </div>
                                        <div className="form-group item-total">
                                            <label className="hide-desktop">Total</label>
                                            <div className="total-display">{total}</div>
                                        </div>
                                        <button type="button" className="delete-item-btn" onClick={() => remove(index)}>
                                            <Trash size={20} />
                                        </button>
                                    </div>
                                );
                            })}

                            <Button
                                type="button"
                                variant="secondary"
                                className="add-item-btn"
                                onClick={() => append({ name: '', quantity: 1, price: 0, total: 0 })}
                            >
                                + Add New Item
                            </Button>
                        </section>

                    </form>
                </div>

                <div className="form-actions">
                    <Button variant="edit" onClick={onClose}>Discard</Button>
                    <Button variant="draft" onClick={handleSaveAsDraft}>Save as Draft</Button>
                    <Button variant="primary" type="submit" form="invoice-form">{invoiceToEdit ? 'Save Changes' : 'Save & Send'}</Button>
                </div>
            </div>
        </>
    );
}
