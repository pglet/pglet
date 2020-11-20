import { createSlice } from '@reduxjs/toolkit'

const userSlice = createSlice({
    name: 'users',
    initialState: {
        details: null,
        loading: false,
        error: null
      },
    reducers: {
        fetchUserBegin(state, action) {
            state.loading = true;
            state.error = null;
        },
        fetchUserSuccess(state, action) {
            state.loading = false;
            state.details = action.payload;
        },
        fetchUserError(state, action) {
            state.loading = false;
            state.error = action.payload;
        }
    }
})

const fetchUserBegin = userSlice.actions.fetchUserBegin
const fetchUserSuccess = userSlice.actions.fetchUserSuccess
const fetchUserError = userSlice.actions.fetchUserError

export const fetchUser = (userId) => dispatch => {
    console.log("fetchUser()");
    dispatch(fetchUserBegin());
    fetch(`/api/users/${userId}`)
        .then(handleErrors)
        .then(res => res.json())
        .then(json => {
            console.log(json);
            dispatch(fetchUserSuccess(json));
            return json;
        })
        .catch(error => { dispatch(fetchUserError(error)) });
}

// Handle HTTP errors since fetch won't.
function handleErrors(response) {
    console.log(response);
    if (!response.ok) {
        throw Error(response.statusText);
    }
    return response;
}

export default userSlice.reducer