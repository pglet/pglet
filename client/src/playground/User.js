import React from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { fetchUser } from '../slices/userSlice'

const User = ({userId}) => {

    const user = useSelector(state => state.user.details);
    const loading = useSelector(state => state.user.loading);
    const error = useSelector(state => state.user.error);

    //console.log(user);

    const dispatch = useDispatch();

    React.useEffect(() => {
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