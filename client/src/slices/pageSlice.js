import { createSlice } from '@reduxjs/toolkit'
import Cookies from 'universal-cookie';
//import { current } from '@reduxjs/toolkit'

const cookies = new Cookies();

const initialState = {
    "error": null,
    "controls": {
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
            cookies.set(`sid-${action.payload.pageName}#${action.payload.pageHash}`, action.payload.session.id, { path: '/', sameSite: true });
        },
        registerWebClientError(state, action) {
            state.loading = false;
            state.error = action.payload.error;
            state.signinOptions = action.payload.signinOptions;
        },
        appBecomeInactive(state, action) {
            state.error = action.payload.message;
            //cookies.remove(`sid-${action.payload.pageName}#${action.payload.pageHash}`);
        },
        sessionCrashed(state, action) {
            state.error = action.payload.message;
            //cookies.remove(`sid-${action.payload.pageName}#${action.payload.pageHash}`);
        },
        signout(state, action) {
            var redirectUrl = encodeURIComponent(window.location.pathname);
            window.location.replace("/api/auth/signout?redirect_url=" + redirectUrl);
        },         
        addPageControlsSuccess(state, action) {
            const { controls, trimIDs } = action.payload
            addControls(state, controls);
            removeControls(state, trimIDs);
        },
        addPageControlsError(state, action) {
            state.error = action.payload;
        },
        replacePageControlsSuccess(state, action) {
            const { ids, remove, controls } = action.payload

            // clean or remove controls
            if (remove) {
                removeControls(state, ids);
            } else {
                cleanControls(state, ids);
            }

            // add controls
            addControls(state, controls);
        },        
        replacePageControlsError(state, action) {
            state.error = action.payload;
        },
        pageControlsBatchSuccess(state, action) {
            action.payload.forEach(message => {
                
                //console.log(message);

                if (message.action === 'addPageControls') {
                    const { controls, trimIDs } = message.payload
                    addControls(state, controls);
                    removeControls(state, trimIDs);
                } else if (message.action === 'updateControlProps') {
                    changePropsInternal(state, message.payload.props)
                } else if (message.action === 'cleanControl') {
                    cleanControls(state, message.payload.ids)
                } else if (message.action === 'removeControl') {
                    removeControls(state, message.payload.ids)
                }
            })
        },        
        pageControlsBatchError(state, action) {
            state.error = action.payload;
        },        
        changeProps(state, action) {
            changePropsInternal(state, action.payload)
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
            cleanControls(state, ids)
        },        
        removeControl(state, action) {
            const { ids } = action.payload
            removeControls(state, ids)
        }
    }
})

const changePropsInternal = (state, allProps) => {
    allProps.forEach(props => {
        const ctrl = state.controls[props.i];
        if (ctrl) {
            Object.assign(ctrl, props)
        }
    })    
}

const addControls = (state, controls) => {
    let firstParentId = null;
    controls.forEach(ctrl => {
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
}

const cleanControls = (state, ids) => {
    ids.forEach(id => {
        // remove all children
        const descendantIds = getAllDescendantIds(state.controls, id)
        descendantIds.forEach(descId => delete state.controls[descId])

        // cleanup children collection
        state.controls[id].c = []
    })
}

const removeControls = (state, ids) => {
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
    appBecomeInactive,
    sessionCrashed,
    signout,
    addPageControlsSuccess,
    addPageControlsError,
    replacePageControlsSuccess,
    replacePageControlsError,
    pageControlsBatchSuccess,
    pageControlsBatchError,
    changeProps,
    appendProps,
    cleanControl,
    removeControl
} = pageSlice.actions

export default pageSlice.reducer