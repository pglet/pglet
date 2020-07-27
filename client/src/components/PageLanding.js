import React, { useEffect, useContext } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import Page from './Page'
import { WebSocketContext } from '../WebSocket';
import { changeProps, removeChild, deleteNode } from '../features/page/pageSlice'

const PageLanding = () => {

    const root = useSelector(state => state.page.controls[0]);

    const ws = useContext(WebSocketContext);
  
    var dispatch = useDispatch();
  
    useEffect(() => {
  
      ws.registerWebClient("account1/page1");
  
      setTimeout(() => {
        console.log("change control prop!");
        dispatch(changeProps({
          nodeId: "myTxt",
          newProps: {
            "text": "Another text!"
          }
        }));
        dispatch(removeChild({
          nodeId: "1",
          childId: "3"
        }))
        dispatch(deleteNode("3"))
      }, 4000)
    })

    return <Page control={root} />
}

export default PageLanding