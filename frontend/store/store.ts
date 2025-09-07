import { configureStore } from "@reduxjs/toolkit";
import contactReducer from "@/lib/redux/slices/contactSlice";
import ContactApi from "@/lib/redux/api/ContactApi"; 
import { profileApi } from "@/lib/redux/slices/ProfileApiSlice";

const store = configureStore({
  reducer: {
    contact: contactReducer,
    [ContactApi.reducerPath]: ContactApi.reducer,
    [profileApi.reducerPath]: profileApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(ContactApi.middleware).concat(profileApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;