import { createSlice } from '@reduxjs/toolkit'
import Cookies from 'universal-cookie';
//import { current } from '@reduxjs/toolkit'

const cookies = new Cookies();

const initialState = {
    "name": "test-1",
    "error": null,
    "controls": {
        "page": {
            "c": [],
            "i": "page",
            "p": "",
            "t": "page"
        }
    }
}

const pageSlice = createSlice({
    name: 'page',
    initialState: initialState,
    reducers: {
        registerWebClientSuccess(state, action) {
            state.loading = false;
            state.sessionId = action.payload.session.id;
            state.controls = action.payload.session.controls;
            cookies.set(`sid-${action.payload.pageName}`, action.payload.session.id, { path: '/' });
        },
        registerWebClientError(state, action) {
            state.loading = false;
            state.error = action.payload;
        },
        addPageControlsSuccess(state, action) {
            let firstParentId = null;
            action.payload.forEach(ctrl => {
                if (firstParentId == null) {
                    firstParentId = ctrl.p;
                }

                if (!state.controls[ctrl.i]) {
                    state.controls[ctrl.i] = ctrl;

                    if (ctrl.p === firstParentId) {
                        // root control
                        if (typeof ctrl.at === 'undefined') {
                            // append to the end
                            state.controls[ctrl.p].c.push(ctrl.i)
                        } else {
                            // insert at specified position
                            state.controls[ctrl.p].c.splice(ctrl.at, 0, ctrl.i)
                        }
                    }
                }
            })
            //console.log("After addPageControlsSuccess:", current(state))
        },
        addPageControlsError(state, action) {
            state.error = action.payload;
        },        
        changeProps(state, action) {

            action.payload.forEach(props => {
                const ctrl = state.controls[props.i];
                if (ctrl) {
                    Object.assign(ctrl, props)
                }
            })
            //console.log(current(state))
        },
        appendProps(state, action) {

            action.payload.forEach(props => {
                const ctrl = state.controls[props.i];
                if (ctrl) {
                    Object.getOwnPropertyNames(props).forEach(propName => {
                        if (propName === 'i') {
                            return
                        }
                        let v = ctrl[propName]
                        if (!v) {
                            v = ""
                        }
                        ctrl[propName] = v + props[propName]
                    })
                }
            })
            //console.log(current(state))
        },  
        cleanControl(state, action) {
            const { ids } = action.payload

            ids.forEach(id => {
                // remove all children
                const descendantIds = getAllDescendantIds(state.controls, id)
                descendantIds.forEach(descId => delete state.controls[descId])

                // cleanup children collection
                state.controls[id].c = []
            })
        },        
        removeControl(state, action) {
            const { ids } = action.payload

            ids.forEach(id => {
                const ctrl = state.controls[id]

                // remove all children
                const descendantIds = getAllDescendantIds(state.controls, id)
                descendantIds.forEach(descId => delete state.controls[descId])
    
                // delete control itself
                delete state.controls[id]
    
                // remove ID from parent's children collection
                const parent = state.controls[ctrl.p]
                parent.c = parent.c.filter(childId => childId !== id)  
            })          
        }
    }
})

const getAllDescendantIds = (controls, nodeId) => {
    if (controls[nodeId].c) {
        return controls[nodeId].c.reduce((acc, childId) => (
            [...acc, childId, ...getAllDescendantIds(controls, childId)]
        ), [])
    }
    return []
}

export const {
    registerWebClientSuccess,
    registerWebClientError,
    addPageControlsSuccess,
    addPageControlsError,
    changeProps,
    appendProps,
    cleanControl,
    removeControl
} = pageSlice.actions

export default pageSlice.reducer