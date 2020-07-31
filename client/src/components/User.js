import React, { useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { fetchUser } from '../features/users/userSlice'

const User = ({userId}) => {

    const user = useSelector(state => state.user.details);
    const loading = useSelector(state => state.user.loading);
    const error = useSelector(state => state.user.error);

    //console.log(user);

    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(fetchUser(userId));
    }, [userId, dispatch]);

    if (error) {
        return <div>Error! {error.message}</div>;
    }

    if (loading || user == null) {
        return <div>Loading...</div>;
    }

    return <div>User: {user.username}</div>;
}

export default User