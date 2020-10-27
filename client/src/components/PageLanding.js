import React, { useEffect, useContext } from 'react';
import { useParams } from "react-router-dom";
import { useSelector, useDispatch } from 'react-redux';
import Page from './Page'
import { WebSocketContext } from '../WebSocket';
import { changeProps, removeChild, deleteNode } from '../features/page/pageSlice'

const PageLanding = () => {

    let { accountName, pageName } = useParams();

    let fullPageName = `${accountName}/${pageName}`;

    const ws = useContext(WebSocketContext);

    var dispatch = useDispatch();


    useEffect(() => {

        ws.registerWebClient(fullPageName);

        // setTimeout(() => {
        //     console.log("change control prop!");
        //     dispatch(changeProps({
        //         nodeId: "myTxt",
        //         newProps: {
        //             "text": "Another text!"
        //         }
        //     }));
        //     dispatch(removeChild({
        //         nodeId: "1",
        //         childId: "3"
        //     }))
        //     dispatch(deleteNode("3"))
        // }, 4000)
    }, [])

    const err = useSelector(state => state.page.error);
    const root = useSelector(state => state.page.controls['0']);

    if (err) {
        return <div>{err}</div>;
    } else {
        return <Page control={root} />
    }
}

export default PageLanding