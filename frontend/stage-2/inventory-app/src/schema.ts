import { z } from 'zod';

const addressSchema = z.object({
    street: z.string().min(1, "can't be empty"),
    city: z.string().min(1, "can't be empty"),
    postCode: z.string().min(1, "can't be empty"),
    country: z.string().min(1, "can't be empty"),
});

const itemSchema = z.object({
    name: z.string().min(1, "can't be empty"),
    quantity: z.coerce.number().min(1, "must be at least 1"),
    price: z.coerce.number().min(0.01, "must be positive"),
    total: z.number() // Calculated automatically
});

export const invoiceSchema = z.object({
    senderAddress: addressSchema,
    clientName: z.string().min(1, "can't be empty"),
    clientEmail: z.string().email("invalid email"),
    clientAddress: addressSchema,
    createdAt: z.string().min(1, "can't be empty"), // YYYY-MM-DD
    paymentTerms: z.coerce.number(),
    description: z.string().min(1, "can't be empty"),
    items: z.array(itemSchema).min(1, "An item must be added"),
    status: z.enum(['draft', 'pending', 'paid']).default('pending')
});

export type InvoiceFormData = z.infer<typeof invoiceSchema>;
