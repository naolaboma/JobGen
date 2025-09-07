import { configureStore } from "@reduxjs/toolkit";
import contactReducer from "@/lib/redux/slices/contactSlice";
import ContactApi from "@/lib/redux/api/ContactApi"; 

const store = configureStore({
  reducer: {
    contact: contactReducer,
    [ContactApi.reducerPath]: ContactApi.reducer, 
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(ContactApi.middleware), 
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;