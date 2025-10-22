import { configureStore } from "@reduxjs/toolkit";
import contactReducer from "@/lib/redux/slices/contactSlice";
import jobReducer from "@/lib/redux/slices/jobSlice";

import ContactApi from "@/lib/redux/api/ContactApi";
import JobApi from "@/lib/redux/api/JobApi";
import { profileApi } from "@/lib/redux/slices/ProfileApiSlice";

const store = configureStore({
  reducer: {
    contact: contactReducer,
    job: jobReducer,
    [ContactApi.reducerPath]: ContactApi.reducer,
    [JobApi.reducerPath]: JobApi.reducer,
    [profileApi.reducerPath]: profileApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // RTK Query stores non-serializable cache entries under these slices
        ignoredPaths: [
          ContactApi.reducerPath,
          JobApi.reducerPath,
          profileApi.reducerPath,
        ],
        // Ignore RTK Query action types that may carry non-serializable payloads
        ignoredActions: [
          `${ContactApi.reducerPath}/executeMutation/pending`,
          `${ContactApi.reducerPath}/executeMutation/fulfilled`,
          `${ContactApi.reducerPath}/executeMutation/rejected`,
          `${JobApi.reducerPath}/executeQuery/pending`,
          `${JobApi.reducerPath}/executeQuery/fulfilled`,
          `${JobApi.reducerPath}/executeQuery/rejected`,
          `${profileApi.reducerPath}/executeQuery/pending`,
          `${profileApi.reducerPath}/executeQuery/fulfilled`,
          `${profileApi.reducerPath}/executeQuery/rejected`,
        ],
      },
    }).concat(ContactApi.middleware, JobApi.middleware, profileApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
