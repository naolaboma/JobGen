import { createApi,fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const ContactApi = createApi({
    reducerPath: 'contactApi',
    baseQuery: fetchBaseQuery({ baseUrl: 'http://localhost:8080/api/v1' }),
    endpoints: (builder) => ({
        
        // submit a contact form
        submitContactForm: builder.mutation({
            query: (formData) => ({
                url: '/contact',
                method: 'POST',
                body: {
                name: formData.name,
                email: formData.email,
                subject: formData.subject,
                message: formData.message,
                },
            })
        }),


   }),
});

export const { useSubmitContactFormMutation } = ContactApi;

export default ContactApi;