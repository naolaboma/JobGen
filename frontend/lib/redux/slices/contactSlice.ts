import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ContactState {
    lastSubmission: IFormInput | null; 
    isLoading: boolean;
    error: string | null;
    isSubmitted: boolean;
}

const initialState: ContactState = {
    lastSubmission: null, // Corrected typo
    isLoading: false,
    error: null,
    isSubmitted: false,
};

const contactSlice = createSlice({
    name: "contact",
    initialState,
    reducers: {
        submitContactFormStart: (state) => {
            state.isLoading = true;
            state.error = null;
            state.isSubmitted = false;
        },

        submitContactFormSuccess: (state, action: PayloadAction<IFormInput>) => {
            state.isLoading = false;
            state.isSubmitted = true;
            state.lastSubmission = action.payload; // Corrected typo
        },

        submitContactFormFailure: (state, action: PayloadAction<string>) => {
            state.isLoading = false;
            state.error = action.payload;
            state.isSubmitted = false;
        },

        clearContactFormStatus: (state) => {
            state.isLoading = false;
            state.error = null;
            state.isSubmitted = false;
            state.lastSubmission = null; // Corrected typo
        },
    },
});

export const {
    submitContactFormStart,
    submitContactFormSuccess,
    submitContactFormFailure,
    clearContactFormStatus,
} = contactSlice.actions;

export default contactSlice.reducer;